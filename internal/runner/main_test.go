package runner

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeExecCommand(t *testing.T) {
	t.Run("should create correct docker execution command", func(t *testing.T) {
		langNode := LangDefinitions["javascript"]
		tempFileName := "tmpFileName"

		got := makeExecCommand(langNode, tempFileName)
		want := []string{"run", "--rm", "--pull=never", "-v", "tmpFileName:/app/main.js", "node:20-alpine", "node", "/app/main.js"}

		assertState(t, got, want)
	})
}

func assertState[T any](t testing.TB, got, want T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
