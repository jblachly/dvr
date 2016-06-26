package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/julienschmidt/httprouter"
)

func replyJSON(w http.ResponseWriter, status, message string) {
	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	resp := Response{status, message}
	// TODO error checking
	jsonresp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresp)
	fmt.Fprint(w, "\n")
}

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

	// TODO: move definition
	type Request struct {
		Date     time.Time `json:"date"`
		Duration uint      `json:"duration"`
	}

	// Read the posted JSON, or return error
	//if r.Header.Get("Content-Type") != "application/json" {
	//}
	body := new(Request)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		replyJSON(w, "error", "Could not parse POST request into JSON")
		return
	}

	// Ensure posted JSON has right fields
	// or return error
	if body.Date.IsZero() {
		replyJSON(w, "error", "date was not set")
		return
	}
	if body.Duration < 1 {
		replyJSON(w, "error", "duration was less than 1")
		return
	}

	rec := new(Recording)
	rec.Init() // give it a UUID

	// Need to fill this from the posted data (*http.Request)
	//rec.Date = time.Date(2016, time.June, 23, 21, 33, 00, 0000, time.Local)
	//rec.Date = time.Now().Add(time.Duration(30) * time.Second) // 30 seconds in the future
	rec.Date = body.Date
	rec.Duration = body.Duration // 60 s/min * 30 min = 1800 sec
	rec.Scheduled = true

	//requestedTime := time.Date(2018, time.June, 13, 19, 00, 00, 0000, time.UTC)

	if rec.Date.Add(time.Duration(rec.Duration) * time.Second).Before(time.Now()) {

		replyJSON(w, "error", "Requested recording window is in the past")

	} else {

		// post rec to database

		// call the scheduler
		ScheduleRecording(rec)

		msg := fmt.Sprintf("Scheduled recording on %s for %d seconds", rec.Date.Local(), rec.Duration)
		replyJSON(w, "ok", msg)
	}

}
