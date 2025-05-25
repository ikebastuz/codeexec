package runner

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunWithDeps(t *testing.T) {
	t.Run("should return error for unknown language", func(t *testing.T) {
		fs := &TempDirStub{t: t}
		executor := &CommandExecutorStub{}
		result := NewRunner("qwe", "").runWithDeps(fs, executor)

		assertState(t, result, Result{
			Stdout:        "",
			Stderr:        "",
			BuildDuration: 0,
			ExecDuration:  0,
			Error:         "unknown language: qwe",
		})
	})
}

func TestDockerBaseCommand(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("should create correct base docker command", func(t *testing.T) {
		got := mkDockerBaseCommand(tempDir)
		want := []string{
			"run", "--rm",
			"--pull=never",
			"-w", WORKDIR,
			"-v", fmt.Sprintf("%s:%s", tempDir, WORKDIR),
		}

		assertState(t, got, want)
	})
}

func TestBuildCommand(t *testing.T) {
	tempDir := t.TempDir()
	t.Run("should create correct build command for compiled language", func(t *testing.T) {
		lang := LangDefinitions["go"]
		got := mkBuildCommand(lang, tempDir)
		want := append(
			mkDockerBaseCommand(tempDir),
			[]string{lang.image, "go", "build", "-o", "main", "main.go"}...,
		)

		assertState(t, got, want)
	})

	t.Run("should return nil for interpreted language", func(t *testing.T) {
		lang := LangDefinitions["javascript"]
		got := mkBuildCommand(lang, tempDir)
		if got != nil {
			t.Errorf("expected nil for interpreted language, got %v", got)
		}
	})
}

func TestExecCommand(t *testing.T) {
	tempDir := t.TempDir()
	t.Run("should create correct exec command for compiled language", func(t *testing.T) {
		lang := LangDefinitions["go"]
		got := mkExecCommand(lang, tempDir)
		want := append(
			mkDockerBaseCommand(tempDir),
			[]string{lang.image, "./main"}...,
		)

		assertState(t, got, want)
	})

	t.Run("should create correct exec command for interpreted language", func(t *testing.T) {
		lang := LangDefinitions["javascript"]
		got := mkExecCommand(lang, tempDir)
		want := append(
			mkDockerBaseCommand(tempDir),
			[]string{lang.image, "node", "main.js"}...,
		)

		assertState(t, got, want)
	})
}

func assertState[T any](t testing.TB, got, want T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

type TempDirStub struct {
	t          *testing.T
	createdDir bool
	cleanedDir bool
}

func (t *TempDirStub) Create(fileName string, sourceCode string) (string, error) {
	tempDir := t.t.TempDir()
	t.createdDir = true
	return tempDir, nil
}

func (t *TempDirStub) Cleanup() error {
	t.cleanedDir = true
	return nil
}

type CommandExecutorStub struct{}

func (c *CommandExecutorStub) Run(name string, args ...string) (string, string, float64, error) {
	return "", "", 0, nil
}
