package database

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (db *Database) ApplyMigrations(ctx context.Context) error {
	log.Info().Msg("Applying migrations...")

	return nil
}
