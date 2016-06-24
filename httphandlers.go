package main

import (
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func RecordingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintln(w, "Recordings handler function")
}

func NewRecordingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintln(w, "<h1>Scheduling New recording</h1>")

	rec := new(Recording)
	rec.Init() // give it a UUID

	// Need to fill this from the posted data (*http.Request)
	//rec.Date = time.Date(2016, time.June, 23, 21, 33, 00, 0000, time.Local)
	rec.Date = time.Now().Add(time.Duration(30) * time.Second) // 30 seconds in the future
	rec.Duration = 15                                          // 60 s/min * 30 min = 1800 sec
	rec.Scheduled = true

	//requestedTime := time.Date(2018, time.June, 13, 19, 00, 00, 0000, time.UTC)

	if rec.Date.Add(time.Duration(rec.Duration) * time.Second).Before(time.Now()) {

		fmt.Fprintln(w, "ERROR: Requested recording window is in the past")

	} else {

		fmt.Fprintf(w, "Scheduled recording on %s for %d seconds\n", rec.Date.Local(), rec.Duration)

		// post rec to database

		// call the scheduler
		ScheduleRecording(rec)
	}
}
