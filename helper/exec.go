package helper

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type ExecResult struct {
	Stdout    string
	Stderr    string
	ExitCode  int
	IsTimeout bool
}

func GetCommandOutput(cmdAndArgs []string, timeout time.Duration) (*ExecResult, error) {
	if len(cmdAndArgs) == 0 {
		return nil, fmt.Errorf("command cannot be empty")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create command with context
	var cmd *exec.Cmd
	if len(cmdAndArgs) == 1 {
		cmd = exec.CommandContext(ctx, cmdAndArgs[0])
	} else {
		cmd = exec.CommandContext(ctx, cmdAndArgs[0], cmdAndArgs[1:]...)
	}

	// Create buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()

	result := &ExecResult{
		Stdout:    strings.TrimSuffix(stdout.String(), "\n"),
		Stderr:    strings.TrimSuffix(stderr.String(), "\n"),
		ExitCode:  0,
		IsTimeout: false,
	}

	// Check if the command timed out
	if ctx.Err() == context.DeadlineExceeded {
		result.IsTimeout = true
		result.ExitCode = -1
		return result, nil
	}

	// Handle other errors and extract exit code
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Command ran but failed
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				result.ExitCode = status.ExitStatus()
			} else {
				result.ExitCode = 1
			}
		} else {
			// Command couldn't be started
			result.ExitCode = -1
			return result, err
		}
	}

	return result, nil
}
