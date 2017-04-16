package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

func loog(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Recieved a log")
	w.Header().Set("Server", "loog")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.WriteHeader(http.StatusNoContent)

	var view Loog
	// The ID number of the site within our system
	view.ID = r.FormValue("s")
	// The time as of now, in human readable format
	now := time.Now().UTC()
	view.Time = now.Format("2006-01-02 15:04:05")
	// The date, for time-series storage in Redis
	date := now.Format("060102")
	// The URL the visitor is looking upon
	view.Location = r.FormValue("l")
	// The URL that referred this page view
	view.Referrer = r.FormValue("r")
	// The user's IP address
	view.IP, _, _ = net.SplitHostPort(r.RemoteAddr)
	// The width of the browser at the time of loading
	view.ScreenWidth, _ = strconv.Atoi(r.FormValue("w"))
	// The height of the browser at the time of loading
	view.ScreenHeight, _ = strconv.Atoi(r.FormValue("h"))
	// The visitor's browser string.
	view.UserAgent = r.Header.Get("User-Agent")
	// How many pages on the site the user has visited before and including this view
	view.ViewNumber, _ = strconv.Atoi(r.FormValue("v"))

	// Print that out for ease of debugging
	/*fmt.Printf(
	"At %s, user visited site with ID %s, at page %s, from referrer %s with the IP %s and their window at %dx%d their UA is %s and this was page number %d\n",
	view.Time,
	view.ID,
	view.Location,
	view.Referrer,
	view.IP,
	view.ScreenWidth,
	view.ScreenHeight,
	view.UserAgent,
	view.ViewNumber)*/

	// Created the JSON needed for insertation
	ins, err := json.Marshal(view)
	if err != nil {
		w.Write([]byte("Oops, could not encode JSON"))
	}
	//fmt.Println(string(ins))

	//Push the JSON into the relevant Redis list. Todo: check site is allowed
	_ = rds.LPush("l:"+view.ID, string(ins))
	// Limit the list to only the 200 most recent page views
	_ = rds.LTrim("l:"+view.ID, 0, 199)

	// Always increment the number of hits
	_ = rds.HIncrBy("h:"+view.ID, "h:"+date, 1)

	// If this is the first visit, incrememnt the number of visitors.
	if view.ViewNumber == 1 {
		_ = rds.HIncrBy("h:"+view.ID, "v:"+date, 1)
		//fmt.Println("This was a first visit, incremented visit number")
	} else {
		//fmt.Println("This was not a first visit, ignored visit count.")
	}

	w.Write([]byte("OK")) // For testing purposes naturally.
}
