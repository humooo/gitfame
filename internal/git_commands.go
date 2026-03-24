package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

func GitCommitExists(revision, execDir string) error {
	gitShowCmd := exec.Command("git", "rev-parse", "--verify", revision+"^{commit}")
	gitShowCmd.Dir = execDir

	return gitShowCmd.Run()
}

func GitLsTree(commit, execDir string) ([]string, error) {
	gitLsTreeCmd := exec.Command("git", "ls-tree", "-r", "--name-only", commit)
	gitLsTreeCmd.Dir = execDir

	output, err := gitLsTreeCmd.Output()
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []string{}, nil
	}

	return strings.Split(trimmed, "\n"), nil
}

func GitBlame(commit, pathToFile, execDir string) (string, error) {
	gitBlameCmd := exec.Command("git", "blame", "--porcelain", commit, "--", pathToFile)
	gitBlameCmd.Dir = execDir

	output, err := gitBlameCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GitLogLastIdentity(commit, pathToFile, execDir string) (string, string, string, error) {
	gitLogCmd := exec.Command("git", "log", "-1", "--format=%H%x00%an%x00%cn", commit, "--", pathToFile)
	gitLogCmd.Dir = execDir

	output, err := gitLogCmd.Output()
	if err != nil {
		return "", "", "", err
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return "", "", "", fmt.Errorf("no commits found for file %q", pathToFile)
	}

	parts := strings.SplitN(trimmed, "\x00", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("unexpected git log output for file %q", pathToFile)
	}

	return parts[0], parts[1], parts[2], nil
}
