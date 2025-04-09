package brew

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
)

const brewBinary = "brew"

func commandOutputContext(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("open stderr pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("open stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start command: %w", err)
	}

	scanner := bufio.NewScanner(stderr)
	errorMsg := ""
	for scanner.Scan() {
		errorMsg += scanner.Text()
	}

	scanner = bufio.NewScanner(stdout)
	output := ""
	for scanner.Scan() {
		output += scanner.Text()
	}

	if err := cmd.Wait(); err != nil {
		if len(errorMsg) > 0 {
			// nolint: err113
			return "", errors.New(errorMsg)
		}
		return "", fmt.Errorf("finish command :%w", err)
	}

	return output, nil
}
