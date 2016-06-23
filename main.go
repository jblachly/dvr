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
	"time"
)

// External dependencies
import (
	"github.com/julienschmidt/httprouter"
)

// Schedule a recording
// can be called by the http handler for POST to /recordings/
// or by Record() when issuing the next instance of a scheduled series
func ScheduleRecording(r *Recording) bool {

	// Check to see if the program end time has already passed
	// If start time + duration >= current time, panic
	// This condition should be prevented elsewhere
	duration := time.Duration(r.Duration) * time.Second
	durationUntilStart := r.Date.Sub(time.Now())

	if r.Date.Add(duration).After(time.Now()) {
		// TODO enhance error message with more information
		panic("Attempted to schedule recording entirely in the past")
	}

	// OK, now check to see if the start time has already passed
	// This could happen in at least two situations
	// 1) User selects to record something already in progress
	// 2) Program restarts and looks through list of previously
	//    scheduled recordings.
	if r.Date.After(time.Now()) {
		// start recording immediatley - program in progress
		// Record(r)
		// TODO: decide if this should be a channel, go routine, both, etc.
		go func() {
			_ = Record(r) // if make Record() return void could just call go Record(r)
		}()
	}

	// Recording is in the future, so record sometime in the future
	// option 1
	// schedule using time.After
	if false {
		go func() {
			<-time.After(durationUntilStart)
			_ = Record(r)
		}()
	}

	// option 2
	// schedule using timer.AfterFunc
	if true {
		timer := time.AfterFunc(durationUntilStart, func() { Record(r) }) // closure over r
		// consider _ = time.AfterFunc, so long as the timer is actually allocated
		log.Println("Timer scheduled", timer)
	}

	log.Printf("Recording scheduled for %s (%s from now) for duration of %s ", r.Date, durationUntilStart, r.Duration)
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
