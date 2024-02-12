package web

import (
	"log/slog"
	"net/http"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/bnadland/rhizome/internal/views"
	"github.com/go-chi/chi/v5"
)

func PageHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")

		page, err := q.GetPageBySlug(req.Context(), slug)
		if err != nil {
			slog.Warn(err.Error(), "GetPageBySlug", slug)
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
