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

	if !db.Client.IndexExists("Posts", "UserIndex") {
		indexes := []TableAttributes{
			{
				Name:          "User",
				AttributeType: "string",
			},
			{
				Name:          "CreatedAt",
				AttributeType: "string",
			},
		}
		db.Client.CreateIndexesOnTable("Posts", "UserIndex", &indexes, ctx)
	}

	if !db.Client.IndexExists("Posts", "TypeIndex") {
		indexes := []TableAttributes{
			{
				Name:          "Type",
				AttributeType: "string",
			},
			{
				Name:          "CreatedAt",
				AttributeType: "string",
			},
		}
		db.Client.CreateIndexesOnTable("Posts", "TypeIndex", &indexes, ctx)
	}

	return nil
}
