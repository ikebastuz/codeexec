package main

type Response struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Duration float64 `json:"duration"`
	Error    string  `json:"error"`
}
