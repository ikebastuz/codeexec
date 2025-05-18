package runner

type Lang string

type LangDefinition struct {
	image   string
	command []string
}
