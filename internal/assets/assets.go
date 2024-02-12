package assets

import (
	"embed"
	"net/http"
)

//go:embed static/*.js static/*.css
var assetsFS embed.FS

func AssetHandler() http.Handler {
	return http.FileServer(http.FS(assetsFS))
}
