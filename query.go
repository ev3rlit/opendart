package opendart

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
