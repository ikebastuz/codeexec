package runner

import (
	"fmt"
	"os"
	"path/filepath"
)

type TempDirFS struct {
	dir string
}

func (t *TempDirFS) Create(fileName string, sourceCode string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "tmp-app-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	t.dir = tmpDir
	sourceFilePath := filepath.Join(tmpDir, fileName)
	if err := os.WriteFile(sourceFilePath, []byte(sourceCode), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to write code to source file: %w", err)
	}
	return tmpDir, nil
}

func (t *TempDirFS) Cleanup() error {
	if t.dir == "" {
		return nil
	}
	err := os.RemoveAll(t.dir)
	t.dir = ""
	return err
}
