package opendart

import (
	"context"
	"encoding/json"
	"strings"
)

func getJSON(ctx context.Context, config apiCallerConfig, endpoint string, params map[string]string, method string, op string, out any) error {
	resp, err := config.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(config.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return requestError(method, endpoint, op, err, config.apiKey)
	}
	if err := checkHTTP(resp, method, endpoint, op); err != nil {
		return err
	}

	body := resp.Body()
	var envelope statusEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return decodeError(method, endpoint, op, "json", err)
	}
	if err := openDARTError(envelope.Status, envelope.Message, method, endpoint, op); err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return decodeError(method, endpoint, op, "json", err)
	}
	return nil
}

func getFile(ctx context.Context, config apiCallerConfig, endpoint string, params map[string]string, method string, op string) (*FileResponse, error) {
	resp, err := config.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(config.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return nil, requestError(method, endpoint, op, err, config.apiKey)
	}
	if err := checkHTTP(resp, method, endpoint, op); err != nil {
		return nil, err
	}
	if err := decodeBusinessError(resp.Body(), method, endpoint, op); err != nil {
		return nil, err
	}
	return &FileResponse{ContentType: resp.Header().Get("Content-Type"), Body: resp.Body()}, nil
}

func withAPIKey(apiKey string, params map[string]string) map[string]string {
	result := make(map[string]string, len(params)+1)
	result["crtfc_key"] = apiKey
	for key, value := range params {
		if strings.TrimSpace(value) != "" {
			result[key] = value
		}
	}
	return result
}

func endpointOp(endpoint string) string {
	return strings.TrimPrefix(endpoint, "/api/")
}
