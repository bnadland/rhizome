package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"log"
	"net/http"
	"time"
)

type Page struct {
	Id        int
	PageId    string
	VersionId int
	CmsPage   string
	Content   string
	Published time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func getHtml(markdown string) string {
	unsafe := blackfriday.MarkdownCommon([]byte(markdown))
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

func getDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "rhizome.db")
	if err != nil {
		log.Print(err)
	}
	db.LogMode(true)
	return &db
}

func getCmsPage(w http.ResponseWriter, req *http.Request) {
	cmsPage := "/" + req.URL.Query().Get(":CmsPage")
	page := &Page{}
	db := getDb()

	if db.Where("cms_page = ?", cmsPage).First(page).RecordNotFound() {
		http.Error(w, "Page not found.", 404)
		return
	}

	fmt.Fprintf(w, getHtml(page.Content))
}

func main() {
	// database
	db := getDb()
	db.AutoMigrate(&Page{})

	/// make sure index page is there
	var page Page
	db.Where(Page{CmsPage: "/"}).Attrs(Page{VersionId: 1, PageId: uuid.New(), Content: "# hello, world"}).FirstOrInit(&page)
	db.Save(page)
	db.Close()

	// routes
	routes := pat.New()
	routes.Get("/:CmsPage", http.HandlerFunc(getCmsPage))
	routes.Get("/", http.HandlerFunc(getCmsPage))
	http.Handle("/", routes)

	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
