package db

import (
	"embed"
	"io"
	"log/slog"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(databaseURL string) error {
	dsn, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}
	db := dbmate.New(dsn)
	db.AutoDumpSchema = false
	db.Log = io.Discard
	db.FS = migrations
	db.MigrationsDir = []string{"migrations"}
	slog.Info("applying migrations")
	return db.CreateAndMigrate()
}
