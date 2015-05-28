package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type Rhizome struct {
	DestinationDir string
}

func NewRhizome(destDir string, themeDir string) *Rhizome {
	r := &Rhizome{
		DestinationDir: destDir,
	}

	err := os.MkdirAll(r.DestinationDir, 0777)
	if err != nil {
		panic(err)
	}

	r.renderPage("/index.html", "hello, world")

	return r
}

func (r *Rhizome) renderPage(filePath string, content string) error {
	sitePath := path.Join(r.DestinationDir, filePath)

	err := ioutil.WriteFile(sitePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rhizome) serve() {
	renderedSite := http.NewServeMux()
	renderedSite.Handle("/", http.FileServer(http.Dir(r.DestinationDir)))

	fmt.Println("Listening on port 8080 for rendered site")
	http.ListenAndServe(":8080", renderedSite)
}

func main() {
	destinationDir := flag.String("destination", "_site", "Destination directory for rendered site")
	flag.Parse()

	r := NewRhizome(*destinationDir, *themeDir)
	r.serve()
}
