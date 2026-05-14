package opendart

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type corpCodeResult struct {
	Status  string     `xml:"status"`
	Message string     `xml:"message"`
	List    []CorpCode `xml:"list"`
}

// CorpCodes returns OpenDART company codes.
func (client *Client) CorpCodes(ctx context.Context) ([]CorpCode, error) {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParam("crtfc_key", client.apiKey).
		Get("/api/corpCode.xml")
	if err != nil {
		return nil, err
	}
	if err := checkHTTP(resp); err != nil {
		return nil, err
	}

	codes, err := decodeCorpCodeZIP(resp.Body())
	if err != nil {
		if apiErr := decodeBusinessError(resp.Body()); apiErr != nil {
			return nil, apiErr
		}
		return nil, err
	}
	return codes, nil
}

func decodeCorpCodeZIP(body []byte) ([]CorpCode, error) {
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, &DecodeError{Op: "corpCode.zip", Err: err}
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(file.Name), ".xml") {
			continue
		}

		xmlFile, err := file.Open()
		if err != nil {
			return nil, &DecodeError{Op: "corpCode.zip.open", Err: err}
		}
		defer xmlFile.Close()

		data, err := io.ReadAll(xmlFile)
		if err != nil {
			return nil, &DecodeError{Op: "corpCode.xml.read", Err: err}
		}

		var result corpCodeResult
		if err := xml.Unmarshal(data, &result); err != nil {
			return nil, &DecodeError{Op: "corpCode.xml", Err: err}
		}
		if err := openDARTError(result.Status, result.Message); err != nil {
			return nil, err
		}
		return result.List, nil
	}

	return nil, &DecodeError{Op: "corpCode.zip", Err: fmt.Errorf("zip has no XML entries")}
}
