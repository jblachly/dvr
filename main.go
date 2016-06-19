package main

// BUG(james): hardcoded host, port, channel
// BUG(james): UUID is hardcoded (!!!)
// BUG(james): duration is hardcoded
// BUG(james): outfileName is hardcoded

// Standard library
import (
	"fmt"
	//	"encoding/json"
	"log"
	"net/http"
)

// External dependencies
import (
	"github.com/julienschmidt/httprouter"
)

func Record() bool {

	uri := buildURI("192.168.1.111", 5004, "10.1")
	uuid := "286e792c-e9ab-4983-a17f-36e75f129572"
	duration := uint(60)

	if true {
		go streamToDisk(uri, uuid, duration, "./uuid.mp4")
	} else {
		go transcode(uri, uuid, duration, "./uuid.mp4")
	}

	return true
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
