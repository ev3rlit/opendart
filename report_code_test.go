package opendart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportCodesMatchOpenDARTValues(t *testing.T) {
	tests := map[string]string{
		"ReportCodeFirstQuarter":                 ReportCodeFirstQuarter,
		"ReportCodeHalfYear":                     ReportCodeHalfYear,
		"ReportCodeThirdQuarter":                 ReportCodeThirdQuarter,
		"ReportCodeAnnual":                       ReportCodeAnnual,
		"FinancialStatementDivisionSeparate":     FinancialStatementDivisionSeparate,
		"FinancialStatementDivisionConsolidated": FinancialStatementDivisionConsolidated,
		"StatementDivisionBalanceSheet":          StatementDivisionBalanceSheet,
		"StatementDivisionIncomeStatement":       StatementDivisionIncomeStatement,
		"StatementDivisionComprehensiveIncome":   StatementDivisionComprehensiveIncome,
		"StatementDivisionCashFlow":              StatementDivisionCashFlow,
		"StatementDivisionChangesInEquity":       StatementDivisionChangesInEquity,
		"CorpClassKOSPI":                         CorpClassKOSPI,
		"CorpClassKOSDAQ":                        CorpClassKOSDAQ,
		"CorpClassKONEX":                         CorpClassKONEX,
		"CorpClassOther":                         CorpClassOther,
		"DisclosureTypePeriodic":                 DisclosureTypePeriodic,
		"DisclosureTypeMaterial":                 DisclosureTypeMaterial,
		"DisclosureTypeIssuance":                 DisclosureTypeIssuance,
		"DisclosureTypeOwnership":                DisclosureTypeOwnership,
		"DisclosureTypeOther":                    DisclosureTypeOther,
		"DisclosureTypeExternalAudit":            DisclosureTypeExternalAudit,
		"DisclosureTypeFund":                     DisclosureTypeFund,
		"DisclosureTypeAssetBacked":              DisclosureTypeAssetBacked,
		"DisclosureTypeExchange":                 DisclosureTypeExchange,
		"DisclosureTypeFairTrade":                DisclosureTypeFairTrade,
		"DisclosureSortDate":                     DisclosureSortDate,
		"DisclosureSortCompany":                  DisclosureSortCompany,
		"DisclosureSortReport":                   DisclosureSortReport,
		"SortMethodAscending":                    SortMethodAscending,
		"SortMethodDescending":                   SortMethodDescending,
	}

	expected := map[string]string{
		"ReportCodeFirstQuarter":                 "11013",
		"ReportCodeHalfYear":                     "11012",
		"ReportCodeThirdQuarter":                 "11014",
		"ReportCodeAnnual":                       "11011",
		"FinancialStatementDivisionSeparate":     "OFS",
		"FinancialStatementDivisionConsolidated": "CFS",
		"StatementDivisionBalanceSheet":          "BS",
		"StatementDivisionIncomeStatement":       "IS",
		"StatementDivisionComprehensiveIncome":   "CIS",
		"StatementDivisionCashFlow":              "CF",
		"StatementDivisionChangesInEquity":       "SCE",
		"CorpClassKOSPI":                         "Y",
		"CorpClassKOSDAQ":                        "K",
		"CorpClassKONEX":                         "N",
		"CorpClassOther":                         "E",
		"DisclosureTypePeriodic":                 "A",
		"DisclosureTypeMaterial":                 "B",
		"DisclosureTypeIssuance":                 "C",
		"DisclosureTypeOwnership":                "D",
		"DisclosureTypeOther":                    "E",
		"DisclosureTypeExternalAudit":            "F",
		"DisclosureTypeFund":                     "G",
		"DisclosureTypeAssetBacked":              "H",
		"DisclosureTypeExchange":                 "I",
		"DisclosureTypeFairTrade":                "J",
		"DisclosureSortDate":                     "date",
		"DisclosureSortCompany":                  "crp",
		"DisclosureSortReport":                   "rpt",
		"SortMethodAscending":                    "asc",
		"SortMethodDescending":                   "desc",
	}

	assert.Equal(t, expected, tests)
}
