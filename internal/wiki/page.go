package wiki

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/bnadland/rhizome/internal/assets"
	"github.com/bnadland/rhizome/internal/db"
	"github.com/go-chi/chi/v5"
)

type pager interface {
	GetPageBySlug(context.Context, string) (db.Page, error)
}

func PageHandler(q pager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")

		page, err := q.GetPageBySlug(req.Context(), slug)
		if err != nil {
			slog.Warn(err.Error(), "GetPageBySlug", slug)
			w.WriteHeader(http.StatusNotFound)
			if err := assets.NotFound().Render(req.Context(), w); err != nil {
				slog.Error(err.Error())
			}
			return
		}

		if err := Page(page).Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}
