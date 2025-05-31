package runner

import (
	"strings"
)

type Lang string

type LangDefinition struct {
	image                    string
	execCommand              []string
	buildCommand             []string
	sourceFileName           string
	SampleCode               string
	nonDetermenisticKeywords []string
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

type Runner struct {
	lang Lang
	code string
}

type FS interface {
	Create(fileName string, sourceCode string) (string, error)
	Cleanup() error
}

type CommandExecutor interface {
	Run(name string, args ...string) (string, string, float64, error)
}

func (r *Runner) containsNonDetermenisticKeywords() bool {
	ld, ok := LangDefinitions[r.lang]
	if !ok {
		return false
	}
	if ld.nonDetermenisticKeywords == nil {
		return false
	}

	for _, substr := range ld.nonDetermenisticKeywords {
		if strings.Contains(r.code, substr) {
			return true
		}
	}

	return false
}
