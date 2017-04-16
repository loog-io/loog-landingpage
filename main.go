package main

import (
	"encoding/json"
	"fmt"
	"hash/adler32"
	"html/template"
	"net/http"

	"gopkg.in/redis.v5"
	"xojoc.pw/useragent"
)

var rds *redis.Client

func main() {
	// Grab config

	// Set globals
	rds = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//Create mux
	mux := http.NewServeMux()

	// Handle looging of page views
	mux.HandleFunc("/loog", loog)
	// Handle home page and stats viewing
	mux.HandleFunc("/", HomeHandler)
	// Handle static assets
	mux.HandleFunc("/static/", HandleStaticAssets)
	//Handle Favicon
	mux.HandleFunc("/favicon.ico", HandleFavicon)

	// Listen and serve
	fmt.Println("Serving...")
	http.ListenAndServe(":5050", mux)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse template ready to insert data into
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		w.Write([]byte("Error: Could not load template."))
	}

	// Fetch example content from Redis
	res := rds.HMGet("h:w7xAW", "h:170114", "v:170114")

	// Set data to be inserted into HTMl
	var feed FeedData
	feed.Views = res.Val()[0].(string)
	feed.Visits = res.Val()[1].(string)

	// Fetch all the looging information from Redis
	logs := rds.LRange("l:w7xAW", 0, -1)
	//fmt.Printf("%v", logs)
	for _, log := range logs.Val() {
		var view Loog

		// Read the json
		err := json.Unmarshal([]byte(log), &view)
		if err != nil {
			panic(err)
		}

		browser := useragent.Parse(view.UserAgent)

		//browser := parser.Parse(view.UserAgent)

		// Assemble the entry for the feed
		var entry FeedEntry
		entry.URL = view.Location
		entry.ViewNumber = view.ViewNumber
		entry.IP = view.IP
		entry.Color = fmt.Sprintf("%x", int(adler32.Checksum([]byte(view.IP))))[0:6]
		entry.Country = "Canada"
		//entry.Browser = fmt.Sprintf("%s %s.%s", browser.UserAgent.Family, browser.UserAgent.Major, browser.UserAgent.Minor)
		entry.Browser = browser.Name + " " + browser.Version.String()
		entry.UA = view.UserAgent
		entry.OS = browser.OS
		entry.Width = view.ScreenWidth
		entry.Height = view.ScreenHeight
		entry.Time = view.Time

		feed.Feed = append(feed.Feed, entry)
	}

	// Execute the template and send it to the waiting client
	tmpl.Execute(w, feed)
	fmt.Println("Recieved request for home page.")
}
