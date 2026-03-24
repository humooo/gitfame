package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

func (stats *Stats) PrintTabular() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(w, "Name\tLines\tCommits\tFiles"); err != nil {
		return err
	}

	for _, row := range stats.SortedData {
		if _, err := fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", row.Name, row.Lines, row.Commits, row.Files); err != nil {
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
	for _, row := range stats.SortedData {
		rows = append(rows, []string{
			row.Name,
			strconv.Itoa(row.Lines),
			strconv.Itoa(row.Commits),
			strconv.Itoa(row.Files),
		})
	}

	return writer.WriteAll(rows)
}

func (stats *Stats) PrintJSON() error {
	jsonData, err := json.Marshal(stats.SortedData)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(jsonData))
	return err
}

func (stats *Stats) PrintJSONLines() error {
	for _, row := range stats.SortedData {
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
