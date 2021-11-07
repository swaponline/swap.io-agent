package queueEvents

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/config"
)

type QueueEvents struct {
	controllerConn *kafka.Conn
	kafkaWriter    *kafka.Writer
}

func InitializeQueueEvents() *QueueEvents {
	controllerConn, err := kafka.Dial("tcp", config.KAFKA_ADDR)
	if err != nil {
		log.Panicln(err)
	}
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(config.KAFKA_ADDR),
		Balancer: &kafka.LeastBytes{},
	}

	return &QueueEvents{
		controllerConn: controllerConn,
		kafkaWriter:    kafkaWriter,
	}
}

func (q *QueueEvents) WriteTxsEvents(data map[string][]*blockchain.Transaction) error {
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
				Value: bytes,
			})
		}
	}

	if len(kafkaMessages) == 0 {
		return nil
	}

	return q.kafkaWriter.WriteMessages(
		context.Background(),
		kafkaMessages...,
	)
}
func (q *QueueEvents) GetTxEventNotifier(
	ctx context.Context,
	agentUserId string,
) (
	<-chan blockchain.Transaction,
	chan<- struct{},
) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{config.KAFKA_ADDR},
		Topic:     agentUserId,
		GroupID:   "agentId",
		Partition: 0,
		MinBytes:  0,
		MaxBytes:  1e6, // 10MBit
	})

	notifier := make(chan blockchain.Transaction)
	isOk := make(chan struct{})

	go func() {
		defer close(notifier)
		defer close(isOk)

		for ctx.Err() == nil {
			m, err := r.FetchMessage(ctx)
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				log.Println(err)
				<-time.After(time.Second)
				continue
			}

			var tx blockchain.Transaction
			if err := json.Unmarshal(m.Value, &tx); err != nil {
				log.Panicln(err, string(m.Value))
			}

			select {
			case <-ctx.Done(): return
			case notifier <- tx:
			}

			select {
			case <-ctx.Done(): return
			case <-isOk:
			}

			for err := r.CommitMessages(context.Background(), m); err != nil; {
				log.Println("failed to commit messages:", err)
				<-time.After(time.Second)
			}
		}
	}()

	return notifier, isOk
}
func (q *QueueEvents) ReserveQueueForUser(agentUserId string) error {
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

func (*QueueEvents) Start() {}
func (*QueueEvents) Stop() error {
	return nil
}
func (*QueueEvents) Status() error {
	return nil
}
