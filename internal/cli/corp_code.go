package cli

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"strings"

	"github.com/samber/oops"
)

type corpCodeResult struct {
	Status  string     `xml:"status"`
	Message string     `xml:"message"`
	List    []corpCode `xml:"list"`
}

type corpCode struct {
	CorpCode    string `xml:"corp_code" json:"corp_code"`
	CorpName    string `xml:"corp_name" json:"corp_name"`
	CorpEngName string `xml:"corp_eng_name" json:"corp_eng_name"`
	StockCode   string `xml:"stock_code" json:"stock_code"`
	ModifyDate  string `xml:"modify_date" json:"modify_date"`
}

func decodeCorpCodeZIP(body []byte) ([]corpCode, error) {
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, oops.In("opendart_cli").
			With("endpoint", "/api/corpCode.xml", "format", "zip").
			Wrapf(err, "opendart cli: decode corpCode.zip")
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(file.Name), ".xml") {
			continue
		}

		xmlFile, err := file.Open()
		if err != nil {
			return nil, oops.In("opendart_cli").
				With("endpoint", "/api/corpCode.xml", "file", file.Name).
				Wrapf(err, "opendart cli: open corpCode XML")
		}
		defer xmlFile.Close()

		data, err := io.ReadAll(xmlFile)
		if err != nil {
			return nil, oops.In("opendart_cli").
				With("endpoint", "/api/corpCode.xml", "file", file.Name).
				Wrapf(err, "opendart cli: read corpCode XML")
		}

		var result corpCodeResult
		if err := xml.Unmarshal(data, &result); err != nil {
			return nil, oops.In("opendart_cli").
				With("endpoint", "/api/corpCode.xml", "format", "xml").
				Wrapf(err, "opendart cli: decode corpCode XML")
		}
		if result.Status != "" && result.Status != "000" {
			return nil, oops.In("opendart_cli").
				With("endpoint", "/api/corpCode.xml", "status", result.Status, "message", result.Message).
				Errorf("opendart cli: api error: status=%s message=%s", result.Status, result.Message)
		}
		return result.List, nil
	}

	return nil, oops.In("opendart_cli").
		With("endpoint", "/api/corpCode.xml", "format", "zip").
		New("opendart cli: corpCode ZIP has no XML entries")
}
