package cli

import "io"

const (
	defaultBaseURL = "https://opendart.fss.or.kr"
	envAPIKey      = "OPENDART_API_KEY"
	outputJSON     = "json"
	outputRaw      = "raw"
	outputTable    = "table"
	outputCSV      = "csv"
	viewSummary    = "summary"
	viewDetail     = "detail"
	viewSource     = "source"
)

type getenvFunc func(string) string

type rootOptions struct {
	apiKey  string
	baseURL string
	output  string

	out    io.Writer
	errOut io.Writer
	getenv getenvFunc
}

type apiSpec struct {
	Group       string
	Command     string
	Verb        string
	Resource    string
	APIID       string
	OperationID string
	Name        string
	Endpoint    string
	Params      []paramSpec
}

type paramSpec struct {
	Name        string
	Description string
	Required    bool
}
