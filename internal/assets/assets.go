package assets

import (
	"embed"
	"net/http"
)

//go:embed static/*
var assetsFS embed.FS

func Assets() http.Handler {
	return http.FileServer(http.FS(assetsFS))
}
