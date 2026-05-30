package cli

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/awuzag/opendart"
	"github.com/samber/oops"
)

func newSDKClient(options *rootOptions) (*opendart.Client, error) {
	apiKey, err := options.resolvedAPIKey()
	if err != nil {
		return nil, err
	}
	return opendart.New(opendart.Config{APIKey: apiKey}, opendart.WithBaseURL(options.baseURL))
}

func requestGeneric(ctx context.Context, options *rootOptions, spec apiSpec, values map[string]string) ([]byte, string, error) {
	if body, contentType, handled, err := requestTypedFile(ctx, options, spec, values); handled || err != nil {
		return body, contentType, err
	}

	apiKey, err := options.resolvedAPIKey()
	if err != nil {
		return nil, "", err
	}

	baseURL := strings.TrimRight(options.baseURL, "/")
	endpoint := strings.TrimLeft(spec.Endpoint, "/")
	requestURL, err := url.Parse(baseURL + "/" + endpoint)
	if err != nil {
		return nil, "", oops.In("opendart_cli").
			With("endpoint", spec.Endpoint).
			Wrapf(err, "opendart cli: invalid request URL")
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
		return nil, "", oops.In("opendart_cli").
			With("endpoint", spec.Endpoint).
			Wrapf(err, "opendart cli: create request")
	}
	req.Header.Set("User-Agent", "github.com/awuzag/opendart/cmd/opendart")

	resp, err := openDARTHTTPClient().Do(req)
	if err != nil {
		return nil, "", oops.In("opendart_cli").
			With("endpoint", spec.Endpoint).
			Errorf("opendart cli: request %s: %s", spec.Endpoint, redactAPIKey(err.Error(), apiKey))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", oops.In("opendart_cli").
			With("endpoint", spec.Endpoint).
			Wrapf(err, "opendart cli: read response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", oops.In("opendart_cli").
			With("endpoint", spec.Endpoint, "status_code", resp.StatusCode, "status", resp.Status).
			Errorf("opendart cli: http error: status=%d %s", resp.StatusCode, resp.Status)
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func openDARTHTTPClient() *http.Client {
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

func redactAPIKey(message string, apiKey string) string {
	if strings.TrimSpace(apiKey) == "" {
		return message
	}
	return strings.ReplaceAll(message, apiKey, "[REDACTED]")
}

func requestTypedFile(ctx context.Context, options *rootOptions, spec apiSpec, values map[string]string) ([]byte, string, bool, error) {
	if spec.Endpoint != "/api/document.xml" && spec.Endpoint != "/api/fnlttXbrl.xml" {
		return nil, "", false, nil
	}

	client, err := newSDKClient(options)
	if err != nil {
		return nil, "", true, err
	}

	switch spec.Endpoint {
	case "/api/document.xml":
		file, err := client.DocumentRaw(ctx, opendart.DocumentParams{RceptNo: values["rcept_no"]})
		if err != nil {
			return nil, "", true, err
		}
		return file.Body, file.ContentType, true, nil
	case "/api/fnlttXbrl.xml":
		file, err := client.FnlttXbrl(ctx, opendart.FnlttXbrlParams{
			RceptNo:   values["rcept_no"],
			ReprtCode: values["reprt_code"],
		})
		if err != nil {
			return nil, "", true, err
		}
		return file.Body, file.ContentType, true, nil
	default:
		return nil, "", false, nil
	}
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
	return "", oops.In("opendart_cli").
		With("flag", "api-key", "env", envAPIKey).
		Errorf("opendart cli: --api-key or %s is required", envAPIKey)
}
