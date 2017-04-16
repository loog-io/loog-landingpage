package main

type Loog struct {
	ID           string `json:"-"`
	Time         string `json:"t"`
	Location     string `json:"l"`
	Referrer     string `json:"r,omitempty"`
	IP           string `json:"i"`
	ScreenWidth  int    `json:"w"`
	ScreenHeight int    `json:"h"`
	UserAgent    string `json:"a"`
	ViewNumber   int    `json:"n"`
}
