package queueEvents

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/env"
)

type QueueEvents struct {
	controllerConn *kafka.Conn
	kafkaWriter    *kafka.Writer
}

func InitializeQueueEvents() *QueueEvents {
	controllerConn, err := kafka.Dial("tcp", env.KAFKA_ADDR)
	if err != nil {
		log.Panicln(err)
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(env.KAFKA_ADDR),
		Balancer: &kafka.LeastBytes{},
	}

	return &QueueEvents{
		controllerConn: controllerConn,
		kafkaWriter:    kafkaWriter,
	}
}

func (q *QueueEvents) WriteTxsEvents(data map[string][]blockchain.Transaction) error {
	kafkaMessages := make([]kafka.Message, 0)
	for agentUserId, txs := range data {
		for _, tx := range txs {
			bytes, err := json.Marshal(tx)
			if err != nil {
				return err
			}
			kafkaMessages = append(kafkaMessages, kafka.Message{
				Topic: agentUserId,
				Key:   []byte(tx.Hash),
				Value: []byte(bytes),
			})
		}
	}

	return q.kafkaWriter.WriteMessages(
		context.Background(),
		kafkaMessages...,
	)
}
func (q *QueueEvents) GetTxEventNotifier(
	agentUserId string,
) (<-chan blockchain.Transaction, chan<- bool, context.CancelFunc) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{env.KAFKA_ADDR},
		Topic:     agentUserId,
		GroupID:   "agentId",
		Partition: 0,
		MinBytes:  0,
		MaxBytes:  10e6, // 10MBit
	})

	notifier := make(chan blockchain.Transaction)
	isOk := make(chan bool)

	stopCtx, stopFn := context.WithCancel(context.Background())
	go func() {
		for stopCtx.Err() == nil {
			m, err := r.FetchMessage(stopCtx)
			if err != nil {
				log.Println(err)
				<-time.After(time.Second)
				continue
			}

			var tx blockchain.Transaction
			if err := json.Unmarshal(m.Value, &tx); err != nil {
				log.Panicln(err, string(m.Value))
			}
			notifier <- tx
			<-isOk

			for err := r.CommitMessages(context.Background(), m); err != nil; {
				log.Println("failed to commit messages:", err)
				<-time.After(time.Second)
			}
		}
	}()
	go func() {
		<-stopCtx.Done()
		if err := r.Close(); err != nil {
			log.Fatal("failed to close connection:", err)
		}
	}()

	return notifier, isOk, stopFn
}
func (q *QueueEvents) ReseiveQueueForUser(agentUserId string) error {
	err := q.controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             agentUserId,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return err
	}
	return nil
}
