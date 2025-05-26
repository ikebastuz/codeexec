package runner

import (
	"codeexec/internal/config"
	"codeexec/internal/db"
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func NewRunner(lang Lang, code string) *Runner {
	return &Runner{lang: lang, code: code}
}

func (r *Runner) Run() Result {
	fs := &TempDirFS{}
	executor := &CommandExecutorOS{timeout: config.PROCESS_TIMEOUT}

	encoded := EncodeSource(r.code)
	queries := db.GetDB().GetQueries()

	// TODO: check if we can use request context
	cached, err := queries.GetCodeExecutionResult(context.Background(), db.GetCodeExecutionResultParams{
		EncodedCode: encoded,
		Language:    string(r.lang),
	})

	if err == nil {
		return Result{
			Stdout:        db.NullStringToString(cached.Stdout),
			Stderr:        db.NullStringToString(cached.Stderr),
			ExecDuration:  db.NullFloatToFloat(cached.ExecDuration),
			BuildDuration: db.NullFloatToFloat(cached.BuildDuration),
			Error:         db.NullStringToString(cached.Error),
		}
	}

	result := r.runWithDeps(fs, executor)

	// TODO: check if we can use request context
	queries.InsertCodeExecutionResult(context.Background(), db.InsertCodeExecutionResultParams{
		Code:          r.code,
		Language:      string(r.lang),
		EncodedCode:   encoded,
		Stdout:        db.ToNullString(result.Stdout),
		Stderr:        db.ToNullString(result.Stderr),
		Error:         db.ToNullString(result.Error),
		BuildDuration: sql.NullFloat64{Float64: result.BuildDuration, Valid: true},
		ExecDuration:  sql.NullFloat64{Float64: result.ExecDuration, Valid: true},
	})

	return result
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
