package opendart

import (
	"crypto/tls"
	"net/http"

	opendartapi "github.com/awuzag/opendart/internal/generated/opendartapi"
	"github.com/go-resty/resty/v2"
)

// Client calls OpenDART APIs.
type Client struct {
	apiKey    string
	resty     *resty.Client
	apiCaller opendartapi.Caller
}

// New creates an OpenDART client.
func New(config Config, opts ...Option) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	options := buildOptions(opts)
	httpClient := resty.NewWithClient(newOpenDARTHTTPClient())
	if options.httpClient != nil {
		httpClient = resty.NewWithClient(options.httpClient)
	}

	httpClient.
		SetBaseURL(options.baseURL).
		SetHeader("User-Agent", "github.com/awuzag/opendart").
		SetTimeout(options.timeout)

	callerConfig := apiCallerConfig{
		apiKey: config.APIKey,
		resty:  httpClient,
	}
	return &Client{
		apiKey:    callerConfig.apiKey,
		resty:     callerConfig.resty,
		apiCaller: generatedAPICaller{config: callerConfig},
	}, nil
}

func newOpenDARTHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// OpenDART 서버는 Go 기본 TLS 설정과 협상하지 못하는 경우가 있어
	// TLS 1.2와 RSA 계열 cipher suite를 명시한다.
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return &http.Client{Transport: transport}
}
