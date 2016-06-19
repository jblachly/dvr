package main

// BUG(james): if all tuners in use, will fail; need to retry other IPs for other devices, if available

import (
	"io/ioutil"
	"net/http"
)

// stream directly from URI without software transcoding
// stub function
// should be called as goroutine: go streamToDisk(...)
func streamToDisk(uri string, uuid string, duration uint, recordingFilename string) bool {

	resp, err := http.Get(uri)
	if err != nil {
		// TO DO handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TO DO handle error
	}
	_ = body // eliminate compiler warning
	return true
}
