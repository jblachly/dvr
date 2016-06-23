package main

// Standard library
import (
	//"fmt"
	//	"encoding/json"
	"log"
	"time"
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

	if r.Date.Add(duration).Before(time.Now()) {
		// TODO enhance error message with more information
		//log.Printf("r.Date: %s \n r.Duration: %s \n duration: %s \n time.Now(): %s", r.Date, r.Duration, duration, time.Now())
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

	log.Printf("Recording scheduled for %s (%s from now) for duration of %d ", r.Date, durationUntilStart, r.Duration)
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
