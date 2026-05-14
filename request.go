package opendart

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/samber/oops"
)

func getJSON(ctx context.Context, client *Client, endpoint string, params map[string]string, method string, op string, out any) error {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(client.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return requestError(method, endpoint, op, err, client.apiKey)
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

func getList[T any](ctx context.Context, client *Client, endpoint string, params map[string]string, method string, op string) ([]T, error) {
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		List    []T    `json:"list"`
	}
	if err := getJSON(ctx, client, endpoint, params, method, op, &result); err != nil {
		return nil, err
	}
	return result.List, nil
}

func getFile(ctx context.Context, client *Client, endpoint string, params map[string]string, method string, op string) (*FileResponse, error) {
	resp, err := client.resty.R().
		SetContext(ctx).
		SetQueryParams(withAPIKey(client.apiKey, params)).
		Get(endpoint)
	if err != nil {
		return nil, requestError(method, endpoint, op, err, client.apiKey)
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
		return nil, requiredQueryFieldError("CorpCodeQuery", "CorpCode")
	}
	return map[string]string{"corp_code": query.CorpCode}, nil
}

func documentParams(query DocumentQuery) (map[string]string, error) {
	if strings.TrimSpace(query.ReceiptNo) == "" {
		return nil, requiredQueryFieldError("DocumentQuery", "ReceiptNo")
	}
	return map[string]string{"rcept_no": query.ReceiptNo}, nil
}

func periodicReportParams(query PeriodicReportQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, requiredQueryFieldError("PeriodicReportQuery", "CorpCode")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, requiredQueryFieldError("PeriodicReportQuery", "BusinessYear")
	}
	if query.ReportCode == "" {
		return nil, requiredQueryFieldError("PeriodicReportQuery", "ReportCode")
	}
	return map[string]string{
		"corp_code":  query.CorpCode,
		"bsns_year":  query.BusinessYear,
		"reprt_code": string(query.ReportCode),
	}, nil
}

func receiptReportParams(query ReceiptReportQuery) (map[string]string, error) {
	if strings.TrimSpace(query.ReceiptNo) == "" {
		return nil, requiredQueryFieldError("ReceiptReportQuery", "ReceiptNo")
	}
	if query.ReportCode == "" {
		return nil, requiredQueryFieldError("ReceiptReportQuery", "ReportCode")
	}
	return map[string]string{
		"rcept_no":   query.ReceiptNo,
		"reprt_code": string(query.ReportCode),
	}, nil
}

func fullFinancialStatementParams(query FullFinancialStatementQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, requiredQueryFieldError("FullFinancialStatementQuery", "CorpCode")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, requiredQueryFieldError("FullFinancialStatementQuery", "BusinessYear")
	}
	if query.ReportCode == "" {
		return nil, requiredQueryFieldError("FullFinancialStatementQuery", "ReportCode")
	}
	if query.FinancialStatementDiv == "" {
		return nil, requiredQueryFieldError("FullFinancialStatementQuery", "FinancialStatementDiv")
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
		return nil, requiredQueryFieldError("TaxonomyQuery", "StatementDiv")
	}
	return map[string]string{"sj_div": string(query.StatementDiv)}, nil
}

func financialIndexParams(query FinancialIndexQuery) (map[string]string, error) {
	if strings.TrimSpace(query.CorpCode) == "" {
		return nil, requiredQueryFieldError("FinancialIndexQuery", "CorpCode")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return nil, requiredQueryFieldError("FinancialIndexQuery", "BusinessYear")
	}
	if query.ReportCode == "" {
		return nil, requiredQueryFieldError("FinancialIndexQuery", "ReportCode")
	}
	if strings.TrimSpace(query.IndexClassCode) == "" {
		return nil, requiredQueryFieldError("FinancialIndexQuery", "IndexClassCode")
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
		return nil, requiredQueryFieldError("MaterialReportQuery", "CorpCode")
	}
	if strings.TrimSpace(query.BeginDate) == "" {
		return nil, requiredQueryFieldError("MaterialReportQuery", "BeginDate")
	}
	if strings.TrimSpace(query.EndDate) == "" {
		return nil, requiredQueryFieldError("MaterialReportQuery", "EndDate")
	}
	return map[string]string{
		"corp_code": query.CorpCode,
		"bgn_de":    query.BeginDate,
		"end_de":    query.EndDate,
	}, nil
}

type queryFieldError struct {
	Query string
	Field string
}

func (err *queryFieldError) Error() string {
	return "opendart: " + err.Query + "." + err.Field + " is required"
}

func requiredQueryFieldError(query string, field string) error {
	return oops.In("query").
		With("query", query, "field", field).
		Wrap(&queryFieldError{Query: query, Field: field})
}

func requiredQueryError(method string, err error) error {
	if err == nil {
		return nil
	}
	builder := oops.In("query").With("method", method)
	var fieldErr *queryFieldError
	if errors.As(err, &fieldErr) {
		builder = builder.With("query", fieldErr.Query, "field", fieldErr.Field)
	}
	return builder.Wrapf(err, "opendart: %s query", method)
}

func getPeriodic[T any](ctx context.Context, client *Client, method string, endpoint string, query PeriodicReportQuery) ([]T, error) {
	params, err := periodicReportParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, method, endpointOp(endpoint))
}

func getCorpList[T any](ctx context.Context, client *Client, method string, endpoint string, query CorpCodeQuery) ([]T, error) {
	params, err := corpCodeParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, method, endpointOp(endpoint))
}

func getMaterial[T any](ctx context.Context, client *Client, method string, endpoint string, query MaterialReportQuery) ([]T, error) {
	params, err := materialReportParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, method, endpointOp(endpoint))
}

func endpointOp(endpoint string) string {
	return strings.TrimPrefix(endpoint, "/api/")
}
