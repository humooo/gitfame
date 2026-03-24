package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runCommand(t *testing.T, dir string, extraEnv map[string]string, name string, args ...string) string {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for key, value := range extraEnv {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %v\noutput: %s\nerror: %v", name, args, string(output), err)
	}

	return string(output)
}

func createRepoWithSingleCommit(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	runCommand(t, repoDir, nil, "git", "init", "-q")

	filePath := filepath.Join(repoDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("line\n"), 0o644); err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	runCommand(t, repoDir, nil, "git", "add", "file.txt")
	runCommand(t, repoDir, map[string]string{
		"GIT_AUTHOR_NAME":     "Author",
		"GIT_AUTHOR_EMAIL":    "author@example.com",
		"GIT_COMMITTER_NAME":  "Committer",
		"GIT_COMMITTER_EMAIL": "committer@example.com",
	}, "git", "commit", "-q", "-m", "init")

	return repoDir
}
