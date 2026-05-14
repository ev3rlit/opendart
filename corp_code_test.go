package opendart

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorpCodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/corpCode.xml", r.URL.Path)
		assert.Equal(t, "test-key", r.URL.Query().Get("crtfc_key"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(corpCodeFixture(t))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	codes, err := client.CorpCodes(context.Background())
	require.NoError(t, err)
	require.Len(t, codes, 2)

	assert.Equal(t, "00126380", codes[0].CorpCode)
	assert.Equal(t, "삼성전자", codes[0].CorpName)
	assert.Equal(t, "005930", codes[0].StockCode)
	assert.Equal(t, "20240101", codes[0].ModifyDate)
	assert.Equal(t, "00164779", codes[1].CorpCode)
}

func TestCorpCodesReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<result><status>010</status><message>등록되지 않은 키입니다.</message></result>`))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.CorpCodes(context.Background())
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, "010", apiErr.Status)
	assert.Equal(t, "등록되지 않은 키입니다.", apiErr.Message)
}

func TestCorpCodesReturnsDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not a zip"))
	}))
	t.Cleanup(server.Close)

	client, err := New(Config{APIKey: "test-key"}, WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = client.CorpCodes(context.Background())
	require.Error(t, err)

	var decodeErr *DecodeError
	require.True(t, errors.As(err, &decodeErr))
	assert.Equal(t, "corpCode.zip", decodeErr.Op)
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
  <list>
    <corp_code>00164779</corp_code>
    <corp_name>현대자동차</corp_name>
    <corp_eng_name>HYUNDAI MOTOR COMPANY</corp_eng_name>
    <stock_code>005380</stock_code>
    <modify_date>20240102</modify_date>
  </list>
</result>`))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())

	return buf.Bytes()
}
