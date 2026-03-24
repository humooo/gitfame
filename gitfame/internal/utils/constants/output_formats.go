package constants

type OutputFormat string

const (
	Tabular    OutputFormat = "tabular"
	CSV        OutputFormat = "csv"
	SimpleJSON OutputFormat = "json"
	JSONLines  OutputFormat = "json-lines"
)
