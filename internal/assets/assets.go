package assets

import (
	"embed"
	"log/slog"
	"net/http"
)

//go:embed static/*.js static/*.css
var assetsFS embed.FS

func AssetHandler() http.Handler {
	return http.FileServer(http.FS(assetsFS))
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := NotFound().Render(req.Context(), w); err != nil {
			slog.Error(err.Error())
		}
	}
}
