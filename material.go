package opendart

import "context"

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
