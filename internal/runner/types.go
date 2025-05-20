package runner

type Lang string

type LangDefinition struct {
	image      string
	command    []string
	fileName   string
	SampleCode string
}

type Job struct {
	Stdout string
	Stderr string
	Error  error
}
