package opendart

import "context"

// DisclosureList is the paged response from the disclosure search API.
type DisclosureList struct {
	TotalCount string       `json:"total_count"`
	TotalPage  string       `json:"total_page"`
	PageNo     string       `json:"page_no"`
	PageCount  string       `json:"page_count"`
	Items      []Disclosure `json:"list"`
}

// Disclosure is a disclosure search result item.
type Disclosure struct {
	CorpClass   string `json:"corp_cls"`
	CorpName    string `json:"corp_name"`
	CorpCode    string `json:"corp_code"`
	StockCode   string `json:"stock_code"`
	ReportName  string `json:"report_nm"`
	ReceiptNo   string `json:"rcept_no"`
	FlrName     string `json:"flr_nm"`
	ReceiptDate string `json:"rcept_dt"`
	Remark      string `json:"rm"`
}

// CompanyProfile is the response from the company overview API.
type CompanyProfile struct {
	Status             string `json:"status"`
	Message            string `json:"message"`
	CorpCode           string `json:"corp_code"`
	CorpName           string `json:"corp_name"`
	CorpNameEnglish    string `json:"corp_name_eng"`
	StockName          string `json:"stock_name"`
	StockCode          string `json:"stock_code"`
	CEOName            string `json:"ceo_nm"`
	CorpClass          string `json:"corp_cls"`
	JurisdictionOffice string `json:"jurir_no"`
	BusinessNumber     string `json:"bizr_no"`
	Address            string `json:"adres"`
	Homepage           string `json:"hm_url"`
	IRHomepage         string `json:"ir_url"`
	Phone              string `json:"phn_no"`
	Fax                string `json:"fax_no"`
	IndustryCode       string `json:"induty_code"`
	EstablishDate      string `json:"est_dt"`
	AccountMonth       string `json:"acc_mt"`
}

type disclosureListResponse struct {
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	TotalCount string       `json:"total_count"`
	TotalPage  string       `json:"total_page"`
	PageNo     string       `json:"page_no"`
	PageCount  string       `json:"page_count"`
	List       []Disclosure `json:"list"`
}

// Disclosures searches disclosure reports.
func (client *Client) Disclosures(ctx context.Context, query DisclosureListQuery) (*DisclosureList, error) {
	var result disclosureListResponse
	if err := getJSON(ctx, client, "/api/list.json", disclosureListParams(query), "Disclosures", "list.json", &result); err != nil {
		return nil, err
	}
	return &DisclosureList{
		TotalCount: result.TotalCount,
		TotalPage:  result.TotalPage,
		PageNo:     result.PageNo,
		PageCount:  result.PageCount,
		Items:      result.List,
	}, nil
}

// Company returns a company overview.
func (client *Client) Company(ctx context.Context, query CorpCodeQuery) (*CompanyProfile, error) {
	params, err := corpCodeParams(query)
	if err != nil {
		return nil, requiredQueryError("Company", err)
	}
	var result CompanyProfile
	if err := getJSON(ctx, client, "/api/company.json", params, "Company", "company.json", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Document returns the original disclosure document file.
func (client *Client) Document(ctx context.Context, query DocumentQuery) (*FileResponse, error) {
	params, err := documentParams(query)
	if err != nil {
		return nil, requiredQueryError("Document", err)
	}
	return getFile(ctx, client, "/api/document.xml", params, "Document", "document.xml")
}
