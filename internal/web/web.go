package web

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/bnadland/rhizome/internal/assets"
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
			slog.Warn(err.Error(), "GetPageBySlug", slug)
			notFound(w, req)
			return
		}

		if err := views.Page(page).Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := views.NotFound().Render(req.Context(), w); err != nil {
		slog.Error(err.Error())
	}
}

func notFoundHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		notFound(w, req)
	}
}

func NewRouter(q *db.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "text/html", "text/css", "text/javascript"))

	r.NotFound(notFoundHandlerFunc())

	r.Get("/p/{slug}", page(q))
	r.Handle("/static/*", assets.Assets())
	r.Handle("/", http.RedirectHandler("/p/home", http.StatusMovedPermanently))
	return r
}

func NewServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: router,

		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
}

func Run(ctx context.Context, addr string, databaseURL string) error {
	if err := db.Migrate(databaseURL); err != nil {
		return err
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	router := NewRouter(db.New(pool))
	server := NewServer(addr, router)

	go func() {
		slog.Info("listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error(err.Error())
		}
	}()
	wg.Wait()
	slog.Info("shutdown")
	return nil
}
