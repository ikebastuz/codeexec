package runner

const (
	LANG_GO         Lang = "go"
	LANG_JAVASCRIPT Lang = "javascript"
	LANG_PYTHON     Lang = "python"
)

const MAIN_EXEC_FILE string = "/app/main"

var langDefinition = map[Lang]LangDefinition{
	LANG_GO: {
		image:   "golang:1.21-alpine",
		command: []string{"go", "run"},
	},
	LANG_JAVASCRIPT: {
		image:   "node:20-alpine",
		command: []string{"node"},
	},
	LANG_PYTHON: {
		image:   "python:3.11-alpine",
		command: []string{"python3"},
	},
}
