package provider

import (
	"context"
	awsClients "postservice/infrastructure/aws"
	"postservice/infrastructure/kafka"
	"postservice/internal/api"
	"postservice/internal/bus"
	database "postservice/internal/db"
	"postservice/internal/features/create_post"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/rs/zerolog/log"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
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

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
	brokers := p.kafkaBrokers()

	return kafka.NewKafkaConsumer(brokers, eventBus)
}

func (p *Provider) ProvideApiEndpoint(database *database.Database) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database))
}

func (p *Provider) ProvideApiControllers(database *database.Database) []api.Controller {
	return []api.Controller{
		create_post.NewCreatePostController(create_post.CreatePostRepository(*database)),
	}
}

func (p *Provider) ProvideDb(ctx context.Context) (*database.Database, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig(ctx)
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load aws configuration")
		return nil, err
	}

	return database.NewDatabase(awsClients.NewDynamodbClient(cfg)), nil
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

func provideAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-3"))
}

func provideDevEnvironmentDbConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
}
