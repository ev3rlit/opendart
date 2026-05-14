package opendart

import "context"

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
