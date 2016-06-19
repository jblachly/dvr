package main

// BUG(james): if all tuners in use, will fail; need to retry other IPs for other devices, if available

import (
	"log"
	"os/exec"
	"strconv"
)

var HDTVResolution = map[int]string{
	480:  "854:480",
	720:  "1280:720",
	1080: "1920:1080",
}

// Software transcode with FFmpeg
// should be called as a goroutine: go transcode
// as it will be blocking (I think)
// Called by record()
func transcode(uri string, uuid string, duration uint, recordingFilename string) bool {

	// source/ffmpeg-3.0.2/ffmpeg -t 60 -i "http://192.168.1.111:5004/auto/v10.1" -c:v libx264 -b:v 1024k -vf scale=854:480 output.mp4

	// BUG(james): codec is hard coded as x264
	// BUG(james): resolution is hard coded as 480p
	// BUG(james): bitrate is hard coded
	// BUG(james): audio is not transcode
	durationStr := strconv.Itoa(int(duration))
	cmd := exec.Command("ffmpeg", "-t", durationStr, "-i", uri, "-codec:v", "libx264", "-b:v", "512k", "-vf", "scale=854:480", recordingFilename)
	err := cmd.Start()

	if err != nil {
		log.Fatalln(err)
	}

	// not clear to me if it is mandatory to now/later call cmd.Wait()
	// to release resources
	return true
}
