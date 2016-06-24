package main

import (
	"time"
)

import (
	"github.com/satori/go.uuid"
)

type Recording struct {
	Id  string `json:"_id"`
	Rev string `json:"_rev"`

	Date     time.Time `json:"date"`
	Duration uint      `json:"duration"`
	// would love to use time.Duration instead of uint in seconds,
	// but time.Duration does not implement json>Marshaler interface :|
	// https://github.com/golang/go/issues/4712

	Device        string `json:"device"`
	ChannelNumber string `json:"channel_number"`
	ChannelName   string `json:"channel_name"`

	ShowName    string `json:"show_name"`
	ShowSeason  string `json:"show_season"`
	ShowEpisode string `json:"show_episode"`

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
	r.Id = u.String()
}
