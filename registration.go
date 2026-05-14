package opendart

import "context"

// RegistrationReport contains securities registration report groups and rows.
type RegistrationReport struct {
	Title string              `json:"title,omitempty"`
	Group []RegistrationGroup `json:"group,omitempty"`
	List  []RegistrationItem  `json:"list,omitempty"`
}

// RegistrationGroup contains a grouped securities registration report table.
type RegistrationGroup struct {
	Title string             `json:"title,omitempty"`
	List  []RegistrationItem `json:"list,omitempty"`
}

// RegistrationItem is a securities registration report row.
type RegistrationItem ReportFields

// EquitySecuritiesRegistrationReport is a 지분증권 신고서 주요정보 response.
type EquitySecuritiesRegistrationReport = RegistrationReport

// DebtSecuritiesRegistrationReport is a 채무증권 신고서 주요정보 response.
type DebtSecuritiesRegistrationReport = RegistrationReport

// DepositaryReceiptRegistrationReport is a 증권예탁증권 신고서 주요정보 response.
type DepositaryReceiptRegistrationReport = RegistrationReport

// MergerRegistrationReport is a 합병 신고서 주요정보 response.
type MergerRegistrationReport = RegistrationReport

// ShareExchangeTransferRegistrationReport is a 주식의포괄적교환·이전 신고서 주요정보 response.
type ShareExchangeTransferRegistrationReport = RegistrationReport

// DivisionRegistrationReport is a 분할 신고서 주요정보 response.
type DivisionRegistrationReport = RegistrationReport

func getRegistration(ctx context.Context, client *Client, method string, endpoint string, query MaterialReportQuery) (*RegistrationReport, error) {
	params, err := materialReportParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	var result struct {
		Status  string              `json:"status"`
		Message string              `json:"message"`
		Title   string              `json:"title"`
		Group   []RegistrationGroup `json:"group"`
		List    []RegistrationItem  `json:"list"`
	}
	if err := getJSON(ctx, client, endpoint, params, method, endpointOp(endpoint), &result); err != nil {
		return nil, err
	}
	return &RegistrationReport{Title: result.Title, Group: result.Group, List: result.List}, nil
}

// EquitySecuritiesRegistration returns 지분증권 신고서 주요정보.
func (client *Client) EquitySecuritiesRegistration(ctx context.Context, query MaterialReportQuery) (*EquitySecuritiesRegistrationReport, error) {
	return getRegistration(ctx, client, "EquitySecuritiesRegistration", "/api/estkRs.json", query)
}

// DebtSecuritiesRegistration returns 채무증권 신고서 주요정보.
func (client *Client) DebtSecuritiesRegistration(ctx context.Context, query MaterialReportQuery) (*DebtSecuritiesRegistrationReport, error) {
	return getRegistration(ctx, client, "DebtSecuritiesRegistration", "/api/bdRs.json", query)
}

// DepositaryReceiptRegistration returns 증권예탁증권 신고서 주요정보.
func (client *Client) DepositaryReceiptRegistration(ctx context.Context, query MaterialReportQuery) (*DepositaryReceiptRegistrationReport, error) {
	return getRegistration(ctx, client, "DepositaryReceiptRegistration", "/api/stkdpRs.json", query)
}

// MergerRegistration returns 합병 신고서 주요정보.
func (client *Client) MergerRegistration(ctx context.Context, query MaterialReportQuery) (*MergerRegistrationReport, error) {
	return getRegistration(ctx, client, "MergerRegistration", "/api/mgRs.json", query)
}

// ShareExchangeTransferRegistration returns 주식의포괄적교환·이전 신고서 주요정보.
func (client *Client) ShareExchangeTransferRegistration(ctx context.Context, query MaterialReportQuery) (*ShareExchangeTransferRegistrationReport, error) {
	return getRegistration(ctx, client, "ShareExchangeTransferRegistration", "/api/extrRs.json", query)
}

// DivisionRegistration returns 분할 신고서 주요정보.
func (client *Client) DivisionRegistration(ctx context.Context, query MaterialReportQuery) (*DivisionRegistrationReport, error) {
	return getRegistration(ctx, client, "DivisionRegistration", "/api/dvRs.json", query)
}
