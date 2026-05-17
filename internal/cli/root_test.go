package cli

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommandBuildsVerbFirstCommands(t *testing.T) {
	cmd := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(nil))

	assert.NotNil(t, findCommand(t, cmd, "search", "disclosures"))
	assert.NotNil(t, findCommand(t, cmd, "get", "company-profile"))
	assert.NotNil(t, findCommand(t, cmd, "download", "document"))
	assert.NotNil(t, findCommand(t, cmd, "list", "corp-codes"))
	assert.NotNil(t, findCommand(t, cmd, "get", "financial-statement"))
	assert.NotNil(t, findCommand(t, cmd, "get", "merger-decision"))
	assert.NotNil(t, findCommand(t, cmd, "get", "registration-equity"))
	assert.NotNil(t, findCommand(t, cmd, "get", "quarter-performance"))
	assert.NotNil(t, findCommand(t, cmd, "get", "annual-performance"))
	assert.NotNil(t, findCommand(t, cmd, "get", "financial-position"))
	assert.NotNil(t, findCommand(t, cmd, "get", "cash-flow"))
	assert.NotNil(t, findCommand(t, cmd, "get", "financial-metric"))

	for _, name := range []string{"summarize", "compare", "inspect"} {
		legacy := findCommand(t, cmd, name)
		require.NotNil(t, legacy)
		assert.True(t, legacy.Hidden)
		assert.NotEmpty(t, legacy.Deprecated)
		assert.NotContains(t, publicCommandNames(cmd), name)
	}

	legacy := findCommand(t, cmd, "financial", "single-account")
	require.NotNil(t, legacy)
	assert.True(t, legacy.Hidden)
}

func TestListCorpCodesCommandUsesEnvAPIKeyAndSDKClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/corpCode.xml", r.URL.Path)
		assert.Equal(t, "env-key", r.URL.Query().Get("crtfc_key"))
		_, _ = w.Write(corpCodeFixture(t))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{"--base-url", server.URL, "list", "corp-codes"})

	require.NoError(t, cmd.Execute())

	var codes []map[string]string
	require.NoError(t, json.Unmarshal(out.Bytes(), &codes))
	require.Len(t, codes, 1)
	assert.Equal(t, "00126380", codes[0]["corp_code"])
	assert.Equal(t, "005930", codes[0]["stock_code"])
}

func TestCommandRequiresAPIKeyWithoutLeakingSecret(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := NewRootCommand(&out, &errOut, testGetenv(nil))
	cmd.SetArgs([]string{"list", "corp-codes"})

	err := cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), envAPIKey)
	assert.NotContains(t, err.Error(), "secret")
	assert.Empty(t, out.String())
	assertOopsContext(t, err, map[string]any{
		"flag": "api-key",
		"env":  envAPIKey,
	})
}

func TestGenericCommandValidatesRequiredFlags(t *testing.T) {
	cmd := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{"get", "financial-statement"})

	err := cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--corp-code is required")
	assertOopsContext(t, err, map[string]any{
		"command": "get financial-statement",
		"flag":    "corp-code",
	})
}

func TestGenericCommandMapsEndpointAndFlags(t *testing.T) {
	var gotPath string
	gotQuery := map[string]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		for key, values := range r.URL.Query() {
			gotQuery[key] = values[0]
		}
		_, _ = w.Write([]byte(`{"status":"000","message":"정상","list":[]}`))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"get",
		"merger-decision",
		"--corp-code", "00126380",
		"--bgn-de", "20250101",
		"--end-de", "20251231",
	})

	require.NoError(t, cmd.Execute())
	assert.Equal(t, "/api/cmpMgDecsn.json", gotPath)
	assert.Equal(t, "env-key", gotQuery["crtfc_key"])
	assert.Equal(t, "00126380", gotQuery["corp_code"])
	assert.Equal(t, "20250101", gotQuery["bgn_de"])
	assert.Equal(t, "20251231", gotQuery["end_de"])
	assert.JSONEq(t, `{"status":"000","message":"정상","list":[]}`, strings.TrimSpace(out.String()))
}

func TestGenericBinaryCommandDefaultsToJSONEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/document.xml", r.URL.Path)
		assert.Equal(t, "20260330000001", r.URL.Query().Get("rcept_no"))
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write([]byte("zip-bytes"))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"download",
		"document",
		"--rcept-no", "20260330000001",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), `"endpoint":"/api/document.xml"`)
	assert.Contains(t, out.String(), `"content_base64":"emlwLWJ5dGVz"`)
}

func TestGenericBinaryCommandSupportsRawOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fnlttXbrl.xml", r.URL.Path)
		assert.Equal(t, "20260330000001", r.URL.Query().Get("rcept_no"))
		assert.Equal(t, "11011", r.URL.Query().Get("reprt_code"))
		_, _ = w.Write([]byte("xbrl-bytes"))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"--output", "raw",
		"download",
		"financial-xbrl",
		"--rcept-no", "20260330000001",
		"--reprt-code", "11011",
	})

	require.NoError(t, cmd.Execute())
	assert.Equal(t, "xbrl-bytes", out.String())
}

func TestAllCatalogCommandsHaveHelpAndRequiredFlagValidation(t *testing.T) {
	root := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))

	for _, spec := range apiCatalog {
		t.Run(spec.Verb+"/"+spec.Resource, func(t *testing.T) {
			assert.NotEmpty(t, spec.APIID)
			assert.NotEmpty(t, spec.Endpoint)
			assert.NotEmpty(t, spec.Verb)
			assert.NotEmpty(t, spec.Resource)
			cmd := findCommand(t, root, spec.Verb, spec.Resource)
			require.NotNil(t, cmd)
			assert.NoError(t, cmd.Help())
			for _, param := range spec.Params {
				flag := cmd.Flags().Lookup(flagName(param.Name))
				require.NotNil(t, flag, "missing flag for %s", param.Name)
			}

			if firstRequired := firstRequiredParam(spec); firstRequired != "" {
				var out bytes.Buffer
				testRoot := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
				testRoot.SetArgs([]string{spec.Verb, spec.Resource})
				err := testRoot.Execute()
				require.Error(t, err)
				assert.Contains(t, err.Error(), "--"+flagName(firstRequired)+" is required")
			}
		})
	}
}

func TestAllCatalogCommandsRouteRequestsAndWriteOutputShape(t *testing.T) {
	for _, spec := range apiCatalog {
		t.Run(spec.Verb+"/"+spec.Resource, func(t *testing.T) {
			var gotPath string
			gotQuery := map[string]string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				for key, values := range r.URL.Query() {
					gotQuery[key] = values[0]
				}
				switch spec.Endpoint {
				case "/api/corpCode.xml":
					_, _ = w.Write(corpCodeFixture(t))
				case "/api/document.xml", "/api/fnlttXbrl.xml":
					w.Header().Set("Content-Type", "application/zip")
					_, _ = w.Write([]byte("file-bytes"))
				default:
					_, _ = w.Write([]byte(`{"status":"000","message":"정상","list":[]}`))
				}
			}))
			t.Cleanup(server.Close)

			var out bytes.Buffer
			cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
			args := []string{"--base-url", server.URL, spec.Verb, spec.Resource}
			for _, param := range spec.Params {
				args = append(args, "--"+flagName(param.Name), sampleParamValue(param.Name))
			}
			cmd.SetArgs(args)

			require.NoError(t, cmd.Execute())
			assert.Equal(t, spec.Endpoint, gotPath)
			assert.Equal(t, "env-key", gotQuery["crtfc_key"])
			for _, param := range spec.Params {
				assert.Equal(t, sampleParamValue(param.Name), gotQuery[param.Name])
			}

			switch spec.Endpoint {
			case "/api/corpCode.xml":
				var codes []map[string]string
				require.NoError(t, json.Unmarshal(out.Bytes(), &codes))
				require.Len(t, codes, 1)
			case "/api/document.xml", "/api/fnlttXbrl.xml":
				assert.Contains(t, out.String(), `"content_base64":"ZmlsZS1ieXRlcw=="`)
			default:
				assert.JSONEq(t, `{"status":"000","message":"정상","list":[]}`, strings.TrimSpace(out.String()))
			}
		})
	}
}

func TestViewCommandsHaveHelpAndRequiredFlagValidation(t *testing.T) {
	paths := [][]string{
		{"get", "quarter-performance"},
		{"get", "annual-performance"},
		{"get", "financial-position"},
		{"get", "cash-flow"},
		{"get", "financial-metric"},
	}
	root := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	for _, path := range paths {
		t.Run(strings.Join(path, "/"), func(t *testing.T) {
			cmd := findCommand(t, root, path...)
			require.NotNil(t, cmd)
			assert.NoError(t, cmd.Help())

			testRoot := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
			testRoot.SetArgs(path)
			err := testRoot.Execute()
			require.Error(t, err)
		})
	}
}

func TestGetQuarterPerformanceOutputsJSONAndCalculatesFourthQuarter(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"get",
		"quarter-performance",
		"--corp-code", "00126380",
		"--year", "2025",
		"--fs-div", "CFS",
	})

	require.NoError(t, cmd.Execute())

	var result quarterPerformanceResult
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))
	require.Len(t, result.Quarters, 4)
	assert.Equal(t, "1Q", result.Quarters[0].Quarter)
	assert.Equal(t, int64(100), result.Quarters[0].Revenue.Amount)
	assert.Equal(t, "4Q", result.Quarters[3].Quarter)
	assert.Equal(t, int64(250), result.Quarters[3].Revenue.Amount)
	assert.Equal(t, int64(50), result.Quarters[3].OperatingIncome.Amount)
	assert.Equal(t, int64(30), result.Quarters[3].NetIncome.Amount)
}

func TestGetQuarterPerformanceOutputsTable(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"--output", "table",
		"get",
		"quarter-performance",
		"--corp-code", "00126380",
		"--year", "2025",
		"--fs-div", "CFS",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "quarter")
	assert.Contains(t, out.String(), "4Q")
	assert.Contains(t, out.String(), "250")
}

func TestGetQuarterPerformanceOutputsCSV(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"--output", "csv",
		"get",
		"quarter-performance",
		"--corp-code", "00126380",
		"--year", "2025",
		"--fs-div", "CFS",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), "corp_code,quarter,report_code,revenue,gross_profit,operating_income,net_income,currency")
	assert.Contains(t, out.String(), "00126380,4Q,11011,250,120,50,30,KRW")
}

func TestGetQuarterPerformanceDetailIncludesSourceMetrics(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"get",
		"quarter-performance",
		"--corp-code", "00126380",
		"--year", "2025",
		"--fs-div", "CFS",
		"--view", "detail",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), `"reports":`)
	assert.Contains(t, out.String(), `"source_row":`)
}

func TestGetFinancialMetricOutputsSourceRow(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"get",
		"financial-metric",
		"--corp-code", "00126380",
		"--year", "2025",
		"--fs-div", "CFS",
		"--metric", "revenue",
		"--view", "source",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), `"metric_code":"revenue"`)
	assert.Contains(t, out.String(), `"source_account_id":"ifrs-full_Revenue"`)
	assert.Contains(t, out.String(), `"source_row":`)
}

func TestGetAnnualPerformanceAcceptsCorpCodes(t *testing.T) {
	server := newFinancialStatementServer(t)
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"get",
		"annual-performance",
		"--corp-codes", "00126380,00126371",
		"--year", "2025",
		"--fs-div", "CFS",
	})

	require.NoError(t, cmd.Execute())

	var result []metricSummaryResult
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))
	require.Len(t, result, 2)
	assert.Equal(t, "00126371", result[1].CorpCode)
	assert.NotEmpty(t, result[0].Metrics)
}

func firstRequiredParam(spec apiSpec) string {
	for _, param := range spec.Params {
		if param.Required {
			return param.Name
		}
	}
	return ""
}

func findCommand(t *testing.T, root *cobra.Command, path ...string) *cobra.Command {
	t.Helper()
	current := root
	for _, name := range path {
		var next *cobra.Command
		for _, child := range current.Commands() {
			if child.Name() == name {
				next = child
				break
			}
		}
		if next == nil {
			return nil
		}
		current = next
	}
	return current
}

func publicCommandNames(root *cobra.Command) []string {
	names := make([]string, 0)
	for _, child := range root.Commands() {
		if !child.Hidden {
			names = append(names, child.Name())
		}
	}
	return names
}

func testGetenv(values map[string]string) getenvFunc {
	return func(key string) string {
		return values[key]
	}
}

func sampleParamValue(name string) string {
	switch name {
	case "corp_code":
		return "00126380"
	case "bsns_year":
		return "2025"
	case "reprt_code":
		return "11011"
	case "rcept_no":
		return "20260330000001"
	case "bgn_de":
		return "20250101"
	case "end_de":
		return "20251231"
	case "sj_div":
		return "BS"
	case "idx_cl_code":
		return "M210000"
	case "page_no":
		return "1"
	case "page_count":
		return "10"
	default:
		return "sample"
	}
}

func corpCodeFixture(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	file, err := zipWriter.Create("CORPCODE.xml")
	require.NoError(t, err)
	_, err = file.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<result>
  <list>
    <corp_code>00126380</corp_code>
    <corp_name>삼성전자</corp_name>
    <corp_eng_name>SAMSUNG ELECTRONICS CO,.LTD</corp_eng_name>
    <stock_code>005930</stock_code>
    <modify_date>20240101</modify_date>
  </list>
</result>`))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())
	return buf.Bytes()
}

func newFinancialStatementServer(t *testing.T) *httptest.Server {
	t.Helper()
	cumulative := map[string]map[string]int64{
		"11013": {"revenue": 100, "gross": 50, "operating": 20, "net": 10},
		"11012": {"revenue": 250, "gross": 120, "operating": 50, "net": 25},
		"11014": {"revenue": 450, "gross": 200, "operating": 80, "net": 40},
		"11011": {"revenue": 700, "gross": 320, "operating": 130, "net": 70},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fnlttSinglAcntAll.json", r.URL.Path)
		assert.Equal(t, "env-key", r.URL.Query().Get("crtfc_key"))
		assert.NotEmpty(t, r.URL.Query().Get("corp_code"))
		assert.Equal(t, "2025", r.URL.Query().Get("bsns_year"))
		assert.Equal(t, "CFS", r.URL.Query().Get("fs_div"))
		reportCode := r.URL.Query().Get("reprt_code")
		values, ok := cumulative[reportCode]
		require.True(t, ok, "unexpected report code %s", reportCode)
		_, _ = w.Write([]byte(financialStatementFixture(reportCode, r.URL.Query().Get("corp_code"), values)))
	}))
}

func financialStatementFixture(reportCode string, corpCode string, values map[string]int64) string {
	amountField := "thstrm_add_amount"
	if reportCode == "11011" {
		amountField = "thstrm_amount"
	}
	rows := []string{
		financialStatementRow(reportCode, corpCode, "ifrs-full_Revenue", "매출액", amountField, values["revenue"]),
		financialStatementRow(reportCode, corpCode, "ifrs-full_GrossProfit", "매출총이익", amountField, values["gross"]),
		financialStatementRow(reportCode, corpCode, "dart_OperatingIncomeLoss", "영업이익", amountField, values["operating"]),
		financialStatementRow(reportCode, corpCode, "ifrs-full_ProfitLoss", "당기순이익", amountField, values["net"]),
	}
	return `{"status":"000","message":"정상","list":[` + strings.Join(rows, ",") + `]}`
}

func financialStatementRow(reportCode string, corpCode string, accountID string, accountName string, amountField string, amount int64) string {
	row := map[string]string{
		"rcept_no":       "20260330000001",
		"reprt_code":     reportCode,
		"bsns_year":      "2025",
		"corp_code":      corpCode,
		"sj_div":         "IS",
		"sj_nm":          "손익계산서",
		"account_id":     accountID,
		"account_nm":     accountName,
		"account_detail": accountName,
		"thstrm_nm":      "제 57 기",
		"ord":            "1",
		"currency":       "KRW",
		amountField:      strconvFormatInt(amount),
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func strconvFormatInt(value int64) string {
	return strconv.FormatInt(value, 10)
}

func assertOopsContext(t *testing.T, err error, expected map[string]any) {
	t.Helper()

	oopsErr, ok := oops.AsOops(err)
	require.True(t, ok, "expected oops error context")

	context := oopsErr.Context()
	for key, value := range expected {
		assert.Equal(t, value, context[key], "oops context %q", key)
	}
}
