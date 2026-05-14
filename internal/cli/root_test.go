package cli

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommandBuildsExpectedCommands(t *testing.T) {
	cmd := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(nil))

	assert.NotNil(t, findCommand(t, cmd, "corp-codes"))
	assert.NotNil(t, findCommand(t, cmd, "financial-statement"))
	assert.NotNil(t, findCommand(t, cmd, "disclosure", "list"))
	assert.NotNil(t, findCommand(t, cmd, "financial", "single-account"))
	assert.NotNil(t, findCommand(t, cmd, "material", "cmp-mg-decsn"))
	assert.NotNil(t, findCommand(t, cmd, "registration", "equity"))
}

func TestCorpCodesCommandUsesEnvAPIKeyAndSDKClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/corpCode.xml", r.URL.Path)
		assert.Equal(t, "env-key", r.URL.Query().Get("crtfc_key"))
		_, _ = w.Write(corpCodeFixture(t))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{"--base-url", server.URL, "corp-codes"})

	require.NoError(t, cmd.Execute())

	var codes []map[string]string
	require.NoError(t, json.Unmarshal(out.Bytes(), &codes))
	require.Len(t, codes, 1)
	assert.Equal(t, "00126380", codes[0]["corp_code"])
	assert.Equal(t, "005930", codes[0]["stock_code"])
}

func TestFinancialStatementCommandUsesSDKClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fnlttSinglAcnt.json", r.URL.Path)
		assert.Equal(t, "flag-key", r.URL.Query().Get("crtfc_key"))
		assert.Equal(t, "00126380", r.URL.Query().Get("corp_code"))
		assert.Equal(t, "2025", r.URL.Query().Get("bsns_year"))
		assert.Equal(t, "11011", r.URL.Query().Get("reprt_code"))
		_, _ = w.Write([]byte(`{"status":"000","message":"정상","list":[{"rcept_no":"20260330000001","reprt_code":"11011"}]}`))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(nil))
	cmd.SetArgs([]string{
		"--api-key", "flag-key",
		"--base-url", server.URL,
		"financial-statement",
		"--corp-code", "00126380",
		"--business-year", "2025",
		"--report-code", "11011",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), `"rcept_no":"20260330000001"`)
}

func TestCommandRequiresAPIKeyWithoutLeakingSecret(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd := NewRootCommand(&out, &errOut, testGetenv(nil))
	cmd.SetArgs([]string{"corp-codes"})

	err := cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), envAPIKey)
	assert.NotContains(t, err.Error(), "secret")
	assert.Empty(t, out.String())
}

func TestGenericCommandValidatesRequiredFlags(t *testing.T) {
	cmd := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{"financial", "single-account"})

	err := cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--corp-code is required")
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
		"material",
		"cmp-mg-decsn",
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
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write([]byte("zip-bytes"))
	}))
	t.Cleanup(server.Close)

	var out bytes.Buffer
	cmd := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
	cmd.SetArgs([]string{
		"--base-url", server.URL,
		"disclosure",
		"document",
		"--rcept-no", "20260330000001",
	})

	require.NoError(t, cmd.Execute())
	assert.Contains(t, out.String(), `"endpoint":"/api/document.xml"`)
	assert.Contains(t, out.String(), `"content_base64":"emlwLWJ5dGVz"`)
}

func TestAllCatalogCommandsHaveHelpAndRequiredFlagValidation(t *testing.T) {
	root := NewRootCommand(&bytes.Buffer{}, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))

	for _, spec := range apiCatalog {
		t.Run(spec.Group+"/"+spec.Command, func(t *testing.T) {
			assert.NotEmpty(t, spec.APIID)
			assert.NotEmpty(t, spec.Endpoint)
			cmd := findCommand(t, root, spec.Group, spec.Command)
			require.NotNil(t, cmd)
			for _, param := range spec.Params {
				flag := cmd.Flags().Lookup(flagName(param.Name))
				require.NotNil(t, flag, "missing flag for %s", param.Name)
			}

			if firstRequired := firstRequiredParam(spec); firstRequired != "" {
				var out bytes.Buffer
				testRoot := NewRootCommand(&out, &bytes.Buffer{}, testGetenv(map[string]string{envAPIKey: "env-key"}))
				testRoot.SetArgs([]string{spec.Group, spec.Command})
				err := testRoot.Execute()
				require.Error(t, err)
				assert.Contains(t, err.Error(), "--"+flagName(firstRequired)+" is required")
			}
		})
	}
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

func testGetenv(values map[string]string) getenvFunc {
	return func(key string) string {
		return values[key]
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
