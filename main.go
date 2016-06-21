package main

// BUG(james): hardcoded host, port, channel
// BUG(james): UUID is hardcoded (!!!)
// BUG(james): duration is hardcoded
// BUG(james): outfileName is hardcoded

// Standard library
import (
	//"fmt"
	//	"encoding/json"
	"log"
	"net/http"
)

// External dependencies
import (
	"github.com/julienschmidt/httprouter"
)

// Schedule a recording
// can be called by the http handler for POST to /recordings/
// or by Record() when issuing the next instance of a scheduled series
func ScheduleRecording(r *Recording) bool {
	// schedule using time.After
	return true
}

// maybe rename this RecordCallback?
// should be triggered by the scheduler (more correctly, the wait timer)
func Record(r *Recording) bool {

	// 1. for recording object/doc, check if "recurring"
	// 2. If yes, copy to new object; if no, skip to 4
	if r.Recurring {
		s := new(Recording) // allocate memory
		*s = *r             // deep copy (https://play.golang.org/p/R8uStEApCb)

		// 3. new object:
		// a) decrement remaining counter
		// b) post to database
		// c) schedule recording with cron or at or whatever time.After
		s.RecurrenceData.Remaining--
		// post to database
		_ = ScheduleRecording(s)
	}

	// 4. old object: unset scheduled flag; set in progress flag
	// (TO DO: how to enforce in progress flag being turned off later?)
	r.InProgress = true
	r.Scheduled = false

	// 5. update old object in database

	// 6. Record()!

	uri := buildURI("192.168.1.111", 5004, "10.1")
	uuid := "286e792c-e9ab-4983-a17f-36e75f129572"

	if true {
		go streamToDisk(uri, uuid, r.Duration, "./"+uuid+".mp4")
	} else {
		go transcode(uri, uuid, r.Duration, "./"+uuid+".mp4")
	}

	return true
}

func main() {

	log.Println("Starting DVR")

	log.Println("Connecting to database: http://localhost:5984/dvr")
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

	log.Fatal(http.ListenAndServe(":8080", router))
}
