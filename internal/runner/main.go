package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Run(lang Lang, code string) (output string, err error) {
	var lg, ok = langDefinition[lang]
	if !ok {
		return "", fmt.Errorf("unknown language: %s", lang)
	}

	// File
	tmpFile, err := os.CreateTemp("", "tmp-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %s", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(code); err != nil {
		return "", fmt.Errorf("failed to write code to temp file: %s", err)
	}

	var stdout, stderr bytes.Buffer
	var command = []string{"run", "--rm", "-v", fmt.Sprintf("%s:%s", tmpFile.Name(), MAIN_EXEC_FILE)}
	command = append(command, lg.image)
	command = append(command, lg.command...)
	command = append(command, MAIN_EXEC_FILE)

	cmd := exec.Command("docker", command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println("cmd error: ", err)
		return "", fmt.Errorf("stderr: %s", stderr.String())
	}

	return stdout.String(), nil
}
