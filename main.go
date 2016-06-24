package main

// BUG(james): hardcoded host, port, channel
// BUG(james): UUID is hardcoded (!!!)
// BUG(james): duration is hardcoded
// BUG(james): outfileName is hardcoded

// Standard library
import (
	//"fmt"
	//	"encoding/json"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// External dependencies
import (
	"github.com/julienschmidt/httprouter"
)

type Startup struct {
	Bind         string `json:"bind"`
	Port         int    `json:"port"`
	DatabaseHost string `json:"database_host"`
	DatabasePort int    `json:"database_port"`
	DatabaseUser string `json:"database_user"`
	DatabasePass string `json:"database_pass"`
}

// loadStartup parses startup.json to learn basic configuration,
// especially how to connect to the dvr database, which holds
// more advanced configuration
func (s *Startup) loadStartup() error {

	// TODO: parameterize startup config file name
	jsonFile, err := ioutil.ReadFile("startup.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonFile, s); err != nil {
		return err
	}

	return nil // no error
}
func main() {

	log.Println("Starting DVR")

	log.Println("Reading startup.json")
	s := new(Startup)
	err := s.loadStartup()
	if err != nil {
		log.Println("Unfortunately, there was an error opening, reading, or parsing startup.json")
		log.Println("Specific failure message follows:")
		log.Fatalln(err)
	}

	log.Printf("Connecting to database: http://%s:%d/dvr\n", s.DatabaseHost, s.DatabasePort)
	// connect to couchdb

	log.Println("Reading configuration")
	// read configuration from couchdb
	// if first time, run database setup
	if false {
		databaseInitialize()
	}

	// if corrupt, offer to repair

	log.Println("Database consistency check")
	// ensure nothing is "in progress" on startup (?)
	// that being said, if doing a software transcode it coudl still be working unless it got a HUP when dvr closed?
	// streaming to disk would ahve stopped however
	// perhaps the recording function, since running in a goroutine, can wait for finish to turn off the inprogress flag
	// ALSO: if ever > 1 copy of dvr running at same time / on different hosts this logic breaks
	//databaseCheckConsistency()

	log.Println("Scheduling recordings...")
	// read from couchdb
	// 1. look for any type:recording documents with "scheduled": true (shouldnt be any "in progress" now)
	// 2. look for any type:scheduled?
	// use cron-like library or goroutine with time.After() to schedule them
	log.Println("...0 recordings scheduled")

	log.Println("Starting HTTP server on 127.0.0.1:8080")
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	router.GET("/recordings", RecordingsHandler)
	router.GET("/recordings/:id", RecordingsHandler)
	router.POST("/recordings", NewRecordingHandler)
	router.GET("/testrec", NewRecordingHandler) // TESTING

	log.Fatal(http.ListenAndServe(":8080", router))
}
