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
	fmt.Fprintln(w, "Scheduling New recording")

	// Need to upll this from the posted data (*http.Request)
	requestedTime := time.Date(2018, time.June, 13, 19, 00, 00, 0000, time.UTC)

	if requestedTime.Before(time.Now()) {
		fmt.Fprintln(w, "ERROR: Requested time is in the past")
	}

	fmt.Fprintln(w, "Scheduled recording for %s", requestedTime.Local())

	// allocate and populate instance of Recording struct
	// TO DO: take data from http.Request
	rec := Recording{
		Date:      requestedTime,
		Duration:  300,
		Scheduled: true,
	}
	_ = rec // silence compiler warning

	// post rec to database

}
