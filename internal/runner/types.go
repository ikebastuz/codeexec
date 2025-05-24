package runner

type Lang string

type LangDefinition struct {
	image          string
	execCommand    []string
	buildCommand   []string
	sourceFileName string
	SampleCode     string
}

type Job struct {
	Stdout string
	Stderr string
	Error  error
}

type Result struct {
	Stdout        string  `json:"stdout"`
	Stderr        string  `json:"stderr"`
	ExecDuration  float64 `json:"exec_duration"`
	BuildDuration float64 `json:"build_duration"`
	Error         string  `json:"error"`
}
