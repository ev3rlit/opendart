package opendart

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeFnlttSinglAcntAllHandlesRepresentativeCompanyShapes(t *testing.T) {
	tests := map[string]struct {
		source FinancialMetricSource
		rows   []FnlttSinglAcntAllItem
	}{
		"samsung electronics has IS and CIS": {
			source: FinancialMetricSource{
				CorpCode:     "00126380",
				BusinessYear: "2024",
				ReportCode:   ReportCodeAnnual,
				FsDiv:        FinancialStatementDivisionConsolidated,
			},
			rows: []FnlttSinglAcntAllItem{
				financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "455,905,980,000,000"),
				financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "ifrs-full_Liabilities", "부채총계", "120,188,333,000,000"),
				financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "ifrs-full_Equity", "자본총계", "335,717,647,000,000"),
				financialMetricFixtureRow("00126380", StatementDivisionIncomeStatement, "ifrs-full_Revenue", "매출액", "300,870,903,000,000"),
				financialMetricFixtureRow("00126380", StatementDivisionIncomeStatement, "dart_OperatingIncomeLoss", "영업이익", "32,726,016,000,000"),
				financialMetricFixtureRow("00126380", StatementDivisionIncomeStatement, "ifrs-full_ProfitLoss", "당기순이익", "34,451,430,000,000"),
			},
		},
		"samsung electro mechanics reports income metrics under CIS": {
			source: FinancialMetricSource{
				CorpCode:     "00126371",
				BusinessYear: "2024",
				ReportCode:   ReportCodeAnnual,
				FsDiv:        FinancialStatementDivisionConsolidated,
			},
			rows: []FnlttSinglAcntAllItem{
				financialMetricFixtureRow("00126371", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "12,345,000,000,000"),
				financialMetricFixtureRow("00126371", StatementDivisionComprehensiveIncome, "ifrs-full_Revenue", "매출액", "10,294,000,000,000"),
				financialMetricFixtureRow("00126371", StatementDivisionComprehensiveIncome, "dart_OperatingIncomeLoss", "영업이익", "735,000,000,000"),
				financialMetricFixtureRow("00126371", StatementDivisionComprehensiveIncome, "ifrs-full_ProfitLoss", "당기순이익", "623,000,000,000"),
			},
		},
		"hyundai keeps industry specific rows unmapped": {
			source: FinancialMetricSource{
				CorpCode:     "00164742",
				BusinessYear: "2024",
				ReportCode:   ReportCodeAnnual,
				FsDiv:        FinancialStatementDivisionConsolidated,
			},
			rows: []FnlttSinglAcntAllItem{
				financialMetricFixtureRow("00164742", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "309,000,000,000,000"),
				financialMetricFixtureRow("00164742", StatementDivisionIncomeStatement, "ifrs-full_Revenue", "매출액", "175,000,000,000,000"),
				financialMetricFixtureRow("00164742", StatementDivisionIncomeStatement, "dart_OperatingIncomeLoss", "영업이익", "14,000,000,000,000"),
				financialMetricFixtureRow("00164742", StatementDivisionBalanceSheet, "-표준계정코드 미사용-", "금융업채권", "25,000,000,000,000"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			metrics, err := NormalizeFnlttSinglAcntAll(tt.rows, tt.source)
			require.NoError(t, err)
			assert.Empty(t, metrics.Issues)

			revenue, ok := metrics.Find(FinancialMetricRevenue)
			require.True(t, ok)
			assert.Equal(t, tt.source.CorpCode, revenue.CorpCode)
			assert.Equal(t, tt.source.FsDiv, revenue.FsDiv)
			assert.Contains(t, []string{StatementDivisionIncomeStatement, StatementDivisionComprehensiveIncome}, revenue.StatementDiv)
			assert.Equal(t, FinancialMetricMatchAccountIDExact, revenue.MatchMethod)
			assert.Equal(t, FinancialMetricConfidenceHigh, revenue.Confidence)

			assets, ok := metrics.Find(FinancialMetricAssets)
			require.True(t, ok)
			assert.Equal(t, StatementDivisionBalanceSheet, assets.StatementDiv)
			assert.Equal(t, "account_id", assets.MatchedField)

			if tt.source.CorpCode == "00164742" {
				require.Len(t, metrics.UnmappedRows, 1)
				assert.Equal(t, "금융업채권", metrics.UnmappedRows[0].SourceRow.AccountNm)
				assert.Equal(t, "-표준계정코드 미사용-", metrics.UnmappedRows[0].SourceRow.AccountId)
			}
		})
	}
}

func TestNormalizeFnlttSinglAcntAllUsesAliasAfterAccountID(t *testing.T) {
	rows := []FnlttSinglAcntAllItem{
		financialMetricFixtureRow("00126380", StatementDivisionComprehensiveIncome, "-표준계정코드 미사용-", "영업수익", "10,000"),
		financialMetricFixtureRow("00126380", StatementDivisionComprehensiveIncome, "-표준계정코드 미사용-", "연결당기순이익", "(1,234)"),
	}

	metrics, err := NormalizeFnlttSinglAcntAll(rows, FinancialMetricSource{
		CorpCode:     "00126380",
		BusinessYear: "2024",
		ReportCode:   ReportCodeAnnual,
		FsDiv:        FinancialStatementDivisionConsolidated,
	})
	require.NoError(t, err)
	require.Empty(t, metrics.Issues)
	require.Len(t, metrics.Metrics, 2)

	revenue, ok := metrics.Find(FinancialMetricRevenue)
	require.True(t, ok)
	assert.Equal(t, int64(10000), revenue.Amount)
	assert.Equal(t, FinancialMetricMatchAccountNameAlias, revenue.MatchMethod)
	assert.Equal(t, FinancialMetricConfidenceMedium, revenue.Confidence)

	netIncome, ok := metrics.Find(FinancialMetricNetIncome)
	require.True(t, ok)
	assert.Equal(t, int64(-1234), netIncome.Amount)
	assert.Equal(t, "account_nm", netIncome.MatchedField)
}

func TestNormalizeFnlttSinglAcntAllSupportsOverrideRules(t *testing.T) {
	rows := []FnlttSinglAcntAllItem{
		financialMetricFixtureRow("00164742", StatementDivisionBalanceSheet, "-표준계정코드 미사용-", "금융업채권", "25,000"),
	}

	metrics, err := NormalizeFnlttSinglAcntAll(rows, FinancialMetricSource{
		CorpCode:     "00164742",
		BusinessYear: "2024",
		ReportCode:   ReportCodeAnnual,
		FsDiv:        FinancialStatementDivisionConsolidated,
	}, WithFinancialMetricOverrideRules(FinancialMetricRule{
		ID:                 "hyundai_finance_receivables",
		MetricCode:         "finance_receivables",
		Label:              "Finance receivables",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountNameAliases: []string{"금융업채권"},
	}))
	require.NoError(t, err)
	require.Empty(t, metrics.UnmappedRows)
	require.Empty(t, metrics.Issues)
	require.Len(t, metrics.Metrics, 1)
	assert.Equal(t, FinancialMetricCode("finance_receivables"), metrics.Metrics[0].MetricCode)
	assert.Equal(t, FinancialMetricMatchOverride, metrics.Metrics[0].MatchMethod)
	assert.Equal(t, FinancialMetricConfidenceManual, metrics.Metrics[0].Confidence)
	assert.Equal(t, "hyundai_finance_receivables", metrics.Metrics[0].RuleID)
}

func TestNormalizeFnlttSinglAcntAllReportsAmountIssuesWithoutFailing(t *testing.T) {
	rows := []FnlttSinglAcntAllItem{
		financialMetricFixtureRow("00126380", StatementDivisionIncomeStatement, "ifrs-full_Revenue", "매출액", "not-a-number"),
		financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "-표준계정코드 미사용-", "회사고유계정", "123"),
	}

	metrics, err := NormalizeFnlttSinglAcntAll(rows, FinancialMetricSource{})
	require.NoError(t, err)
	require.Empty(t, metrics.Metrics)
	require.Len(t, metrics.Issues, 1)
	require.Len(t, metrics.UnmappedRows, 1)
	assert.Contains(t, metrics.Issues[0].Reason, "invalid thstrm_amount")
	assert.Equal(t, "회사고유계정", metrics.UnmappedRows[0].SourceRow.AccountNm)
}

func TestNormalizeFnlttSinglAcntAllResponseUsesParamsForSource(t *testing.T) {
	response := &FnlttSinglAcntAllResponse{
		Status:  "000",
		Message: "정상",
		List: []FnlttSinglAcntAllItem{
			financialMetricFixtureRow("", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "455,000"),
		},
	}
	params := FnlttSinglAcntAllParams{
		CorpCode:  "00126380",
		BsnsYear:  "2024",
		ReprtCode: ReportCodeAnnual,
		FsDiv:     FinancialStatementDivisionConsolidated,
	}

	metrics, err := NormalizeFnlttSinglAcntAllResponse(response, params)
	require.NoError(t, err)
	assets, ok := metrics.Find(FinancialMetricAssets)
	require.True(t, ok)
	assert.Equal(t, "00126380", assets.CorpCode)
	assert.Equal(t, "2024", assets.BusinessYear)
	assert.Equal(t, FinancialStatementDivisionConsolidated, assets.FsDiv)
}

func TestAnalyzeFnlttSinglAcntAllRowsUsesNullablePeriodAndAmountFields(t *testing.T) {
	row := financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "455,000")
	row.ThstrmAddAmount = ""
	row.FrmtrmNm = "제 55 기"
	row.FrmtrmAmount = ""
	row.BfefrmtrmNm = ""
	row.BfefrmtrmAmount = "123,000"

	analysis := AnalyzeFnlttSinglAcntAllRows([]FnlttSinglAcntAllItem{row})
	require.Empty(t, analysis.Issues)
	require.Len(t, analysis.Rows, 1)

	analysisRow := analysis.Rows[0]
	assert.Equal(t, "00126380", analysisRow.CorpCode)
	assert.Equal(t, "ifrs-full_Assets", analysisRow.AccountID)
	require.NotNil(t, analysisRow.CurrentAmount)
	assert.Equal(t, "455,000", analysisRow.CurrentAmount.Raw)
	assert.Equal(t, int64(455000), analysisRow.CurrentAmount.Value)
	assert.Nil(t, analysisRow.CurrentCumulativeAmount)
	require.NotNil(t, analysisRow.PreviousName)
	assert.Equal(t, "제 55 기", *analysisRow.PreviousName)
	assert.Nil(t, analysisRow.PreviousAmount)
	assert.Nil(t, analysisRow.BeforePreviousName)
	require.NotNil(t, analysisRow.BeforePreviousAmount)
	assert.Equal(t, int64(123000), analysisRow.BeforePreviousAmount.Value)
	assert.Equal(t, row, analysisRow.SourceRow)
}

func TestAnalyzeFnlttSinglAcntAllRowsReportsMissingKeysAndInvalidAmounts(t *testing.T) {
	row := financialMetricFixtureRow("00126380", StatementDivisionBalanceSheet, "ifrs-full_Assets", "자산총계", "455,000")
	row.AccountDetail = ""
	row.Ord = ""
	row.FrmtrmAmount = "not-a-number"

	analysis := AnalyzeFnlttSinglAcntAllRows([]FnlttSinglAcntAllItem{row})
	require.Len(t, analysis.Issues, 3)
	assert.Equal(t, "account_detail", analysis.Issues[0].Field)
	assert.Equal(t, "ord", analysis.Issues[1].Field)
	assert.Equal(t, "frmtrm_amount", analysis.Issues[2].Field)
	assert.Contains(t, analysis.Issues[2].Reason, "invalid amount")
	assert.Nil(t, analysis.Rows[0].PreviousAmount)
}

func TestFnlttSinglAcntAllMetricsFetchesRawThenNormalizes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fnlttSinglAcntAll.json", r.URL.Path)
		assert.Equal(t, "test-key", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, FinancialStatementDivisionConsolidated, r.URL.Query().Get("fs_div"))
		_, _ = w.Write([]byte(`{
			"status":"000",
			"message":"정상",
			"list":[
				{"corp_code":"00126380","bsns_year":"2024","reprt_code":"11011","sj_div":"BS","account_id":"ifrs-full_Assets","account_nm":"자산총계","thstrm_amount":"455,000","currency":"KRW"}
			]
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	metrics, err := client.FnlttSinglAcntAllMetrics(context.Background(), FnlttSinglAcntAllParams{
		CorpCode:  "00126380",
		BsnsYear:  "2024",
		ReprtCode: ReportCodeAnnual,
		FsDiv:     FinancialStatementDivisionConsolidated,
	})
	require.NoError(t, err)
	assets, ok := metrics.Find(FinancialMetricAssets)
	require.True(t, ok)
	assert.Equal(t, int64(455000), assets.Amount)
	assert.Equal(t, "KRW", assets.Currency)
	assert.Equal(t, FinancialMetricAmountCurrent, assets.AmountSource)
}

func financialMetricFixtureRow(corpCode string, sjDiv string, accountID string, accountName string, amount string) FnlttSinglAcntAllItem {
	return FnlttSinglAcntAllItem{
		RceptNo:       "20250331000000",
		ReprtCode:     ReportCodeAnnual,
		BsnsYear:      "2024",
		CorpCode:      corpCode,
		SjDiv:         sjDiv,
		SjNm:          sjDiv,
		AccountId:     accountID,
		AccountNm:     accountName,
		AccountDetail: accountName,
		ThstrmNm:      "제 56 기",
		ThstrmAmount:  amount,
		Ord:           "1",
		Currency:      "KRW",
	}
}
