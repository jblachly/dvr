package main

import (
	"encoding/json"
	"log"

	"github.com/jblachly/ddoc-builder"
	"github.com/jblachly/go-couchdb"
)

type DesignDoc struct {
	couchdb.BasicDocumentWithMtime

	Language string                 `json:"language"`
	Views    map[string]interface{} `json:"views"`
}

// load design document into new database
func populateDatabase(db *couchdb.CouchDB) error {

	b, err := ddoc.Build("design", "couchdb")

	// unmarshal into struct
	log.Print(string(b))
	doc := new(DesignDoc)
	err = json.Unmarshal(b, doc)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("pre PostDocument")
	_, err = db.PostDocument(doc)

	return nil
}
