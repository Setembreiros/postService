package main

import (
	"context"
	"os"
	"os/signal"
	"postservice/cmd/provider"
	"postservice/infrastructure/kafka"
	"postservice/internal/api"
	"postservice/internal/bus"
	database "postservice/internal/db"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	ctx              context.Context
	cancel           context.CancelFunc
	configuringTasks sync.WaitGroup
	runningTasks     sync.WaitGroup
	env              string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))

	app := &app{
		ctx:    ctx,
		cancel: cancel,
		env:    env,
	}

	app.configuringLog()

	log.Info().Msgf("Starting PostService service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(env)
	database, err := provider.ProvideDb(ctx)
	if err != nil {
		os.Exit(1)
	}
	eventBus, err := provider.ProvideEventBus()
	if err != nil {
		os.Exit(1)
	}
	subscriptions := provider.ProvideSubscriptions()
	apiEnpoint := provider.ProvideApiEndpoint(database)
	kafkaConsumer, err := provider.ProvideKafkaConsumer(eventBus)
	if err != nil {
		os.Exit(1)
	}

	app.runConfigurationTasks(database, subscriptions, eventBus)
	app.runServerTasks(kafkaConsumer, apiEnpoint)
}

func (app *app) configuringLog() {
	if app.env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.With().Caller().Logger()
}

func (app *app) runConfigurationTasks(database *database.Database, subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	app.configuringTasks.Add(2)
	go app.applyMigrations(database)
	go app.subcribeEvents(subscriptions, eventBus) // Always subscribe event before init Kafka
	app.configuringTasks.Wait()
}

func (app *app) runServerTasks(kafkaConsumer *kafka.KafkaConsumer, apiEnpoint *api.Api) {
	app.runningTasks.Add(2)
	go app.initKafkaConsumption(kafkaConsumer)
	go app.runApiEndpoint(apiEnpoint)

	blockForever()

	app.shutdown()
}

func (app *app) applyMigrations(database *database.Database) {
	defer app.configuringTasks.Done()

	err := database.ApplyMigrations(app.ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Migrations failed")
	}
	log.Info().Msg("Migrations finished")
}

func (app *app) subcribeEvents(subscriptions *[]bus.EventSubscription, eventBus *bus.EventBus) {
	defer app.configuringTasks.Done()

	log.Info().Msg("Subscribing events...")

	for _, subscription := range *subscriptions {
		eventBus.Subscribe(&subscription, app.ctx)
		log.Info().Msgf("%s subscribed", subscription.EventType)
	}

	log.Info().Msg("All events subscribed")
}

func (app *app) initKafkaConsumption(kafkaConsumer *kafka.KafkaConsumer) {
	defer app.runningTasks.Done()

	err := kafkaConsumer.InitConsumption(app.ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Kafka Consumption failed")
	}
	log.Info().Msg("Kafka Consumer Group stopped")
}

func (app *app) runApiEndpoint(apiEnpoint *api.Api) {
	defer app.runningTasks.Done()

	err := apiEnpoint.Run(app.ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Closing PostService Api failed")
	}
	log.Info().Msg("PostService Api stopped")
}

func blockForever() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func (app *app) shutdown() {
	app.cancel()
	log.Info().Msg("Shutting down PostService Service...")
	app.runningTasks.Wait()
	log.Info().Msg("PostService Service stopped")
}
