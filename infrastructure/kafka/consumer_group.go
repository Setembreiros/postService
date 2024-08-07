package kafka

import (
	"context"
	"errors"
	"postservice/internal/bus"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaConsumer struct {
	ConsumerGroup sarama.ConsumerGroup
	eventBus      *bus.EventBus
}

func NewKafkaConsumer(brokers []string, eventBus *bus.EventBus) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	groupId := "postervice-group"

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupId, config)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error creating consumer group client")
		return nil, err
	}

	return &KafkaConsumer{
		ConsumerGroup: consumerGroup,
		eventBus:      eventBus,
	}, nil
}

func (k *KafkaConsumer) InitConsumption(ctx context.Context) error {
	consumer := Consumer{
		ready:    make(chan bool),
		eventBus: k.eventBus,
	}

	log.Info().Msg("Initiating Kafka Consumer Group...")

	var wg sync.WaitGroup
	wg.Add(1)
	go k.runConsumerGroup(ctx, &wg, &consumer)

	<-consumer.ready // Await till the consumer has been set up
	log.Info().Msg("Kafka Consumer up and running...")

	<-ctx.Done()
	log.Info().Msg("Terminating Kafka Consumer: context cancelled")

	wg.Wait()
	if err := k.ConsumerGroup.Close(); err != nil {
		log.Error().Stack().Err(err).Msg("Error closing Kafka Consumer Group")
		return err
	}

	return nil
}

func (k *KafkaConsumer) runConsumerGroup(ctx context.Context, wg *sync.WaitGroup, consumer *Consumer) {
	defer wg.Done()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.ConsumerGroup.Consume(ctx, getTopics(), consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				log.Error().Stack().Err(err).Msg("Consumer Group was closed")
				return
			}
			log.Panic().Stack().Err(err).Msg("Error from consumer")
		}
		// check if context was cancelled, signaling that the consumer should stop
		err := ctx.Err()
		if err != nil {
			return
		}
		consumer.ready = make(chan bool)
	}
}
