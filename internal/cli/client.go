package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ev3rlit/opendart"
)

func newSDKClient(options *rootOptions) (*opendart.Client, error) {
	apiKey, err := options.resolvedAPIKey()
	if err != nil {
		return nil, err
	}
	return opendart.New(opendart.Config{APIKey: apiKey}, opendart.WithBaseURL(options.baseURL))
}

func requestGeneric(ctx context.Context, options *rootOptions, spec apiSpec, values map[string]string) ([]byte, string, error) {
	apiKey, err := options.resolvedAPIKey()
	if err != nil {
		return nil, "", err
	}

	baseURL := strings.TrimRight(options.baseURL, "/")
	endpoint := strings.TrimLeft(spec.Endpoint, "/")
	requestURL, err := url.Parse(baseURL + "/" + endpoint)
	if err != nil {
		return nil, "", fmt.Errorf("opendart cli: invalid request URL: %w", err)
	}

	query := requestURL.Query()
	query.Set("crtfc_key", apiKey)
	for name, value := range values {
		if strings.TrimSpace(value) != "" {
			query.Set(name, value)
		}
	}
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("opendart cli: create request: %w", err)
	}
	req.Header.Set("User-Agent", "github.com/ev3rlit/opendart/cmd/opendart")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("opendart cli: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("opendart cli: http error: status=%d %s", resp.StatusCode, resp.Status)
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func (options *rootOptions) resolvedAPIKey() (string, error) {
	if strings.TrimSpace(options.apiKey) != "" {
		return strings.TrimSpace(options.apiKey), nil
	}
	if options.getenv != nil {
		if apiKey := strings.TrimSpace(options.getenv(envAPIKey)); apiKey != "" {
			return apiKey, nil
		}
	}
	return "", fmt.Errorf("opendart cli: --api-key or %s is required", envAPIKey)
}
