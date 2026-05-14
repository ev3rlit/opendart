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

func TestDisclosureTypedMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("crtfc_key"))
		switch r.URL.Path {
		case "/api/list.json":
			assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
			assert.Equal(t, "20250101", r.URL.Query().Get("bgn_de"))
			assert.Equal(t, "2", r.URL.Query().Get("page_no"))
			_, _ = w.Write([]byte(`{
				"status":"000",
				"message":"정상",
				"total_count":"1",
				"total_page":"1",
				"page_no":"2",
				"page_count":"10",
				"list":[{"corp_code":"00126380","corp_name":"삼성전자","rcept_no":"20260330000001"}]
			}`))
		case "/api/company.json":
			assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
			_, _ = w.Write([]byte(`{"status":"000","message":"정상","corp_code":"00126380","corp_name":"삼성전자"}`))
		case "/api/document.xml":
			assert.Equal(t, "20260330000001", r.URL.Query().Get("rcept_no"))
			w.Header().Set("Content-Type", "application/zip")
			_, _ = w.Write([]byte("zip-bytes"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	disclosures, err := client.Disclosures(context.Background(), DisclosureListQuery{
		CorpCode:  "00126380",
		BeginDate: "20250101",
		PageNo:    2,
		PageCount: 10,
	})
	require.NoError(t, err)
	require.Len(t, disclosures.Items, 1)
	assert.Equal(t, "1", disclosures.TotalCount)
	assert.Equal(t, "20260330000001", disclosures.Items[0].ReceiptNo)

	company, err := client.Company(context.Background(), CorpCodeQuery{CorpCode: "00126380"})
	require.NoError(t, err)
	assert.Equal(t, "삼성전자", company.CorpName)

	document, err := client.Document(context.Background(), DocumentQuery{ReceiptNo: "20260330000001"})
	require.NoError(t, err)
	assert.Equal(t, "application/zip", document.ContentType)
	assert.Equal(t, []byte("zip-bytes"), document.Body)
}

func TestTypedMethodsMapRepresentativeEndpoints(t *testing.T) {
	expectedPaths := map[string]bool{
		"/api/irdsSttus.json":         false,
		"/api/fnlttSinglAcntAll.json": false,
		"/api/fnlttXbrl.xml":          false,
		"/api/xbrlTaxonomy.json":      false,
		"/api/fnlttSinglIndx.json":    false,
		"/api/majorstock.json":        false,
		"/api/cmpMgDecsn.json":        false,
		"/api/estkRs.json":            false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := expectedPaths[r.URL.Path]; !ok {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		expectedPaths[r.URL.Path] = true
		assert.Equal(t, "test-key", r.URL.Query().Get("crtfc_key"))

		switch r.URL.Path {
		case "/api/fnlttXbrl.xml":
			assert.Equal(t, "20260330000001", r.URL.Query().Get("rcept_no"))
			assert.Equal(t, string(ReportAnnual), r.URL.Query().Get("reprt_code"))
			_, _ = w.Write([]byte("xbrl-bytes"))
		case "/api/estkRs.json":
			assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
			assert.Equal(t, "20250101", r.URL.Query().Get("bgn_de"))
			_, _ = w.Write([]byte(`{
				"status":"000",
				"message":"정상",
				"title":"지분증권",
				"group":[{"title":"모집","list":[{"rcept_no":"20260330000001"}]}]
			}`))
		default:
			_, _ = w.Write([]byte(`{"status":"000","message":"정상","list":[{"rcept_no":"20260330000001","account_nm":"자산총계","custom_field":"kept"}]}`))
		}
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)
	ctx := context.Background()

	periodic := PeriodicReportQuery{CorpCode: "00126380", BusinessYear: "2025", ReportCode: ReportAnnual}
	rows, err := client.CapitalIncreaseDecreaseStatus(ctx, periodic)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "kept", rows[0]["custom_field"])

	full, err := client.FullFinancialStatement(ctx, FullFinancialStatementQuery{
		CorpCode:              "00126380",
		BusinessYear:          "2025",
		ReportCode:            ReportAnnual,
		FinancialStatementDiv: FinancialStatementConsolidated,
	})
	require.NoError(t, err)
	require.Len(t, full, 1)
	assert.Equal(t, "자산총계", full[0].AccountName)

	xbrl, err := client.FinancialStatementXBRL(ctx, ReceiptReportQuery{
		ReceiptNo:  "20260330000001",
		ReportCode: ReportAnnual,
	})
	require.NoError(t, err)
	assert.Equal(t, []byte("xbrl-bytes"), xbrl.Body)

	taxonomy, err := client.XBRLTaxonomy(ctx, TaxonomyQuery{StatementDiv: StatementBalanceSheet})
	require.NoError(t, err)
	require.Len(t, taxonomy, 1)
	assert.Equal(t, "자산총계", taxonomy[0].AccountName)

	index, err := client.FinancialIndex(ctx, FinancialIndexQuery{
		CorpCode:       "00126380",
		BusinessYear:   "2025",
		ReportCode:     ReportAnnual,
		IndexClassCode: "M210000",
	})
	require.NoError(t, err)
	require.Len(t, index, 1)
	assert.Equal(t, "자산총계", index[0].AccountName)

	ownership, err := client.MajorStock(ctx, CorpCodeQuery{CorpCode: "00126380"})
	require.NoError(t, err)
	require.Len(t, ownership, 1)
	assert.Equal(t, "20260330000001", ownership[0]["rcept_no"])

	material, err := client.CompanyMergerDecision(ctx, MaterialReportQuery{
		CorpCode:  "00126380",
		BeginDate: "20250101",
		EndDate:   "20251231",
	})
	require.NoError(t, err)
	require.Len(t, material, 1)
	assert.Equal(t, "kept", material[0]["custom_field"])

	registration, err := client.EquitySecuritiesRegistration(ctx, MaterialReportQuery{
		CorpCode:  "00126380",
		BeginDate: "20250101",
		EndDate:   "20251231",
	})
	require.NoError(t, err)
	require.Len(t, registration.Group, 1)
	assert.Equal(t, "지분증권", registration.Title)

	for path, called := range expectedPaths {
		assert.True(t, called, "path was not called: %s", path)
	}
}

func TestTypedMethodsValidateRequiredQueries(t *testing.T) {
	client, err := New(Config{APIKey: "test-key"})
	require.NoError(t, err)

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "corp", run: func() error {
			_, err := client.Company(context.Background(), CorpCodeQuery{})
			return err
		}},
		{name: "document", run: func() error {
			_, err := client.Document(context.Background(), DocumentQuery{})
			return err
		}},
		{name: "periodic", run: func() error {
			_, err := client.CapitalIncreaseDecreaseStatus(context.Background(), PeriodicReportQuery{})
			return err
		}},
		{name: "receipt", run: func() error {
			_, err := client.FinancialStatementXBRL(context.Background(), ReceiptReportQuery{})
			return err
		}},
		{name: "full financial", run: func() error {
			_, err := client.FullFinancialStatement(context.Background(), FullFinancialStatementQuery{})
			return err
		}},
		{name: "taxonomy", run: func() error {
			_, err := client.XBRLTaxonomy(context.Background(), TaxonomyQuery{})
			return err
		}},
		{name: "financial index", run: func() error {
			_, err := client.FinancialIndex(context.Background(), FinancialIndexQuery{})
			return err
		}},
		{name: "material", run: func() error {
			_, err := client.CompanyMergerDecision(context.Background(), MaterialReportQuery{})
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "required")
		})
	}
}

func TestTypedMethodsReturnAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"013","message":"조회된 데이타가 없습니다."}`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.CompanyMergerDecision(context.Background(), MaterialReportQuery{
		CorpCode:  "00126380",
		BeginDate: "20250101",
		EndDate:   "20251231",
	})
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "013", apiErr.Status)
}

func TestRequestErrorRedactsAPIKey(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"000","message":"정상"}`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "secret-test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.Company(context.Background(), CorpCodeQuery{CorpCode: "00126380"})
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "secret-test-key")
	assert.Contains(t, err.Error(), "[REDACTED]")
}
