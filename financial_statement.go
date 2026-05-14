package opendart

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// FinancialStatementQuery is the query for the single-company major accounts API.
type FinancialStatementQuery struct {
	CorpCode     string
	BusinessYear string
	ReportCode   ReportCode
}

type financialStatementResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	List    []FinancialStatement `json:"list"`
}

// FinancialStatement returns major account financial statements for one company.
func (client *Client) FinancialStatement(ctx context.Context, query FinancialStatementQuery) ([]FinancialStatement, error) {
	if err := validateFinancialStatementQuery(query); err != nil {
		return nil, err
	}

	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"crtfc_key":  client.apiKey,
			"corp_code":  query.CorpCode,
			"bsns_year":  query.BusinessYear,
			"reprt_code": string(query.ReportCode),
		}).
		Get("/api/fnlttSinglAcnt.json")
	if err != nil {
		return nil, requestError("/api/fnlttSinglAcnt.json", err, client.apiKey)
	}
	if err := checkHTTP(resp); err != nil {
		return nil, err
	}

	var result financialStatementResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, &DecodeError{Op: "fnlttSinglAcnt.json", Err: err}
	}
	if err := openDARTError(result.Status, result.Message); err != nil {
		return nil, err
	}
	return result.List, nil
}

func validateFinancialStatementQuery(query FinancialStatementQuery) error {
	if strings.TrimSpace(query.CorpCode) == "" {
		return errors.New("opendart: FinancialStatementQuery.CorpCode is required")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return errors.New("opendart: FinancialStatementQuery.BusinessYear is required")
	}
	if query.ReportCode == "" {
		return errors.New("opendart: FinancialStatementQuery.ReportCode is required")
	}
	return nil
}
