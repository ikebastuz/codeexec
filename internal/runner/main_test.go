package runner

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tempDir = "tmpDir"
var tempDirMount = fmt.Sprintf("%s:%s", tempDir, "/app")
var baseCommand = []string{"run", "--rm", "--pull=never", "-w", WORKDIR, "-v", tempDirMount}

func TestMakeCommandForCompiled(t *testing.T) {
	lang := LangDefinitions["go"]

	t.Run("should create correct docker build command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, true)
		want := append(
			baseCommand,
			[]string{lang.image, "go", "build", "-o", "main", "main.go"}...,
		)

		assertState(t, got, want)
	})
	t.Run("should create correct docker execution command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, false)
		want := append(
			baseCommand,
			[]string{lang.image, "./main"}...,
		)

		assertState(t, got, want)
	})
}

func TestMakeCommandForInterpreted(t *testing.T) {
	lang := LangDefinitions["javascript"]

	t.Run("should create correct docker execution command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, false)
		want := append(
			baseCommand,
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
