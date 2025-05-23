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
