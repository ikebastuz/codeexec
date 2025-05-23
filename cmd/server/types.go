package main

type Response struct {
	Stdout        string  `json:"stdout"`
	Stderr        string  `json:"stderr"`
	ExecDuration  float64 `json:"exec_duration"`
	BuildDuration float64 `json:"build_duration"`
	Error         string  `json:"error"`
}
