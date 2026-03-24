package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var binaryPath string

func TestMain(m *testing.M) {
	moduleRoot, err := findModuleRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tmpDir, err := os.MkdirTemp("", "gitfame-bin-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	binaryPath = filepath.Join(tmpDir, "gitfame")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/gitfame")
	buildCmd.Dir = moduleRoot
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", errors.New("could not find module root with go.mod")
		}

		dir = parentDir
	}
}

func TestGitFame(t *testing.T) {
	bundlesDir := filepath.Join("testdata", "bundles")
	testsDir := filepath.Join("testdata", "tests")
	testDirs := listTestDirs(t, testsDir)

	for _, dir := range testDirs {
		tc := readTestCase(t, filepath.Join(testsDir, dir))

		t.Run(dir+"/"+tc.Name, func(t *testing.T) {
			repoDir := t.TempDir()
			unbundle(t, filepath.Join(bundlesDir, tc.Bundle), repoDir)

			headBefore := getHEADRef(t, repoDir)
			args := append([]string{"--repository", repoDir}, tc.Args...)

			cmd := exec.Command(binaryPath, args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			output, err := cmd.Output()
			if tc.Error {
				if err == nil {
					t.Fatalf("expected non-zero exit, got success")
				}

				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("expected ExitError, got %T (%v)", err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected command error: %v\nstderr:\n%s", err, stderr.String())
				}

				compareResults(t, tc.Expected, output, tc.Format)
			}

			headAfter := getHEADRef(t, repoDir)
			if headBefore != headAfter {
				t.Fatalf("repository HEAD changed: before=%q after=%q", headBefore, headAfter)
			}
		})
	}
}

func listTestDirs(t *testing.T, path string) []string {
	t.Helper()

	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatalf("could not read test dirs: %v", err)
	}

	names := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		names = append(names, entry.Name())
	}

	toInt := func(name string) int {
		value, convErr := strconv.Atoi(name)
		if convErr != nil {
			t.Fatalf("test dir %q is not numeric: %v", name, convErr)
		}
		return value
	}

	sort.Slice(names, func(i, j int) bool {
		return toInt(names[i]) < toInt(names[j])
	})

	return names
}

type testCase struct {
	*testDescription
	Expected []byte
}

type testDescription struct {
	Name   string   `yaml:"name"`
	Args   []string `yaml:"args"`
	Bundle string   `yaml:"bundle"`
	Error  bool     `yaml:"error"`
	Format string   `yaml:"format,omitempty"`
}

func readTestCase(t *testing.T, path string) *testCase {
	t.Helper()

	desc := readTestDescription(t, path)
	expected, err := os.ReadFile(filepath.Join(path, "expected.out"))
	if err != nil {
		t.Fatalf("could not read expected.out: %v", err)
	}

	return &testCase{testDescription: desc, Expected: expected}
}

func readTestDescription(t *testing.T, path string) *testDescription {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(path, "description.yaml"))
	if err != nil {
		t.Fatalf("could not read description.yaml: %v", err)
	}

	var desc testDescription
	if err := yaml.Unmarshal(data, &desc); err != nil {
		t.Fatalf("could not parse yaml description: %v", err)
	}

	return &desc
}

func unbundle(t *testing.T, srcBundlePath, dstRepoPath string) {
	t.Helper()

	cmd := exec.Command("git", "clone", srcBundlePath, dstRepoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git clone bundle failed: %v", err)
	}
}

func compareResults(t *testing.T, expected, actual []byte, format string) {
	t.Helper()

	switch format {
	case "json":
		compareJSON(t, expected, actual)
	case "json-lines":
		compareJSONLines(t, expected, actual)
	default:
		if string(expected) != string(actual) {
			t.Fatalf("unexpected output\n--- expected ---\n%s\n--- actual ---\n%s", string(expected), string(actual))
		}
	}
}

func compareJSON(t *testing.T, expected, actual []byte) {
	t.Helper()

	expectedTrimmed := strings.TrimSpace(string(expected))
	actualTrimmed := strings.TrimSpace(string(actual))

	if expectedTrimmed == "" {
		if actualTrimmed != "" {
			t.Fatalf("expected empty output, got %q", actualTrimmed)
		}
		return
	}

	var expectedValue any
	if err := json.Unmarshal([]byte(expectedTrimmed), &expectedValue); err != nil {
		t.Fatalf("could not parse expected json: %v", err)
	}

	var actualValue any
	if err := json.Unmarshal([]byte(actualTrimmed), &actualValue); err != nil {
		t.Fatalf("could not parse actual json: %v\nactual: %s", err, actualTrimmed)
	}

	if !deepJSONEqual(expectedValue, actualValue) {
		t.Fatalf("json mismatch\nexpected: %s\nactual:   %s", expectedTrimmed, actualTrimmed)
	}
}

func compareJSONLines(t *testing.T, expected, actual []byte) {
	t.Helper()

	expectedLines := parseJSONLines(expected)
	actualLines := parseJSONLines(actual)
	if len(expectedLines) != len(actualLines) {
		t.Fatalf("json-lines length mismatch: expected=%d actual=%d", len(expectedLines), len(actualLines))
	}

	for i := range expectedLines {
		compareJSON(t, expectedLines[i], actualLines[i])
	}
}

func parseJSONLines(data []byte) [][]byte {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil
	}

	lines := bytes.Split(trimmed, []byte("\n"))
	out := make([][]byte, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		out = append(out, line)
	}

	return out
}

func deepJSONEqual(a, b any) bool {
	left, errLeft := json.Marshal(a)
	right, errRight := json.Marshal(b)
	if errLeft != nil || errRight != nil {
		return false
	}

	return string(left) == string(right)
}

func getHEADRef(t *testing.T, repoPath string) string {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("could not read HEAD ref: %v", err)
	}

	return strings.TrimSpace(string(output))
}
