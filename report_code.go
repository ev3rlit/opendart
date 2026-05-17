package opendart

const (
	// ReportCodeFirstQuarter is the OpenDART report code for first-quarter reports.
	ReportCodeFirstQuarter = "11013"
	// ReportCodeHalfYear is the OpenDART report code for half-year reports.
	ReportCodeHalfYear = "11012"
	// ReportCodeThirdQuarter is the OpenDART report code for third-quarter reports.
	ReportCodeThirdQuarter = "11014"
	// ReportCodeAnnual is the OpenDART report code for annual business reports.
	ReportCodeAnnual = "11011"
)

const (
	// FinancialStatementDivisionSeparate is the OpenDART code for separate financial statements.
	FinancialStatementDivisionSeparate = "OFS"
	// FinancialStatementDivisionConsolidated is the OpenDART code for consolidated financial statements.
	FinancialStatementDivisionConsolidated = "CFS"
)

const (
	// StatementDivisionBalanceSheet is the OpenDART code for balance sheets.
	StatementDivisionBalanceSheet = "BS"
	// StatementDivisionIncomeStatement is the OpenDART code for income statements.
	StatementDivisionIncomeStatement = "IS"
	// StatementDivisionComprehensiveIncome is the OpenDART code for comprehensive income statements.
	StatementDivisionComprehensiveIncome = "CIS"
	// StatementDivisionCashFlow is the OpenDART code for cash flow statements.
	StatementDivisionCashFlow = "CF"
	// StatementDivisionChangesInEquity is the OpenDART code for statements of changes in equity.
	StatementDivisionChangesInEquity = "SCE"
)

const (
	// CorpClassKOSPI is the OpenDART corporation class code for KOSPI-listed companies.
	CorpClassKOSPI = "Y"
	// CorpClassKOSDAQ is the OpenDART corporation class code for KOSDAQ-listed companies.
	CorpClassKOSDAQ = "K"
	// CorpClassKONEX is the OpenDART corporation class code for KONEX-listed companies.
	CorpClassKONEX = "N"
	// CorpClassOther is the OpenDART corporation class code for other companies.
	CorpClassOther = "E"
)

const (
	// DisclosureTypePeriodic is the OpenDART disclosure type code for periodic disclosures.
	DisclosureTypePeriodic = "A"
	// DisclosureTypeMaterial is the OpenDART disclosure type code for material event reports.
	DisclosureTypeMaterial = "B"
	// DisclosureTypeIssuance is the OpenDART disclosure type code for issuance disclosures.
	DisclosureTypeIssuance = "C"
	// DisclosureTypeOwnership is the OpenDART disclosure type code for ownership disclosures.
	DisclosureTypeOwnership = "D"
	// DisclosureTypeOther is the OpenDART disclosure type code for other disclosures.
	DisclosureTypeOther = "E"
	// DisclosureTypeExternalAudit is the OpenDART disclosure type code for external audit disclosures.
	DisclosureTypeExternalAudit = "F"
	// DisclosureTypeFund is the OpenDART disclosure type code for fund disclosures.
	DisclosureTypeFund = "G"
	// DisclosureTypeAssetBacked is the OpenDART disclosure type code for asset-backed securities disclosures.
	DisclosureTypeAssetBacked = "H"
	// DisclosureTypeExchange is the OpenDART disclosure type code for exchange disclosures.
	DisclosureTypeExchange = "I"
	// DisclosureTypeFairTrade is the OpenDART disclosure type code for fair trade disclosures.
	DisclosureTypeFairTrade = "J"
)

const (
	// DisclosureSortDate is the OpenDART disclosure search sort code for receipt date.
	DisclosureSortDate = "date"
	// DisclosureSortCompany is the OpenDART disclosure search sort code for company name.
	DisclosureSortCompany = "crp"
	// DisclosureSortReport is the OpenDART disclosure search sort code for report name.
	DisclosureSortReport = "rpt"
)

const (
	// SortMethodAscending is the OpenDART sort method code for ascending order.
	SortMethodAscending = "asc"
	// SortMethodDescending is the OpenDART sort method code for descending order.
	SortMethodDescending = "desc"
)
