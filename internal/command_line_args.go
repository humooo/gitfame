package internal

import (
	"fmt"
	"os"

	"github.com/humooo/urlshortener/internal/utils/constants"
	flag "github.com/spf13/pflag"
)

type CommandLineArgs struct {
	RepositoryPath *string
	CommitPointer  *string
	SortOrderKey   constants.OrderKey
	UseCommiter    *bool
	OutputFormat   constants.OutputFormat
	Extensions     *[]string
	Languages      *[]string
	Exclude        *[]string
	Restricted     *[]string
}

func NewCommandLineArgs() *CommandLineArgs {
	return &CommandLineArgs{}
}

func isValidFilePath(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func commitExists(commitHash, gitDir string) bool {
	return GitShow(commitHash, gitDir) == nil
}

func isValidOrderKey(orderKey constants.OrderKey) bool {
	switch orderKey {
	case constants.Lines, constants.Commits, constants.Files:
		return true
	default:
		return false
	}
}

func isValidOutputFormat(outputFormat constants.OutputFormat) bool {
	switch outputFormat {
	case constants.Tabular, constants.CSV, constants.SimpleJSON, constants.JSONLines:
		return true
	default:
		return false
	}
}

func (cla *CommandLineArgs) GetCommandLineArgs() error {

	cla.RepositoryPath = flag.String("repository", "./", "Repository path.")
	cla.CommitPointer = flag.String("revision", "HEAD", "Pointer to a commit.")
	orderKey := flag.String("order-by", "lines", "Sort key.")
	cla.UseCommiter = flag.Bool("use-committer", false, "Use commiter.")
	format := flag.String("format", "tabular", "Output format.")
	cla.Extensions = flag.StringSlice("extensions", []string{}, "List of Extensions to search.")
	cla.Languages = flag.StringSlice("languages", []string{}, "List of Languages to search.")
	cla.Exclude = flag.StringSlice("exclude", []string{}, "Glob-patterns to Exclude.")
	cla.Restricted = flag.StringSlice("restrict-to", []string{}, "List of restrictions to match.")

	flag.Parse()

	switch *orderKey {
	case "lines":
		cla.SortOrderKey = constants.Lines
	case "commits":
		cla.SortOrderKey = constants.Commits
	case "files":
		cla.SortOrderKey = constants.Files
	}

	switch *format {
	case "tabular":
		cla.OutputFormat = constants.Tabular
	case "csv":
		cla.OutputFormat = constants.CSV
	case "json":
		cla.OutputFormat = constants.SimpleJSON
	case "json-lines":
		cla.OutputFormat = constants.JSONLines
	}

	if !isValidFilePath(*cla.RepositoryPath) {
		return fmt.Errorf("invalid repository path: %s", *cla.RepositoryPath)
	}
	if !commitExists(*cla.CommitPointer, *cla.RepositoryPath) {
		return fmt.Errorf("no such commit: %s", *cla.CommitPointer)
	}
	if !isValidOrderKey(cla.SortOrderKey) {
		return fmt.Errorf("invalid sort order key: %s. Should be one of: 'lines', 'commits', 'files'", cla.SortOrderKey)
	}
	if !isValidOutputFormat(cla.OutputFormat) {
		return fmt.Errorf("invalid output format: %s. Should be one of: 'tabular', 'csv', 'json', 'json-lines'", cla.OutputFormat)
	}

	return nil
}
