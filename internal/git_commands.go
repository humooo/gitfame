package internal

import (
	"os/exec"
	"strings"
)

func GitShow(commit, execDir string) error {
	gitShowCmd := exec.Command("git", "show", commit)
	gitShowCmd.Dir = execDir

	err := gitShowCmd.Run()
	return err
}

func GitLsTree(commit, execDir string) (string, error) {
	gitLsTreeCmd := exec.Command("git", "ls-tree", "-r", "--name-only", commit)
	gitLsTreeCmd.Dir = execDir

	var gitLsTreeCmdOutput strings.Builder
	gitLsTreeCmd.Stdout = &gitLsTreeCmdOutput

	err := gitLsTreeCmd.Run()
	return gitLsTreeCmdOutput.String(), err
}

func GitBlame(commit, pathToFile, execDir string) (string, error) {
	gitBlameCmd := exec.Command("git", "blame", "--line-porcelain", "-b", commit, pathToFile)
	gitBlameCmd.Dir = execDir

	var gitBlameCmdOutput strings.Builder
	gitBlameCmd.Stdout = &gitBlameCmdOutput

	err := gitBlameCmd.Run()
	return gitBlameCmdOutput.String(), err
}

func GitLog(commit, pathToFile, execDir string) (string, error) {
	gitLogCmd := exec.Command("git", "log", "-p", commit, "--follow", "--", pathToFile)
	gitLogCmd.Dir = execDir

	var gitLogCmdOutput strings.Builder
	gitLogCmd.Stdout = &gitLogCmdOutput

	err := gitLogCmd.Run()
	return gitLogCmdOutput.String(), err
}
