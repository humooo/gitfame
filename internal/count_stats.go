package internal

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/humooo/gitfame/internal/utils/constants"
)

type Stats struct {
	UserToLines      map[string]int
	UserToCommits    map[string]map[string]bool
	UserToNumCommits map[string]int
	UserToFiles      map[string]map[string]bool
	UserToNumFiles   map[string]int
	CombinedData     map[string][3]int
	SortedData       []outputRow
}

type blameLine struct {
	CommitHash string
	Person     string
}

type outputRow struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func CountStatistics(fp *FilesParams) (Stats, error) {
	stats := Stats{
		UserToLines:      make(map[string]int),
		UserToCommits:    make(map[string]map[string]bool),
		UserToNumCommits: make(map[string]int),
		UserToFiles:      make(map[string]map[string]bool),
		UserToNumFiles:   make(map[string]int),
		CombinedData:     make(map[string][3]int),
	}

	for _, path := range fp.FilesList {
		if err := stats.ProcessFile(path, fp.Cla.RepositoryPath, fp.Cla.CommitPointer, fp.Cla.UseCommitter); err != nil {
			return Stats{}, err
		}
	}

	stats.CombineResults()
	return stats, nil
}

func (stats *Stats) AddLine(person string) {
	stats.UserToLines[person]++
}

func (stats *Stats) AddCommit(person, commit string) {
	if _, ok := stats.UserToCommits[person]; !ok {
		stats.UserToCommits[person] = make(map[string]bool)
	}

	if _, ok := stats.UserToCommits[person][commit]; !ok {
		stats.UserToCommits[person][commit] = true
		stats.UserToNumCommits[person]++
	}
}

func (stats *Stats) AddFile(person, path string) {
	if _, ok := stats.UserToFiles[person]; !ok {
		stats.UserToFiles[person] = make(map[string]bool)
	}

	if _, ok := stats.UserToFiles[person][path]; !ok {
		stats.UserToFiles[person][path] = true
		stats.UserToNumFiles[person]++
	}
}

func (stats *Stats) ProcessFile(path, gitDir, commitPointer string, useCommitter bool) error {
	blameOutput, err := GitBlame(commitPointer, path, gitDir)
	if err != nil {
		return fmt.Errorf("git blame failed for %q: %w", path, err)
	}

	if strings.TrimSpace(blameOutput) == "" {
		commitHash, author, committer, err := GitLogLastIdentity(commitPointer, path, gitDir)
		if err != nil {
			return fmt.Errorf("git log failed for empty file %q: %w", path, err)
		}

		person := author
		if useCommitter {
			person = committer
		}

		stats.AddCommit(person, commitHash)
		stats.AddFile(person, path)
		return nil
	}

	blameLines, err := parseBlamePorcelain(blameOutput, useCommitter)
	if err != nil {
		return fmt.Errorf("could not parse git blame output for %q: %w", path, err)
	}

	for _, blameLine := range blameLines {
		stats.AddLine(blameLine.Person)
		stats.AddCommit(blameLine.Person, blameLine.CommitHash)
		stats.AddFile(blameLine.Person, path)
	}

	return nil
}

func parseBlamePorcelain(blameOutput string, useCommitter bool) ([]blameLine, error) {
	lines := strings.Split(blameOutput, "\n")
	commitToPerson := make(map[string]string)
	parsed := make([]blameLine, 0)

	for i := 0; i < len(lines); {
		line := lines[i]
		if line == "" {
			i++
			continue
		}

		commitHash, isHeader := parseBlameHeader(line)
		if !isHeader {
			return nil, fmt.Errorf("unexpected blame line: %q", line)
		}

		i++
		person, cached := commitToPerson[commitHash]
		sawContentLine := false
		for i < len(lines) {
			metadataLine := lines[i]
			if strings.HasPrefix(metadataLine, "\t") {
				sawContentLine = true
				i++
				break
			}

			if !cached {
				if !useCommitter && strings.HasPrefix(metadataLine, "author ") {
					person = strings.TrimPrefix(metadataLine, "author ")
				}
				if useCommitter && strings.HasPrefix(metadataLine, "committer ") {
					person = strings.TrimPrefix(metadataLine, "committer ")
				}
			}

			i++
		}

		if !sawContentLine {
			return nil, fmt.Errorf("blame block for commit %s does not contain source line", commitHash)
		}

		if person == "" {
			if cachedPerson, ok := commitToPerson[commitHash]; ok {
				person = cachedPerson
			} else {
				return nil, fmt.Errorf("missing %s for commit %s", blameIdentityLabel(useCommitter), commitHash)
			}
		}

		commitToPerson[commitHash] = person
		parsed = append(parsed, blameLine{CommitHash: commitHash, Person: person})
	}

	return parsed, nil
}

func parseBlameHeader(line string) (string, bool) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return "", false
	}

	if !isHexHash(fields[0]) {
		return "", false
	}

	if _, err := strconv.Atoi(fields[1]); err != nil {
		return "", false
	}

	if _, err := strconv.Atoi(fields[2]); err != nil {
		return "", false
	}

	return fields[0], true
}

func isHexHash(value string) bool {
	if len(value) != 40 {
		return false
	}

	for _, symbol := range value {
		isDigit := symbol >= '0' && symbol <= '9'
		isLowerHex := symbol >= 'a' && symbol <= 'f'
		if !isDigit && !isLowerHex {
			return false
		}
	}

	return true
}

func blameIdentityLabel(useCommitter bool) string {
	if useCommitter {
		return "committer"
	}

	return "author"
}

func (stats *Stats) CombineResults() {
	for name, numCommits := range stats.UserToNumCommits {
		numLines := 0

		if actualNumLines, ok := stats.UserToLines[name]; ok {
			numLines = actualNumLines
		}

		stats.CombinedData[name] = [3]int{
			numLines,
			numCommits,
			stats.UserToNumFiles[name],
		}
	}
}

func (stats *Stats) SortResults(sortKey constants.OrderKey) {
	users := make([]string, 0, len(stats.UserToNumCommits))
	for user := range stats.UserToNumCommits {
		if user != "Not Committed Yet" {
			users = append(users, user)
		}
	}

	priority := sortPriority(sortKey)
	sort.SliceStable(users, func(i, j int) bool {
		return stats.userLess(users[i], users[j], priority)
	})

	sortedStats := make([]outputRow, 0, len(users))
	for _, user := range users {
		userData := stats.CombinedData[user]
		sortedStats = append(sortedStats, outputRow{
			Name:    user,
			Lines:   userData[linesIdx],
			Commits: userData[commitsIdx],
			Files:   userData[filesIdx],
		})
	}

	stats.SortedData = sortedStats
}

const (
	linesIdx = iota
	commitsIdx
	filesIdx
)

func sortPriority(sortKey constants.OrderKey) [3]int {
	switch sortKey {
	case constants.Commits:
		return [3]int{commitsIdx, linesIdx, filesIdx}
	case constants.Files:
		return [3]int{filesIdx, linesIdx, commitsIdx}
	default:
		return [3]int{linesIdx, commitsIdx, filesIdx}
	}
}

func (stats *Stats) userLess(leftUser, rightUser string, priority [3]int) bool {
	left := stats.CombinedData[leftUser]
	right := stats.CombinedData[rightUser]

	for _, metric := range priority {
		if left[metric] != right[metric] {
			return left[metric] > right[metric]
		}
	}

	return leftUser < rightUser
}
