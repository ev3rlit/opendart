package opendart

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// HTTPError reports a non-2xx HTTP response.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("opendart: http error: status=%d %s", err.StatusCode, err.Status)
}

// DecodeError reports a response decoding failure.
type DecodeError struct {
	Op  string
	Err error
}

func (err *DecodeError) Error() string {
	return fmt.Sprintf("opendart: decode %s: %v", err.Op, err.Err)
}

func (err *DecodeError) Unwrap() error {
	return err.Err
}

// APIError reports an OpenDART business error response.
type APIError struct {
	Status  string
	Message string
}

func (err *APIError) Error() string {
	return fmt.Sprintf("opendart: api error: status=%s message=%s", err.Status, err.Message)
}

func checkHTTP(resp *resty.Response) error {
	if resp.IsSuccess() {
		return nil
	}
	return &HTTPError{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Body:       string(resp.Body()),
	}
}

type statusEnvelope struct {
	Status  string `json:"status" xml:"status"`
	Message string `json:"message" xml:"message"`
}

func openDARTError(status, message string) error {
	if status == "" || status == "000" {
		return nil
	}
	return &APIError{Status: status, Message: message}
}

func decodeBusinessError(body []byte) error {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return nil
	}

	var jsonEnvelope statusEnvelope
	if strings.HasPrefix(trimmed, "{") {
		if err := json.Unmarshal(body, &jsonEnvelope); err == nil {
			return openDARTError(jsonEnvelope.Status, jsonEnvelope.Message)
		}
		return nil
	}

	var xmlEnvelope statusEnvelope
	if err := xml.Unmarshal(body, &xmlEnvelope); err == nil {
		return openDARTError(xmlEnvelope.Status, xmlEnvelope.Message)
	}
	return nil
}
