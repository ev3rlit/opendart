package opendart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DisclosureListQuery is the query for the disclosure search API.
type DisclosureListQuery struct {
	CorpCode          string
	BeginDate         string
	EndDate           string
	LastReportOnly    string
	DisclosureType    string
	DisclosureSubType string
	CorpClass         string
	Sort              string
	SortMethod        string
	PageNo            int
	PageCount         int
}

// CorpCodeQuery is the query for APIs that require only corp_code.
type CorpCodeQuery struct {
	CorpCode string
}

// DocumentQuery is the query for disclosure document file APIs.
type DocumentQuery struct {
	ReceiptNo string
}

// PeriodicReportQuery is the query for periodic report information APIs.
type PeriodicReportQuery struct {
	CorpCode     string
	BusinessYear string
	ReportCode   ReportCode
}

// ReceiptReportQuery is the query for receipt-number report file APIs.
type ReceiptReportQuery struct {
	ReceiptNo  string
	ReportCode ReportCode
}

// FullFinancialStatementQuery is the query for all accounts of one company.
type FullFinancialStatementQuery struct {
	CorpCode              string
	BusinessYear          string
	ReportCode            ReportCode
	FinancialStatementDiv FinancialStatementDivision
}

// TaxonomyQuery is the query for XBRL taxonomy statement forms.
type TaxonomyQuery struct {
	StatementDiv StatementDivision
}

// FinancialIndexQuery is the query for financial index APIs.
type FinancialIndexQuery struct {
	CorpCode       string
	BusinessYear   string
	ReportCode     ReportCode
	IndexClassCode string
}

// MaterialReportQuery is the query for material and registration report APIs.
type MaterialReportQuery struct {
	CorpCode  string
	BeginDate string
	EndDate   string
}

// FileResponse contains bytes returned by OpenDART file APIs.
type FileResponse struct {
	ContentType string
	Body        []byte
}

func getJSON(ctx context.Context, client *Client, endpoint string, params map[string]string, op string, out any) error {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(client.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return requestError(endpoint, err, client.apiKey)
	}
	if err := checkHTTP(resp); err != nil {
		return err
	}

	body := resp.Body()
	var envelope statusEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return &DecodeError{Op: op, Err: err}
	}
	if err := openDARTError(envelope.Status, envelope.Message); err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return &DecodeError{Op: op, Err: err}
	}
	return nil
}

func getList[T any](ctx context.Context, client *Client, endpoint string, params map[string]string, op string) ([]T, error) {
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		List    []T    `json:"list"`
	}
	if err := getJSON(ctx, client, endpoint, params, op, &result); err != nil {
		return nil, err
	}
	return result.List, nil
}

func getFile(ctx context.Context, client *Client, endpoint string, params map[string]string) (*FileResponse, error) {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(client.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return nil, requestError(endpoint, err, client.apiKey)
	}
	if err := checkHTTP(resp); err != nil {
		return nil, err
	}
	if err := decodeBusinessError(resp.Body()); err != nil {
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

func disclosureListParams(query DisclosureListQuery) map[string]string {
	params := map[string]string{
		"corp_code":        query.CorpCode,
		"bgn_de":           query.BeginDate,
		"end_de":           query.EndDate,
		"last_reprt_at":    query.LastReportOnly,
		"pblntf_ty":        query.DisclosureType,
		"pblntf_detail_ty": query.DisclosureSubType,
		"corp_cls":         query.CorpClass,
		"sort":             query.Sort,
		"sort_mth":         query.SortMethod,
	}
	if query.PageNo > 0 {
		params["page_no"] = strconv.Itoa(query.PageNo)
	}
	if query.PageCount > 0 {
		params["page_count"] = strconv.Itoa(query.PageCount)
	}
	return params
}

func corpCodeParams(query CorpCodeQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, errors.New("opendart: CorpCodeQuery.CorpCode is required")
	}
	return map[string]string{"corp_code": query.CorpCode}, nil
}

func documentParams(query DocumentQuery) (map[string]string, error) {
	if strings.TrimSpace(query.ReceiptNo) == "" {
		return nil, errors.New("opendart: DocumentQuery.ReceiptNo is required")
	}
	return map[string]string{"rcept_no": query.ReceiptNo}, nil
}

func periodicReportParams(query PeriodicReportQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, errors.New("opendart: PeriodicReportQuery.CorpCode is required")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, errors.New("opendart: PeriodicReportQuery.BusinessYear is required")
	}
	if query.ReportCode == "" {
		return nil, errors.New("opendart: PeriodicReportQuery.ReportCode is required")
	}
	return map[string]string{
		"corp_code":  query.CorpCode,
		"bsns_year":  query.BusinessYear,
		"reprt_code": string(query.ReportCode),
	}, nil
}

func receiptReportParams(query ReceiptReportQuery) (map[string]string, error) {
	if strings.TrimSpace(query.ReceiptNo) == "" {
		return nil, errors.New("opendart: ReceiptReportQuery.ReceiptNo is required")
	}
	if query.ReportCode == "" {
		return nil, errors.New("opendart: ReceiptReportQuery.ReportCode is required")
	}
	return map[string]string{
		"rcept_no":   query.ReceiptNo,
		"reprt_code": string(query.ReportCode),
	}, nil
}

func fullFinancialStatementParams(query FullFinancialStatementQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, errors.New("opendart: FullFinancialStatementQuery.CorpCode is required")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, errors.New("opendart: FullFinancialStatementQuery.BusinessYear is required")
	}
	if query.ReportCode == "" {
		return nil, errors.New("opendart: FullFinancialStatementQuery.ReportCode is required")
	}
	if query.FinancialStatementDiv == "" {
		return nil, errors.New("opendart: FullFinancialStatementQuery.FinancialStatementDiv is required")
	}
	return map[string]string{
		"corp_code":  query.CorpCode,
		"bsns_year":  query.BusinessYear,
		"reprt_code": string(query.ReportCode),
		"fs_div":     string(query.FinancialStatementDiv),
	}, nil
}

func taxonomyParams(query TaxonomyQuery) (map[string]string, error) {
	if query.StatementDiv == "" {
		return nil, errors.New("opendart: TaxonomyQuery.StatementDiv is required")
	}
	return map[string]string{"sj_div": string(query.StatementDiv)}, nil
}

func financialIndexParams(query FinancialIndexQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, errors.New("opendart: FinancialIndexQuery.CorpCode is required")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, errors.New("opendart: FinancialIndexQuery.BusinessYear is required")
	}
	if query.ReportCode == "" {
		return nil, errors.New("opendart: FinancialIndexQuery.ReportCode is required")
	}
	if strings.TrimSpace(query.IndexClassCode) == "" {
		return nil, errors.New("opendart: FinancialIndexQuery.IndexClassCode is required")
	}
	return map[string]string{
		"corp_code":   query.CorpCode,
		"bsns_year":   query.BusinessYear,
		"reprt_code":  string(query.ReportCode),
		"idx_cl_code": query.IndexClassCode,
	}, nil
}

func materialReportParams(query MaterialReportQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, errors.New("opendart: MaterialReportQuery.CorpCode is required")
	}
	if strings.TrimSpace(query.BeginDate) == "" {
		return nil, errors.New("opendart: MaterialReportQuery.BeginDate is required")
	}
	if strings.TrimSpace(query.EndDate) == "" {
		return nil, errors.New("opendart: MaterialReportQuery.EndDate is required")
	}
	return map[string]string{
		"corp_code": query.CorpCode,
		"bgn_de":    query.BeginDate,
		"end_de":    query.EndDate,
	}, nil
}

func requiredQueryError(method string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("opendart: %s query: %w", method, err)
}
