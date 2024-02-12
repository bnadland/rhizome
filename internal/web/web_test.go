package web_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/bnadland/rhizome/internal/db"
	"github.com/bnadland/rhizome/internal/web"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestHomepage(t *testing.T) {
	t.Parallel()
	r := web.NewRouter(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusFound, resp.Code)
}

func TestMain(m *testing.M) {
	databaseURL := os.Getenv("DATABASE_URL_TEST")
	if databaseURL != "" {
		dsn, err := url.Parse(databaseURL)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		testDB := dbmate.New(dsn)
		testDB.AutoDumpSchema = false
		testDB.FS = db.Migrations
		testDB.MigrationsDir = []string{"migrations"}
		if err := testDB.Drop(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if err := testDB.CreateAndMigrate(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		ctx := context.Background()
		conn, err := pgx.Connect(ctx, databaseURL)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		q := db.New(conn)

		if err = q.SavePage(ctx, db.SavePageParams{
			Title:   "Home",
			Slug:    "home",
			Content: "hello, world",
		}); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		conn.Close(ctx)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}
