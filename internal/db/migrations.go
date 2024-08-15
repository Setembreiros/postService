package database

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (db *Database) ApplyMigrations(ctx context.Context) error {
	log.Info().Msg("Applying migrations...")

	if !db.Client.TableExists("Posts") {
		keys := []TableAttributes{
			{
				Name:          "PostId",
				AttributeType: "string",
			},
		}
		err := db.Client.CreateTable("Posts", &keys, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
