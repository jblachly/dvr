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
	"strconv"

	"github.com/jblachly/go-couchdb"
	"github.com/julienschmidt/httprouter"
)

type Context struct {
	db *couchdb.CouchDB
}

// Global context
// TODO: need to pass this around
var ctx Context

type Startup struct {
	Bind         string `json:"bind"`
	Port         int    `json:"port"`
	DatabaseHost string `json:"database_host"`
	DatabasePort int    `json:"database_port"`
	DatabaseUser string `json:"database_user"`
	DatabasePass string `json:"database_pass"`
}

// Load parses startup.json to learn basic configuration,
// especially how to connect to the dvr database, which holds
// more advanced configuration
func (s *Startup) Load() error {

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
	err := s.Load()
	if err != nil {
		log.Println("Unfortunately, there was an error opening, reading, or parsing startup.json")
		log.Println("Specific failure message follows:")
		log.Fatalln(err)
	}

	// connect to couchdb
	log.Printf("Connecting to database: http://%s:%d/dvr\n", s.DatabaseHost, s.DatabasePort)
	ctx.db, err = couchdb.Database("http://couchdb.local:5984", "dvr", s.DatabaseUser, s.DatabasePass)
	if err != nil {
		log.Println("couchdb.Database()")
		log.Fatalln(err)
	}

	info, err := ctx.db.Info()
	if err != nil {
		log.Println("db.Info() error")
		log.Println(info)
		if couchErr, ok := err.(*couchdb.CouchError); ok {
			// examine error
			if couchErr.Reason == "no_db_file" {

				// if no_db_file this may be the first time we've run this
				log.Println("Creating new database...")
				ctx.db, err = couchdb.CreateDatabase("http://couchdb.local:5984", "dvr", s.DatabaseUser, s.DatabasePass)
				if err != nil {
					log.Fatalln(err)
				}

				log.Println("Initializing database...")
				err = populateDatabase(ctx.db)
				if err != nil {
					log.Fatalln(err)
				}

			} // else handle no authorization, give helpful error message?
		} else {
			log.Fatalln(err)
		}
	}

	log.Println("Reading configuration from database")
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

	log.Println("Queuing previously scheduled recordings...")
	// read from couchdb
	// 1. look for any type:recording documents with "scheduled": true (shouldnt be any "in progress" now)
	// 2. look for any type:scheduled?
	// use cron-like library or goroutine with time.After() to schedule them
	log.Println("...0 recordings scheduled")

	// Start web server
	webServerAddress := s.Bind + ":" + strconv.Itoa(s.Port)
	log.Printf("Starting HTTP server on %s\n", webServerAddress)
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	router.GET("/devices", DevicesHandler)
	router.GET("/devices/:id", DevicesHandler)
	router.PUT("/devices/:id", DevicesHandler)
	router.DELETE("/devices/:id", DevicesHandler)

	router.GET("/channels", ChannelsHandler)
	//router.GET("/channels/:id", ChannelsHandler)

	router.GET("/recordings", RecordingsHandler)
	router.GET("/recordings/:id", RecordingsHandler)
	router.POST("/recordings", NewRecordingHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
