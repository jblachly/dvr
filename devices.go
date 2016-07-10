package main

import (
	"github.com/jblachly/go-couchdb"
)

// should I rename this, and structs like it, DeviceDoc (xxxDoc) ?
type Device struct {
	couchdb.BasicDocumentWithMtime

	Type string `json:"type,omitempty"`

	Host string `json:"host,omitempty"` // Host as entered in Web GUI e.g. 192.168.1.111
	// TODO: instead of host, shoudl I instead just use BaseURL?

	HDHomeRun // embed anonymously
}
