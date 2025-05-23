package runner

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var tempDir = "tmpDir"
var tempDirMount = fmt.Sprintf("%s:%s", tempDir, "/app")

func TestMakeCommandForCompiled(t *testing.T) {
	lang := LangDefinitions["go"]

	t.Run("should create correct docker build command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, true)
		want := []string{
			"run", "--rm", "--pull=never", "-w", "/app", "-v", tempDirMount, lang.image, "go", "build", "-o", "main", "main.go",
		}

		assertState(t, got, want)
	})
	t.Run("should create correct docker execution command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, false)
		want := []string{
			"run", "--rm", "--pull=never", "-w", "/app", "-v", tempDirMount, lang.image, "./main",
		}

		assertState(t, got, want)
	})
}

func TestMakeCommandForInterpreted(t *testing.T) {
	lang := LangDefinitions["javascript"]

	t.Run("should create correct docker execution command", func(t *testing.T) {
		got := mkCommand(lang, tempDir, false)
		want := []string{
			"run", "--rm", "--pull=never", "-w", "/app", "-v", tempDirMount, lang.image, "node", "main.js",
		}

		assertState(t, got, want)
	})
}

func assertState[T any](t testing.TB, got, want T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
