package opendart

import (
	"github.com/go-resty/resty/v2"
)

// Client calls OpenDART APIs.
type Client struct {
	apiKey string
	resty  *resty.Client
}

// New creates an OpenDART client.
func New(config Config, opts ...Option) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	options := buildOptions(opts)
	httpClient := resty.New()
	if options.httpClient != nil {
		httpClient = resty.NewWithClient(options.httpClient)
	}

	httpClient.
		SetBaseURL(options.baseURL).
		SetHeader("User-Agent", "github.com/ev3rlit/opendart").
		SetTimeout(options.timeout)

	return &Client{
		apiKey: config.APIKey,
		resty:  httpClient,
	}, nil
}
