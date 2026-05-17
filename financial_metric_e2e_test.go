//go:build e2e

package opendart

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envFinancialMetricAuditTargets      = "OPENDART_E2E_TARGETS"
	envFinancialMetricAuditTargetsFile  = "OPENDART_E2E_TARGETS_FILE"
	envFinancialMetricAuditBusinessYear = "OPENDART_E2E_BSNS_YEAR"
	envFinancialMetricAuditReportCode   = "OPENDART_E2E_REPRT_CODE"
	envFinancialMetricAuditFSDiv        = "OPENDART_E2E_FS_DIV"
	envFinancialMetricAuditMinTargets   = "OPENDART_E2E_MIN_TARGETS"
	envFinancialMetricAuditMinFetched   = "OPENDART_E2E_MIN_FETCHED"
	envFinancialMetricAuditMaxTargets   = "OPENDART_E2E_MAX_TARGETS"
	envFinancialMetricAuditAccountLimit = "OPENDART_E2E_ACCOUNT_REPORT_LIMIT"
)

func TestE2EFinancialMetricAuditForMajorCompanies(t *testing.T) {
	apiKey := strings.TrimSpace(os.Getenv("OPENDART_API_KEY"))
	if apiKey == "" {
		t.Skip("OPENDART_API_KEY is required for live OpenDART e2e financial metric audit")
	}

	targets := loadFinancialMetricAuditTargets(t)
	if len(targets) == 0 {
		t.Skip("set OPENDART_E2E_TARGETS or OPENDART_E2E_TARGETS_FILE for live financial metric audit targets")
	}

	minTargets := financialMetricAuditEnvInt(envFinancialMetricAuditMinTargets, 1)
	require.GreaterOrEqual(t, len(targets), minTargets, "not enough e2e audit targets")

	maxTargets := financialMetricAuditEnvInt(envFinancialMetricAuditMaxTargets, 200)
	if len(targets) > maxTargets {
		targets = targets[:maxTargets]
	}

	client, err := New(Config{APIKey: apiKey})
	require.NoError(t, err)

	ctx := context.Background()
	corpCodes := loadFinancialMetricAuditCorpCodes(ctx, t, client, targets)
	targets = resolveFinancialMetricAuditTargets(t, targets, corpCodes)

	params := buildFinancialMetricAuditParams()
	report := newFinancialMetricAuditReport(targets, params)
	for _, target := range targets {
		response, err := client.FnlttSinglAcntAll(ctx, FnlttSinglAcntAllParams{
			CorpCode:  target.CorpCode,
			BsnsYear:  params.BusinessYear,
			ReprtCode: params.ReportCode,
			FsDiv:     params.FsDiv,
		})
		if err != nil {
			if noData, ok := financialMetricAuditNoDataResult(target, err); ok {
				report.NoData = append(report.NoData, noData)
				continue
			}
			report.Failures = append(report.Failures, financialMetricAuditFailure{
				Target: target,
				Error:  err.Error(),
			})
			continue
		}

		metrics, err := NormalizeFnlttSinglAcntAllResponse(response, FnlttSinglAcntAllParams{
			CorpCode:  target.CorpCode,
			BsnsYear:  params.BusinessYear,
			ReprtCode: params.ReportCode,
			FsDiv:     params.FsDiv,
		})
		require.NoError(t, err)
		report.AddCompany(target, response.List, metrics)

		assert.NotEmpty(t, response.List, "raw rows should not be empty for %s", target.CorpCode)
		for _, metric := range metrics.Metrics {
			assert.NotEmpty(t, metric.MetricCode)
			assert.NotEmpty(t, metric.SourceAccountName)
			assert.NotEmpty(t, metric.MatchMethod)
			assert.GreaterOrEqual(t, metric.SourceRowIndex, 0)
		}
	}

	finished := report.Finish(financialMetricAuditEnvInt(envFinancialMetricAuditAccountLimit, 50))
	encoded, err := json.MarshalIndent(finished, "", "  ")
	require.NoError(t, err)
	t.Logf("financial_metric_audit_report=%s", string(encoded))
	require.Empty(t, finished.Failures, "live OpenDART calls failed")
	require.GreaterOrEqual(t, finished.FetchedCount, financialMetricAuditEnvInt(envFinancialMetricAuditMinFetched, 1), "not enough companies returned financial statement rows")
}

type financialMetricAuditTarget struct {
	CorpCode  string `json:"corp_code,omitempty"`
	StockCode string `json:"stock_code,omitempty"`
	Name      string `json:"name,omitempty"`
}

type financialMetricAuditParams struct {
	BusinessYear string `json:"business_year"`
	ReportCode   string `json:"report_code"`
	FsDiv        string `json:"fs_div"`
}

type financialMetricAuditReport struct {
	Params              financialMetricAuditParams            `json:"params"`
	TargetCount         int                                   `json:"target_count"`
	FetchedCount        int                                   `json:"fetched_count"`
	FailureCount        int                                   `json:"failure_count"`
	NoDataCount         int                                   `json:"no_data_count"`
	TotalRawRows        int                                   `json:"total_raw_rows"`
	TotalMetrics        int                                   `json:"total_metrics"`
	TotalUnmappedRows   int                                   `json:"total_unmapped_rows"`
	TotalIssues         int                                   `json:"total_issues"`
	StatementDivs       map[string]int                        `json:"statement_divs"`
	RawFieldCoverage    []financialMetricAuditCoverage        `json:"raw_field_coverage"`
	MetricCoverage      []financialMetricAuditMetricCoverage  `json:"metric_coverage"`
	AccountIDCoverage   []financialMetricAuditAccountCoverage `json:"account_id_coverage"`
	AccountNameCoverage []financialMetricAuditAccountCoverage `json:"account_name_coverage"`
	Companies           []financialMetricAuditCompany         `json:"companies"`
	Failures            []financialMetricAuditFailure         `json:"failures,omitempty"`
	NoData              []financialMetricAuditNoData          `json:"no_data,omitempty"`
	rawFields           map[string]*financialMetricAuditCounter
	metrics             map[FinancialMetricCode]*financialMetricAuditMetricCounter
	accountIDs          map[string]*financialMetricAuditAccountCounter
	accountNames        map[string]*financialMetricAuditAccountCounter
	targets             []financialMetricAuditTarget
	fetchedTargets      map[string]financialMetricAuditTarget
}

type financialMetricAuditCompany struct {
	Target         financialMetricAuditTarget `json:"target"`
	RawRows        int                        `json:"raw_rows"`
	Metrics        int                        `json:"metrics"`
	UnmappedRows   int                        `json:"unmapped_rows"`
	Issues         int                        `json:"issues"`
	StatementDivs  []string                   `json:"statement_divs"`
	MissingMetrics []FinancialMetricCode      `json:"missing_metrics,omitempty"`
}

type financialMetricAuditFailure struct {
	Target financialMetricAuditTarget `json:"target"`
	Error  string                     `json:"error"`
}

type financialMetricAuditNoData struct {
	Target  financialMetricAuditTarget `json:"target"`
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
}

type financialMetricAuditCoverage struct {
	Field        string  `json:"field"`
	Companies    int     `json:"companies"`
	Rows         int     `json:"rows"`
	CompanyRatio float64 `json:"company_ratio"`
	RowRatio     float64 `json:"row_ratio"`
}

type financialMetricAuditMetricCoverage struct {
	MetricCode       FinancialMetricCode          `json:"metric_code"`
	Companies        int                          `json:"companies"`
	CompanyRatio     float64                      `json:"company_ratio"`
	AccountIDExact   int                          `json:"account_id_exact"`
	AccountNameAlias int                          `json:"account_name_alias"`
	Override         int                          `json:"override"`
	MissingTargets   []financialMetricAuditTarget `json:"missing_targets,omitempty"`
}

type financialMetricAuditAccountCoverage struct {
	Value        string  `json:"value"`
	Companies    int     `json:"companies"`
	Rows         int     `json:"rows"`
	CompanyRatio float64 `json:"company_ratio"`
}

type financialMetricAuditCounter struct {
	rows      int
	companies map[string]struct{}
}

type financialMetricAuditMetricCounter struct {
	companies map[string]struct{}
	methods   map[FinancialMetricMatchMethod]int
}

type financialMetricAuditAccountCounter struct {
	rows      int
	companies map[string]struct{}
}

type financialMetricAuditCorpCodeXML struct {
	Items []financialMetricAuditCorpCodeItem `xml:"list"`
}

type financialMetricAuditCorpCodeItem struct {
	CorpCode  string `xml:"corp_code"`
	CorpName  string `xml:"corp_name"`
	StockCode string `xml:"stock_code"`
}

func loadFinancialMetricAuditTargets(t *testing.T) []financialMetricAuditTarget {
	t.Helper()

	var raw string
	if inline := strings.TrimSpace(os.Getenv(envFinancialMetricAuditTargets)); inline != "" {
		raw = inline
	} else if filePath := strings.TrimSpace(os.Getenv(envFinancialMetricAuditTargetsFile)); filePath != "" {
		body, err := os.ReadFile(filePath)
		require.NoError(t, err)
		raw = string(body)
	}

	targets := make([]financialMetricAuditTarget, 0)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.ReplaceAll(line, ",", " ")
		fields := strings.Fields(line)
		targets = append(targets, parseFinancialMetricAuditTarget(t, fields))
	}
	return targets
}

func parseFinancialMetricAuditTarget(t *testing.T, fields []string) financialMetricAuditTarget {
	t.Helper()
	require.NotEmpty(t, fields)

	target := financialMetricAuditTarget{}
	switch {
	case isFinancialMetricAuditCode(fields[0], 8):
		target.CorpCode = fields[0]
	case isFinancialMetricAuditCode(fields[0], 6):
		target.StockCode = fields[0]
	default:
		t.Fatalf("invalid target identifier %q", fields[0])
	}

	if len(fields) >= 2 {
		switch {
		case isFinancialMetricAuditCode(fields[1], 8) && target.CorpCode == "":
			target.CorpCode = fields[1]
		case isFinancialMetricAuditCode(fields[1], 6) && target.StockCode == "":
			target.StockCode = fields[1]
		default:
			target.Name = fields[1]
		}
	}
	if len(fields) >= 3 {
		target.Name = strings.Join(fields[2:], " ")
	}
	return target
}

func isFinancialMetricAuditCode(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func loadFinancialMetricAuditCorpCodes(ctx context.Context, t *testing.T, client *Client, targets []financialMetricAuditTarget) map[string]financialMetricAuditCorpCodeItem {
	t.Helper()
	needsStockMapping := false
	for _, target := range targets {
		if target.CorpCode == "" {
			needsStockMapping = true
			break
		}
	}
	if !needsStockMapping {
		return nil
	}

	file, err := client.CorpCode(ctx)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(file.Body), int64(len(file.Body)))
	require.NoError(t, err)
	require.NotEmpty(t, reader.File)

	var xmlBody []byte
	for _, zippedFile := range reader.File {
		if strings.HasSuffix(strings.ToLower(zippedFile.Name), ".xml") {
			rc, err := zippedFile.Open()
			require.NoError(t, err)
			xmlBody, err = ioReadAllAndClose(rc)
			require.NoError(t, err)
			break
		}
	}
	require.NotEmpty(t, xmlBody)

	var parsed financialMetricAuditCorpCodeXML
	require.NoError(t, xml.Unmarshal(xmlBody, &parsed))

	byStockCode := make(map[string]financialMetricAuditCorpCodeItem, len(parsed.Items))
	for _, item := range parsed.Items {
		if strings.TrimSpace(item.StockCode) != "" {
			byStockCode[item.StockCode] = item
		}
	}
	return byStockCode
}

func ioReadAllAndClose(rc interface {
	Read([]byte) (int, error)
	Close() error
}) ([]byte, error) {
	defer rc.Close()
	return io.ReadAll(rc)
}

func resolveFinancialMetricAuditTargets(t *testing.T, targets []financialMetricAuditTarget, corpCodes map[string]financialMetricAuditCorpCodeItem) []financialMetricAuditTarget {
	t.Helper()
	resolved := make([]financialMetricAuditTarget, len(targets))
	for index, target := range targets {
		if target.CorpCode == "" {
			item, ok := corpCodes[target.StockCode]
			require.True(t, ok, "stock code %s not found in OpenDART corp code list", target.StockCode)
			target.CorpCode = item.CorpCode
			if target.Name == "" {
				target.Name = item.CorpName
			}
		}
		resolved[index] = target
	}
	return resolved
}

func buildFinancialMetricAuditParams() financialMetricAuditParams {
	return financialMetricAuditParams{
		BusinessYear: financialMetricAuditEnv(envFinancialMetricAuditBusinessYear, "2024"),
		ReportCode:   financialMetricAuditEnv(envFinancialMetricAuditReportCode, ReportCodeAnnual),
		FsDiv:        financialMetricAuditEnv(envFinancialMetricAuditFSDiv, FinancialStatementDivisionConsolidated),
	}
}

func newFinancialMetricAuditReport(targets []financialMetricAuditTarget, params financialMetricAuditParams) *financialMetricAuditReport {
	report := &financialMetricAuditReport{
		Params:         params,
		TargetCount:    len(targets),
		StatementDivs:  map[string]int{},
		rawFields:      map[string]*financialMetricAuditCounter{},
		metrics:        map[FinancialMetricCode]*financialMetricAuditMetricCounter{},
		accountIDs:     map[string]*financialMetricAuditAccountCounter{},
		accountNames:   map[string]*financialMetricAuditAccountCounter{},
		targets:        targets,
		fetchedTargets: map[string]financialMetricAuditTarget{},
	}
	for _, rule := range DefaultFinancialMetricRules() {
		report.metrics[rule.MetricCode] = &financialMetricAuditMetricCounter{
			companies: map[string]struct{}{},
			methods:   map[FinancialMetricMatchMethod]int{},
		}
	}
	return report
}

func (report *financialMetricAuditReport) AddCompany(target financialMetricAuditTarget, rows []FnlttSinglAcntAllItem, metrics *FinancialMetricSet) {
	report.FetchedCount++
	report.fetchedTargets[target.CorpCode] = target
	report.TotalRawRows += len(rows)
	report.TotalMetrics += len(metrics.Metrics)
	report.TotalUnmappedRows += len(metrics.UnmappedRows)
	report.TotalIssues += len(metrics.Issues)

	companyKey := target.CorpCode
	statementDivs := map[string]struct{}{}
	for _, row := range rows {
		if strings.TrimSpace(row.SjDiv) != "" {
			report.StatementDivs[row.SjDiv]++
			statementDivs[row.SjDiv] = struct{}{}
		}
		report.addRawFieldCoverage(companyKey, row)
		report.addAccountCoverage(companyKey, row)
	}
	for _, metric := range metrics.Metrics {
		counter := report.metricCounter(metric.MetricCode)
		counter.companies[companyKey] = struct{}{}
		counter.methods[metric.MatchMethod]++
	}

	report.Companies = append(report.Companies, financialMetricAuditCompany{
		Target:         target,
		RawRows:        len(rows),
		Metrics:        len(metrics.Metrics),
		UnmappedRows:   len(metrics.UnmappedRows),
		Issues:         len(metrics.Issues),
		StatementDivs:  sortedMapKeys(statementDivs),
		MissingMetrics: report.missingMetrics(companyKey),
	})
}

func (report *financialMetricAuditReport) Finish(accountLimit int) *financialMetricAuditReport {
	report.FailureCount = len(report.Failures)
	report.NoDataCount = len(report.NoData)
	report.RawFieldCoverage = report.finishRawFieldCoverage()
	report.MetricCoverage = report.finishMetricCoverage()
	report.AccountIDCoverage = report.finishAccountCoverage(report.accountIDs, accountLimit)
	report.AccountNameCoverage = report.finishAccountCoverage(report.accountNames, accountLimit)
	report.rawFields = nil
	report.metrics = nil
	report.accountIDs = nil
	report.accountNames = nil
	report.targets = nil
	report.fetchedTargets = nil
	return report
}

func financialMetricAuditNoDataResult(target financialMetricAuditTarget, err error) (financialMetricAuditNoData, bool) {
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != "013" {
		return financialMetricAuditNoData{}, false
	}
	return financialMetricAuditNoData{
		Target:  target,
		Status:  apiErr.Status,
		Message: apiErr.Message,
	}, true
}

func (report *financialMetricAuditReport) addRawFieldCoverage(companyKey string, row FnlttSinglAcntAllItem) {
	value := reflect.ValueOf(row)
	valueType := value.Type()
	for index := range value.NumField() {
		field := valueType.Field(index)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		if strings.TrimSpace(value.Field(index).String()) == "" {
			continue
		}
		counter := report.rawFieldCounter(name)
		counter.rows++
		counter.companies[companyKey] = struct{}{}
	}
}

func (report *financialMetricAuditReport) addAccountCoverage(companyKey string, row FnlttSinglAcntAllItem) {
	if accountID := strings.TrimSpace(row.AccountId); accountID != "" {
		counter := report.accountCounter(report.accountIDs, accountID)
		counter.rows++
		counter.companies[companyKey] = struct{}{}
	}
	if accountName := strings.TrimSpace(row.AccountNm); accountName != "" {
		counter := report.accountCounter(report.accountNames, accountName)
		counter.rows++
		counter.companies[companyKey] = struct{}{}
	}
}

func (report *financialMetricAuditReport) rawFieldCounter(name string) *financialMetricAuditCounter {
	counter, ok := report.rawFields[name]
	if !ok {
		counter = &financialMetricAuditCounter{companies: map[string]struct{}{}}
		report.rawFields[name] = counter
	}
	return counter
}

func (report *financialMetricAuditReport) metricCounter(code FinancialMetricCode) *financialMetricAuditMetricCounter {
	counter, ok := report.metrics[code]
	if !ok {
		counter = &financialMetricAuditMetricCounter{
			companies: map[string]struct{}{},
			methods:   map[FinancialMetricMatchMethod]int{},
		}
		report.metrics[code] = counter
	}
	return counter
}

func (report *financialMetricAuditReport) accountCounter(counters map[string]*financialMetricAuditAccountCounter, value string) *financialMetricAuditAccountCounter {
	counter, ok := counters[value]
	if !ok {
		counter = &financialMetricAuditAccountCounter{companies: map[string]struct{}{}}
		counters[value] = counter
	}
	return counter
}

func (report *financialMetricAuditReport) finishRawFieldCoverage() []financialMetricAuditCoverage {
	fields := make([]financialMetricAuditCoverage, 0, len(report.rawFields))
	for field, counter := range report.rawFields {
		fields = append(fields, financialMetricAuditCoverage{
			Field:        field,
			Companies:    len(counter.companies),
			Rows:         counter.rows,
			CompanyRatio: financialMetricAuditRatio(len(counter.companies), report.FetchedCount),
			RowRatio:     financialMetricAuditRatio(counter.rows, report.TotalRawRows),
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Companies == fields[j].Companies {
			return fields[i].Field < fields[j].Field
		}
		return fields[i].Companies > fields[j].Companies
	})
	return fields
}

func (report *financialMetricAuditReport) finishMetricCoverage() []financialMetricAuditMetricCoverage {
	coverages := make([]financialMetricAuditMetricCoverage, 0, len(report.metrics))
	for code, counter := range report.metrics {
		coverages = append(coverages, financialMetricAuditMetricCoverage{
			MetricCode:       code,
			Companies:        len(counter.companies),
			CompanyRatio:     financialMetricAuditRatio(len(counter.companies), report.FetchedCount),
			AccountIDExact:   counter.methods[FinancialMetricMatchAccountIDExact],
			AccountNameAlias: counter.methods[FinancialMetricMatchAccountNameAlias],
			Override:         counter.methods[FinancialMetricMatchOverride],
			MissingTargets:   report.missingTargetsForMetric(counter.companies),
		})
	}
	sort.Slice(coverages, func(i, j int) bool {
		if coverages[i].Companies == coverages[j].Companies {
			return coverages[i].MetricCode < coverages[j].MetricCode
		}
		return coverages[i].Companies > coverages[j].Companies
	})
	return coverages
}

func (report *financialMetricAuditReport) finishAccountCoverage(counters map[string]*financialMetricAuditAccountCounter, limit int) []financialMetricAuditAccountCoverage {
	coverages := make([]financialMetricAuditAccountCoverage, 0, len(counters))
	for value, counter := range counters {
		coverages = append(coverages, financialMetricAuditAccountCoverage{
			Value:        value,
			Companies:    len(counter.companies),
			Rows:         counter.rows,
			CompanyRatio: financialMetricAuditRatio(len(counter.companies), report.FetchedCount),
		})
	}
	sort.Slice(coverages, func(i, j int) bool {
		if coverages[i].Companies == coverages[j].Companies {
			if coverages[i].Rows == coverages[j].Rows {
				return coverages[i].Value < coverages[j].Value
			}
			return coverages[i].Rows > coverages[j].Rows
		}
		return coverages[i].Companies > coverages[j].Companies
	})
	if limit > 0 && len(coverages) > limit {
		return coverages[:limit]
	}
	return coverages
}

func (report *financialMetricAuditReport) missingMetrics(companyKey string) []FinancialMetricCode {
	missing := make([]FinancialMetricCode, 0)
	for code, counter := range report.metrics {
		if _, ok := counter.companies[companyKey]; !ok {
			missing = append(missing, code)
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i] < missing[j] })
	return missing
}

func (report *financialMetricAuditReport) missingTargetsForMetric(companies map[string]struct{}) []financialMetricAuditTarget {
	missing := make([]financialMetricAuditTarget, 0)
	for _, target := range report.targets {
		if _, ok := report.fetchedTargets[target.CorpCode]; !ok {
			continue
		}
		if _, ok := companies[target.CorpCode]; !ok {
			missing = append(missing, target)
		}
	}
	return missing
}

func financialMetricAuditEnv(name string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func financialMetricAuditEnvInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func financialMetricAuditRatio(numerator int, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

func sortedMapKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
