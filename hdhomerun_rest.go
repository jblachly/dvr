package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"url"
)

func getJSON(host, endpoint) ([]byte, error) {

	u := url.URL{Scheme: "http", Host: host, Path: endpoint}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err) // HD HomeRun not correctly configured should not panic or abort program
		return nil, err
	}

	j, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err) // HD HomeRun not correctly configured should not panic or abort program
		return nil, err
	}

	return j, nil
}

// DISCOVER
// http://deviceip:80/discover.json

type HDHomeRun struct {
}

func discoverHDHR(host string) (*HDHomeRun, error) {

	return nil, nil
}

// LINEUP STATUS
// http://deviceip:80/lineup_status.json
// Example result:
// {"ScanInProgress":0,"ScanPossible":1,"Source":"Antenna","SourceList":["Antenna","Cable"]}

type LineupStatus struct {
	ScanInProgress int
	ScanPossible   int
	Source         string
	SourceList     []string
}

func getLineupStatus(host string) (*LineupStatus, error) {

	j, err := getJSON(host, "lineup_status.json")
	if err != nil {
		return nil, err
	}

	ls := new(LineupStatus{})
	json.Unmarshal(j, &ls)

	return ls, nil
}

/* LINEUP

http://deviceip:80/lineup.json

Initially, this is an empty array:
[]

After running a channel scan, it is an array of detected channels:
[
	{
	"GuideNumber": "4.1",
	"GuideName": "WCMH-DT",
	"HD": 1,
	"URL": "http://192.168.1.111:5004/auto/v4.1"
	},
	...
]
*/

type Lineup struct {
	// TODO: look for field set when channel disabled
	GuideNumber string
	GuideName   string
	HD          int
	Favorite    int
	Subscribed  int // undocumented, unclear if int
	DRM         int // undocumented, unclear if int
	URL         string
}

func getLineup(host string) ([]Lineup, error) {

	j, err := getJSON(host, "lineup.json")
	if err != nil {
		return nil, err
	}

	lineups := make([]Lineup, 0, 0)
	json.Unmarshal(j, &lineups)

	return lineups, nil

}
