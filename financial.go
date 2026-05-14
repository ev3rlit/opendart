package opendart

import (
	"context"
	"encoding/json"
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
		return nil, requiredQueryError("FinancialStatement", err)
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
		return nil, requestError("FinancialStatement", "/api/fnlttSinglAcnt.json", "fnlttSinglAcnt.json", err, client.apiKey)
	}
	if err := checkHTTP(resp, "FinancialStatement", "/api/fnlttSinglAcnt.json", "fnlttSinglAcnt.json"); err != nil {
		return nil, err
	}

	var result financialStatementResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, decodeError("FinancialStatement", "/api/fnlttSinglAcnt.json", "fnlttSinglAcnt.json", "json", err)
	}
	if err := openDARTError(result.Status, result.Message, "FinancialStatement", "/api/fnlttSinglAcnt.json", "fnlttSinglAcnt.json"); err != nil {
		return nil, err
	}
	return result.List, nil
}

func validateFinancialStatementQuery(query FinancialStatementQuery) error {
	if strings.TrimSpace(query.CorpCode) == "" {
		return requiredQueryFieldError("FinancialStatementQuery", "CorpCode")
	}
	if strings.TrimSpace(query.BusinessYear) == "" {
		return requiredQueryFieldError("FinancialStatementQuery", "BusinessYear")
	}
	if query.ReportCode == "" {
		return requiredQueryFieldError("FinancialStatementQuery", "ReportCode")
	}
	return nil
}

// FullFinancialStatement is an all-account financial statement row.
type FullFinancialStatement struct {
	ReceiptNo             string                     `json:"rcept_no"`
	ReportCode            ReportCode                 `json:"reprt_code"`
	BusinessYear          string                     `json:"bsns_year"`
	CorpCode              string                     `json:"corp_code"`
	StatementDiv          StatementDivision          `json:"sj_div"`
	StatementName         string                     `json:"sj_nm"`
	AccountID             string                     `json:"account_id"`
	AccountName           string                     `json:"account_nm"`
	AccountDetail         string                     `json:"account_detail"`
	CurrentTermName       string                     `json:"thstrm_nm"`
	CurrentTermAmount     string                     `json:"thstrm_amount"`
	CurrentTermAddAmount  string                     `json:"thstrm_add_amount"`
	PreviousTermName      string                     `json:"frmtrm_nm"`
	PreviousTermAmount    string                     `json:"frmtrm_amount"`
	PreviousTermQuarter   string                     `json:"frmtrm_q_nm"`
	PreviousTermQAmount   string                     `json:"frmtrm_q_amount"`
	PreviousTermAddAmount string                     `json:"frmtrm_add_amount"`
	FinancialStatementDiv FinancialStatementDivision `json:"fs_div"`
	FinancialStatement    string                     `json:"fs_nm"`
	Order                 string                     `json:"ord"`
	Currency              string                     `json:"currency"`
}

// XBRLTaxonomyItem is an XBRL taxonomy financial statement form row.
type XBRLTaxonomyItem struct {
	StatementDiv  StatementDivision `json:"sj_div"`
	AccountID     string            `json:"account_id"`
	AccountName   string            `json:"account_nm"`
	Language      string            `json:"lang"`
	DataType      string            `json:"data_tp"`
	Calculation   string            `json:"calculation"`
	Reference     string            `json:"ref"`
	Label         string            `json:"label"`
	ParentAccount string            `json:"parent_account_id"`
}

// FinancialIndexItem is a financial index row.
type FinancialIndexItem struct {
	BusinessYear   string     `json:"bsns_year"`
	CorpCode       string     `json:"corp_code"`
	StockCode      string     `json:"stock_code"`
	ReportCode     ReportCode `json:"reprt_code"`
	AccountName    string     `json:"account_nm"`
	CurrentTerm    string     `json:"thstrm"`
	PreviousTerm   string     `json:"frmtrm"`
	BeforePrevious string     `json:"bfefrmtrm"`
	IndexClassCode string     `json:"idx_cl_code"`
	IndexClassName string     `json:"idx_cl_nm"`
}

// MultiCompanyFinancialStatements returns major account statements for multiple companies.
func (client *Client) MultiCompanyFinancialStatements(ctx context.Context, query PeriodicReportQuery) ([]FinancialStatement, error) {
	return getPeriodic[FinancialStatement](ctx, client, "MultiCompanyFinancialStatements", "/api/fnlttMultiAcnt.json", query)
}

// FinancialStatementXBRL returns the XBRL original financial statement file.
func (client *Client) FinancialStatementXBRL(ctx context.Context, query ReceiptReportQuery) (*FileResponse, error) {
	params, err := receiptReportParams(query)
	if err != nil {
		return nil, requiredQueryError("FinancialStatementXBRL", err)
	}
	return getFile(ctx, client, "/api/fnlttXbrl.xml", params, "FinancialStatementXBRL", "fnlttXbrl.xml")
}

// FullFinancialStatement returns all accounts for one company.
func (client *Client) FullFinancialStatement(ctx context.Context, query FullFinancialStatementQuery) ([]FullFinancialStatement, error) {
	params, err := fullFinancialStatementParams(query)
	if err != nil {
		return nil, requiredQueryError("FullFinancialStatement", err)
	}
	return getList[FullFinancialStatement](ctx, client, "/api/fnlttSinglAcntAll.json", params, "FullFinancialStatement", "fnlttSinglAcntAll.json")
}

// XBRLTaxonomy returns XBRL taxonomy financial statement forms.
func (client *Client) XBRLTaxonomy(ctx context.Context, query TaxonomyQuery) ([]XBRLTaxonomyItem, error) {
	params, err := taxonomyParams(query)
	if err != nil {
		return nil, requiredQueryError("XBRLTaxonomy", err)
	}
	return getList[XBRLTaxonomyItem](ctx, client, "/api/xbrlTaxonomy.json", params, "XBRLTaxonomy", "xbrlTaxonomy.json")
}

// FinancialIndex returns financial index rows for one company.
func (client *Client) FinancialIndex(ctx context.Context, query FinancialIndexQuery) ([]FinancialIndexItem, error) {
	params, err := financialIndexParams(query)
	if err != nil {
		return nil, requiredQueryError("FinancialIndex", err)
	}
	return getList[FinancialIndexItem](ctx, client, "/api/fnlttSinglIndx.json", params, "FinancialIndex", "fnlttSinglIndx.json")
}

// CompanyFinancialIndex returns financial index rows for multiple companies.
func (client *Client) CompanyFinancialIndex(ctx context.Context, query FinancialIndexQuery) ([]FinancialIndexItem, error) {
	params, err := financialIndexParams(query)
	if err != nil {
		return nil, requiredQueryError("CompanyFinancialIndex", err)
	}
	return getList[FinancialIndexItem](ctx, client, "/api/fnlttCmpnyIndx.json", params, "CompanyFinancialIndex", "fnlttCmpnyIndx.json")
}
