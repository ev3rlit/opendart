package opendart

import "context"

// DisclosureList is the paged response from the disclosure search API.
type DisclosureList struct {
	TotalCount string       `json:"total_count"`
	TotalPage  string       `json:"total_page"`
	PageNo     string       `json:"page_no"`
	PageCount  string       `json:"page_count"`
	Items      []Disclosure `json:"list"`
}

// Disclosure is a disclosure search result item.
type Disclosure struct {
	CorpClass   string `json:"corp_cls"`
	CorpName    string `json:"corp_name"`
	CorpCode    string `json:"corp_code"`
	StockCode   string `json:"stock_code"`
	ReportName  string `json:"report_nm"`
	ReceiptNo   string `json:"rcept_no"`
	FlrName     string `json:"flr_nm"`
	ReceiptDate string `json:"rcept_dt"`
	Remark      string `json:"rm"`
}

// CompanyProfile is the response from the company overview API.
type CompanyProfile struct {
	Status             string `json:"status"`
	Message            string `json:"message"`
	CorpCode           string `json:"corp_code"`
	CorpName           string `json:"corp_name"`
	CorpNameEnglish    string `json:"corp_name_eng"`
	StockName          string `json:"stock_name"`
	StockCode          string `json:"stock_code"`
	CEOName            string `json:"ceo_nm"`
	CorpClass          string `json:"corp_cls"`
	JurisdictionOffice string `json:"jurir_no"`
	BusinessNumber     string `json:"bizr_no"`
	Address            string `json:"adres"`
	Homepage           string `json:"hm_url"`
	IRHomepage         string `json:"ir_url"`
	Phone              string `json:"phn_no"`
	Fax                string `json:"fax_no"`
	IndustryCode       string `json:"induty_code"`
	EstablishDate      string `json:"est_dt"`
	AccountMonth       string `json:"acc_mt"`
}

// ReportFields is a lossless field map for OpenDART report rows.
type ReportFields map[string]string

// CapitalIncreaseDecreaseStatus is a 증자(감자) 현황 row.
type CapitalIncreaseDecreaseStatus ReportFields

// DividendMatter is a 배당에 관한 사항 row.
type DividendMatter ReportFields

// TreasuryStockAcquisitionDisposalStatus is a 자기주식 취득 및 처분 현황 row.
type TreasuryStockAcquisitionDisposalStatus ReportFields

// LargestShareholderStatus is a 최대주주 현황 row.
type LargestShareholderStatus ReportFields

// LargestShareholderChangeStatus is a 최대주주 변동현황 row.
type LargestShareholderChangeStatus ReportFields

// MinorityShareholderStatus is a 소액주주 현황 row.
type MinorityShareholderStatus ReportFields

// ExecutiveStatus is an 임원 현황 row.
type ExecutiveStatus ReportFields

// EmployeeStatus is a 직원 현황 row.
type EmployeeStatus ReportFields

// DirectorAuditorIndividualCompensationStatus is an 이사·감사 개인별 보수 row.
type DirectorAuditorIndividualCompensationStatus ReportFields

// DirectorAuditorTotalCompensationStatus is an 이사·감사 전체 보수 row.
type DirectorAuditorTotalCompensationStatus ReportFields

// IndividualCompensation is a 개인별 보수지급 금액 row.
type IndividualCompensation ReportFields

// OtherCorporationInvestmentStatus is a 타법인 출자현황 row.
type OtherCorporationInvestmentStatus ReportFields

// StockTotalQuantityStatus is a 주식의 총수 현황 row.
type StockTotalQuantityStatus ReportFields

// DebtSecuritiesIssuanceResult is a 채무증권 발행실적 row.
type DebtSecuritiesIssuanceResult ReportFields

// CommercialPaperOutstandingBalance is a 기업어음증권 미상환 잔액 row.
type CommercialPaperOutstandingBalance ReportFields

// ShortTermBondOutstandingBalance is a 단기사채 미상환 잔액 row.
type ShortTermBondOutstandingBalance ReportFields

// CorporateBondOutstandingBalance is a 회사채 미상환 잔액 row.
type CorporateBondOutstandingBalance ReportFields

// NewCapitalSecuritiesOutstandingBalance is a 신종자본증권 미상환 잔액 row.
type NewCapitalSecuritiesOutstandingBalance ReportFields

// ContingentCapitalSecuritiesOutstandingBalance is a 조건부 자본증권 미상환 잔액 row.
type ContingentCapitalSecuritiesOutstandingBalance ReportFields

// AuditorNameAndOpinion is a 회계감사인의 명칭 및 감사의견 row.
type AuditorNameAndOpinion ReportFields

// AuditServiceContractStatus is a 감사용역체결현황 row.
type AuditServiceContractStatus ReportFields

// NonAuditServiceContractStatus is a 비감사용역 계약체결 현황 row.
type NonAuditServiceContractStatus ReportFields

// OutsideDirectorChangeStatus is a 사외이사 및 그 변동현황 row.
type OutsideDirectorChangeStatus ReportFields

// UnregisteredExecutiveCompensationStatus is a 미등기임원 보수현황 row.
type UnregisteredExecutiveCompensationStatus ReportFields

// DirectorAuditorShareholderMeetingApprovedAmount is a 주주총회 승인금액 row.
type DirectorAuditorShareholderMeetingApprovedAmount ReportFields

// DirectorAuditorCompensationByType is a 보수지급금액 유형별 row.
type DirectorAuditorCompensationByType ReportFields

// PublicOfferingCapitalUseDetail is a 공모자금의 사용내역 row.
type PublicOfferingCapitalUseDetail ReportFields

// PrivatePlacementCapitalUseDetail is a 사모자금의 사용내역 row.
type PrivatePlacementCapitalUseDetail ReportFields

// FullFinancialStatement is an all-account financial statement row.
type FullFinancialStatement struct {
	ReceiptNo             string                     `json:"rcept_no"`
	ReportCode            ReportCode                 `json:"reprt_code"`
	BusinessYear          string                     `json:"bsns_year"`
	CorpCode              string                     `json:"corp_code"`
	StatementDiv          StatementDivision          `json:"sj_div"`
	StatementName         string                     `json:"sj_nm"`
	AccountID             string                     `json:"account_id"`
	AccountName           string                     `json:"account_nm"`
	AccountDetail         string                     `json:"account_detail"`
	CurrentTermName       string                     `json:"thstrm_nm"`
	CurrentTermAmount     string                     `json:"thstrm_amount"`
	CurrentTermAddAmount  string                     `json:"thstrm_add_amount"`
	PreviousTermName      string                     `json:"frmtrm_nm"`
	PreviousTermAmount    string                     `json:"frmtrm_amount"`
	PreviousTermQuarter   string                     `json:"frmtrm_q_nm"`
	PreviousTermQAmount   string                     `json:"frmtrm_q_amount"`
	PreviousTermAddAmount string                     `json:"frmtrm_add_amount"`
	FinancialStatementDiv FinancialStatementDivision `json:"fs_div"`
	FinancialStatement    string                     `json:"fs_nm"`
	Order                 string                     `json:"ord"`
	Currency              string                     `json:"currency"`
}

// XBRLTaxonomyItem is an XBRL taxonomy financial statement form row.
type XBRLTaxonomyItem struct {
	StatementDiv  StatementDivision `json:"sj_div"`
	AccountID     string            `json:"account_id"`
	AccountName   string            `json:"account_nm"`
	Language      string            `json:"lang"`
	DataType      string            `json:"data_tp"`
	Calculation   string            `json:"calculation"`
	Reference     string            `json:"ref"`
	Label         string            `json:"label"`
	ParentAccount string            `json:"parent_account_id"`
}

// FinancialIndexItem is a financial index row.
type FinancialIndexItem struct {
	BusinessYear   string     `json:"bsns_year"`
	CorpCode       string     `json:"corp_code"`
	StockCode      string     `json:"stock_code"`
	ReportCode     ReportCode `json:"reprt_code"`
	AccountName    string     `json:"account_nm"`
	CurrentTerm    string     `json:"thstrm"`
	PreviousTerm   string     `json:"frmtrm"`
	BeforePrevious string     `json:"bfefrmtrm"`
	IndexClassCode string     `json:"idx_cl_code"`
	IndexClassName string     `json:"idx_cl_nm"`
}

// MajorStockReport is a 대량보유 상황보고 row.
type MajorStockReport ReportFields

// ExecutiveStockReport is an 임원ㆍ주요주주 소유보고 row.
type ExecutiveStockReport ReportFields

// MaterialReportItem is a lossless row for material report APIs.
type MaterialReportItem ReportFields

// AssetTransferEtcPutbackOptionReport is a 자산양수도(기타), 풋백옵션 row.
type AssetTransferEtcPutbackOptionReport ReportFields

// DefaultOccurrenceReport is a 부도발생 row.
type DefaultOccurrenceReport ReportFields

// BusinessSuspensionReport is an 영업정지 row.
type BusinessSuspensionReport ReportFields

// RehabilitationProcedureApplicationReport is a 회생절차 개시신청 row.
type RehabilitationProcedureApplicationReport ReportFields

// DissolutionReasonOccurrenceReport is a 해산사유 발생 row.
type DissolutionReasonOccurrenceReport ReportFields

// PaidInCapitalIncreaseDecisionReport is a 유상증자 결정 row.
type PaidInCapitalIncreaseDecisionReport ReportFields

// BonusIssueDecisionReport is a 무상증자 결정 row.
type BonusIssueDecisionReport ReportFields

// PaidInBonusIssueDecisionReport is a 유무상증자 결정 row.
type PaidInBonusIssueDecisionReport ReportFields

// CapitalReductionDecisionReport is a 감자 결정 row.
type CapitalReductionDecisionReport ReportFields

// CreditorBankManagementProcedureStartReport is a 채권은행 등의 관리절차 개시 row.
type CreditorBankManagementProcedureStartReport ReportFields

// LawsuitFilingReport is a 소송 등의 제기 row.
type LawsuitFilingReport ReportFields

// OverseasListingDecisionReport is a 해외 증권시장 주권등 상장 결정 row.
type OverseasListingDecisionReport ReportFields

// OverseasDelistingDecisionReport is a 해외 증권시장 주권등 상장폐지 결정 row.
type OverseasDelistingDecisionReport ReportFields

// OverseasListingReport is a 해외 증권시장 주권등 상장 row.
type OverseasListingReport ReportFields

// OverseasDelistingReport is a 해외 증권시장 주권등 상장폐지 row.
type OverseasDelistingReport ReportFields

// ConvertibleBondIssueDecisionReport is a 전환사채권 발행결정 row.
type ConvertibleBondIssueDecisionReport ReportFields

// BondWithWarrantIssueDecisionReport is a 신주인수권부사채권 발행결정 row.
type BondWithWarrantIssueDecisionReport ReportFields

// ExchangeableBondIssueDecisionReport is a 교환사채권 발행결정 row.
type ExchangeableBondIssueDecisionReport ReportFields

// CreditorBankManagementProcedureStopReport is a 채권은행 등의 관리절차 중단 row.
type CreditorBankManagementProcedureStopReport ReportFields

// WriteDownContingentCapitalBondIssueDecisionReport is a 상각형 조건부자본증권 발행결정 row.
type WriteDownContingentCapitalBondIssueDecisionReport ReportFields

// TreasuryStockAcquisitionDecisionReport is a 자기주식 취득 결정 row.
type TreasuryStockAcquisitionDecisionReport ReportFields

// TreasuryStockDisposalDecisionReport is a 자기주식 처분 결정 row.
type TreasuryStockDisposalDecisionReport ReportFields

// TreasuryStockTrustContractConclusionDecisionReport is a 자기주식취득 신탁계약 체결 결정 row.
type TreasuryStockTrustContractConclusionDecisionReport ReportFields

// TreasuryStockTrustContractCancellationDecisionReport is a 자기주식취득 신탁계약 해지 결정 row.
type TreasuryStockTrustContractCancellationDecisionReport ReportFields

// BusinessTransferInDecisionReport is an 영업양수 결정 row.
type BusinessTransferInDecisionReport ReportFields

// BusinessTransferOutDecisionReport is an 영업양도 결정 row.
type BusinessTransferOutDecisionReport ReportFields

// TangibleAssetTransferInDecisionReport is a 유형자산 양수 결정 row.
type TangibleAssetTransferInDecisionReport ReportFields

// TangibleAssetTransferOutDecisionReport is a 유형자산 양도 결정 row.
type TangibleAssetTransferOutDecisionReport ReportFields

// OtherCorporationStockInvestmentSecurityTransferInDecisionReport is a 타법인 주식 및 출자증권 양수결정 row.
type OtherCorporationStockInvestmentSecurityTransferInDecisionReport ReportFields

// OtherCorporationStockInvestmentSecurityTransferOutDecisionReport is a 타법인 주식 및 출자증권 양도결정 row.
type OtherCorporationStockInvestmentSecurityTransferOutDecisionReport ReportFields

// StockRelatedBondTransferInDecisionReport is a 주권 관련 사채권 양수 결정 row.
type StockRelatedBondTransferInDecisionReport ReportFields

// StockRelatedBondTransferOutDecisionReport is a 주권 관련 사채권 양도 결정 row.
type StockRelatedBondTransferOutDecisionReport ReportFields

// CompanyMergerDecisionReport is a 회사합병 결정 row.
type CompanyMergerDecisionReport ReportFields

// CompanyDivisionDecisionReport is a 회사분할 결정 row.
type CompanyDivisionDecisionReport ReportFields

// CompanyDivisionMergerDecisionReport is a 회사분할합병 결정 row.
type CompanyDivisionMergerDecisionReport ReportFields

// StockExchangeTransferDecisionReport is a 주식교환·이전 결정 row.
type StockExchangeTransferDecisionReport ReportFields

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

type disclosureListResponse struct {
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	TotalCount string       `json:"total_count"`
	TotalPage  string       `json:"total_page"`
	PageNo     string       `json:"page_no"`
	PageCount  string       `json:"page_count"`
	List       []Disclosure `json:"list"`
}

// Disclosures searches disclosure reports.
func (client *Client) Disclosures(ctx context.Context, query DisclosureListQuery) (*DisclosureList, error) {
	var result disclosureListResponse
	if err := getJSON(ctx, client, "/api/list.json", disclosureListParams(query), "list.json", &result); err != nil {
		return nil, err
	}
	return &DisclosureList{
		TotalCount: result.TotalCount,
		TotalPage:  result.TotalPage,
		PageNo:     result.PageNo,
		PageCount:  result.PageCount,
		Items:      result.List,
	}, nil
}

// Company returns a company overview.
func (client *Client) Company(ctx context.Context, query CorpCodeQuery) (*CompanyProfile, error) {
	params, err := corpCodeParams(query)
	if err != nil {
		return nil, requiredQueryError("Company", err)
	}
	var result CompanyProfile
	if err := getJSON(ctx, client, "/api/company.json", params, "company.json", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Document returns the original disclosure document file.
func (client *Client) Document(ctx context.Context, query DocumentQuery) (*FileResponse, error) {
	params, err := documentParams(query)
	if err != nil {
		return nil, requiredQueryError("Document", err)
	}
	return getFile(ctx, client, "/api/document.xml", params)
}

func getPeriodic[T any](ctx context.Context, client *Client, method string, endpoint string, query PeriodicReportQuery) ([]T, error) {
	params, err := periodicReportParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, endpoint[5:])
}

func getCorpList[T any](ctx context.Context, client *Client, method string, endpoint string, query CorpCodeQuery) ([]T, error) {
	params, err := corpCodeParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, endpoint[5:])
}

func getMaterial[T any](ctx context.Context, client *Client, method string, endpoint string, query MaterialReportQuery) ([]T, error) {
	params, err := materialReportParams(query)
	if err != nil {
		return nil, requiredQueryError(method, err)
	}
	return getList[T](ctx, client, endpoint, params, endpoint[5:])
}

// CapitalIncreaseDecreaseStatus returns 증자(감자) 현황 rows.
func (client *Client) CapitalIncreaseDecreaseStatus(ctx context.Context, query PeriodicReportQuery) ([]CapitalIncreaseDecreaseStatus, error) {
	return getPeriodic[CapitalIncreaseDecreaseStatus](ctx, client, "CapitalIncreaseDecreaseStatus", "/api/irdsSttus.json", query)
}

// DividendMatters returns 배당에 관한 사항 rows.
func (client *Client) DividendMatters(ctx context.Context, query PeriodicReportQuery) ([]DividendMatter, error) {
	return getPeriodic[DividendMatter](ctx, client, "DividendMatters", "/api/alotMatter.json", query)
}

// TreasuryStockAcquisitionDisposalStatus returns 자기주식 취득 및 처분 현황 rows.
func (client *Client) TreasuryStockAcquisitionDisposalStatus(ctx context.Context, query PeriodicReportQuery) ([]TreasuryStockAcquisitionDisposalStatus, error) {
	return getPeriodic[TreasuryStockAcquisitionDisposalStatus](ctx, client, "TreasuryStockAcquisitionDisposalStatus", "/api/tesstkAcqsDspsSttus.json", query)
}

// LargestShareholderStatus returns 최대주주 현황 rows.
func (client *Client) LargestShareholderStatus(ctx context.Context, query PeriodicReportQuery) ([]LargestShareholderStatus, error) {
	return getPeriodic[LargestShareholderStatus](ctx, client, "LargestShareholderStatus", "/api/hyslrSttus.json", query)
}

// LargestShareholderChangeStatus returns 최대주주 변동현황 rows.
func (client *Client) LargestShareholderChangeStatus(ctx context.Context, query PeriodicReportQuery) ([]LargestShareholderChangeStatus, error) {
	return getPeriodic[LargestShareholderChangeStatus](ctx, client, "LargestShareholderChangeStatus", "/api/hyslrChgSttus.json", query)
}

// MinorityShareholderStatus returns 소액주주 현황 rows.
func (client *Client) MinorityShareholderStatus(ctx context.Context, query PeriodicReportQuery) ([]MinorityShareholderStatus, error) {
	return getPeriodic[MinorityShareholderStatus](ctx, client, "MinorityShareholderStatus", "/api/mrhlSttus.json", query)
}

// ExecutiveStatus returns 임원 현황 rows.
func (client *Client) ExecutiveStatus(ctx context.Context, query PeriodicReportQuery) ([]ExecutiveStatus, error) {
	return getPeriodic[ExecutiveStatus](ctx, client, "ExecutiveStatus", "/api/exctvSttus.json", query)
}

// EmployeeStatus returns 직원 현황 rows.
func (client *Client) EmployeeStatus(ctx context.Context, query PeriodicReportQuery) ([]EmployeeStatus, error) {
	return getPeriodic[EmployeeStatus](ctx, client, "EmployeeStatus", "/api/empSttus.json", query)
}

// DirectorAuditorIndividualCompensationStatus returns 이사·감사의 개인별 보수현황 rows.
func (client *Client) DirectorAuditorIndividualCompensationStatus(ctx context.Context, query PeriodicReportQuery) ([]DirectorAuditorIndividualCompensationStatus, error) {
	return getPeriodic[DirectorAuditorIndividualCompensationStatus](ctx, client, "DirectorAuditorIndividualCompensationStatus", "/api/hmvAuditIndvdlBySttus.json", query)
}

// DirectorAuditorTotalCompensationStatus returns 이사·감사 전체의 보수현황 rows.
func (client *Client) DirectorAuditorTotalCompensationStatus(ctx context.Context, query PeriodicReportQuery) ([]DirectorAuditorTotalCompensationStatus, error) {
	return getPeriodic[DirectorAuditorTotalCompensationStatus](ctx, client, "DirectorAuditorTotalCompensationStatus", "/api/hmvAuditAllSttus.json", query)
}

// IndividualCompensation returns 개인별 보수지급 금액 rows.
func (client *Client) IndividualCompensation(ctx context.Context, query PeriodicReportQuery) ([]IndividualCompensation, error) {
	return getPeriodic[IndividualCompensation](ctx, client, "IndividualCompensation", "/api/indvdlByPay.json", query)
}

// OtherCorporationInvestmentStatus returns 타법인 출자현황 rows.
func (client *Client) OtherCorporationInvestmentStatus(ctx context.Context, query PeriodicReportQuery) ([]OtherCorporationInvestmentStatus, error) {
	return getPeriodic[OtherCorporationInvestmentStatus](ctx, client, "OtherCorporationInvestmentStatus", "/api/otrCprInvstmntSttus.json", query)
}

// StockTotalQuantityStatus returns 주식의 총수 현황 rows.
func (client *Client) StockTotalQuantityStatus(ctx context.Context, query PeriodicReportQuery) ([]StockTotalQuantityStatus, error) {
	return getPeriodic[StockTotalQuantityStatus](ctx, client, "StockTotalQuantityStatus", "/api/stockTotqySttus.json", query)
}

// DebtSecuritiesIssuanceResults returns 채무증권 발행실적 rows.
func (client *Client) DebtSecuritiesIssuanceResults(ctx context.Context, query PeriodicReportQuery) ([]DebtSecuritiesIssuanceResult, error) {
	return getPeriodic[DebtSecuritiesIssuanceResult](ctx, client, "DebtSecuritiesIssuanceResults", "/api/detScritsIsuAcmslt.json", query)
}

// CommercialPaperOutstandingBalance returns 기업어음증권 미상환 잔액 rows.
func (client *Client) CommercialPaperOutstandingBalance(ctx context.Context, query PeriodicReportQuery) ([]CommercialPaperOutstandingBalance, error) {
	return getPeriodic[CommercialPaperOutstandingBalance](ctx, client, "CommercialPaperOutstandingBalance", "/api/entrprsBilScritsNrdmpBlce.json", query)
}

// ShortTermBondOutstandingBalance returns 단기사채 미상환 잔액 rows.
func (client *Client) ShortTermBondOutstandingBalance(ctx context.Context, query PeriodicReportQuery) ([]ShortTermBondOutstandingBalance, error) {
	return getPeriodic[ShortTermBondOutstandingBalance](ctx, client, "ShortTermBondOutstandingBalance", "/api/srtpdPsndbtNrdmpBlce.json", query)
}

// CorporateBondOutstandingBalance returns 회사채 미상환 잔액 rows.
func (client *Client) CorporateBondOutstandingBalance(ctx context.Context, query PeriodicReportQuery) ([]CorporateBondOutstandingBalance, error) {
	return getPeriodic[CorporateBondOutstandingBalance](ctx, client, "CorporateBondOutstandingBalance", "/api/cprndNrdmpBlce.json", query)
}

// NewCapitalSecuritiesOutstandingBalance returns 신종자본증권 미상환 잔액 rows.
func (client *Client) NewCapitalSecuritiesOutstandingBalance(ctx context.Context, query PeriodicReportQuery) ([]NewCapitalSecuritiesOutstandingBalance, error) {
	return getPeriodic[NewCapitalSecuritiesOutstandingBalance](ctx, client, "NewCapitalSecuritiesOutstandingBalance", "/api/newCaplScritsNrdmpBlce.json", query)
}

// ContingentCapitalSecuritiesOutstandingBalance returns 조건부 자본증권 미상환 잔액 rows.
func (client *Client) ContingentCapitalSecuritiesOutstandingBalance(ctx context.Context, query PeriodicReportQuery) ([]ContingentCapitalSecuritiesOutstandingBalance, error) {
	return getPeriodic[ContingentCapitalSecuritiesOutstandingBalance](ctx, client, "ContingentCapitalSecuritiesOutstandingBalance", "/api/cndlCaplScritsNrdmpBlce.json", query)
}

// AuditorNameAndOpinion returns 회계감사인의 명칭 및 감사의견 rows.
func (client *Client) AuditorNameAndOpinion(ctx context.Context, query PeriodicReportQuery) ([]AuditorNameAndOpinion, error) {
	return getPeriodic[AuditorNameAndOpinion](ctx, client, "AuditorNameAndOpinion", "/api/accnutAdtorNmNdAdtOpinion.json", query)
}

// AuditServiceContractStatus returns 감사용역체결현황 rows.
func (client *Client) AuditServiceContractStatus(ctx context.Context, query PeriodicReportQuery) ([]AuditServiceContractStatus, error) {
	return getPeriodic[AuditServiceContractStatus](ctx, client, "AuditServiceContractStatus", "/api/adtServcCnclsSttus.json", query)
}

// NonAuditServiceContractStatus returns 비감사용역 계약체결 현황 rows.
func (client *Client) NonAuditServiceContractStatus(ctx context.Context, query PeriodicReportQuery) ([]NonAuditServiceContractStatus, error) {
	return getPeriodic[NonAuditServiceContractStatus](ctx, client, "NonAuditServiceContractStatus", "/api/accnutAdtorNonAdtServcCnclsSttus.json", query)
}

// OutsideDirectorChangeStatus returns 사외이사 및 그 변동현황 rows.
func (client *Client) OutsideDirectorChangeStatus(ctx context.Context, query PeriodicReportQuery) ([]OutsideDirectorChangeStatus, error) {
	return getPeriodic[OutsideDirectorChangeStatus](ctx, client, "OutsideDirectorChangeStatus", "/api/outcmpnyDrctrNdChangeSttus.json", query)
}

// UnregisteredExecutiveCompensationStatus returns 미등기임원 보수현황 rows.
func (client *Client) UnregisteredExecutiveCompensationStatus(ctx context.Context, query PeriodicReportQuery) ([]UnregisteredExecutiveCompensationStatus, error) {
	return getPeriodic[UnregisteredExecutiveCompensationStatus](ctx, client, "UnregisteredExecutiveCompensationStatus", "/api/unrstExctvMendngSttus.json", query)
}

// DirectorAuditorShareholderMeetingApprovedAmount returns 주주총회 승인금액 rows.
func (client *Client) DirectorAuditorShareholderMeetingApprovedAmount(ctx context.Context, query PeriodicReportQuery) ([]DirectorAuditorShareholderMeetingApprovedAmount, error) {
	return getPeriodic[DirectorAuditorShareholderMeetingApprovedAmount](ctx, client, "DirectorAuditorShareholderMeetingApprovedAmount", "/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json", query)
}

// DirectorAuditorCompensationByType returns 보수지급금액 유형별 rows.
func (client *Client) DirectorAuditorCompensationByType(ctx context.Context, query PeriodicReportQuery) ([]DirectorAuditorCompensationByType, error) {
	return getPeriodic[DirectorAuditorCompensationByType](ctx, client, "DirectorAuditorCompensationByType", "/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json", query)
}

// PublicOfferingCapitalUseDetails returns 공모자금의 사용내역 rows.
func (client *Client) PublicOfferingCapitalUseDetails(ctx context.Context, query PeriodicReportQuery) ([]PublicOfferingCapitalUseDetail, error) {
	return getPeriodic[PublicOfferingCapitalUseDetail](ctx, client, "PublicOfferingCapitalUseDetails", "/api/pssrpCptalUseDtls.json", query)
}

// PrivatePlacementCapitalUseDetails returns 사모자금의 사용내역 rows.
func (client *Client) PrivatePlacementCapitalUseDetails(ctx context.Context, query PeriodicReportQuery) ([]PrivatePlacementCapitalUseDetail, error) {
	return getPeriodic[PrivatePlacementCapitalUseDetail](ctx, client, "PrivatePlacementCapitalUseDetails", "/api/prvsrpCptalUseDtls.json", query)
}

// MultiCompanyFinancialStatements returns major account statements for multiple companies.
func (client *Client) MultiCompanyFinancialStatements(ctx context.Context, query PeriodicReportQuery) ([]FinancialStatement, error) {
	return getPeriodic[FinancialStatement](ctx, client, "MultiCompanyFinancialStatements", "/api/fnlttMultiAcnt.json", query)
}

// FinancialStatementXBRL returns the XBRL original financial statement file.
func (client *Client) FinancialStatementXBRL(ctx context.Context, query ReceiptReportQuery) (*FileResponse, error) {
	params, err := receiptReportParams(query)
	if err != nil {
		return nil, requiredQueryError("FinancialStatementXBRL", err)
	}
	return getFile(ctx, client, "/api/fnlttXbrl.xml", params)
}

// FullFinancialStatement returns all accounts for one company.
func (client *Client) FullFinancialStatement(ctx context.Context, query FullFinancialStatementQuery) ([]FullFinancialStatement, error) {
	params, err := fullFinancialStatementParams(query)
	if err != nil {
		return nil, requiredQueryError("FullFinancialStatement", err)
	}
	return getList[FullFinancialStatement](ctx, client, "/api/fnlttSinglAcntAll.json", params, "fnlttSinglAcntAll.json")
}

// XBRLTaxonomy returns XBRL taxonomy financial statement forms.
func (client *Client) XBRLTaxonomy(ctx context.Context, query TaxonomyQuery) ([]XBRLTaxonomyItem, error) {
	params, err := taxonomyParams(query)
	if err != nil {
		return nil, requiredQueryError("XBRLTaxonomy", err)
	}
	return getList[XBRLTaxonomyItem](ctx, client, "/api/xbrlTaxonomy.json", params, "xbrlTaxonomy.json")
}

// FinancialIndex returns financial index rows for one company.
func (client *Client) FinancialIndex(ctx context.Context, query FinancialIndexQuery) ([]FinancialIndexItem, error) {
	params, err := financialIndexParams(query)
	if err != nil {
		return nil, requiredQueryError("FinancialIndex", err)
	}
	return getList[FinancialIndexItem](ctx, client, "/api/fnlttSinglIndx.json", params, "fnlttSinglIndx.json")
}

// CompanyFinancialIndex returns financial index rows for multiple companies.
func (client *Client) CompanyFinancialIndex(ctx context.Context, query FinancialIndexQuery) ([]FinancialIndexItem, error) {
	params, err := financialIndexParams(query)
	if err != nil {
		return nil, requiredQueryError("CompanyFinancialIndex", err)
	}
	return getList[FinancialIndexItem](ctx, client, "/api/fnlttCmpnyIndx.json", params, "fnlttCmpnyIndx.json")
}

// MajorStock returns 대량보유 상황보고 rows.
func (client *Client) MajorStock(ctx context.Context, query CorpCodeQuery) ([]MajorStockReport, error) {
	return getCorpList[MajorStockReport](ctx, client, "MajorStock", "/api/majorstock.json", query)
}

// ExecutiveStock returns 임원ㆍ주요주주 소유보고 rows.
func (client *Client) ExecutiveStock(ctx context.Context, query CorpCodeQuery) ([]ExecutiveStockReport, error) {
	return getCorpList[ExecutiveStockReport](ctx, client, "ExecutiveStock", "/api/elestock.json", query)
}

// AssetTransferEtcPutbackOption returns 자산양수도(기타), 풋백옵션 rows.
func (client *Client) AssetTransferEtcPutbackOption(ctx context.Context, query MaterialReportQuery) ([]AssetTransferEtcPutbackOptionReport, error) {
	return getMaterial[AssetTransferEtcPutbackOptionReport](ctx, client, "AssetTransferEtcPutbackOption", "/api/astInhtrfEtcPtbkOpt.json", query)
}

// DefaultOccurrence returns 부도발생 rows.
func (client *Client) DefaultOccurrence(ctx context.Context, query MaterialReportQuery) ([]DefaultOccurrenceReport, error) {
	return getMaterial[DefaultOccurrenceReport](ctx, client, "DefaultOccurrence", "/api/dfOcr.json", query)
}

// BusinessSuspension returns 영업정지 rows.
func (client *Client) BusinessSuspension(ctx context.Context, query MaterialReportQuery) ([]BusinessSuspensionReport, error) {
	return getMaterial[BusinessSuspensionReport](ctx, client, "BusinessSuspension", "/api/bsnSp.json", query)
}

// RehabilitationProcedureApplication returns 회생절차 개시신청 rows.
func (client *Client) RehabilitationProcedureApplication(ctx context.Context, query MaterialReportQuery) ([]RehabilitationProcedureApplicationReport, error) {
	return getMaterial[RehabilitationProcedureApplicationReport](ctx, client, "RehabilitationProcedureApplication", "/api/ctrcvsBgrq.json", query)
}

// DissolutionReasonOccurrence returns 해산사유 발생 rows.
func (client *Client) DissolutionReasonOccurrence(ctx context.Context, query MaterialReportQuery) ([]DissolutionReasonOccurrenceReport, error) {
	return getMaterial[DissolutionReasonOccurrenceReport](ctx, client, "DissolutionReasonOccurrence", "/api/dsRsOcr.json", query)
}

// PaidInCapitalIncreaseDecision returns 유상증자 결정 rows.
func (client *Client) PaidInCapitalIncreaseDecision(ctx context.Context, query MaterialReportQuery) ([]PaidInCapitalIncreaseDecisionReport, error) {
	return getMaterial[PaidInCapitalIncreaseDecisionReport](ctx, client, "PaidInCapitalIncreaseDecision", "/api/piicDecsn.json", query)
}

// BonusIssueDecision returns 무상증자 결정 rows.
func (client *Client) BonusIssueDecision(ctx context.Context, query MaterialReportQuery) ([]BonusIssueDecisionReport, error) {
	return getMaterial[BonusIssueDecisionReport](ctx, client, "BonusIssueDecision", "/api/fricDecsn.json", query)
}

// PaidInBonusIssueDecision returns 유무상증자 결정 rows.
func (client *Client) PaidInBonusIssueDecision(ctx context.Context, query MaterialReportQuery) ([]PaidInBonusIssueDecisionReport, error) {
	return getMaterial[PaidInBonusIssueDecisionReport](ctx, client, "PaidInBonusIssueDecision", "/api/pifricDecsn.json", query)
}

// CapitalReductionDecision returns 감자 결정 rows.
func (client *Client) CapitalReductionDecision(ctx context.Context, query MaterialReportQuery) ([]CapitalReductionDecisionReport, error) {
	return getMaterial[CapitalReductionDecisionReport](ctx, client, "CapitalReductionDecision", "/api/crDecsn.json", query)
}

// CreditorBankManagementProcedureStart returns 채권은행 등의 관리절차 개시 rows.
func (client *Client) CreditorBankManagementProcedureStart(ctx context.Context, query MaterialReportQuery) ([]CreditorBankManagementProcedureStartReport, error) {
	return getMaterial[CreditorBankManagementProcedureStartReport](ctx, client, "CreditorBankManagementProcedureStart", "/api/bnkMngtPcbg.json", query)
}

// LawsuitFiling returns 소송 등의 제기 rows.
func (client *Client) LawsuitFiling(ctx context.Context, query MaterialReportQuery) ([]LawsuitFilingReport, error) {
	return getMaterial[LawsuitFilingReport](ctx, client, "LawsuitFiling", "/api/lwstLg.json", query)
}

// OverseasListingDecision returns 해외 증권시장 주권등 상장 결정 rows.
func (client *Client) OverseasListingDecision(ctx context.Context, query MaterialReportQuery) ([]OverseasListingDecisionReport, error) {
	return getMaterial[OverseasListingDecisionReport](ctx, client, "OverseasListingDecision", "/api/ovLstDecsn.json", query)
}

// OverseasDelistingDecision returns 해외 증권시장 주권등 상장폐지 결정 rows.
func (client *Client) OverseasDelistingDecision(ctx context.Context, query MaterialReportQuery) ([]OverseasDelistingDecisionReport, error) {
	return getMaterial[OverseasDelistingDecisionReport](ctx, client, "OverseasDelistingDecision", "/api/ovDlstDecsn.json", query)
}

// OverseasListing returns 해외 증권시장 주권등 상장 rows.
func (client *Client) OverseasListing(ctx context.Context, query MaterialReportQuery) ([]OverseasListingReport, error) {
	return getMaterial[OverseasListingReport](ctx, client, "OverseasListing", "/api/ovLst.json", query)
}

// OverseasDelisting returns 해외 증권시장 주권등 상장폐지 rows.
func (client *Client) OverseasDelisting(ctx context.Context, query MaterialReportQuery) ([]OverseasDelistingReport, error) {
	return getMaterial[OverseasDelistingReport](ctx, client, "OverseasDelisting", "/api/ovDlst.json", query)
}

// ConvertibleBondIssueDecision returns 전환사채권 발행결정 rows.
func (client *Client) ConvertibleBondIssueDecision(ctx context.Context, query MaterialReportQuery) ([]ConvertibleBondIssueDecisionReport, error) {
	return getMaterial[ConvertibleBondIssueDecisionReport](ctx, client, "ConvertibleBondIssueDecision", "/api/cvbdIsDecsn.json", query)
}

// BondWithWarrantIssueDecision returns 신주인수권부사채권 발행결정 rows.
func (client *Client) BondWithWarrantIssueDecision(ctx context.Context, query MaterialReportQuery) ([]BondWithWarrantIssueDecisionReport, error) {
	return getMaterial[BondWithWarrantIssueDecisionReport](ctx, client, "BondWithWarrantIssueDecision", "/api/bdwtIsDecsn.json", query)
}

// ExchangeableBondIssueDecision returns 교환사채권 발행결정 rows.
func (client *Client) ExchangeableBondIssueDecision(ctx context.Context, query MaterialReportQuery) ([]ExchangeableBondIssueDecisionReport, error) {
	return getMaterial[ExchangeableBondIssueDecisionReport](ctx, client, "ExchangeableBondIssueDecision", "/api/exbdIsDecsn.json", query)
}

// CreditorBankManagementProcedureStop returns 채권은행 등의 관리절차 중단 rows.
func (client *Client) CreditorBankManagementProcedureStop(ctx context.Context, query MaterialReportQuery) ([]CreditorBankManagementProcedureStopReport, error) {
	return getMaterial[CreditorBankManagementProcedureStopReport](ctx, client, "CreditorBankManagementProcedureStop", "/api/bnkMngtPcsp.json", query)
}

// WriteDownContingentCapitalBondIssueDecision returns 상각형 조건부자본증권 발행결정 rows.
func (client *Client) WriteDownContingentCapitalBondIssueDecision(ctx context.Context, query MaterialReportQuery) ([]WriteDownContingentCapitalBondIssueDecisionReport, error) {
	return getMaterial[WriteDownContingentCapitalBondIssueDecisionReport](ctx, client, "WriteDownContingentCapitalBondIssueDecision", "/api/wdCocobdIsDecsn.json", query)
}

// TreasuryStockAcquisitionDecision returns 자기주식 취득 결정 rows.
func (client *Client) TreasuryStockAcquisitionDecision(ctx context.Context, query MaterialReportQuery) ([]TreasuryStockAcquisitionDecisionReport, error) {
	return getMaterial[TreasuryStockAcquisitionDecisionReport](ctx, client, "TreasuryStockAcquisitionDecision", "/api/tsstkAqDecsn.json", query)
}

// TreasuryStockDisposalDecision returns 자기주식 처분 결정 rows.
func (client *Client) TreasuryStockDisposalDecision(ctx context.Context, query MaterialReportQuery) ([]TreasuryStockDisposalDecisionReport, error) {
	return getMaterial[TreasuryStockDisposalDecisionReport](ctx, client, "TreasuryStockDisposalDecision", "/api/tsstkDpDecsn.json", query)
}

// TreasuryStockTrustContractConclusionDecision returns 자기주식취득 신탁계약 체결 결정 rows.
func (client *Client) TreasuryStockTrustContractConclusionDecision(ctx context.Context, query MaterialReportQuery) ([]TreasuryStockTrustContractConclusionDecisionReport, error) {
	return getMaterial[TreasuryStockTrustContractConclusionDecisionReport](ctx, client, "TreasuryStockTrustContractConclusionDecision", "/api/tsstkAqTrctrCnsDecsn.json", query)
}

// TreasuryStockTrustContractCancellationDecision returns 자기주식취득 신탁계약 해지 결정 rows.
func (client *Client) TreasuryStockTrustContractCancellationDecision(ctx context.Context, query MaterialReportQuery) ([]TreasuryStockTrustContractCancellationDecisionReport, error) {
	return getMaterial[TreasuryStockTrustContractCancellationDecisionReport](ctx, client, "TreasuryStockTrustContractCancellationDecision", "/api/tsstkAqTrctrCcDecsn.json", query)
}

// BusinessTransferInDecision returns 영업양수 결정 rows.
func (client *Client) BusinessTransferInDecision(ctx context.Context, query MaterialReportQuery) ([]BusinessTransferInDecisionReport, error) {
	return getMaterial[BusinessTransferInDecisionReport](ctx, client, "BusinessTransferInDecision", "/api/bsnInhDecsn.json", query)
}

// BusinessTransferOutDecision returns 영업양도 결정 rows.
func (client *Client) BusinessTransferOutDecision(ctx context.Context, query MaterialReportQuery) ([]BusinessTransferOutDecisionReport, error) {
	return getMaterial[BusinessTransferOutDecisionReport](ctx, client, "BusinessTransferOutDecision", "/api/bsnTrfDecsn.json", query)
}

// TangibleAssetTransferInDecision returns 유형자산 양수 결정 rows.
func (client *Client) TangibleAssetTransferInDecision(ctx context.Context, query MaterialReportQuery) ([]TangibleAssetTransferInDecisionReport, error) {
	return getMaterial[TangibleAssetTransferInDecisionReport](ctx, client, "TangibleAssetTransferInDecision", "/api/tgastInhDecsn.json", query)
}

// TangibleAssetTransferOutDecision returns 유형자산 양도 결정 rows.
func (client *Client) TangibleAssetTransferOutDecision(ctx context.Context, query MaterialReportQuery) ([]TangibleAssetTransferOutDecisionReport, error) {
	return getMaterial[TangibleAssetTransferOutDecisionReport](ctx, client, "TangibleAssetTransferOutDecision", "/api/tgastTrfDecsn.json", query)
}

// OtherCorporationStockInvestmentSecurityTransferInDecision returns 타법인 주식 및 출자증권 양수결정 rows.
func (client *Client) OtherCorporationStockInvestmentSecurityTransferInDecision(ctx context.Context, query MaterialReportQuery) ([]OtherCorporationStockInvestmentSecurityTransferInDecisionReport, error) {
	return getMaterial[OtherCorporationStockInvestmentSecurityTransferInDecisionReport](ctx, client, "OtherCorporationStockInvestmentSecurityTransferInDecision", "/api/otcprStkInvscrInhDecsn.json", query)
}

// OtherCorporationStockInvestmentSecurityTransferOutDecision returns 타법인 주식 및 출자증권 양도결정 rows.
func (client *Client) OtherCorporationStockInvestmentSecurityTransferOutDecision(ctx context.Context, query MaterialReportQuery) ([]OtherCorporationStockInvestmentSecurityTransferOutDecisionReport, error) {
	return getMaterial[OtherCorporationStockInvestmentSecurityTransferOutDecisionReport](ctx, client, "OtherCorporationStockInvestmentSecurityTransferOutDecision", "/api/otcprStkInvscrTrfDecsn.json", query)
}

// StockRelatedBondTransferInDecision returns 주권 관련 사채권 양수 결정 rows.
func (client *Client) StockRelatedBondTransferInDecision(ctx context.Context, query MaterialReportQuery) ([]StockRelatedBondTransferInDecisionReport, error) {
	return getMaterial[StockRelatedBondTransferInDecisionReport](ctx, client, "StockRelatedBondTransferInDecision", "/api/stkrtbdInhDecsn.json", query)
}

// StockRelatedBondTransferOutDecision returns 주권 관련 사채권 양도 결정 rows.
func (client *Client) StockRelatedBondTransferOutDecision(ctx context.Context, query MaterialReportQuery) ([]StockRelatedBondTransferOutDecisionReport, error) {
	return getMaterial[StockRelatedBondTransferOutDecisionReport](ctx, client, "StockRelatedBondTransferOutDecision", "/api/stkrtbdTrfDecsn.json", query)
}

// CompanyMergerDecision returns 회사합병 결정 rows.
func (client *Client) CompanyMergerDecision(ctx context.Context, query MaterialReportQuery) ([]CompanyMergerDecisionReport, error) {
	return getMaterial[CompanyMergerDecisionReport](ctx, client, "CompanyMergerDecision", "/api/cmpMgDecsn.json", query)
}

// CompanyDivisionDecision returns 회사분할 결정 rows.
func (client *Client) CompanyDivisionDecision(ctx context.Context, query MaterialReportQuery) ([]CompanyDivisionDecisionReport, error) {
	return getMaterial[CompanyDivisionDecisionReport](ctx, client, "CompanyDivisionDecision", "/api/cmpDvDecsn.json", query)
}

// CompanyDivisionMergerDecision returns 회사분할합병 결정 rows.
func (client *Client) CompanyDivisionMergerDecision(ctx context.Context, query MaterialReportQuery) ([]CompanyDivisionMergerDecisionReport, error) {
	return getMaterial[CompanyDivisionMergerDecisionReport](ctx, client, "CompanyDivisionMergerDecision", "/api/cmpDvmgDecsn.json", query)
}

// StockExchangeTransferDecision returns 주식교환·이전 결정 rows.
func (client *Client) StockExchangeTransferDecision(ctx context.Context, query MaterialReportQuery) ([]StockExchangeTransferDecisionReport, error) {
	return getMaterial[StockExchangeTransferDecisionReport](ctx, client, "StockExchangeTransferDecision", "/api/stkExtrDecsn.json", query)
}

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
	if err := getJSON(ctx, client, endpoint, params, endpoint[5:], &result); err != nil {
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
