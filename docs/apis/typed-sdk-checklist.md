# 과거 SDK friendly method 대응표

- 기준 문서: `docs/apis/official-inventory.md`
- 확인 날짜: 2026-05-14
- 범위: 공식 OpenDART API 85개 중 과거 수작업 friendly method 83개
- SDK package: `github.com/awuzag/opendart`
- 현재 root SDK는 OpenAPI `operationId` 기반 generated raw method를 기본으로 제공한다.
- 이 표의 friendly 이름은 `docs/apis/sdk-names.yaml`에 보존되어 있다.

## 요약

| Group | APIs | Friendly names |
| --- | ---: | ---: |
| disclosure | 4 | 4 |
| company | 30 | 28 |
| financial | 7 | 7 |
| ownership | 2 | 2 |
| material | 36 | 36 |
| registration | 6 | 6 |
| total | 85 | 83 |

## 대응표

| Group | API | Endpoint | SDK method |
| --- | --- | --- | --- |
| disclosure | 공시검색 | `/api/list.json` | `Client.Disclosures(ctx, query)` |
| disclosure | 기업개황 | `/api/company.json` | `Client.Company(ctx, query)` |
| disclosure | 공시서류원본파일 | `/api/document.xml` | `Client.Document(ctx, query)` |
| disclosure | 고유번호 | `/api/corpCode.xml` | `Client.CorpCodes(ctx)` |
| company | 증자(감자) 현황 | `/api/irdsSttus.json` | `Client.CapitalIncreaseDecreaseStatus(ctx, query)` |
| company | 배당에 관한 사항 | `/api/alotMatter.json` | `Client.DividendMatters(ctx, query)` |
| company | 자기주식 취득 및 처분 현황 | `/api/tesstkAcqsDspsSttus.json` | `Client.TreasuryStockAcquisitionDisposalStatus(ctx, query)` |
| company | 최대주주 현황 | `/api/hyslrSttus.json` | `Client.LargestShareholderStatus(ctx, query)` |
| company | 최대주주 변동현황 | `/api/hyslrChgSttus.json` | `Client.LargestShareholderChangeStatus(ctx, query)` |
| company | 소액주주 현황 | `/api/mrhlSttus.json` | `Client.MinorityShareholderStatus(ctx, query)` |
| company | 임원 현황 | `/api/exctvSttus.json` | `Client.ExecutiveStatus(ctx, query)` |
| company | 직원 현황 | `/api/empSttus.json` | `Client.EmployeeStatus(ctx, query)` |
| company | 이사·감사의 개인별 보수현황 | `/api/hmvAuditIndvdlBySttus.json` | `Client.DirectorAuditorIndividualCompensationStatus(ctx, query)` |
| company | 이사·감사 전체의 보수현황 | `/api/hmvAuditAllSttus.json` | `Client.DirectorAuditorTotalCompensationStatus(ctx, query)` |
| company | 개인별 보수지급 금액 | `/api/indvdlByPay.json` | `Client.IndividualCompensation(ctx, query)` |
| company | 타법인 출자현황 | `/api/otrCprInvstmntSttus.json` | `Client.OtherCorporationInvestmentStatus(ctx, query)` |
| company | 주식의 총수 현황 | `/api/stockTotqySttus.json` | `Client.StockTotalQuantityStatus(ctx, query)` |
| company | 채무증권 발행실적 | `/api/detScritsIsuAcmslt.json` | `Client.DebtSecuritiesIssuanceResults(ctx, query)` |
| company | 기업어음증권 미상환 잔액 | `/api/entrprsBilScritsNrdmpBlce.json` | `Client.CommercialPaperOutstandingBalance(ctx, query)` |
| company | 단기사채 미상환 잔액 | `/api/srtpdPsndbtNrdmpBlce.json` | `Client.ShortTermBondOutstandingBalance(ctx, query)` |
| company | 회사채 미상환 잔액 | `/api/cprndNrdmpBlce.json` | `Client.CorporateBondOutstandingBalance(ctx, query)` |
| company | 신종자본증권 미상환 잔액 | `/api/newCaplScritsNrdmpBlce.json` | `Client.NewCapitalSecuritiesOutstandingBalance(ctx, query)` |
| company | 조건부 자본증권 미상환 잔액 | `/api/cndlCaplScritsNrdmpBlce.json` | `Client.ContingentCapitalSecuritiesOutstandingBalance(ctx, query)` |
| company | 회계감사인의 명칭 및 감사의견 | `/api/accnutAdtorNmNdAdtOpinion.json` | `Client.AuditorNameAndOpinion(ctx, query)` |
| company | 감사용역체결현황 | `/api/adtServcCnclsSttus.json` | `Client.AuditServiceContractStatus(ctx, query)` |
| company | 회계감사인과의 비감사용역 계약체결 현황 | `/api/accnutAdtorNonAdtServcCnclsSttus.json` | `Client.NonAuditServiceContractStatus(ctx, query)` |
| company | 사외이사 및 그 변동현황 | `/api/outcmpnyDrctrNdChangeSttus.json` | `Client.OutsideDirectorChangeStatus(ctx, query)` |
| company | 미등기임원 보수현황 | `/api/unrstExctvMendngSttus.json` | `Client.UnregisteredExecutiveCompensationStatus(ctx, query)` |
| company | 이사·감사 전체의 보수현황(주주총회 승인금액) | `/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json` | `Client.DirectorAuditorShareholderMeetingApprovedAmount(ctx, query)` |
| company | 이사·감사 전체의 보수현황(보수지급금액 - 유형별) | `/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json` | `Client.DirectorAuditorCompensationByType(ctx, query)` |
| company | 공모자금의 사용내역 | `/api/pssrpCptalUseDtls.json` | `Client.PublicOfferingCapitalUseDetails(ctx, query)` |
| company | 사모자금의 사용내역 | `/api/prvsrpCptalUseDtls.json` | `Client.PrivatePlacementCapitalUseDetails(ctx, query)` |
| company | 이사·감사의 개인별 보수현황(5억원 이상) (Ver 2.0) | `/api/hmvAuditIndvdlBySttusV2.json` | 미구현 |
| company | 개인별 보수지급 금액(5억이상 상위5인) (Ver 2.0) | `/api/indvdlByPayV2.json` | 미구현 |
| financial | 단일회사 주요계정 | `/api/fnlttSinglAcnt.json` | `Client.FinancialStatement(ctx, query)` |
| financial | 다중회사 주요계정 | `/api/fnlttMultiAcnt.json` | `Client.MultiCompanyFinancialStatements(ctx, query)` |
| financial | 재무제표 원본파일(XBRL) | `/api/fnlttXbrl.xml` | `Client.FinancialStatementXBRL(ctx, query)` |
| financial | 단일회사 전체 재무제표 | `/api/fnlttSinglAcntAll.json` | `Client.FullFinancialStatement(ctx, query)` |
| financial | XBRL택사노미재무제표양식 | `/api/xbrlTaxonomy.json` | `Client.XBRLTaxonomy(ctx, query)` |
| financial | 단일회사 주요 재무지표 | `/api/fnlttSinglIndx.json` | `Client.FinancialIndex(ctx, query)` |
| financial | 다중회사 주요 재무지표 | `/api/fnlttCmpnyIndx.json` | `Client.CompanyFinancialIndex(ctx, query)` |
| ownership | 대량보유 상황보고 | `/api/majorstock.json` | `Client.MajorStock(ctx, query)` |
| ownership | 임원ㆍ주요주주 소유보고 | `/api/elestock.json` | `Client.ExecutiveStock(ctx, query)` |
| material | 자산양수도(기타), 풋백옵션 | `/api/astInhtrfEtcPtbkOpt.json` | `Client.AssetTransferEtcPutbackOption(ctx, query)` |
| material | 부도발생 | `/api/dfOcr.json` | `Client.DefaultOccurrence(ctx, query)` |
| material | 영업정지 | `/api/bsnSp.json` | `Client.BusinessSuspension(ctx, query)` |
| material | 회생절차 개시신청 | `/api/ctrcvsBgrq.json` | `Client.RehabilitationProcedureApplication(ctx, query)` |
| material | 해산사유 발생 | `/api/dsRsOcr.json` | `Client.DissolutionReasonOccurrence(ctx, query)` |
| material | 유상증자 결정 | `/api/piicDecsn.json` | `Client.PaidInCapitalIncreaseDecision(ctx, query)` |
| material | 무상증자 결정 | `/api/fricDecsn.json` | `Client.BonusIssueDecision(ctx, query)` |
| material | 유무상증자 결정 | `/api/pifricDecsn.json` | `Client.PaidInBonusIssueDecision(ctx, query)` |
| material | 감자 결정 | `/api/crDecsn.json` | `Client.CapitalReductionDecision(ctx, query)` |
| material | 채권은행 등의 관리절차 개시 | `/api/bnkMngtPcbg.json` | `Client.CreditorBankManagementProcedureStart(ctx, query)` |
| material | 소송 등의 제기 | `/api/lwstLg.json` | `Client.LawsuitFiling(ctx, query)` |
| material | 해외 증권시장 주권등 상장 결정 | `/api/ovLstDecsn.json` | `Client.OverseasListingDecision(ctx, query)` |
| material | 해외 증권시장 주권등 상장폐지 결정 | `/api/ovDlstDecsn.json` | `Client.OverseasDelistingDecision(ctx, query)` |
| material | 해외 증권시장 주권등 상장 | `/api/ovLst.json` | `Client.OverseasListing(ctx, query)` |
| material | 해외 증권시장 주권등 상장폐지 | `/api/ovDlst.json` | `Client.OverseasDelisting(ctx, query)` |
| material | 전환사채권 발행결정 | `/api/cvbdIsDecsn.json` | `Client.ConvertibleBondIssueDecision(ctx, query)` |
| material | 신주인수권부사채권 발행결정 | `/api/bdwtIsDecsn.json` | `Client.BondWithWarrantIssueDecision(ctx, query)` |
| material | 교환사채권 발행결정 | `/api/exbdIsDecsn.json` | `Client.ExchangeableBondIssueDecision(ctx, query)` |
| material | 채권은행 등의 관리절차 중단 | `/api/bnkMngtPcsp.json` | `Client.CreditorBankManagementProcedureStop(ctx, query)` |
| material | 상각형 조건부자본증권 발행결정 | `/api/wdCocobdIsDecsn.json` | `Client.WriteDownContingentCapitalBondIssueDecision(ctx, query)` |
| material | 자기주식 취득 결정 | `/api/tsstkAqDecsn.json` | `Client.TreasuryStockAcquisitionDecision(ctx, query)` |
| material | 자기주식 처분 결정 | `/api/tsstkDpDecsn.json` | `Client.TreasuryStockDisposalDecision(ctx, query)` |
| material | 자기주식취득 신탁계약 체결 결정 | `/api/tsstkAqTrctrCnsDecsn.json` | `Client.TreasuryStockTrustContractConclusionDecision(ctx, query)` |
| material | 자기주식취득 신탁계약 해지 결정 | `/api/tsstkAqTrctrCcDecsn.json` | `Client.TreasuryStockTrustContractCancellationDecision(ctx, query)` |
| material | 영업양수 결정 | `/api/bsnInhDecsn.json` | `Client.BusinessTransferInDecision(ctx, query)` |
| material | 영업양도 결정 | `/api/bsnTrfDecsn.json` | `Client.BusinessTransferOutDecision(ctx, query)` |
| material | 유형자산 양수 결정 | `/api/tgastInhDecsn.json` | `Client.TangibleAssetTransferInDecision(ctx, query)` |
| material | 유형자산 양도 결정 | `/api/tgastTrfDecsn.json` | `Client.TangibleAssetTransferOutDecision(ctx, query)` |
| material | 타법인 주식 및 출자증권 양수결정 | `/api/otcprStkInvscrInhDecsn.json` | `Client.OtherCorporationStockInvestmentSecurityTransferInDecision(ctx, query)` |
| material | 타법인 주식 및 출자증권 양도결정 | `/api/otcprStkInvscrTrfDecsn.json` | `Client.OtherCorporationStockInvestmentSecurityTransferOutDecision(ctx, query)` |
| material | 주권 관련 사채권 양수 결정 | `/api/stkrtbdInhDecsn.json` | `Client.StockRelatedBondTransferInDecision(ctx, query)` |
| material | 주권 관련 사채권 양도 결정 | `/api/stkrtbdTrfDecsn.json` | `Client.StockRelatedBondTransferOutDecision(ctx, query)` |
| material | 회사합병 결정 | `/api/cmpMgDecsn.json` | `Client.CompanyMergerDecision(ctx, query)` |
| material | 회사분할 결정 | `/api/cmpDvDecsn.json` | `Client.CompanyDivisionDecision(ctx, query)` |
| material | 회사분할합병 결정 | `/api/cmpDvmgDecsn.json` | `Client.CompanyDivisionMergerDecision(ctx, query)` |
| material | 주식교환·이전 결정 | `/api/stkExtrDecsn.json` | `Client.StockExchangeTransferDecision(ctx, query)` |
| registration | 지분증권 | `/api/estkRs.json` | `Client.EquitySecuritiesRegistration(ctx, query)` |
| registration | 채무증권 | `/api/bdRs.json` | `Client.DebtSecuritiesRegistration(ctx, query)` |
| registration | 증권예탁증권 | `/api/stkdpRs.json` | `Client.DepositaryReceiptRegistration(ctx, query)` |
| registration | 합병 | `/api/mgRs.json` | `Client.MergerRegistration(ctx, query)` |
| registration | 주식의포괄적교환·이전 | `/api/extrRs.json` | `Client.ShareExchangeTransferRegistration(ctx, query)` |
| registration | 분할 | `/api/dvRs.json` | `Client.DivisionRegistration(ctx, query)` |
