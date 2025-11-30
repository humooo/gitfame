package internal

import (
	"sort"
	"strconv"
	"strings"

	"github.com/humooo/urlshortener/internal/utils/constants"
)

type Stats struct {
	UserToLines      map[string]int
	UserToCommits    map[string]map[string]bool
	UserToNumCommits map[string]int
	UserToFiles      map[string]map[string]bool
	UserToNumFiles   map[string]int
	CombinedData     map[string][3]int
	SortedData       [][4]string
}

var totalStats Stats

func CountStatistics(fp *FilesParams) Stats {
	totalStats = Stats{
		UserToLines:      make(map[string]int),
		UserToCommits:    make(map[string]map[string]bool),
		UserToNumCommits: make(map[string]int),
		UserToFiles:      make(map[string]map[string]bool),
		UserToNumFiles:   make(map[string]int),
		CombinedData:     make(map[string][3]int),
	}

	for _, path := range fp.FilesList {
		ProcessFile(path, *fp.Cla.RepositoryPath, *fp.Cla.CommitPointer, *fp.Cla.UseCommiter)
	}

	totalStats.CombineResults()
	return totalStats
}

func AddLine(author string) {
	totalStats.UserToLines[author]++
}

func AddCommit(author, commit string) {
	if _, ok := totalStats.UserToCommits[author]; !ok {
		totalStats.UserToCommits[author] = make(map[string]bool)
	}
	if _, ok := totalStats.UserToCommits[author][commit]; !ok {
		totalStats.UserToCommits[author][commit] = true
		totalStats.UserToNumCommits[author]++
	}
}

func AddFile(author, path string) {
	if _, ok := totalStats.UserToFiles[author]; !ok {
		totalStats.UserToFiles[author] = make(map[string]bool)
	}
	if _, ok := totalStats.UserToFiles[author][path]; !ok {
		totalStats.UserToFiles[author][path] = true
		totalStats.UserToNumFiles[author]++
	}
}

func ProcessFile(path, gitDir, commitPointer string, useCommiter bool) {
	commitersLog, err := GitBlame(commitPointer, path, gitDir)
	ProcessError(err, "ProcessFile")

	statLines := strings.Split(commitersLog, "\n")
	author := ""
	commitHash := ""

	if len(statLines) == 1 && statLines[0] == "" {
		// Empty file
		gitLog, err := GitLog(commitPointer, path, gitDir)
		ProcessError(err, "ProcessFile")

		logLines := strings.Split(gitLog, "\n")
		commitHash = strings.Split(logLines[0], " ")[1]
		words := strings.Split(logLines[1], " ")
		author = strings.Join(words[1:len(words)-1], " ")

		AddCommit(author, commitHash)
		AddFile(author, path)
	}

	for i := 0; i < len(statLines); i++ {
		words := strings.Split(statLines[i], " ")

		if useCommiter {
			if words[0] == "committer" {
				commitHash = strings.Split(statLines[i-5], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		} else {
			if words[0] == "author" {
				commitHash = strings.Split(statLines[i-1], " ")[0]
				author = strings.Join(words[1:], " ")
			} else {
				continue
			}
		}

		AddLine(author)
		AddCommit(author, commitHash)
		AddFile(author, path)
	}
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
	var users []string
	for user := range stats.UserToNumCommits {
		if user != "Not Committed Yet" {
			users = append(users, user)
		}
	}

	switch sortKey {
	case constants.Lines:
		sort.SliceStable(users, func(i, j int) bool {
			if stats.CombinedData[users[i]][0] == stats.CombinedData[users[j]][0] {
				if stats.CombinedData[users[i]][1] == stats.CombinedData[users[j]][1] {
					if stats.CombinedData[users[i]][2] == stats.CombinedData[users[j]][2] {
						return users[i] < users[j]
					}
					return stats.CombinedData[users[i]][2] > stats.CombinedData[users[j]][2]
				}
				return stats.CombinedData[users[i]][1] > stats.CombinedData[users[j]][1]
			}
			return stats.CombinedData[users[i]][0] > stats.CombinedData[users[j]][0]
		})
	case constants.Commits:
		sort.SliceStable(users, func(i, j int) bool {
			if stats.CombinedData[users[i]][1] == stats.CombinedData[users[j]][1] {
				if stats.CombinedData[users[i]][0] == stats.CombinedData[users[j]][0] {
					if stats.CombinedData[users[i]][2] == stats.CombinedData[users[j]][2] {
						return users[i] < users[j]
					}
					return stats.CombinedData[users[i]][2] > stats.CombinedData[users[j]][2]
				}
				return stats.CombinedData[users[i]][0] > stats.CombinedData[users[j]][0]
			}
			return stats.CombinedData[users[i]][1] > stats.CombinedData[users[j]][1]
		})
	case constants.Files:
		sort.SliceStable(users, func(i, j int) bool {
			if stats.CombinedData[users[i]][2] == stats.CombinedData[users[j]][2] {
				if stats.CombinedData[users[i]][0] == stats.CombinedData[users[j]][0] {
					if stats.CombinedData[users[i]][1] == stats.CombinedData[users[j]][1] {
						return users[i] < users[j]
					}
					return stats.CombinedData[users[i]][1] > stats.CombinedData[users[j]][1]
				}
				return stats.CombinedData[users[i]][0] > stats.CombinedData[users[j]][0]
			}
			return stats.CombinedData[users[i]][2] > stats.CombinedData[users[j]][2]
		})
	}

	var sortedStats [][4]string

	for _, user := range users {
		sortedStats = append(sortedStats, [4]string{user, strconv.Itoa(stats.CombinedData[user][0]), strconv.Itoa(stats.CombinedData[user][1]), strconv.Itoa(stats.CombinedData[user][2])})
	}
	stats.SortedData = sortedStats
}
