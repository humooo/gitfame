package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

func (stats *Stats) PrintTabular() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, err := fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
	ProcessError(err, "PrintTabular")

	for _, line := range stats.SortedData {
		_, err = fmt.Fprintln(w, line[0]+"\t"+line[1]+"\t"+line[2]+"\t"+line[3])
		ProcessError(err, "PrintTabular")
	}

	err = w.Flush()
	ProcessError(err, "PrintTabular")
}

func (stats *Stats) PrintCSV() {
	header := []string{"Name", "Lines", "Commits", "Files"}
	w := csv.NewWriter(os.Stdout)
	var buff [][]string
	buff = append(buff, header)
	for _, line := range stats.SortedData {
		buff = append(buff, []string{line[0], line[1], line[2], line[3]})
	}
	err := w.WriteAll(buff)
	ProcessError(err, "PrintCSV")
}

func (stats *Stats) PrintJSON() {
	var buff []map[string]interface{}
	for _, line := range stats.SortedData {
		lines, err := strconv.Atoi(line[1])
		ProcessError(err, "PrintJSON: could not convert num of lines")

		commits, err := strconv.Atoi(line[2])
		ProcessError(err, "PrintJSON: could not convert num of commits")

		files, err := strconv.Atoi(line[3])
		ProcessError(err, "PrintJSON: could not convert num of files")

		buff = append(buff, map[string]interface{}{
			"name":    line[0],
			"lines":   lines,
			"commits": commits,
			"files":   files,
		})
	}
	jsonData, err := json.Marshal(buff)
	ProcessError(err, "PrintJSON: could not marshal json")

	fmt.Println(string(jsonData))
}

func (stats *Stats) PrintJSONLines() {
	for _, line := range stats.SortedData {
		lines, err := strconv.Atoi(line[1])
		ProcessError(err, "PrintJSONLines: could not convert num of lines")

		commits, err := strconv.Atoi(line[2])
		ProcessError(err, "PrintJSONLines: could not convert num of commits")

		files, err := strconv.Atoi(line[3])
		ProcessError(err, "PrintJSONLines: could not convert num of files")

		jsonLine, err := json.Marshal(map[string]interface{}{
			"name":    line[0],
			"lines":   lines,
			"commits": commits,
			"files":   files,
		})
		ProcessError(err, "PrintJSONLines: could not marshal json")

		fmt.Println(string(jsonLine))
	}
}
