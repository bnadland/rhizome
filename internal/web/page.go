package web

import (
	"log/slog"
	"net/http"

	"github.com/bnadland/rhizome/internal/db"
	"github.com/go-chi/chi/v5"
)

func page(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")

		page, err := q.GetPageBySlug(req.Context(), slug)
		if err != nil {
			slog.Warn(err.Error(), "GetPageBySlug", slug)
			w.WriteHeader(http.StatusNotFound)
			if err := NotFound().Render(req.Context(), w); err != nil {
				slog.Error(err.Error())
			}
			return
		}

		if err := Page(page).Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}
