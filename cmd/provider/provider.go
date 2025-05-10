package provider

import (
	"context"
	awsClients "postservice/infrastructure/aws"
	"postservice/infrastructure/kafka"
	"postservice/internal/api"
	"postservice/internal/bus"
	database "postservice/internal/db"
	"postservice/internal/features/create_post"
	"postservice/internal/features/delete_post"
	"postservice/internal/features/get_post"
	objectstorage "postservice/internal/objectStorage"

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

func (p *Provider) ProvideApiEndpoint(database *database.Database, objectRepository *objectstorage.ObjectStorage, bus *bus.EventBus) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database, objectRepository, bus))
}

func (p *Provider) ProvideApiControllers(database *database.Database, objectRepository *objectstorage.ObjectStorage, bus *bus.EventBus) []api.Controller {
	return []api.Controller{
		create_post.NewCreatePostController(create_post.NewCreatePostService(create_post.NewCreatePostRepository(database, objectRepository), bus), bus),
		get_post.NewGetPostController(get_post.NewGetPostRepository(database, objectRepository)),
		delete_post.NewDeletePostController(delete_post.NewDeletePostRepository(database, objectRepository), bus),
	}
}

func (p *Provider) ProvideDb(ctx context.Context) (*database.Database, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig(ctx, "8000")
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load aws configuration")
		return nil, err
	}

	return database.NewDatabase(awsClients.NewDynamodbClient(cfg)), nil
}

func (p *Provider) ProvideObjectStorage(ctx context.Context) (*objectstorage.ObjectStorage, error) {
	var cfg aws.Config
	var err error

	if p.env == "development" {
		cfg, err = provideDevEnvironmentDbConfig(ctx, "4566")
	} else {
		cfg, err = provideAwsConfig(ctx)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load aws configuration")
		return nil, err
	}

	return objectstorage.NewObjectStorage(awsClients.NewS3Client(cfg, "artis-bucket")), nil
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.0.242:9092",
			"172.31.7.110:9092",
		}
	}
}

func provideAwsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-3"))
}

func provideDevEnvironmentDbConfig(ctx context.Context, port string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:" + port}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
}
