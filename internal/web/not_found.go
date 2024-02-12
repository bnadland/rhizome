package web

import (
	"log/slog"
	"net/http"

	"github.com/bnadland/rhizome/internal/views"
)

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := views.NotFound().Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}
