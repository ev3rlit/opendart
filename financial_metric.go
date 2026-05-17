package opendart

import (
	"context"
	"strconv"
	"strings"

	"github.com/samber/oops"
)

// FinancialMetricCode identifies a normalized financial metric.
type FinancialMetricCode string

const (
	// FinancialMetricRevenue is revenue or operating revenue.
	FinancialMetricRevenue FinancialMetricCode = "revenue"
	// FinancialMetricGrossProfit is gross profit.
	FinancialMetricGrossProfit FinancialMetricCode = "gross_profit"
	// FinancialMetricOperatingIncome is operating income or loss.
	FinancialMetricOperatingIncome FinancialMetricCode = "operating_income"
	// FinancialMetricNetIncome is profit or loss for the period.
	FinancialMetricNetIncome FinancialMetricCode = "net_income"
	// FinancialMetricAssets is total assets.
	FinancialMetricAssets FinancialMetricCode = "assets"
	// FinancialMetricCurrentAssets is current assets.
	FinancialMetricCurrentAssets FinancialMetricCode = "current_assets"
	// FinancialMetricNonCurrentAssets is non-current assets.
	FinancialMetricNonCurrentAssets FinancialMetricCode = "non_current_assets"
	// FinancialMetricLiabilities is total liabilities.
	FinancialMetricLiabilities FinancialMetricCode = "liabilities"
	// FinancialMetricCurrentLiabilities is current liabilities.
	FinancialMetricCurrentLiabilities FinancialMetricCode = "current_liabilities"
	// FinancialMetricNonCurrentLiabilities is non-current liabilities.
	FinancialMetricNonCurrentLiabilities FinancialMetricCode = "non_current_liabilities"
	// FinancialMetricEquity is total equity.
	FinancialMetricEquity FinancialMetricCode = "equity"
	// FinancialMetricCashAndCashEquivalents is cash and cash equivalents.
	FinancialMetricCashAndCashEquivalents FinancialMetricCode = "cash_and_cash_equivalents"
	// FinancialMetricOperatingCashFlow is net cash flows from operating activities.
	FinancialMetricOperatingCashFlow FinancialMetricCode = "operating_cash_flow"
)

// FinancialMetricMatchMethod describes how a source row matched a metric rule.
type FinancialMetricMatchMethod string

const (
	// FinancialMetricMatchAccountIDExact means account_id matched exactly.
	FinancialMetricMatchAccountIDExact FinancialMetricMatchMethod = "account_id_exact"
	// FinancialMetricMatchAccountNameAlias means account_nm matched an alias rule.
	FinancialMetricMatchAccountNameAlias FinancialMetricMatchMethod = "account_name_alias"
	// FinancialMetricMatchOverride means an override rule matched the source row.
	FinancialMetricMatchOverride FinancialMetricMatchMethod = "override"
)

// FinancialMetricConfidence describes the reliability of a normalized match.
type FinancialMetricConfidence string

const (
	// FinancialMetricConfidenceHigh is used for exact account_id matches.
	FinancialMetricConfidenceHigh FinancialMetricConfidence = "high"
	// FinancialMetricConfidenceMedium is used for account_nm alias matches.
	FinancialMetricConfidenceMedium FinancialMetricConfidence = "medium"
	// FinancialMetricConfidenceManual is used for explicit override matches.
	FinancialMetricConfidenceManual FinancialMetricConfidence = "manual"
)

const (
	// FinancialMetricAmountCurrent is the OpenDART thstrm_amount field.
	FinancialMetricAmountCurrent = "thstrm_amount"
	// FinancialMetricAmountCurrentCumulative is the OpenDART thstrm_add_amount field.
	FinancialMetricAmountCurrentCumulative = "thstrm_add_amount"
)

// FinancialMetricSource carries call-level context that is not always present in each raw row.
type FinancialMetricSource struct {
	CorpCode     string `json:"corp_code,omitempty"`
	BusinessYear string `json:"business_year,omitempty"`
	ReportCode   string `json:"report_code,omitempty"`
	FsDiv        string `json:"fs_div,omitempty"`
}

// FinancialMetric is a normalized metric with a lossless pointer back to its source row.
type FinancialMetric struct {
	MetricCode          FinancialMetricCode        `json:"metric_code"`
	Label               string                     `json:"label"`
	Amount              int64                      `json:"amount"`
	AmountRaw           string                     `json:"amount_raw"`
	AmountSource        string                     `json:"amount_source"`
	Currency            string                     `json:"currency,omitempty"`
	CorpCode            string                     `json:"corp_code,omitempty"`
	BusinessYear        string                     `json:"business_year,omitempty"`
	ReportCode          string                     `json:"report_code,omitempty"`
	FsDiv               string                     `json:"fs_div,omitempty"`
	StatementDiv        string                     `json:"statement_div,omitempty"`
	SourceAccountID     string                     `json:"source_account_id,omitempty"`
	SourceAccountName   string                     `json:"source_account_name,omitempty"`
	SourceAccountDetail string                     `json:"source_account_detail,omitempty"`
	SourceRowIndex      int                        `json:"source_row_index"`
	SourceRow           FnlttSinglAcntAllItem      `json:"source_row"`
	MatchMethod         FinancialMetricMatchMethod `json:"match_method"`
	MatchedField        string                     `json:"matched_field,omitempty"`
	MatchedValue        string                     `json:"matched_value,omitempty"`
	Confidence          FinancialMetricConfidence  `json:"confidence"`
	RuleID              string                     `json:"rule_id,omitempty"`
}

// FinancialMetricUnmappedRow preserves a raw row that did not map to a normalized metric.
type FinancialMetricUnmappedRow struct {
	SourceRowIndex int                   `json:"source_row_index"`
	SourceRow      FnlttSinglAcntAllItem `json:"source_row"`
}

// FinancialMetricIssue reports a non-fatal normalization problem for a raw row.
type FinancialMetricIssue struct {
	SourceRowIndex int                   `json:"source_row_index"`
	Reason         string                `json:"reason"`
	SourceRow      FnlttSinglAcntAllItem `json:"source_row"`
}

// FinancialMetricSet is the normalized view over a single-account-all raw response.
type FinancialMetricSet struct {
	Source       FinancialMetricSource        `json:"source"`
	Metrics      []FinancialMetric            `json:"metrics"`
	UnmappedRows []FinancialMetricUnmappedRow `json:"unmapped_rows,omitempty"`
	Issues       []FinancialMetricIssue       `json:"issues,omitempty"`
}

// Find returns the first normalized metric with the requested code.
func (set *FinancialMetricSet) Find(code FinancialMetricCode) (FinancialMetric, bool) {
	if set == nil {
		return FinancialMetric{}, false
	}
	for _, metric := range set.Metrics {
		if metric.MetricCode == code {
			return metric, true
		}
	}
	return FinancialMetric{}, false
}

// FinancialMetricOption customizes financial metric normalization.
type FinancialMetricOption func(*financialMetricConfig)

type financialMetricConfig struct {
	rules         []FinancialMetricRule
	overrideRules []FinancialMetricRule
}

// WithFinancialMetricRules adds account mapping rules after the default rules.
func WithFinancialMetricRules(rules ...FinancialMetricRule) FinancialMetricOption {
	return func(config *financialMetricConfig) {
		config.rules = append(config.rules, cloneFinancialMetricRules(rules)...)
	}
}

// WithFinancialMetricOverrideRules adds manual rules that run before default exact and alias rules.
func WithFinancialMetricOverrideRules(rules ...FinancialMetricRule) FinancialMetricOption {
	return func(config *financialMetricConfig) {
		config.overrideRules = append(config.overrideRules, cloneFinancialMetricRules(rules)...)
	}
}

// NormalizeFnlttSinglAcntAll normalizes single-account-all raw rows into common financial metrics.
func NormalizeFnlttSinglAcntAll(rows []FnlttSinglAcntAllItem, source FinancialMetricSource, opts ...FinancialMetricOption) (*FinancialMetricSet, error) {
	config := financialMetricConfig{
		rules: DefaultFinancialMetricRules(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}

	result := &FinancialMetricSet{
		Source: source,
	}
	for index, row := range rows {
		match, ok := matchFinancialMetric(row, config)
		if !ok {
			result.UnmappedRows = append(result.UnmappedRows, FinancialMetricUnmappedRow{
				SourceRowIndex: index,
				SourceRow:      row,
			})
			continue
		}

		amountRaw, amountSource := financialMetricAmount(row)
		amount, err := parseFinancialMetricAmount(amountRaw)
		if err != nil {
			result.Issues = append(result.Issues, FinancialMetricIssue{
				SourceRowIndex: index,
				Reason:         "invalid " + amountSource + ": " + err.Error(),
				SourceRow:      row,
			})
			continue
		}

		result.Metrics = append(result.Metrics, FinancialMetric{
			MetricCode:          match.rule.MetricCode,
			Label:               match.rule.Label,
			Amount:              amount,
			AmountRaw:           amountRaw,
			AmountSource:        amountSource,
			Currency:            row.Currency,
			CorpCode:            firstNonEmpty(row.CorpCode, source.CorpCode),
			BusinessYear:        firstNonEmpty(row.BsnsYear, source.BusinessYear),
			ReportCode:          firstNonEmpty(row.ReprtCode, source.ReportCode),
			FsDiv:               source.FsDiv,
			StatementDiv:        row.SjDiv,
			SourceAccountID:     row.AccountId,
			SourceAccountName:   row.AccountNm,
			SourceAccountDetail: row.AccountDetail,
			SourceRowIndex:      index,
			SourceRow:           row,
			MatchMethod:         match.method,
			MatchedField:        match.field,
			MatchedValue:        match.value,
			Confidence:          match.confidence,
			RuleID:              match.rule.ID,
		})
	}
	return result, nil
}

// NormalizeFnlttSinglAcntAllResponse normalizes a raw single-account-all response with its request params.
func NormalizeFnlttSinglAcntAllResponse(response *FnlttSinglAcntAllResponse, params FnlttSinglAcntAllParams, opts ...FinancialMetricOption) (*FinancialMetricSet, error) {
	if response == nil {
		return nil, oops.In("financial_metric").
			New("opendart: nil FnlttSinglAcntAllResponse")
	}
	return NormalizeFnlttSinglAcntAll(response.List, FinancialMetricSource{
		CorpCode:     params.CorpCode,
		BusinessYear: params.BsnsYear,
		ReportCode:   params.ReprtCode,
		FsDiv:        params.FsDiv,
	}, opts...)
}

// FnlttSinglAcntAllMetrics fetches raw single-account-all rows and returns their normalized metric view.
func (client *Client) FnlttSinglAcntAllMetrics(ctx context.Context, params FnlttSinglAcntAllParams, opts ...FinancialMetricOption) (*FinancialMetricSet, error) {
	response, err := client.FnlttSinglAcntAll(ctx, params)
	if err != nil {
		return nil, err
	}
	return NormalizeFnlttSinglAcntAllResponse(response, params, opts...)
}

func financialMetricAmount(row FnlttSinglAcntAllItem) (string, string) {
	if strings.TrimSpace(row.ThstrmAmount) != "" {
		return row.ThstrmAmount, FinancialMetricAmountCurrent
	}
	return row.ThstrmAddAmount, FinancialMetricAmountCurrentCumulative
}

func parseFinancialMetricAmount(text string) (int64, error) {
	normalized := strings.TrimSpace(text)
	normalized = strings.ReplaceAll(normalized, ",", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	if strings.HasPrefix(normalized, "(") && strings.HasSuffix(normalized, ")") {
		normalized = "-" + strings.TrimSuffix(strings.TrimPrefix(normalized, "("), ")")
	}
	return strconv.ParseInt(normalized, 10, 64)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
