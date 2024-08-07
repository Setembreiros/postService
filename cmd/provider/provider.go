package provider

import (
	"postservice/infrastructure/kafka"
	"postservice/internal/api"
	"postservice/internal/bus"
)

type Provider struct {
	env     string
	connStr string
}

func NewProvider(env, connStr string) *Provider {
	return &Provider{
		env:     env,
		connStr: connStr,
	}
}

func (p *Provider) ProvideEventBus() (*bus.EventBus, error) {
	kafkaProducer, err := kafka.NewKafkaProducer(p.kafkaBrokers())
	if err != nil {
		return nil, err
	}

	return bus.NewEventBus(kafkaProducer), nil
}

func (p *Provider) ProvideSubscriptions() *[]bus.EventSubscription {
	return &[]bus.EventSubscription{
		{},
	}
}

func (p *Provider) ProvideApiEndpoint() *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers())
}

func (p *Provider) ProvideApiControllers() []api.Controller {
	return []api.Controller{}
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.36.175:9092",
			"172.31.45.255:9092",
		}
	}
}
