package cli

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ev3rlit/opendart"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

type financialViewParams struct {
	corpCode  string
	corpCodes string
	year      string
	fsDiv     string
	view      string
}

type metricValue struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency,omitempty"`
}

type quarterPerformanceResult struct {
	CorpCode string                     `json:"corp_code"`
	Year     string                     `json:"year"`
	FsDiv    string                     `json:"fs_div"`
	Quarters []quarterPerformancePeriod `json:"quarters"`
	Reports  []metricDetailResult       `json:"reports,omitempty"`
}

type quarterPerformancePeriod struct {
	Quarter         string      `json:"quarter"`
	ReportCode      string      `json:"report_code"`
	Revenue         metricValue `json:"revenue"`
	GrossProfit     metricValue `json:"gross_profit"`
	OperatingIncome metricValue `json:"operating_income"`
	NetIncome       metricValue `json:"net_income"`
}

type namedMetricValue struct {
	Metric   opendart.FinancialMetricCode `json:"metric"`
	Label    string                       `json:"label"`
	Amount   int64                        `json:"amount"`
	Currency string                       `json:"currency,omitempty"`
}

type metricSummaryResult struct {
	CorpCode   string             `json:"corp_code"`
	Year       string             `json:"year"`
	FsDiv      string             `json:"fs_div"`
	ReportCode string             `json:"report_code"`
	Metrics    []namedMetricValue `json:"metrics"`
}

type metricDetailResult struct {
	CorpCode   string                     `json:"corp_code"`
	Year       string                     `json:"year"`
	FsDiv      string                     `json:"fs_div"`
	ReportCode string                     `json:"report_code"`
	Metrics    []opendart.FinancialMetric `json:"metrics"`
}

type companyComparisonResult struct {
	Year      string                `json:"year"`
	FsDiv     string                `json:"fs_div"`
	Companies []metricSummaryResult `json:"companies"`
}

func addFinancialViewCommands(parent *cobra.Command, options *rootOptions) {
	parent.AddCommand(newQuarterPerformanceCommand(options, "quarter-performance", runQuarterPerformanceSummary))
	parent.AddCommand(newMetricSummaryCommand(options, "annual-performance", "연간 손익 주요 지표를 조회합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricRevenue,
		opendart.FinancialMetricGrossProfit,
		opendart.FinancialMetricOperatingIncome,
		opendart.FinancialMetricNetIncome,
	}))
	parent.AddCommand(newMetricSummaryCommand(options, "financial-position", "재무상태표 주요 지표를 조회합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricAssets,
		opendart.FinancialMetricCurrentAssets,
		opendart.FinancialMetricNonCurrentAssets,
		opendart.FinancialMetricLiabilities,
		opendart.FinancialMetricCurrentLiabilities,
		opendart.FinancialMetricNonCurrentLiabilities,
		opendart.FinancialMetricEquity,
	}))
	parent.AddCommand(newMetricSummaryCommand(options, "cash-flow", "현금흐름 주요 지표를 조회합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricCashAndCashEquivalents,
		opendart.FinancialMetricOperatingCashFlow,
	}))
	parent.AddCommand(newFinancialMetricCommand(options, "financial-metric", viewSummary))
}

func newSummarizeCommand(options *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "OpenDART 원문 응답을 요약합니다.",
	}
	cmd.AddCommand(newQuarterPerformanceCommand(options, "quarter-performance", runQuarterPerformanceSummary))
	cmd.AddCommand(newMetricSummaryCommand(options, "annual-performance", "연간 손익 주요 지표를 요약합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricRevenue,
		opendart.FinancialMetricGrossProfit,
		opendart.FinancialMetricOperatingIncome,
		opendart.FinancialMetricNetIncome,
	}))
	cmd.AddCommand(newMetricSummaryCommand(options, "financial-position", "재무상태표 주요 지표를 요약합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricAssets,
		opendart.FinancialMetricCurrentAssets,
		opendart.FinancialMetricNonCurrentAssets,
		opendart.FinancialMetricLiabilities,
		opendart.FinancialMetricCurrentLiabilities,
		opendart.FinancialMetricNonCurrentLiabilities,
		opendart.FinancialMetricEquity,
	}))
	cmd.AddCommand(newMetricSummaryCommand(options, "cash-flow", "현금흐름 주요 지표를 요약합니다.", []opendart.FinancialMetricCode{
		opendart.FinancialMetricCashAndCashEquivalents,
		opendart.FinancialMetricOperatingCashFlow,
	}))
	return cmd
}

func newCompareCommand(options *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "회사와 기간별 주요 지표를 비교합니다.",
	}
	cmd.AddCommand(newCompareQuarterPerformanceCommand(options))
	cmd.AddCommand(newCompareCompaniesCommand(options))
	return cmd
}

func newInspectCommand(options *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "OpenDART 가공 지표의 원천 row를 확인합니다.",
	}
	cmd.AddCommand(newFinancialMetricCommand(options, "metric", viewSource))
	return cmd
}

func newQuarterPerformanceCommand(options *rootOptions, use string, run func(context.Context, *opendart.Client, financialViewParams) (*quarterPerformanceResult, error)) *cobra.Command {
	params := financialViewParams{view: viewSummary}
	cmd := &cobra.Command{
		Use:   use,
		Short: "분기 손익 주요 지표를 조회합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := validateFinancialViewParams(use, params)
			if err != nil {
				return err
			}
			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			results := make([]quarterPerformanceResult, 0, len(codes))
			for _, code := range codes {
				target := params
				target.corpCode = code
				result, err := run(cmd.Context(), client, target)
				if err != nil {
					return err
				}
				results = append(results, *result)
			}
			return writeQuarterPerformanceResults(options, results)
		},
	}
	addFinancialViewFlags(cmd, &params)
	return cmd
}

func newMetricSummaryCommand(options *rootOptions, use string, short string, metrics []opendart.FinancialMetricCode) *cobra.Command {
	params := financialViewParams{view: viewSummary}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := validateFinancialViewParams(use, params)
			if err != nil {
				return err
			}
			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			if params.view == viewSummary {
				results := make([]metricSummaryResult, 0, len(codes))
				for _, code := range codes {
					target := params
					target.corpCode = code
					result, err := runMetricSummary(cmd.Context(), client, target, opendart.ReportCodeAnnual, metrics)
					if err != nil {
						return err
					}
					results = append(results, *result)
				}
				return writeMetricSummaryResults(options, results)
			}
			results := make([]metricDetailResult, 0, len(codes))
			for _, code := range codes {
				target := params
				target.corpCode = code
				result, err := runMetricDetail(cmd.Context(), client, target, opendart.ReportCodeAnnual, metrics)
				if err != nil {
					return err
				}
				results = append(results, *result)
			}
			return writeMetricDetailResults(options, results)
		},
	}
	addFinancialViewFlags(cmd, &params)
	return cmd
}

func newCompareQuarterPerformanceCommand(options *rootOptions) *cobra.Command {
	var corpCodes string
	params := financialViewParams{view: viewSummary}
	cmd := &cobra.Command{
		Use:   "quarter-performance",
		Short: "여러 회사의 분기 손익 주요 지표를 비교합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := parseCorpCodes("compare", corpCodes)
			if err != nil {
				return err
			}
			if _, err := validateFinancialViewParams("quarter-performance", financialViewParams{corpCode: codes[0], year: params.year, fsDiv: params.fsDiv, view: viewSummary}); err != nil {
				return err
			}
			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			results := make([]quarterPerformanceResult, 0, len(codes))
			for _, code := range codes {
				result, err := runQuarterPerformanceSummary(cmd.Context(), client, financialViewParams{
					corpCode: code,
					year:     params.year,
					fsDiv:    params.fsDiv,
				})
				if err != nil {
					return err
				}
				results = append(results, *result)
			}
			return writeQuarterPerformanceResults(options, results)
		},
	}
	cmd.Flags().StringVar(&corpCodes, "corp-codes", "", "Comma-separated OpenDART corp_code values.")
	cmd.Flags().StringVar(&params.year, "year", "", "Business year, for example 2025.")
	cmd.Flags().StringVar(&params.fsDiv, "fs-div", "", "Financial statement division: CFS or OFS.")
	return cmd
}

func newCompareCompaniesCommand(options *rootOptions) *cobra.Command {
	var corpCodes string
	params := financialViewParams{view: viewSummary}
	cmd := &cobra.Command{
		Use:   "companies",
		Short: "여러 회사의 연간 손익 주요 지표를 비교합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := parseCorpCodes("compare", corpCodes)
			if err != nil {
				return err
			}
			if _, err := validateFinancialViewParams("companies", financialViewParams{corpCode: codes[0], year: params.year, fsDiv: params.fsDiv, view: viewSummary}); err != nil {
				return err
			}
			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			result := companyComparisonResult{
				Year:      params.year,
				FsDiv:     params.fsDiv,
				Companies: make([]metricSummaryResult, 0, len(codes)),
			}
			for _, code := range codes {
				summary, err := runMetricSummary(cmd.Context(), client, financialViewParams{
					corpCode: code,
					year:     params.year,
					fsDiv:    params.fsDiv,
				}, opendart.ReportCodeAnnual, []opendart.FinancialMetricCode{
					opendart.FinancialMetricRevenue,
					opendart.FinancialMetricOperatingIncome,
					opendart.FinancialMetricNetIncome,
				})
				if err != nil {
					return err
				}
				result.Companies = append(result.Companies, *summary)
			}
			return writeView(options, result, writeCompanyComparisonTable(result), writeCompanyComparisonCSV(result))
		},
	}
	cmd.Flags().StringVar(&corpCodes, "corp-codes", "", "Comma-separated OpenDART corp_code values.")
	cmd.Flags().StringVar(&params.year, "year", "", "Business year, for example 2025.")
	cmd.Flags().StringVar(&params.fsDiv, "fs-div", "", "Financial statement division: CFS or OFS.")
	return cmd
}

func newFinancialMetricCommand(options *rootOptions, use string, defaultView string) *cobra.Command {
	params := financialViewParams{view: defaultView}
	var metricCode string
	var reportCode string
	cmd := &cobra.Command{
		Use:   use,
		Short: "정규화 지표와 원천 row를 확인합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			codes, err := validateFinancialViewParams(use, params)
			if err != nil {
				return err
			}
			if strings.TrimSpace(metricCode) == "" {
				return requiredFlagError(use, "metric")
			}
			if strings.TrimSpace(reportCode) == "" {
				reportCode = opendart.ReportCodeAnnual
			}
			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			results := make([]opendart.FinancialMetric, 0, len(codes))
			for _, code := range codes {
				target := params
				target.corpCode = code
				set, err := fetchNormalizedMetrics(cmd.Context(), client, target, reportCode)
				if err != nil {
					return err
				}
				metric, ok := set.Find(opendart.FinancialMetricCode(metricCode))
				if !ok {
					return oops.In("opendart_cli").
						With("metric", metricCode, "corp_code", code).
						Errorf("opendart cli: metric %q was not found", metricCode)
				}
				results = append(results, metric)
			}
			if params.view == viewSummary {
				summaries := make([]namedMetricValue, 0, len(results))
				for _, metric := range results {
					summaries = append(summaries, namedMetricValue{
						Metric:   metric.MetricCode,
						Label:    metric.Label,
						Amount:   metric.Amount,
						Currency: metric.Currency,
					})
				}
				return writeFinancialMetricSummaryResults(options, summaries)
			}
			return writeFinancialMetricResults(options, results)
		},
	}
	addFinancialViewFlags(cmd, &params)
	cmd.Flags().StringVar(&metricCode, "metric", "", "Metric code, for example revenue.")
	cmd.Flags().StringVar(&reportCode, "report-code", opendart.ReportCodeAnnual, "Report code, for example 11011.")
	return cmd
}

func addFinancialViewFlags(cmd *cobra.Command, params *financialViewParams) {
	cmd.Flags().StringVar(&params.corpCode, "corp-code", "", "OpenDART corp_code.")
	cmd.Flags().StringVar(&params.corpCodes, "corp-codes", "", "Comma-separated OpenDART corp_code values.")
	cmd.Flags().StringVar(&params.year, "year", "", "Business year, for example 2025.")
	cmd.Flags().StringVar(&params.fsDiv, "fs-div", "", "Financial statement division: CFS or OFS.")
	cmd.Flags().StringVar(&params.view, "view", params.view, "View depth: summary, detail, or source.")
}

func validateFinancialViewParams(command string, params financialViewParams) ([]string, error) {
	switch strings.TrimSpace(params.view) {
	case viewSummary, viewDetail, viewSource:
	default:
		return nil, oops.In("opendart_cli").
			With("command", command, "view", params.view).
			Errorf("opendart cli: unsupported view %q", params.view)
	}
	codes := financialViewCorpCodes(params)
	if len(codes) == 0 {
		return nil, requiredFlagError(command, "corp-code")
	}
	if strings.TrimSpace(params.year) == "" {
		return nil, requiredFlagError(command, "year")
	}
	if strings.TrimSpace(params.fsDiv) == "" {
		return nil, requiredFlagError(command, "fs-div")
	}
	return codes, nil
}

func financialViewCorpCodes(params financialViewParams) []string {
	codes := make([]string, 0, 1)
	if code := strings.TrimSpace(params.corpCode); code != "" {
		codes = append(codes, code)
	}
	return append(codes, splitCorpCodes(params.corpCodes)...)
}

func runQuarterPerformanceSummary(ctx context.Context, client *opendart.Client, params financialViewParams) (*quarterPerformanceResult, error) {
	reportCodes := []string{
		opendart.ReportCodeFirstQuarter,
		opendart.ReportCodeHalfYear,
		opendart.ReportCodeThirdQuarter,
		opendart.ReportCodeAnnual,
	}
	cumulative := make([]map[opendart.FinancialMetricCode]metricValue, 0, len(reportCodes))
	reports := make([]metricDetailResult, 0, len(reportCodes))
	for _, reportCode := range reportCodes {
		set, err := fetchNormalizedMetrics(ctx, client, params, reportCode)
		if err != nil {
			return nil, err
		}
		cumulative = append(cumulative, cumulativeMetricValues(set, []opendart.FinancialMetricCode{
			opendart.FinancialMetricRevenue,
			opendart.FinancialMetricGrossProfit,
			opendart.FinancialMetricOperatingIncome,
			opendart.FinancialMetricNetIncome,
		}))
		if params.view != viewSummary {
			reports = append(reports, metricDetailResult{
				CorpCode:   params.corpCode,
				Year:       params.year,
				FsDiv:      params.fsDiv,
				ReportCode: reportCode,
				Metrics: normalizedMetrics(set, []opendart.FinancialMetricCode{
					opendart.FinancialMetricRevenue,
					opendart.FinancialMetricGrossProfit,
					opendart.FinancialMetricOperatingIncome,
					opendart.FinancialMetricNetIncome,
				}),
			})
		}
	}

	result := &quarterPerformanceResult{
		CorpCode: params.corpCode,
		Year:     params.year,
		FsDiv:    params.fsDiv,
		Quarters: make([]quarterPerformancePeriod, 0, 4),
	}
	for index, reportCode := range reportCodes {
		previous := map[opendart.FinancialMetricCode]metricValue{}
		if index > 0 {
			previous = cumulative[index-1]
		}
		current := cumulative[index]
		result.Quarters = append(result.Quarters, quarterPerformancePeriod{
			Quarter:         []string{"1Q", "2Q", "3Q", "4Q"}[index],
			ReportCode:      reportCode,
			Revenue:         subtractMetric(current, previous, opendart.FinancialMetricRevenue),
			GrossProfit:     subtractMetric(current, previous, opendart.FinancialMetricGrossProfit),
			OperatingIncome: subtractMetric(current, previous, opendart.FinancialMetricOperatingIncome),
			NetIncome:       subtractMetric(current, previous, opendart.FinancialMetricNetIncome),
		})
	}
	if params.view != viewSummary {
		result.Reports = reports
	}
	return result, nil
}

func runMetricSummary(ctx context.Context, client *opendart.Client, params financialViewParams, reportCode string, metrics []opendart.FinancialMetricCode) (*metricSummaryResult, error) {
	set, err := fetchNormalizedMetrics(ctx, client, params, reportCode)
	if err != nil {
		return nil, err
	}
	result := &metricSummaryResult{
		CorpCode:   params.corpCode,
		Year:       params.year,
		FsDiv:      params.fsDiv,
		ReportCode: reportCode,
		Metrics:    make([]namedMetricValue, 0, len(metrics)),
	}
	for _, code := range metrics {
		metric, ok := set.Find(code)
		if !ok {
			continue
		}
		result.Metrics = append(result.Metrics, namedMetricValue{
			Metric:   metric.MetricCode,
			Label:    metric.Label,
			Amount:   metric.Amount,
			Currency: metric.Currency,
		})
	}
	return result, nil
}

func runMetricDetail(ctx context.Context, client *opendart.Client, params financialViewParams, reportCode string, metrics []opendart.FinancialMetricCode) (*metricDetailResult, error) {
	set, err := fetchNormalizedMetrics(ctx, client, params, reportCode)
	if err != nil {
		return nil, err
	}
	result := &metricDetailResult{
		CorpCode:   params.corpCode,
		Year:       params.year,
		FsDiv:      params.fsDiv,
		ReportCode: reportCode,
		Metrics:    normalizedMetrics(set, metrics),
	}
	return result, nil
}

func normalizedMetrics(set *opendart.FinancialMetricSet, metrics []opendart.FinancialMetricCode) []opendart.FinancialMetric {
	result := make([]opendart.FinancialMetric, 0, len(metrics))
	for _, code := range metrics {
		metric, ok := set.Find(code)
		if !ok {
			continue
		}
		result = append(result, metric)
	}
	return result
}

func fetchNormalizedMetrics(ctx context.Context, client *opendart.Client, params financialViewParams, reportCode string) (*opendart.FinancialMetricSet, error) {
	return client.FnlttSinglAcntAllMetrics(ctx, opendart.FnlttSinglAcntAllParams{
		CorpCode:  params.corpCode,
		BsnsYear:  params.year,
		ReprtCode: reportCode,
		FsDiv:     params.fsDiv,
	})
}

func cumulativeMetricValues(set *opendart.FinancialMetricSet, metrics []opendart.FinancialMetricCode) map[opendart.FinancialMetricCode]metricValue {
	result := make(map[opendart.FinancialMetricCode]metricValue, len(metrics))
	for _, code := range metrics {
		metric, ok := set.Find(code)
		if !ok {
			continue
		}
		amount := metric.Amount
		if strings.TrimSpace(metric.SourceRow.ThstrmAddAmount) != "" {
			if parsed, err := parseCLIAmount(metric.SourceRow.ThstrmAddAmount); err == nil {
				amount = parsed
			}
		}
		result[code] = metricValue{
			Amount:   amount,
			Currency: metric.Currency,
		}
	}
	return result
}

func subtractMetric(current map[opendart.FinancialMetricCode]metricValue, previous map[opendart.FinancialMetricCode]metricValue, code opendart.FinancialMetricCode) metricValue {
	value := current[code]
	if prior, ok := previous[code]; ok {
		value.Amount -= prior.Amount
	}
	return value
}

func parseCLIAmount(value string) (int64, error) {
	normalized := strings.TrimSpace(value)
	normalized = strings.ReplaceAll(normalized, ",", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	if strings.HasPrefix(normalized, "(") && strings.HasSuffix(normalized, ")") {
		normalized = "-" + strings.TrimSuffix(strings.TrimPrefix(normalized, "("), ")")
	}
	return strconv.ParseInt(normalized, 10, 64)
}

func parseCorpCodes(command string, value string) ([]string, error) {
	codes := splitCorpCodes(value)
	if len(codes) == 0 {
		return nil, requiredFlagError(command, "corp-codes")
	}
	return codes, nil
}

func splitCorpCodes(value string) []string {
	parts := strings.Split(value, ",")
	codes := make([]string, 0, len(parts))
	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code != "" {
			codes = append(codes, code)
		}
	}
	return codes
}

func writeQuarterPerformanceResults(options *rootOptions, results []quarterPerformanceResult) error {
	return writeView(options, singleOrMany(results), writeQuarterPerformanceTable(results), writeQuarterPerformanceCSV(results))
}

func writeQuarterPerformanceTable(results []quarterPerformanceResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = io.WriteString(writer, "corp_code\tquarter\treport_code\trevenue\tgross_profit\toperating_income\tnet_income\tcurrency\n")
		for _, result := range results {
			for _, quarter := range result.Quarters {
				_, _ = io.WriteString(writer, result.CorpCode+"\t"+quarter.Quarter+"\t"+quarter.ReportCode+"\t"+
					strconv.FormatInt(quarter.Revenue.Amount, 10)+"\t"+
					strconv.FormatInt(quarter.GrossProfit.Amount, 10)+"\t"+
					strconv.FormatInt(quarter.OperatingIncome.Amount, 10)+"\t"+
					strconv.FormatInt(quarter.NetIncome.Amount, 10)+"\t"+
					quarter.Revenue.Currency+"\n")
			}
		}
		return writer.Flush()
	}
}

func writeQuarterPerformanceCSV(results []quarterPerformanceResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := csv.NewWriter(out)
		if err := writer.Write([]string{"corp_code", "quarter", "report_code", "revenue", "gross_profit", "operating_income", "net_income", "currency"}); err != nil {
			return err
		}
		for _, result := range results {
			for _, quarter := range result.Quarters {
				if err := writer.Write([]string{
					result.CorpCode,
					quarter.Quarter,
					quarter.ReportCode,
					strconv.FormatInt(quarter.Revenue.Amount, 10),
					strconv.FormatInt(quarter.GrossProfit.Amount, 10),
					strconv.FormatInt(quarter.OperatingIncome.Amount, 10),
					strconv.FormatInt(quarter.NetIncome.Amount, 10),
					quarter.Revenue.Currency,
				}); err != nil {
					return err
				}
			}
		}
		writer.Flush()
		return writer.Error()
	}
}

func writeMetricSummaryResults(options *rootOptions, results []metricSummaryResult) error {
	return writeView(options, singleOrMany(results), writeMetricSummaryTable(results), writeMetricSummaryCSV(results))
}

func writeMetricSummaryTable(results []metricSummaryResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = io.WriteString(writer, "corp_code\treport_code\tmetric\tlabel\tamount\tcurrency\n")
		for _, result := range results {
			for _, metric := range result.Metrics {
				_, _ = io.WriteString(writer, result.CorpCode+"\t"+result.ReportCode+"\t"+string(metric.Metric)+"\t"+metric.Label+"\t"+strconv.FormatInt(metric.Amount, 10)+"\t"+metric.Currency+"\n")
			}
		}
		return writer.Flush()
	}
}

func writeMetricSummaryCSV(results []metricSummaryResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := csv.NewWriter(out)
		if err := writer.Write([]string{"corp_code", "report_code", "metric", "label", "amount", "currency"}); err != nil {
			return err
		}
		for _, result := range results {
			for _, metric := range result.Metrics {
				if err := writer.Write([]string{result.CorpCode, result.ReportCode, string(metric.Metric), metric.Label, strconv.FormatInt(metric.Amount, 10), metric.Currency}); err != nil {
					return err
				}
			}
		}
		writer.Flush()
		return writer.Error()
	}
}

func writeMetricDetailResults(options *rootOptions, results []metricDetailResult) error {
	return writeView(options, singleOrMany(results), writeMetricDetailTable(results), writeMetricDetailCSV(results))
}

func writeMetricDetailTable(results []metricDetailResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = io.WriteString(writer, "corp_code\treport_code\tmetric\tamount\tcurrency\taccount_id\taccount_name\tmatch\n")
		for _, result := range results {
			for _, metric := range result.Metrics {
				_, _ = io.WriteString(writer, result.CorpCode+"\t"+result.ReportCode+"\t"+string(metric.MetricCode)+"\t"+strconv.FormatInt(metric.Amount, 10)+"\t"+metric.Currency+"\t"+metric.SourceAccountID+"\t"+metric.SourceAccountName+"\t"+string(metric.MatchMethod)+"\n")
			}
		}
		return writer.Flush()
	}
}

func writeMetricDetailCSV(results []metricDetailResult) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := csv.NewWriter(out)
		if err := writer.Write([]string{"corp_code", "report_code", "metric", "amount", "currency", "account_id", "account_name", "match_method", "source_row_index"}); err != nil {
			return err
		}
		for _, result := range results {
			for _, metric := range result.Metrics {
				if err := writer.Write([]string{result.CorpCode, result.ReportCode, string(metric.MetricCode), strconv.FormatInt(metric.Amount, 10), metric.Currency, metric.SourceAccountID, metric.SourceAccountName, string(metric.MatchMethod), strconv.Itoa(metric.SourceRowIndex)}); err != nil {
					return err
				}
			}
		}
		writer.Flush()
		return writer.Error()
	}
}

func writeFinancialMetricSummaryResults(options *rootOptions, results []namedMetricValue) error {
	return writeView(options, singleOrMany(results), writeFinancialMetricSummaryTable(results), writeFinancialMetricSummaryCSV(results))
}

func writeFinancialMetricSummaryTable(results []namedMetricValue) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = io.WriteString(writer, "metric\tlabel\tamount\tcurrency\n")
		for _, metric := range results {
			_, _ = io.WriteString(writer, string(metric.Metric)+"\t"+metric.Label+"\t"+strconv.FormatInt(metric.Amount, 10)+"\t"+metric.Currency+"\n")
		}
		return writer.Flush()
	}
}

func writeFinancialMetricSummaryCSV(results []namedMetricValue) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := csv.NewWriter(out)
		if err := writer.Write([]string{"metric", "label", "amount", "currency"}); err != nil {
			return err
		}
		for _, metric := range results {
			if err := writer.Write([]string{string(metric.Metric), metric.Label, strconv.FormatInt(metric.Amount, 10), metric.Currency}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	}
}

func writeFinancialMetricResults(options *rootOptions, results []opendart.FinancialMetric) error {
	return writeView(options, singleOrMany(results), writeFinancialMetricTable(results), writeFinancialMetricCSV(results))
}

func writeFinancialMetricTable(results []opendart.FinancialMetric) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		_, _ = io.WriteString(writer, "metric\tamount\tcurrency\taccount_id\taccount_name\tmatch\n")
		for _, metric := range results {
			_, _ = io.WriteString(writer, string(metric.MetricCode)+"\t"+strconv.FormatInt(metric.Amount, 10)+"\t"+metric.Currency+"\t"+metric.SourceAccountID+"\t"+metric.SourceAccountName+"\t"+string(metric.MatchMethod)+"\n")
		}
		return writer.Flush()
	}
}

func writeFinancialMetricCSV(results []opendart.FinancialMetric) func(io.Writer) error {
	return func(out io.Writer) error {
		writer := csv.NewWriter(out)
		if err := writer.Write([]string{"metric", "amount", "currency", "account_id", "account_name", "match_method", "source_row_index"}); err != nil {
			return err
		}
		for _, metric := range results {
			if err := writer.Write([]string{string(metric.MetricCode), strconv.FormatInt(metric.Amount, 10), metric.Currency, metric.SourceAccountID, metric.SourceAccountName, string(metric.MatchMethod), strconv.Itoa(metric.SourceRowIndex)}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	}
}

func writeCompanyComparisonTable(result companyComparisonResult) func(io.Writer) error {
	return func(out io.Writer) error {
		return writeMetricSummaryTable(result.Companies)(out)
	}
}

func writeCompanyComparisonCSV(result companyComparisonResult) func(io.Writer) error {
	return func(out io.Writer) error {
		return writeMetricSummaryCSV(result.Companies)(out)
	}
}

func singleOrMany[T any](values []T) any {
	if len(values) == 1 {
		return values[0]
	}
	return values
}

func writeView(options *rootOptions, value any, writeTable func(io.Writer) error, writeCSV func(io.Writer) error) error {
	switch options.output {
	case outputJSON:
		return writeJSON(options.out, value)
	case outputTable:
		return writeTable(options.out)
	case outputCSV:
		return writeCSV(options.out)
	case outputRaw:
		return oops.In("opendart_cli").
			With("output", options.output).
			New("opendart cli: raw output is not supported for view commands")
	default:
		return options.validateOutput()
	}
}
