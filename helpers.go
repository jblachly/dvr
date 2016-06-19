package main

import (
	"path/filepath"
	"strconv"
)

// Build the URI for streaming
// This includes host (hostname or IP), port, and virtual channel
// Returns a URI where the stream can be retrieved, e.g.:
// http://192.168.1.222:5004/auto/v10.1
//
// Potential enhancements:
// Later, could include duration or transcoding parameters, or ?DRM
func buildURI(networkHost string, port uint16, virtualChannel string) (uri string) {
	uri = "http://" + networkHost + ":" + strconv.Itoa(int(port)) + "/auto/v" + virtualChannel
	return
}

// Build the output filename
// for now, it concatenates
// 1) recordings directory
// 2) UUID
// 3) extension (mpeg/mp4/whatever) -- this will have to be parameterized
//		according to whether transcoding is planned or not, I suppose
func buildRecordingFilename(recordPath string, uuid string, extension string) (recordingFilename string) {
	absRecordPath, _ := filepath.Abs(recordPath)
	recordingFilename = filepath.Join(absRecordPath, uuid) + "." + extension
	return
}
