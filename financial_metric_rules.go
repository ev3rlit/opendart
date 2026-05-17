package opendart

import "strings"

// FinancialMetricRule maps OpenDART account identifiers or names to one normalized metric.
type FinancialMetricRule struct {
	ID                 string              `json:"id,omitempty"`
	MetricCode         FinancialMetricCode `json:"metric_code"`
	Label              string              `json:"label"`
	StatementDivisions []string            `json:"statement_divisions,omitempty"`
	AccountIDs         []string            `json:"account_ids,omitempty"`
	AccountNameAliases []string            `json:"account_name_aliases,omitempty"`
}

// DefaultFinancialMetricRules returns the built-in common financial metric rules.
func DefaultFinancialMetricRules() []FinancialMetricRule {
	return cloneFinancialMetricRules(defaultFinancialMetricRules)
}

var defaultFinancialMetricRules = []FinancialMetricRule{
	{
		ID:                 "revenue",
		MetricCode:         FinancialMetricRevenue,
		Label:              "Revenue",
		StatementDivisions: []string{StatementDivisionIncomeStatement, StatementDivisionComprehensiveIncome},
		AccountIDs: []string{
			"ifrs-full_Revenue",
			"ifrs-full_RevenueFromContractsWithCustomers",
		},
		AccountNameAliases: []string{"매출액", "매출", "영업수익", "수익"},
	},
	{
		ID:                 "gross_profit",
		MetricCode:         FinancialMetricGrossProfit,
		Label:              "Gross profit",
		StatementDivisions: []string{StatementDivisionIncomeStatement, StatementDivisionComprehensiveIncome},
		AccountIDs:         []string{"ifrs-full_GrossProfit"},
		AccountNameAliases: []string{"매출총이익", "매출총손익"},
	},
	{
		ID:                 "operating_income",
		MetricCode:         FinancialMetricOperatingIncome,
		Label:              "Operating income",
		StatementDivisions: []string{StatementDivisionIncomeStatement, StatementDivisionComprehensiveIncome},
		AccountIDs:         []string{"dart_OperatingIncomeLoss"},
		AccountNameAliases: []string{"영업이익", "영업손익"},
	},
	{
		ID:                 "net_income",
		MetricCode:         FinancialMetricNetIncome,
		Label:              "Net income",
		StatementDivisions: []string{StatementDivisionIncomeStatement, StatementDivisionComprehensiveIncome},
		AccountIDs: []string{
			"ifrs-full_ProfitLoss",
			"ifrs-full_ProfitLossAttributableToOwnersOfParent",
		},
		AccountNameAliases: []string{"당기순이익", "당기순손익", "연결당기순이익", "분기순이익", "반기순이익"},
	},
	{
		ID:                 "assets",
		MetricCode:         FinancialMetricAssets,
		Label:              "Assets",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_Assets"},
		AccountNameAliases: []string{"자산총계", "자산"},
	},
	{
		ID:                 "current_assets",
		MetricCode:         FinancialMetricCurrentAssets,
		Label:              "Current assets",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_CurrentAssets"},
		AccountNameAliases: []string{"유동자산"},
	},
	{
		ID:                 "non_current_assets",
		MetricCode:         FinancialMetricNonCurrentAssets,
		Label:              "Non-current assets",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_NoncurrentAssets"},
		AccountNameAliases: []string{"비유동자산"},
	},
	{
		ID:                 "liabilities",
		MetricCode:         FinancialMetricLiabilities,
		Label:              "Liabilities",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_Liabilities"},
		AccountNameAliases: []string{"부채총계", "부채"},
	},
	{
		ID:                 "current_liabilities",
		MetricCode:         FinancialMetricCurrentLiabilities,
		Label:              "Current liabilities",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_CurrentLiabilities"},
		AccountNameAliases: []string{"유동부채"},
	},
	{
		ID:                 "non_current_liabilities",
		MetricCode:         FinancialMetricNonCurrentLiabilities,
		Label:              "Non-current liabilities",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_NoncurrentLiabilities"},
		AccountNameAliases: []string{"비유동부채"},
	},
	{
		ID:                 "equity",
		MetricCode:         FinancialMetricEquity,
		Label:              "Equity",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_Equity"},
		AccountNameAliases: []string{"자본총계", "자본"},
	},
	{
		ID:                 "cash_and_cash_equivalents",
		MetricCode:         FinancialMetricCashAndCashEquivalents,
		Label:              "Cash and cash equivalents",
		StatementDivisions: []string{StatementDivisionBalanceSheet},
		AccountIDs:         []string{"ifrs-full_CashAndCashEquivalents"},
		AccountNameAliases: []string{"현금및현금성자산", "현금및현금성자산의증가"},
	},
	{
		ID:                 "operating_cash_flow",
		MetricCode:         FinancialMetricOperatingCashFlow,
		Label:              "Operating cash flow",
		StatementDivisions: []string{StatementDivisionCashFlow},
		AccountIDs: []string{
			"ifrs-full_CashFlowsFromUsedInOperatingActivities",
			"ifrs-full_CashFlowsFromUsedInOperations",
		},
		AccountNameAliases: []string{"영업활동현금흐름", "영업활동으로인한현금흐름"},
	},
}

type financialMetricRuleMatch struct {
	rule       FinancialMetricRule
	method     FinancialMetricMatchMethod
	field      string
	value      string
	confidence FinancialMetricConfidence
}

func matchFinancialMetric(row FnlttSinglAcntAllItem, config financialMetricConfig) (financialMetricRuleMatch, bool) {
	if match, ok := matchFinancialMetricRules(row, config.overrideRules, true, true, true, FinancialMetricMatchOverride, FinancialMetricConfidenceManual); ok {
		return match, true
	}
	if match, ok := matchFinancialMetricRules(row, config.rules, true, false, false, FinancialMetricMatchAccountIDExact, FinancialMetricConfidenceHigh); ok {
		return match, true
	}
	return matchFinancialMetricRules(row, config.rules, false, true, false, FinancialMetricMatchAccountNameAlias, FinancialMetricConfidenceMedium)
}

func matchFinancialMetricRules(
	row FnlttSinglAcntAllItem,
	rules []FinancialMetricRule,
	matchAccountID bool,
	matchAccountName bool,
	allowStandardUnusedID bool,
	method FinancialMetricMatchMethod,
	confidence FinancialMetricConfidence,
) (financialMetricRuleMatch, bool) {
	for _, rule := range rules {
		if !financialMetricStatementMatches(row, rule) {
			continue
		}
		if matchAccountID {
			if value, ok := exactStringMatch(row.AccountId, rule.AccountIDs, allowStandardUnusedID); ok {
				return financialMetricRuleMatch{
					rule:       rule,
					method:     method,
					field:      "account_id",
					value:      value,
					confidence: confidence,
				}, true
			}
		}
		if matchAccountName {
			if value, ok := normalizedStringMatch(row.AccountNm, rule.AccountNameAliases); ok {
				return financialMetricRuleMatch{
					rule:       rule,
					method:     method,
					field:      "account_nm",
					value:      value,
					confidence: confidence,
				}, true
			}
		}
	}
	return financialMetricRuleMatch{}, false
}

func financialMetricStatementMatches(row FnlttSinglAcntAllItem, rule FinancialMetricRule) bool {
	if len(rule.StatementDivisions) == 0 {
		return true
	}
	return stringInSlice(row.SjDiv, rule.StatementDivisions)
}

func exactStringMatch(value string, candidates []string, allowStandardUnused bool) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}
	if trimmed == "-표준계정코드 미사용-" && !allowStandardUnused {
		return "", false
	}
	for _, candidate := range candidates {
		if trimmed == strings.TrimSpace(candidate) {
			return candidate, true
		}
	}
	return "", false
}

func normalizedStringMatch(value string, candidates []string) (string, bool) {
	normalized := normalizeFinancialMetricText(value)
	if normalized == "" {
		return "", false
	}
	for _, candidate := range candidates {
		if normalized == normalizeFinancialMetricText(candidate) {
			return candidate, true
		}
	}
	return "", false
}

func normalizeFinancialMetricText(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), ""))
}

func stringInSlice(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if value == candidate {
			return true
		}
	}
	return false
}

func cloneFinancialMetricRules(rules []FinancialMetricRule) []FinancialMetricRule {
	cloned := make([]FinancialMetricRule, len(rules))
	for index, rule := range rules {
		cloned[index] = rule
		cloned[index].StatementDivisions = append([]string(nil), rule.StatementDivisions...)
		cloned[index].AccountIDs = append([]string(nil), rule.AccountIDs...)
		cloned[index].AccountNameAliases = append([]string(nil), rule.AccountNameAliases...)
	}
	return cloned
}
