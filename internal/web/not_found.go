package web

import (
	"log/slog"
	"net/http"
)

func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := NotFound().Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}
