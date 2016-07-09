package main

import (
	"github.com/jblachly/go-couchdb"
)

type Channel struct {
	couchdb.BasicDocumentWithMtime

	// couchdb.BasicDocumentW and bdwMtime include ID and Rev; also Deleted, Attachments, Created, and Modified
	//Id  string `json:"_id"`
	//Rev string `json:"_rev"`

	HDHomeRun string `json:"hdhomerun"`

	Lineup // embed anonymous struct (defined in hdhomeru_rest.go)

}
