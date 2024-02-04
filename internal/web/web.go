package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/bnadland/rhizome/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/graceful"
)

func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := views.NotFound().Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}

}

func page(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")

		page, err := q.GetPageBySlug(req.Context(), slug)
		if err != nil {
			slog.Warn(err.Error(), "GetPageBySlug", slug)
			notFound()(w, req)
			return
		}

		if err := views.Page(page).Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}

func GetRouter(q *db.Queries) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.NotFound(notFound())
	r.Get("/p/{slug}", page(q))
	r.Handle("/", http.RedirectHandler("/p/home", http.StatusMovedPermanently))
	return r
}

func Run(addr string) error {
	if err := db.Migrate(); err != nil {
		return err
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	server := graceful.WithDefaults(&http.Server{
		Addr:    addr,
		Handler: GetRouter(db.New(pool)),
	})

	slog.Info("listening", "addr", addr)
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		return err
	}
	slog.Info("shutdown gracefully")
	return nil
}
