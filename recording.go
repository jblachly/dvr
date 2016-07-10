package main

// Standard library
import (
	//"fmt"
	//	"encoding/json"
	"log"
	"time"

	"github.com/jblachly/go-couchdb"
	"github.com/satori/go.uuid"
)

type Recording struct {
	couchdb.BasicDocumentWithMtime // anonymous

	// couchdb.BasicDocumentW and bdwMtime include ID and Rev; also Deleted, Attachments, Created, and Modified
	//Id  string `json:"_id"`
	//Rev string `json:"_rev"`

	Date     time.Time `json:"date"`
	Duration uint      `json:"duration"`
	// would love to use time.Duration instead of uint in seconds,
	// but time.Duration does not implement json>Marshaler interface :|
	// https://github.com/golang/go/issues/4712

	Device        string `json:"device"`
	ChannelNumber string `json:"channel_number"`
	ChannelName   string `json:"channel_name"`

	ProgramName    string `json:"show_name"`
	ProgramSeason  string `json:"show_season"`
	ProgramEpisode string `json:"show_episode"`

	Filename string `json:"filename"`

	Scheduled  bool `json:"scheduled"` // should this be "future" ?
	InProgress bool `json:"in_progress"`
	Recurring  bool `json:"recurring"`

	RecurrenceData Recurrence `json:"recurrence"`

	// add embedded struct for recurrence info
}

type Recurrence struct {
	Name      string `json:"name"`      // Title for the recurrence, e.g. "Jeopardy" or "CBS Sunday Morning" or "2016 Olympics"
	Frequency string `json:"frequency"` // daily, daily-weekday, weekly, monthly
	Remaining uint   `json:"remaining"` // number of occurences remaining
	Path      string `json:"path"`      // relative path to save recordings belonging to this recurrence
	// decrement this ONLY once the recording starts, otherwise consecutive starts-stops
	// of the server will bring the counter down
}

// Init is a constructor for type Recording
// Presently, it only generates a type 4 (random) UUID
func (r *Recording) Init() {
	u := uuid.NewV4()
	r.ID = u.String()
}

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
		// This should not have reached this point - the httpHandler should not schedule this,
		// nor should the startup function when parsing database
		// to reduce code duplication, should I move those two checks to here only, and return err
		// instead of panic?
		// TODO enhance error message with more information
		//log.Printf("r.Date: %s \n r.Duration: %s \n duration: %s \n time.Now(): %s", r.Date, r.Duration, duration, time.Now())
		panic("Attempted to schedule recording entirely in the past")
	}

	// OK, now check to see if the start time has already passed
	// This could happen in at least two situations
	// 1) User selects to record something already in progress
	// 2) Program restarts and looks through list of previously
	//    scheduled recordings.
	if time.Now().After(r.Date) {
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
		log.Printf("<scheduleRecording> durationUntilStart: %d", durationUntilStart)
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

	if true {
		go streamToDisk(uri, r.ID, r.Duration, "./"+r.ID+".mp4")
	} else {
		go transcode(uri, r.ID, r.Duration, "./"+r.ID+".mp4")
	}

	return true
}
