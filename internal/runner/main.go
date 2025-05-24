package runner

import (
	"codeexec/internal/config"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func NewRunner(lang Lang, code string) *Runner {
	return &Runner{lang: lang, code: code}
}

func (r *Runner) Run() Result {
	fs := &TempDirFS{}
	executor := &CommandExecutorOS{timeout: config.PROCESS_TIMEOUT}
	return r.runWithDeps(fs, executor)
}

func (r *Runner) runWithDeps(fs FS, executor CommandExecutor) Result {
	ld, ok := LangDefinitions[r.lang]
	if !ok {
		return Result{Error: fmt.Sprintf("unknown language: %s", r.lang)}
	}

	tempDir, err := fs.Create(ld.sourceFileName, r.code)
	if err != nil {
		return Result{Error: err.Error()}
	}
	defer fs.Cleanup()

	var buildDuration float64
	var stdout, stderr string

	if ld.buildCommand != nil {
		buildCommand := mkBuildCommand(ld, tempDir)
		stdout, stderr, buildDuration, err = executor.Run("docker", buildCommand...)

		if err != nil {
			log.Errorf("Failed to build: %s", err)
			return Result{
				Stdout:        stdout,
				Stderr:        stderr,
				BuildDuration: buildDuration,
				Error:         err.Error(),
			}
		}
	}

	execCommand := mkExecCommand(ld, tempDir)
	stdout, stderr, execDuration, err := executor.Run("docker", execCommand...)

	result := Result{
		Stdout:        stdout,
		Stderr:        stderr,
		BuildDuration: buildDuration,
		ExecDuration:  execDuration,
	}

	if err != nil {
		log.Errorf("Execution error: %s", err)
		result.Error = err.Error()
	}
	return result
}
