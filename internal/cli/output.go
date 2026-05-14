package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func writeRawJSON(out io.Writer, body []byte) error {
	if !json.Valid(body) {
		return fmt.Errorf("opendart cli: response is not valid JSON")
	}
	if _, err := out.Write(body); err != nil {
		return err
	}
	_, err := out.Write([]byte("\n"))
	return err
}

func writeBinaryJSON(out io.Writer, spec apiSpec, contentType string, body []byte) error {
	return writeJSON(out, map[string]string{
		"api_id":         spec.APIID,
		"api_name":       spec.Name,
		"endpoint":       spec.Endpoint,
		"content_type":   contentType,
		"content_base64": base64.StdEncoding.EncodeToString(body),
	})
}
