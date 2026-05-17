package opendart

import (
	"strings"

	"github.com/samber/oops"
)

// FnlttSinglAcntAllAmount is a parsed amount value with its original OpenDART text.
type FnlttSinglAcntAllAmount struct {
	Raw   string `json:"raw"`
	Value int64  `json:"value"`
}

// FnlttSinglAcntAllAnalysisRow is the analysis schema for one single-account-all raw row.
//
// A 2024 annual consolidated audit over 188 successfully fetched major KOSPI
// companies found the key string fields below present in every raw row. Amount
// and prior-period fields are nullable because OpenDART legitimately leaves
// them blank for some rows.
type FnlttSinglAcntAllAnalysisRow struct {
	SourceRowIndex int `json:"source_row_index"`

	RceptNo       string `json:"rcept_no"`
	ReprtCode     string `json:"reprt_code"`
	BusinessYear  string `json:"bsns_year"`
	CorpCode      string `json:"corp_code"`
	StatementDiv  string `json:"sj_div"`
	StatementName string `json:"sj_nm"`
	AccountID     string `json:"account_id"`
	AccountName   string `json:"account_nm"`
	AccountDetail string `json:"account_detail"`
	CurrentName   string `json:"thstrm_nm"`
	Order         string `json:"ord"`
	Currency      string `json:"currency"`

	CurrentAmount            *FnlttSinglAcntAllAmount `json:"thstrm_amount,omitempty"`
	CurrentCumulativeAmount  *FnlttSinglAcntAllAmount `json:"thstrm_add_amount,omitempty"`
	PreviousName             *string                  `json:"frmtrm_nm,omitempty"`
	PreviousAmount           *FnlttSinglAcntAllAmount `json:"frmtrm_amount,omitempty"`
	PreviousQuarterName      *string                  `json:"frmtrm_q_nm,omitempty"`
	PreviousQuarterAmount    *FnlttSinglAcntAllAmount `json:"frmtrm_q_amount,omitempty"`
	PreviousCumulativeAmount *FnlttSinglAcntAllAmount `json:"frmtrm_add_amount,omitempty"`
	BeforePreviousName       *string                  `json:"bfefrmtrm_nm,omitempty"`
	BeforePreviousAmount     *FnlttSinglAcntAllAmount `json:"bfefrmtrm_amount,omitempty"`

	SourceRow FnlttSinglAcntAllItem `json:"source_row"`
}

// FnlttSinglAcntAllAnalysisIssue reports a row that does not fit the analysis schema cleanly.
type FnlttSinglAcntAllAnalysisIssue struct {
	SourceRowIndex int                   `json:"source_row_index"`
	Field          string                `json:"field"`
	Reason         string                `json:"reason"`
	SourceRow      FnlttSinglAcntAllItem `json:"source_row"`
}

// FnlttSinglAcntAllAnalysisSet is a lossless analysis view over single-account-all rows.
type FnlttSinglAcntAllAnalysisSet struct {
	Rows   []FnlttSinglAcntAllAnalysisRow   `json:"rows"`
	Issues []FnlttSinglAcntAllAnalysisIssue `json:"issues,omitempty"`
}

// AnalyzeFnlttSinglAcntAllRows converts raw single-account-all rows into the analysis schema.
func AnalyzeFnlttSinglAcntAllRows(rows []FnlttSinglAcntAllItem) *FnlttSinglAcntAllAnalysisSet {
	result := &FnlttSinglAcntAllAnalysisSet{
		Rows: make([]FnlttSinglAcntAllAnalysisRow, 0, len(rows)),
	}
	for index, row := range rows {
		analysisRow := FnlttSinglAcntAllAnalysisRow{
			SourceRowIndex:      index,
			RceptNo:             row.RceptNo,
			ReprtCode:           row.ReprtCode,
			BusinessYear:        row.BsnsYear,
			CorpCode:            row.CorpCode,
			StatementDiv:        row.SjDiv,
			StatementName:       row.SjNm,
			AccountID:           row.AccountId,
			AccountName:         row.AccountNm,
			AccountDetail:       row.AccountDetail,
			CurrentName:         row.ThstrmNm,
			Order:               row.Ord,
			Currency:            row.Currency,
			PreviousName:        optionalFnlttSinglAcntAllString(row.FrmtrmNm),
			PreviousQuarterName: optionalFnlttSinglAcntAllString(row.FrmtrmQNm),
			BeforePreviousName:  optionalFnlttSinglAcntAllString(row.BfefrmtrmNm),
			SourceRow:           row,
		}

		result.addMissingKeyIssues(index, row)
		analysisRow.CurrentAmount = result.optionalAmount(index, row, "thstrm_amount", row.ThstrmAmount)
		analysisRow.CurrentCumulativeAmount = result.optionalAmount(index, row, "thstrm_add_amount", row.ThstrmAddAmount)
		analysisRow.PreviousAmount = result.optionalAmount(index, row, "frmtrm_amount", row.FrmtrmAmount)
		analysisRow.PreviousQuarterAmount = result.optionalAmount(index, row, "frmtrm_q_amount", row.FrmtrmQAmount)
		analysisRow.PreviousCumulativeAmount = result.optionalAmount(index, row, "frmtrm_add_amount", row.FrmtrmAddAmount)
		analysisRow.BeforePreviousAmount = result.optionalAmount(index, row, "bfefrmtrm_amount", row.BfefrmtrmAmount)

		result.Rows = append(result.Rows, analysisRow)
	}
	return result
}

// AnalyzeFnlttSinglAcntAllResponse converts a raw single-account-all response into the analysis schema.
func AnalyzeFnlttSinglAcntAllResponse(response *FnlttSinglAcntAllResponse) (*FnlttSinglAcntAllAnalysisSet, error) {
	if response == nil {
		return nil, oops.In("financial_statement_schema").
			New("opendart: nil FnlttSinglAcntAllResponse")
	}
	return AnalyzeFnlttSinglAcntAllRows(response.List), nil
}

func (set *FnlttSinglAcntAllAnalysisSet) addMissingKeyIssues(index int, row FnlttSinglAcntAllItem) {
	for _, field := range requiredFnlttSinglAcntAllAnalysisFields(row) {
		set.Issues = append(set.Issues, FnlttSinglAcntAllAnalysisIssue{
			SourceRowIndex: index,
			Field:          field,
			Reason:         "missing required analysis key field",
			SourceRow:      row,
		})
	}
}

func (set *FnlttSinglAcntAllAnalysisSet) optionalAmount(index int, row FnlttSinglAcntAllItem, field string, value string) *FnlttSinglAcntAllAmount {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parsed, err := parseFinancialMetricAmount(value)
	if err != nil {
		set.Issues = append(set.Issues, FnlttSinglAcntAllAnalysisIssue{
			SourceRowIndex: index,
			Field:          field,
			Reason:         "invalid amount: " + err.Error(),
			SourceRow:      row,
		})
		return nil
	}
	return &FnlttSinglAcntAllAmount{
		Raw:   value,
		Value: parsed,
	}
}

func requiredFnlttSinglAcntAllAnalysisFields(row FnlttSinglAcntAllItem) []string {
	fields := make([]string, 0)
	if strings.TrimSpace(row.AccountId) == "" {
		fields = append(fields, "account_id")
	}
	if strings.TrimSpace(row.AccountNm) == "" {
		fields = append(fields, "account_nm")
	}
	if strings.TrimSpace(row.AccountDetail) == "" {
		fields = append(fields, "account_detail")
	}
	if strings.TrimSpace(row.CorpCode) == "" {
		fields = append(fields, "corp_code")
	}
	if strings.TrimSpace(row.BsnsYear) == "" {
		fields = append(fields, "bsns_year")
	}
	if strings.TrimSpace(row.ReprtCode) == "" {
		fields = append(fields, "reprt_code")
	}
	if strings.TrimSpace(row.SjDiv) == "" {
		fields = append(fields, "sj_div")
	}
	if strings.TrimSpace(row.SjNm) == "" {
		fields = append(fields, "sj_nm")
	}
	if strings.TrimSpace(row.Currency) == "" {
		fields = append(fields, "currency")
	}
	if strings.TrimSpace(row.Ord) == "" {
		fields = append(fields, "ord")
	}
	if strings.TrimSpace(row.RceptNo) == "" {
		fields = append(fields, "rcept_no")
	}
	if strings.TrimSpace(row.ThstrmNm) == "" {
		fields = append(fields, "thstrm_nm")
	}
	return fields
}

func optionalFnlttSinglAcntAllString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}
