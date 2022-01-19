package queueEvents

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"swap.io-agent/src/blockchain"
	"swap.io-agent/src/config"
)

type QueueEvents struct {
	p sarama.SyncProducer
}

func InitializeQueueEvents() *QueueEvents {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = -1

	p, err := sarama.NewSyncProducer([]string{config.KAFKA_ADDR}, kafkaConfig)
	if err != nil {
		log.Panicln(err)
	}

	return &QueueEvents{
		p,
	}
}

func (q *QueueEvents) WriteTxsEvents(data map[string][]*blockchain.Transaction) error {
	kafkaMessages := make([]*sarama.ProducerMessage, 0)
	for agentUserId, txs := range data {
		for _, tx := range txs {
			bytes, err := json.Marshal(tx)
			if err != nil {
				return err
			}
			kafkaMessages = append(kafkaMessages, &sarama.ProducerMessage{
				Topic: config.BLOCKCHAIN + "." + agentUserId,
				Key:   sarama.StringEncoder(tx.Hash),
				Value: sarama.ByteEncoder(bytes),
			})
		}
	}

	if len(kafkaMessages) == 0 {
		return nil
	}

	return q.p.SendMessages(kafkaMessages)
}

type consumerGroupHandler struct {
	ctx      context.Context
	notifier chan blockchain.Transaction
	isOk     chan struct{}
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		var tx blockchain.Transaction
		if err := json.Unmarshal(msg.Value, &tx); err != nil {
			log.Panicln(err, string(msg.Value))
		}

		select {
		case <-h.ctx.Done():
			return nil
		case h.notifier <- tx:
		}

		select {
		case <-h.ctx.Done():
			return nil
		case <-h.isOk:
		}

		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}
func (q *QueueEvents) GetTxEventNotifier(
	ctx context.Context,
	agentUserId string,
) (
	<-chan blockchain.Transaction,
	chan<- struct{},
) {
	consumerGroupConfig := sarama.NewConfig()
	consumerGroupConfig.Consumer.Return.Errors = true
	consumerGroupConfig.Consumer.Offsets.AutoCommit.Enable = false
	consumerGroupConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	groupConsumer, err := sarama.NewConsumerGroup(
		[]string{config.KAFKA_ADDR},
		"agentId",
		consumerGroupConfig,
	)
	if err != nil {
		log.Panicln(err)
	}
	go func() {
		for {
			select {
			case err, ok := <-groupConsumer.Errors():
				{
					if ok {
						log.Println("ERROR", err)
					} else {
						return
					}
				}
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}()

	notifier := make(chan blockchain.Transaction)
	isOk := make(chan struct{})

	go func() {
		defer close(notifier)
		defer close(isOk)

		handler := consumerGroupHandler{
			ctx:      ctx,
			notifier: notifier,
			isOk:     isOk,
		}

		for {
			err := groupConsumer.Consume(
				ctx,
				[]string{config.BLOCKCHAIN + "." + agentUserId},
				handler,
			)
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				log.Println(err)
				<-time.After(time.Second)
				continue
			}
		}
	}()

	return notifier, isOk
}
func (q *QueueEvents) ReserveQueueForUser(agentUserId string) error {
	admin, err := sarama.NewClusterAdmin([]string{config.KAFKA_ADDR}, nil)
	if err != nil {
		log.Println("Error while creating cluster admin: ", err.Error())
		return err
	}
	defer func() { _ = admin.Close() }()
	err = admin.CreateTopic(config.BLOCKCHAIN+"."+agentUserId, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil && !strings.HasPrefix(err.Error(), sarama.ErrTopicAlreadyExists.Error()) {
		log.Fatal("Error while creating topic: ", err.Error())
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
