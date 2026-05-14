package opendart

// ReportCode identifies an OpenDART periodic report.
type ReportCode string

const (
	// ReportFirstQuarter is the first-quarter report code.
	ReportFirstQuarter ReportCode = "11013"
	// ReportHalf is the half-year report code.
	ReportHalf ReportCode = "11012"
	// ReportThirdQuarter is the third-quarter report code.
	ReportThirdQuarter ReportCode = "11014"
	// ReportAnnual is the annual business report code.
	ReportAnnual ReportCode = "11011"
)

// FinancialStatementDivision identifies separate or consolidated statements.
type FinancialStatementDivision string

const (
	// FinancialStatementSeparate means separate financial statements.
	FinancialStatementSeparate FinancialStatementDivision = "OFS"
	// FinancialStatementConsolidated means consolidated financial statements.
	FinancialStatementConsolidated FinancialStatementDivision = "CFS"
)

// StatementDivision identifies a financial statement section.
type StatementDivision string

const (
	// StatementBalanceSheet is a balance sheet.
	StatementBalanceSheet StatementDivision = "BS"
	// StatementIncomeStatement is an income statement.
	StatementIncomeStatement StatementDivision = "IS"
	// StatementComprehensiveIncome is a comprehensive income statement.
	StatementComprehensiveIncome StatementDivision = "CIS"
	// StatementCashFlow is a cash flow statement.
	StatementCashFlow StatementDivision = "CF"
	// StatementChangesInEquity is a statement of changes in equity.
	StatementChangesInEquity StatementDivision = "SCE"
)

// CorpCode is a company code entry from the OpenDART corpCode ZIP XML.
type CorpCode struct {
	CorpCode    string `xml:"corp_code" json:"corp_code"`
	CorpName    string `xml:"corp_name" json:"corp_name"`
	CorpEngName string `xml:"corp_eng_name" json:"corp_eng_name"`
	StockCode   string `xml:"stock_code" json:"stock_code"`
	ModifyDate  string `xml:"modify_date" json:"modify_date"`
}

// FinancialStatement is an item returned by the single-company major accounts API.
type FinancialStatement struct {
	ReceiptNo                string                     `json:"rcept_no"`
	BusinessYear             string                     `json:"bsns_year"`
	StockCode                string                     `json:"stock_code"`
	ReportCode               ReportCode                 `json:"reprt_code"`
	AccountName              string                     `json:"account_nm"`
	FinancialStatementDiv    FinancialStatementDivision `json:"fs_div"`
	FinancialStatementName   string                     `json:"fs_nm"`
	StatementDiv             StatementDivision          `json:"sj_div"`
	StatementName            string                     `json:"sj_nm"`
	CurrentTermName          string                     `json:"thstrm_nm"`
	CurrentTermDate          string                     `json:"thstrm_dt"`
	CurrentTermAmount        string                     `json:"thstrm_amount"`
	CurrentTermAddAmount     string                     `json:"thstrm_add_amount"`
	PreviousTermName         string                     `json:"frmtrm_nm"`
	PreviousTermDate         string                     `json:"frmtrm_dt"`
	PreviousTermAmount       string                     `json:"frmtrm_amount"`
	PreviousTermAddAmount    string                     `json:"frmtrm_add_amount"`
	BeforePreviousTermName   string                     `json:"bfefrmtrm_nm"`
	BeforePreviousTermDate   string                     `json:"bfefrmtrm_dt"`
	BeforePreviousTermAmount string                     `json:"bfefrmtrm_amount"`
	Order                    string                     `json:"ord"`
	Currency                 string                     `json:"currency"`
}
