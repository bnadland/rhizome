package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
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

func (db *Database) UpdateDocument(doc *Document) error {
	docjson, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	err = db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		docid := new(bytes.Buffer)
		docid.WriteString(doc.Id)
		err = b.Put(docid.Bytes(), docjson)
		if err != nil {
			return err
		}
		docid.WriteString(strconv.Itoa(doc.Version))
		err := b.Put(docid.Bytes(), docjson)
		return err
	})
	return err
}

func (db *Database) GetDocument(id string, version ...int) (*Document, error) {
	doc := Document{}
	db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("documents"))
		docid := new(bytes.Buffer)
		docid.WriteString(id)
		if len(version) == 1 {
			docid.WriteString(strconv.Itoa(version[0]))
		}
		err := json.Unmarshal(b.Get(docid.Bytes()), &doc)
		return err
	})
	return &doc, nil
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

	d1 := Document{
		Id:      "1",
		Title:   "Main",
		Text:    "# hello, world",
		Page:    "/Home",
		Version: 1,
		Created: time.Now(),
		Updated: time.Now(),
	}
	database.UpdateDocument(&d1)
	d2, err := database.GetDocument("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", d2)
}
