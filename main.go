package main

import (
	//    "encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	// "github.com/russross/blackfriday"
	// "code.google.com/p/go-uuid/uuid"
)

type Config struct {
	dsn string
}

var config = new(Config)

type Database struct {
	db *bolt.DB
}

func (db *Database) Init(dsn string) (*bolt.DB, error) {
	boltdb, err := bolt.Open(dsn, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db.db = boltdb

	err = db.db.Update(func(tx *bolt.Tx) error {
		for _, bucketname := range []string{"documents"} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketname))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return db.db, err
}

func (db *Database) UpdateDocument(doc *Document) {
	fmt.Printf("%v\n%v\n", db.db, doc)
}

var database = new(Database)

type Document struct {
	Id      string
	Title   string
	Text    string
	Page    string
	Version int
	Created time.Time
	Updated time.Time
}

func main() {
	// config
	flag.StringVar(&config.dsn, "dsn", "rhizome.db", "database filename")
	flag.Parse()

	// db
	db, err := database.Init(config.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// start
	fmt.Println("# rhizome")
	fmt.Printf("using database: %s\n", config.dsn)

	doc := Document{
		Id:      "1",
		Title:   "Main",
		Text:    "# hello, world",
		Page:    "/Home",
		Version: 1,
		Created: time.Now(),
		Updated: time.Now(),
	}
	database.UpdateDocument(&doc)
}
