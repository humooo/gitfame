package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/humooo/gitfame/internal/utils/constants"
	flag "github.com/spf13/pflag"
)

type CommandLineArgs struct {
	RepositoryPath string
	CommitPointer  string
	SortOrderKey   constants.OrderKey
	UseCommitter   bool
	OutputFormat   constants.OutputFormat
	Extensions     []string
	Languages      []string
	Exclude        []string
	Restricted     []string
}

func NewCommandLineArgs() *CommandLineArgs {
	return &CommandLineArgs{}
}

func isExistingDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func parseOrderKey(value string) (constants.OrderKey, error) {
	switch value {
	case string(constants.Lines):
		return constants.Lines, nil
	case string(constants.Commits):
		return constants.Commits, nil
	case string(constants.Files):
		return constants.Files, nil
	default:
		return "", fmt.Errorf("invalid sort order key: %q. should be one of: 'lines', 'commits', 'files'", value)
	}
}

func parseOutputFormat(value string) (constants.OutputFormat, error) {
	switch value {
	case string(constants.Tabular):
		return constants.Tabular, nil
	case string(constants.CSV):
		return constants.CSV, nil
	case string(constants.SimpleJSON):
		return constants.SimpleJSON, nil
	case string(constants.JSONLines):
		return constants.JSONLines, nil
	default:
		return "", fmt.Errorf("invalid output format: %q. should be one of: 'tabular', 'csv', 'json', 'json-lines'", value)
	}
}

func normalizeStringSlice(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		normalized = append(normalized, trimmed)
	}

	return normalized
}

func validateGlobPatterns(patterns []string) error {
	for _, pattern := range patterns {
		if _, err := filepath.Match(pattern, "dummy"); err != nil {
			return fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
		}
	}

	return nil
}

func parseFlags(flagSet *flag.FlagSet, args []string) error {
	var parseOutput bytes.Buffer
	flagSet.SetOutput(&parseOutput)

	if err := flagSet.Parse(args); err != nil {
		output := strings.TrimRight(parseOutput.String(), "\n")
		if errors.Is(err, flag.ErrHelp) {
			if output != "" {
				_, _ = fmt.Fprintln(os.Stdout, output)
			}

			return err
		}

		if output != "" {
			return errors.New(output)
		}

		return err
	}

	return nil
}

func (cla *CommandLineArgs) GetCommandLineArgs(args []string) error {
	flagSet := flag.NewFlagSet("gitfame", flag.ContinueOnError)
	flagSet.SortFlags = false

	repositoryPath := flagSet.String("repository", "./", "Repository path.")
	commitPointer := flagSet.String("revision", "HEAD", "Pointer to a commit.")
	orderKey := flagSet.String("order-by", "lines", "Sort key.")
	useCommitter := flagSet.Bool("use-committer", false, "Use committer.")
	format := flagSet.String("format", "tabular", "Output format.")
	extensions := flagSet.StringSlice("extensions", []string{}, "List of Extensions to search.")
	languages := flagSet.StringSlice("languages", []string{}, "List of Languages to search.")
	exclude := flagSet.StringSlice("exclude", []string{}, "Glob-patterns to Exclude.")
	restricted := flagSet.StringSlice("restrict-to", []string{}, "List of restrictions to match.")

	if err := parseFlags(flagSet, args); err != nil {
		return err
	}

	if flagSet.NArg() > 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(flagSet.Args(), " "))
	}

	sortOrderKey, err := parseOrderKey(*orderKey)
	if err != nil {
		return err
	}

	outputFormat, err := parseOutputFormat(*format)
	if err != nil {
		return err
	}

	normalizedExtensions := normalizeStringSlice(*extensions)
	normalizedLanguages := normalizeStringSlice(*languages)
	normalizedExclude := normalizeStringSlice(*exclude)
	normalizedRestricted := normalizeStringSlice(*restricted)

	if !isExistingDir(*repositoryPath) {
		return fmt.Errorf("invalid repository path: %s", *repositoryPath)
	}

	if err := GitCommitExists(*commitPointer, *repositoryPath); err != nil {
		return fmt.Errorf("invalid revision %q: %w", *commitPointer, err)
	}

	if err := validateGlobPatterns(normalizedExclude); err != nil {
		return err
	}

	if err := validateGlobPatterns(normalizedRestricted); err != nil {
		return err
	}

	cla.RepositoryPath = *repositoryPath
	cla.CommitPointer = *commitPointer
	cla.SortOrderKey = sortOrderKey
	cla.UseCommitter = *useCommitter
	cla.OutputFormat = outputFormat
	cla.Extensions = normalizedExtensions
	cla.Languages = normalizedLanguages
	cla.Exclude = normalizedExclude
	cla.Restricted = normalizedRestricted

	return nil
}
