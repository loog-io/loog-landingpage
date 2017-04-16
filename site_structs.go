package main

// FeedData contains the data needed to be fed to the real time feed template.
type FeedData struct {
	Views  string
	Visits string
	Feed   []FeedEntry
}

// FeedEntry contains the data for each entry in the real-time log
type FeedEntry struct {
	Referrer   string
	URL        string
	ViewNumber int
	IP         string
	Color      string
	Country    string
	Browser    string
	OS         string
	UA         string
	Width      int
	Height     int
	Time       string
}
