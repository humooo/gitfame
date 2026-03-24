package internal

import (
	"strings"
	"testing"
)

func TestGetCommandLineArgs_InvalidExcludePatternReturnsError(t *testing.T) {
	repoDir := createRepoWithSingleCommit(t)

	cla := NewCommandLineArgs()
	err := cla.GetCommandLineArgs([]string{
		"--repository", repoDir,
		"--exclude", "[",
	})
	if err == nil {
		t.Fatalf("expected error for invalid glob pattern")
	}

	if !strings.Contains(err.Error(), "invalid glob pattern") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetCommandLineArgs_InvalidOrderByContainsValue(t *testing.T) {
	repoDir := createRepoWithSingleCommit(t)

	cla := NewCommandLineArgs()
	err := cla.GetCommandLineArgs([]string{
		"--repository", repoDir,
		"--order-by", "wrong",
	})
	if err == nil {
		t.Fatalf("expected error for invalid order-by")
	}

	if !strings.Contains(err.Error(), "wrong") {
		t.Fatalf("expected invalid value in error, got: %v", err)
	}
}
