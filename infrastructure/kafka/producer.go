package kafka

import (
	"postservice/internal/bus"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error creating consumer group client")
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
	}, nil
}

func (kp *KafkaProducer) Publish(event *bus.Event) error {
	msg := &sarama.ProducerMessage{
		Topic: event.Type,
		Value: sarama.StringEncoder(event.Data),
	}

	_, _, err := kp.Producer.SendMessage(msg)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error publishing")
		return err
	}

	log.Info().Msgf("Event %s published on Kafka", event.Type)

	return nil
}
