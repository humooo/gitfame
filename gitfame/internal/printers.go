package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type outputRow struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func parseOutputRow(line [4]string) (outputRow, error) {
	lines, err := strconv.Atoi(line[1])
	if err != nil {
		return outputRow{}, fmt.Errorf("could not convert line count %q: %w", line[1], err)
	}

	commits, err := strconv.Atoi(line[2])
	if err != nil {
		return outputRow{}, fmt.Errorf("could not convert commit count %q: %w", line[2], err)
	}

	files, err := strconv.Atoi(line[3])
	if err != nil {
		return outputRow{}, fmt.Errorf("could not convert file count %q: %w", line[3], err)
	}

	return outputRow{
		Name:    line[0],
		Lines:   lines,
		Commits: commits,
		Files:   files,
	}, nil
}

func (stats *Stats) PrintTabular() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(w, "Name\tLines\tCommits\tFiles"); err != nil {
		return err
	}

	for _, line := range stats.SortedData {
		if _, err := fmt.Fprintln(w, line[0]+"\t"+line[1]+"\t"+line[2]+"\t"+line[3]); err != nil {
			return err
		}
	}

	return w.Flush()
}

func (stats *Stats) PrintCSV() error {
	header := []string{"Name", "Lines", "Commits", "Files"}
	writer := csv.NewWriter(os.Stdout)

	rows := make([][]string, 0, len(stats.SortedData)+1)
	rows = append(rows, header)
	for _, line := range stats.SortedData {
		rows = append(rows, []string{line[0], line[1], line[2], line[3]})
	}

	return writer.WriteAll(rows)
}

func (stats *Stats) PrintJSON() error {
	rows := make([]outputRow, 0, len(stats.SortedData))
	for _, line := range stats.SortedData {
		row, err := parseOutputRow(line)
		if err != nil {
			return err
		}

		rows = append(rows, row)
	}

	jsonData, err := json.Marshal(rows)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(jsonData))
	return err
}

func (stats *Stats) PrintJSONLines() error {
	for _, line := range stats.SortedData {
		row, err := parseOutputRow(line)
		if err != nil {
			return err
		}

		jsonLine, err := json.Marshal(row)
		if err != nil {
			return fmt.Errorf("could not marshal json line: %w", err)
		}

		if _, err := fmt.Fprintln(os.Stdout, string(jsonLine)); err != nil {
			return err
		}
	}

	return nil
}
