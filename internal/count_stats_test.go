package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBlamePorcelain_UsesCachedIdentity(t *testing.T) {
	commitA := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	commitB := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	blameOutput := commitA + " 1 1 1\n" +
		"author Alice\n" +
		"author-mail <alice@example.com>\n" +
		"author-time 1\n" +
		"author-tz +0000\n" +
		"committer Carol\n" +
		"committer-mail <carol@example.com>\n" +
		"committer-time 1\n" +
		"committer-tz +0000\n" +
		"summary msg\n" +
		"filename file.txt\n" +
		"\tline1\n" +
		commitB + " 2 2 2\n" +
		"author Bob\n" +
		"author-mail <bob@example.com>\n" +
		"author-time 2\n" +
		"author-tz +0000\n" +
		"committer Bob\n" +
		"committer-mail <bob@example.com>\n" +
		"committer-time 2\n" +
		"committer-tz +0000\n" +
		"summary msg\n" +
		"filename file.txt\n" +
		"\tline2\n" +
		commitB + " 3 3\n" +
		"\tline3\n" +
		commitA + " 4 4\n" +
		"\tline4\n"

	authorBlameLines, err := parseBlamePorcelain(blameOutput, false)
	if err != nil {
		t.Fatalf("parseBlamePorcelain(author) failed: %v", err)
	}
	if len(authorBlameLines) != 4 {
		t.Fatalf("expected 4 blamed lines, got %d", len(authorBlameLines))
	}

	authorNames := []string{authorBlameLines[0].Person, authorBlameLines[1].Person, authorBlameLines[2].Person, authorBlameLines[3].Person}
	expectedAuthors := []string{"Alice", "Bob", "Bob", "Alice"}
	for index := range expectedAuthors {
		if authorNames[index] != expectedAuthors[index] {
			t.Fatalf("unexpected author at index %d: want %q, got %q", index, expectedAuthors[index], authorNames[index])
		}
	}

	committerBlameLines, err := parseBlamePorcelain(blameOutput, true)
	if err != nil {
		t.Fatalf("parseBlamePorcelain(committer) failed: %v", err)
	}

	committerNames := []string{committerBlameLines[0].Person, committerBlameLines[1].Person, committerBlameLines[2].Person, committerBlameLines[3].Person}
	expectedCommitters := []string{"Carol", "Bob", "Bob", "Carol"}
	for index := range expectedCommitters {
		if committerNames[index] != expectedCommitters[index] {
			t.Fatalf("unexpected committer at index %d: want %q, got %q", index, expectedCommitters[index], committerNames[index])
		}
	}
}

func TestCountStatistics_EmptyFileUsesCommitter(t *testing.T) {
	repoDir := t.TempDir()
	runCommand(t, repoDir, nil, "git", "init", "-q")

	emptyPath := filepath.Join(repoDir, "empty.txt")
	if err := os.WriteFile(emptyPath, []byte{}, 0o644); err != nil {
		t.Fatalf("could not create empty file: %v", err)
	}

	runCommand(t, repoDir, nil, "git", "add", "empty.txt")
	runCommand(t, repoDir, map[string]string{
		"GIT_AUTHOR_NAME":     "Author Empty",
		"GIT_AUTHOR_EMAIL":    "author-empty@example.com",
		"GIT_COMMITTER_NAME":  "Committer Empty",
		"GIT_COMMITTER_EMAIL": "committer-empty@example.com",
	}, "git", "commit", "-q", "-m", "add empty")

	fp := &FilesParams{
		FilesList: []string{"empty.txt"},
		Cla: &CommandLineArgs{
			RepositoryPath: repoDir,
			CommitPointer:  "HEAD",
			UseCommitter:   true,
		},
	}

	stats, err := CountStatistics(fp)
	if err != nil {
		t.Fatalf("CountStatistics failed: %v", err)
	}

	summary, ok := stats.CombinedData["Committer Empty"]
	if !ok {
		t.Fatalf("expected committer entry for empty file")
	}

	if summary != [3]int{0, 1, 1} {
		t.Fatalf("unexpected stats for committer: got %v, want [0 1 1]", summary)
	}

	if _, authorPresent := stats.CombinedData["Author Empty"]; authorPresent {
		t.Fatalf("did not expect author entry when --use-committer is enabled")
	}
}
