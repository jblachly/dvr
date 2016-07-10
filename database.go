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

// removeDesignDoc removes existing design document, if any
func removeDesignDoc(db *couchdb.CouchDB) error {
	ddocTmp := new(couchdb.BasicDocumentWithMtime)
	err := ctx.db.GetDocument(&ddocTmp, "_design/design")
	if err != nil {
		log.Println("ERROR: databaseInitialize()")
	}

	rev := ddocTmp.Rev
	if rev != "" { // if _rev exists, design doc exists
		// TODO deal with success / errors
		_, err = ctx.db.DeleteDocument(ddocTmp.ID, rev)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadDesignDoc loads the design document into new database
func loadDesignDoc(db *couchdb.CouchDB) error {

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

	// TODO: handle error
	return nil
}
