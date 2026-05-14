package opendart

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinancialStatement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fnlttSinglAcnt.json", r.URL.Path)
		assert.Equal(t, "test-key", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
		assert.Equal(t, "2025", r.URL.Query().Get("bsns_year"))
		assert.Equal(t, string(ReportAnnual), r.URL.Query().Get("reprt_code"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "000",
			"message": "정상",
			"list": [
				{
					"rcept_no": "20260330000001",
					"bsns_year": "2025",
					"stock_code": "005930",
					"reprt_code": "11011",
					"account_nm": "자본총계",
					"fs_div": "CFS",
					"fs_nm": "연결재무제표",
					"sj_div": "BS",
					"sj_nm": "재무상태표",
					"thstrm_nm": "제 57 기",
					"thstrm_dt": "2025.12.31 현재",
					"thstrm_amount": "100",
					"thstrm_add_amount": "",
					"frmtrm_nm": "제 56 기",
					"frmtrm_dt": "2024.12.31 현재",
					"frmtrm_amount": "90",
					"frmtrm_add_amount": "",
					"bfefrmtrm_nm": "제 55 기",
					"bfefrmtrm_dt": "2023.12.31 현재",
					"bfefrmtrm_amount": "80",
					"ord": "1",
					"currency": "KRW"
				}
			]
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	statements, err := client.FinancialStatement(context.Background(), FinancialStatementQuery{
		CorpCode:     "00126380",
		BusinessYear: "2025",
		ReportCode:   ReportAnnual,
	})
	require.NoError(t, err)
	require.Len(t, statements, 1)

	statement := statements[0]
	assert.Equal(t, "20260330000001", statement.ReceiptNo)
	assert.Equal(t, ReportAnnual, statement.ReportCode)
	assert.Equal(t, FinancialStatementConsolidated, statement.FinancialStatementDiv)
	assert.Equal(t, StatementBalanceSheet, statement.StatementDiv)
	assert.Equal(t, "KRW", statement.Currency)
}

func TestFinancialStatementReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.FinancialStatement(context.Background(), FinancialStatementQuery{
		CorpCode:     "00126380",
		BusinessYear: "2025",
		ReportCode:   ReportAnnual,
	})
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "013", apiErr.Status)
	assert.Equal(t, "조회된 데이타가 없습니다.", apiErr.Message)
}

func TestFinancialStatementReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.FinancialStatement(context.Background(), FinancialStatementQuery{
		CorpCode:     "00126380",
		BusinessYear: "2025",
		ReportCode:   ReportAnnual,
	})
	require.Error(t, err)

	var httpErr *HTTPError
	require.True(t, errors.As(err, &httpErr))
	assert.Equal(t, http.StatusServiceUnavailable, httpErr.StatusCode)
}

func TestFinancialStatementValidatesQuery(t *testing.T) {
	client, err := New(Config{APIKey: "test-key"})
	require.NoError(t, err)

	_, err = client.FinancialStatement(context.Background(), FinancialStatementQuery{})
	require.Error(t, err)
}
