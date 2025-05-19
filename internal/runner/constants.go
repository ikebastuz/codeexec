package runner

const (
	LANG_GO         Lang = "go"
	LANG_JAVASCRIPT Lang = "javascript"
	LANG_PYTHON     Lang = "python"
	LANG_RUBY       Lang = "ruby"
	LANG_RUST       Lang = "rust"
	LANG_PHP        Lang = "php"
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
	LANG_RUBY: {
		image:   "ruby:3.3-alpine",
		command: []string{"ruby"},
		ext:     "rb",
	},
	LANG_RUST: {
		image:   "rust:1.77-alpine",
		command: []string{"sh", "-c", "rustc /app/main.rs -o /app/main && /app/main"},
		ext:     "rs",
	},
	LANG_PHP: {
		image:   "php:8.3-cli-alpine",
		command: []string{"php"},
		ext:     "php",
	},
}

func GetLangs() []string {
	var languages []string
	for k := range langDefinition {
		languages = append(languages, string(k))
	}
	return languages
}
