package kafka

import (
	"errors"

	"postservice/internal/bus"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	ready    chan bool
	eventBus *bus.EventBus
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				err := errors.New("message channel was closed")
				log.Error().Err(err)
				return nil
			}
			log.Info().Msgf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

			event := bus.Event{
				Type: message.Topic,
				Data: message.Value,
			}
			consumer.eventBus.PublishLocal(event)
		case <-session.Context().Done():
			return nil
		}
	}
}
