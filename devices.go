package main

import (
	"github.com/jblachly/go-couchdb"
)

// should I rename this, and structs like it, DeviceDoc (xxxDoc) ?
type Device struct {
	couchdb.BasicDocumentWithMtime

	Type string `json:"type,omitempty"`

	HDHomeRun // embed anonymously
}
