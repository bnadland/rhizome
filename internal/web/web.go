package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/bnadland/rhizome/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func page(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")
		page, err := q.GetPageBySlug(req.Context(), slug)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			if err := views.NotFound().Render(req.Context(), w); err != nil {
				slog.Error(err.Error())
			}
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
	q := db.New(pool)
	r := GetRouter(q)
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	slog.Info("listening", "addr", addr)
	return s.ListenAndServe()
}
