package runner

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type CommandExecutorOS struct {
	timeout time.Duration
}

func (r *CommandExecutorOS) Run(name string, args ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Infof("Running command: \"%s %s\"", name, strings.Join(args, " "))

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), errors.New(ERR_TIMEOUT_EXCEEDED)
	}

	return stdout.String(), stderr.String(), err
}
