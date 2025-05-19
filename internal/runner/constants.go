package runner

const (
	LANG_GO         Lang = "go"
	LANG_JAVASCRIPT Lang = "javascript"
	LANG_PYTHON     Lang = "python"
)

const WORKDIR = "/app"
const EXEC_FILE = "main"

var langDefinition = map[Lang]LangDefinition{
	LANG_GO: {
		image:   "golang:1.21-alpine",
		command: []string{"go", "run"},
		ext:     "go",
	},
	LANG_JAVASCRIPT: {
		image:   "node:20-alpine",
		command: []string{"node"},
		ext:     "js",
	},
	LANG_PYTHON: {
		image:   "python:3.11-alpine",
		command: []string{"python3"},
		ext:     "py",
	},
}

func GetLangs() []string {
	var languages []string
	for k := range langDefinition {
		languages = append(languages, string(k))
	}
	return languages
}
