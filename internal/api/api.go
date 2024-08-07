package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Api struct {
	port        int
	env         string
	controllers []Controller
}

func NewApiEndpoint(env string, controllers []Controller) *Api {
	return &Api{
		port:        6666,
		env:         env,
		controllers: controllers,
	}
}

func (api *Api) Run(ctx context.Context) error {
	routes := api.routes()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", api.port),
		Handler:           routes,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Info().Msgf("Starting PostService Api Server on port %d", api.port)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("PostService Api Server failed")
		}
	}()

	<-ctx.Done()
	return server.Shutdown(ctx)
}
