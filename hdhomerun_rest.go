package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func getJSON(host, endpoint string) ([]byte, error) {

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
	// the HDHomeRun JSON field names are in camel case as below
	FriendlyName    string
	ModelNumber     string
	FirmwareName    string
	FirmwareVersion string
	DeviceID        string
	DeviceAuth      string
	BaseURL         string
	LineupURL       string
}

func discoverHDHR(host string) (*HDHomeRun, error) {

	j, err := getJSON(host, "discover.json")
	if err != nil {
		return nil, err
	}

	hdhr := new(HDHomeRun)
	json.Unmarshal(j, &hdhr)

	return hdhr, nil
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

	j, err := getJSON(host, "lineup_status.json") // TODO: ?show=found
	if err != nil {
		return nil, err
	}

	ls := new(LineupStatus)
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
	GuideNumber string
	GuideName   string
	HD          int
	Favorite    int
	Enabled     int // zero default is not sensible here / consider ignore
	Subscribed  int // undocumented, unclear if int
	DRM         int // undocumented, unclear if int
	URL         string
}

func getLineup(host string) ([]Lineup, error) {

	j, err := getJSON(host, "lineup.json?show=found")
	if err != nil {
		return nil, err
	}

	lineups := make([]Lineup, 0, 0)
	json.Unmarshal(j, &lineups)

	return lineups, nil

}
