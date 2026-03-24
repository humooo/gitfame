package integration

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIStreamRouting(t *testing.T) {
	t.Run("help_to_stdout", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--help")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("help command failed: %v", err)
		}

		if !strings.Contains(stdout.String(), "Usage of gitfame:") {
			t.Fatalf("help output is missing usage header:\n%s", stdout.String())
		}

		if strings.TrimSpace(stderr.String()) != "" {
			t.Fatalf("expected empty stderr for help, got:\n%s", stderr.String())
		}
	})

	t.Run("parse_errors_to_stderr", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--unknown-flag-for-test")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected parse error exit status, got success")
		}

		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("expected ExitError, got %T (%v)", err, err)
		}

		if strings.TrimSpace(stdout.String()) != "" {
			t.Fatalf("expected empty stdout for parse error, got:\n%s", stdout.String())
		}

		stderrText := stderr.String()
		if !strings.Contains(stderrText, "unknown flag: --unknown-flag-for-test") {
			t.Fatalf("expected unknown flag message in stderr, got:\n%s", stderrText)
		}
	})
}
