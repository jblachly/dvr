package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/jblachly/go-couchdb"
	"github.com/julienschmidt/httprouter"
)

type Page struct {
	PageTitle string
}

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

func replyJSONobj(w http.ResponseWriter, obj interface{}) {

	jsonresp, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresp)
	fmt.Fprintf(w, "\n")

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// fmt.Fprint(w, "Welcome!\n")
	p := Page{PageTitle: "DVR Main"}

	t, err := template.ParseFiles("templates/base.html", "templates/index.html")
	if err != nil {
		replyJSON(w, "error", err.Error())
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		replyJSON(w, "error", err.Error())
		return
	}

}

func DevicesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Check to see if URL includes an /:id/
	// if none set, ByName() returns empty string
	if ps.ByName("id") == "" {

		vargs := couchdb.ViewArgs{Reduce: couchdb.FalsePointer}
		vres, err := ctx.db.View("design", "devicesByDeviceID", vargs)

		// return vres.Rows

		var p = struct {
			Page // anonymous embedding
			Rows []couchdb.ViewRow
		}{
			Page: Page{"Video Devices"},
			Rows: vres.Rows,
		}

		t, err := template.ParseFiles("templates/base.html", "templates/devices.html")
		if err != nil {
			replyJSON(w, "error", err.Error())
		}
		err = t.Execute(w, p)
		if err != nil {
			replyJSON(w, "error", err.Error())
		}

		return
	}

	// else :id included in URL

	switch r.Method {
	case "GET":
		/*
			vargs := couchdb.ViewArgs{Key: ps.ByName("id"), Reduce: couchdb.FalsePointer}
			vres, _ := ctx.db.View("design", "devices", vargs)
			replyJSONobj(w, vres)
		*/
		dev := Device{}
		ctx.db.GetDocument(&dev, ps.ByName("id"))
		replyJSONobj(w, dev)

	case "POST":
		log.Println("POST to /devices/:id")

		// TODO: Add a timeout to the discoverHDHR call
		hdhr, err := discoverHDHR(ps.ByName("id"))
		// TODO: if err != nil Flash failure message
		if err == nil {
			dev := Device{Type: "device", Host: ps.ByName("id"), HDHomeRun: *hdhr}
			ctx.db.PostDocument(dev)
		}

		http.Redirect(w, r, "/devices", http.StatusSeeOther)

	case "PUT":
		panic("PUT to /devices/:id [shouldnt happen]")

	case "DELETE":
		log.Println("DELETE to /devices/:id")

		dev := Device{}
		ctx.db.GetDocument(&dev, ps.ByName("id"))

		if dev.ID == "" {
			replyJSON(w, "error", "cannot locate "+ps.ByName("id"))
			return
		}

		cs, err := ctx.db.DeleteDocument(dev.ID, dev.Rev)
		if err != nil || !cs.OK {
			replyJSON(w, "error", err.Error())
			return
		} else {
			replyJSON(w, "ok", "deletion_succeeded")
		}

	}

}

func ChannelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var p = struct {
		Page     // anonymous embedding
		Channels [3]string
	}{
		Page:     Page{"Channel Lineup"},
		Channels: [3]string{"ABC", "NBC", "CBS"},
	}

	t, err := template.ParseFiles("templates/base.html", "templates/channels.html")
	if err != nil {
		replyJSON(w, "error", err.Error())
	}
	err = t.Execute(w, p)
	if err != nil {
		replyJSON(w, "error", err.Error())
	}
}

func RecordingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var p = struct {
		Page
		Recordings [3]string
	}{
		Page:       Page{"Recordings"},
		Recordings: [3]string{"recording1", "recording2", "recording3"},
	}

	t, err := template.ParseFiles("templates/base.html", "templates/recordings.html")
	if err != nil {
		replyJSON(w, "error", err.Error())
	}
	err = t.Execute(w, p)
	if err != nil {
		replyJSON(w, "error", err.Error())
	}
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

///////////////////////////////////////
// API

func ResetDatabase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {}

// UpdateDesignDoc replaces existing design document with freshly built
// Currently its code is identical to the beginning of databaseInitialize,
// and coudl be factored out into a common call
func UpdateDesignDoc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// remove existing design document, if any
	err := removeDesignDoc(ctx.db)
	if err != nil {
		replyJSON(w, "error", err.Error())
		return
	}

	// load a compiled design doc
	err = loadDesignDoc(ctx.db)
	if err != nil {
		replyJSON(w, "error", err.Error())
		return
	}

	//return nil

}
