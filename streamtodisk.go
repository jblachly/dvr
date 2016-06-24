package main

// BUG(james): if all tuners in use, will fail; need to retry other IPs for other devices, if available

import (
	"io"
	//	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// stream directly from URI without software transcoding
// stub function
// should be called as goroutine: go streamToDisk(...)
func streamToDisk(uri string, uuid string, duration uint, recordingFilename string) bool {

	// Open output file
	outFile, err := os.Create(recordingFilename)
	if err != nil {
		log.Println("Could not open ", recordingFilename)
		return false
	}
	defer outFile.Close()

	log.Printf("<streamToDisk> Opening %s\n", uri)
	resp, err := http.Get(uri)
	if err != nil {
		// TO DO handle error
	}

	defer resp.Body.Close()
	log.Printf("<streamToDisk> reading from %s\n", uri)

	/*
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// TO DO handle error
		}
		_ = body // eliminate compiler warning
	*/

	written, err := io.CopyBuffer(outFile, resp.Body, nil)
	if err != nil {
		log.Printf("<streamtoDisk> %s", err)
		return false
	}

	log.Printf("<streamToDisk> Wrote %d bytes\n", written)

	return true
}
