// This app starts a webservice that shows the last 25 titles of the radio-paradise mellow-mix webradio
// STEK7212 20220627
// LICENSE: MIT
//
//
// todo: code comments

package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// global variable 'output' to store track title
var output = "empty"

// global variable 'rpPlaylist' to store 25 titles
var rpPlaylist = make([]string, 25)

func main() {

	fmt.Println("** Radio Paradise - Mellow Mix Playlist")
	fmt.Println("***************************************")
	titleOld := ""

	// prefill variable 'rpPlaylist' with 25 empty tracks
	for i := 24; i > -1; i-- {
		rpPlaylist[i] = ".<br>"
	}

	// handler for webserver main-route '/'
	http.HandleFunc("/", handler)

	// start webserver on port 7070
	go func() {
		log.Fatal(http.ListenAndServe(":7070", nil))
	}()

	// endless for-loop for main routine
	for {
		titleNew := getTitle()
		if titleNew != titleOld {
			currentTime := time.Now()                          // get current time
			filetime := currentTime.Format("2006.01.02 15:04") // build string with current date from timestamp
			output = filetime + " | " + titleNew + "<br>"
			fmt.Println(output)
			rpPlaylist = leftRotation(rpPlaylist, len(rpPlaylist), 1)
			rpPlaylist[24] = output
		}
		titleOld = titleNew
		time.Sleep(90 * time.Second) // wait for 90 seconds until next title
	}
}

// function handler route "/" Webserver - shows all 25 tracks
func handler(w http.ResponseWriter, r *http.Request) {
	htmloutput := "" +
		"<!DOCTYPE html><html lang=\"de\">" +
		"<head>" +
		"<meta charset=\"UTF-8\">" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">" +
		"<meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">" +
		"<title>RP-Playlist</title>" +
		"</head>" +
		"<body>" +
		"<font face=\"sans-serif\">" +
		"<style>" +
		"body {" +
		"  font-size: 12px;" +
		"}" +
		"a:link { text-decoration: none; }" +
		"</style>" +
		"<a href=\"/\"><h1>Radio Paradise Playlist</h1></a>"

	fmt.Fprintf(w, htmloutput)

	for i := 24; i > -1; i-- {
		fmt.Fprintf(w, rpPlaylist[i])
	}
	fmt.Fprintf(w, "</body></html>")
}

// leftRotation - shifts string variable left - to rotate playlist
func leftRotation(a []string, size int, rotation int) []string {
	var newArray []string
	for i := 0; i < rotation; i++ {
		newArray = a[1:size]
		newArray = append(newArray, a[0])
		a = newArray
	}
	return a
}

// get actual title with the help of the external 'streamripper' - command
func getTitle() string {
	playlistCmd := exec.Command("streamripper", "https://stream.radioparadise.com/mellow-32", "-M", "0")
	playlistOut, err := playlistCmd.Output()
	if err != nil {
		fmt.Println("error.")
	}
	title := "empty title string - maybe an error?"
	s := strings.Split(string(playlistOut), "\n")
	if len(s) > 5 {
		title = s[6]
	}
	titleLen := len(title) - 5
	titleNew := title[17:titleLen]
	return titleNew
}
