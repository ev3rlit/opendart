# 공식 API 인벤토리

- 확인 날짜: 2026-05-14
- 공식 개발가이드 API 수: 85개
- 공식 개발가이드 기준 그룹:
  - 공시정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS001
  - 정기보고서 주요정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS002
  - 정기보고서 재무정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS003
  - 지분공시 종합정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS004
  - 주요사항보고서 주요정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS005
  - 증권신고서 주요정보: https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS006

## 공통 요청/응답 규칙

- 모든 API는 `crtfc_key`를 요구한다. CLI에서는 `--api-key` 또는 `OPENDART_API_KEY`로 받으며, command flag로 직접 노출하지 않는다.
- JSON API의 공통 응답 필드는 `status`, `message`다. 목록형 API는 `list`를 포함한다.
- 공시검색은 paging 필드 `page_no`, `page_count`, `total_count`, `total_page`를 추가로 반환한다.
- 증권신고서 주요정보 API는 `group`, `title`, `list` 구조가 반복될 수 있다.
- 파일/XML API인 `document.xml`, `fnlttXbrl.xml`은 OpenDART 원문 bytes를 반환한다. CLI 기본 JSON 출력에서는 `api_id`, `api_name`, `endpoint`, `content_type`, `content_base64` envelope로 감싼다. `corpCode.xml`은 `opendart list corp-codes`에서 ZIP 내부 XML을 디코드해 `list[]` JSON으로 출력하고, `--output raw`일 때만 원문 bytes를 출력한다.
- business error는 공통 메시지 코드(`status`, `message`)를 따른다. 기본 테스트는 live 호출 없이 fake server로 검증한다.
- 공식 응답 필드 전체는 `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`에 OpenAPI 3.1 형태로 추출했다.
- CLI catalog는 `go generate ./internal/generated/opendartapi` 실행 시 OpenAPI split 문서와 `docs/apis/cli-names.yaml`에서 생성된다. raw API는 `search`, `get`, `list`, `download` verb 아래에 붙고, 가공 view도 top-level 동작 verb가 아니라 `get <business-resource>` 리소스로 노출한다.
- root package `opendart`는 공식 문서나 repo-local OpenAPI scrape에서 목록이 확인되는 반복 코드값을 public const로 제공한다.
- `idx_cl_code`, `pblntf_detail_ty`처럼 전체 목록 확인이 필요한 상세 코드는 이번 public const 범위에서 제외하고 후속 조사 대상으로 둔다.

## 파라미터 프로필

| Profile | 요청 파라미터 |
| --- | --- |
| `corp` | `corp_code` |
| `periodic` | `corp_code`, `bsns_year`, `reprt_code` |
| `material` | `corp_code`, `bgn_de`, `end_de` |
| `receipt-report` | `rcept_no`, `reprt_code` |
| `taxonomy` | `sj_div` |
| `financial-index` | `corp_code`, `bsns_year`, `reprt_code`, `idx_cl_code` |

## Verb-first CLI Command 인벤토리

| Group | API ID | API | Endpoint | CLI command | Params | Response fields |
| --- | --- | --- | --- | --- | --- | --- |
| disclosure | 2019001 | 공시검색 | `/api/list.json` | `opendart search disclosures` | optional `corp_code`, `bgn_de`, `end_de`, `last_reprt_at`, `pblntf_ty`, `pblntf_detail_ty`, `corp_cls`, `sort`, `sort_mth`, `page_no`, `page_count` | `status`, `message`, paging fields, `list[]` disclosure fields |
| disclosure | 2019002 | 기업개황 | `/api/company.json` | `opendart get company-profile` | `corp` | `status`, `message`, company profile fields |
| disclosure | 2019003 | 공시서류원본파일 | `/api/document.xml` | `opendart download document` | `rcept_no` | file bytes or base64 envelope |
| disclosure | 2019018 | 고유번호 | `/api/corpCode.xml` | `opendart list corp-codes` | none beyond auth | `list[]` corp code fields |
| company | 2019004 | 증자(감자) 현황 | `/api/irdsSttus.json` | `opendart get capital-change-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019005 | 배당에 관한 사항 | `/api/alotMatter.json` | `opendart get dividend-matters` | `periodic` | `status`, `message`, `list[]` |
| company | 2019006 | 자기주식 취득 및 처분 현황 | `/api/tesstkAcqsDspsSttus.json` | `opendart get treasury-stock-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019007 | 최대주주 현황 | `/api/hyslrSttus.json` | `opendart get largest-shareholder-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019008 | 최대주주 변동현황 | `/api/hyslrChgSttus.json` | `opendart get largest-shareholder-change-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019009 | 소액주주 현황 | `/api/mrhlSttus.json` | `opendart get minority-shareholder-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019010 | 임원 현황 | `/api/exctvSttus.json` | `opendart get executive-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019011 | 직원 현황 | `/api/empSttus.json` | `opendart get employee-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2019012 | 이사·감사의 개인별 보수현황(5억원 이상) | `/api/hmvAuditIndvdlBySttus.json` | `opendart get director-auditor-pay-individual` | `periodic` | `status`, `message`, `list[]` |
| company | 2019013 | 이사·감사 전체의 보수현황(보수지급금액 - 이사·감사 전체) | `/api/hmvAuditAllSttus.json` | `opendart get director-auditor-pay-total` | `periodic` | `status`, `message`, `list[]` |
| company | 2019014 | 개인별 보수지급 금액(5억이상 상위5인) | `/api/indvdlByPay.json` | `opendart get individual-pay-top5` | `periodic` | `status`, `message`, `list[]` |
| company | 2019015 | 타법인 출자현황 | `/api/otrCprInvstmntSttus.json` | `opendart get other-company-investments` | `periodic` | `status`, `message`, `list[]` |
| company | 2020002 | 주식의 총수 현황 | `/api/stockTotqySttus.json` | `opendart get stock-total-quantity` | `periodic` | `status`, `message`, `list[]` |
| company | 2020003 | 채무증권 발행실적 | `/api/detScritsIsuAcmslt.json` | `opendart get debt-securities-issued` | `periodic` | `status`, `message`, `list[]` |
| company | 2020004 | 기업어음증권 미상환 잔액 | `/api/entrprsBilScritsNrdmpBlce.json` | `opendart get commercial-paper-balance` | `periodic` | `status`, `message`, `list[]` |
| company | 2020005 | 단기사채 미상환 잔액 | `/api/srtpdPsndbtNrdmpBlce.json` | `opendart get short-term-bond-balance` | `periodic` | `status`, `message`, `list[]` |
| company | 2020006 | 회사채 미상환 잔액 | `/api/cprndNrdmpBlce.json` | `opendart get corporate-bond-balance` | `periodic` | `status`, `message`, `list[]` |
| company | 2020007 | 신종자본증권 미상환 잔액 | `/api/newCaplScritsNrdmpBlce.json` | `opendart get hybrid-capital-balance` | `periodic` | `status`, `message`, `list[]` |
| company | 2020008 | 조건부 자본증권 미상환 잔액 | `/api/cndlCaplScritsNrdmpBlce.json` | `opendart get contingent-capital-balance` | `periodic` | `status`, `message`, `list[]` |
| company | 2020009 | 회계감사인의 명칭 및 감사의견 | `/api/accnutAdtorNmNdAdtOpinion.json` | `opendart get auditor-opinion` | `periodic` | `status`, `message`, `list[]` |
| company | 2020010 | 감사용역체결현황 | `/api/adtServcCnclsSttus.json` | `opendart get audit-service-contracts` | `periodic` | `status`, `message`, `list[]` |
| company | 2020011 | 회계감사인과의 비감사용역 계약체결 현황 | `/api/accnutAdtorNonAdtServcCnclsSttus.json` | `opendart get non-audit-service-contracts` | `periodic` | `status`, `message`, `list[]` |
| company | 2020012 | 사외이사 및 그 변동현황 | `/api/outcmpnyDrctrNdChangeSttus.json` | `opendart get outside-director-change-status` | `periodic` | `status`, `message`, `list[]` |
| company | 2020013 | 미등기임원 보수현황 | `/api/unrstExctvMendngSttus.json` | `opendart get unregistered-executive-pay` | `periodic` | `status`, `message`, `list[]` |
| company | 2020014 | 이사·감사 전체의 보수현황(주주총회 승인금액) | `/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json` | `opendart get director-auditor-approved-pay` | `periodic` | `status`, `message`, `list[]` |
| company | 2020015 | 이사·감사 전체의 보수현황(보수지급금액 - 유형별) | `/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json` | `opendart get director-auditor-paid-by-type` | `periodic` | `status`, `message`, `list[]` |
| company | 2020016 | 공모자금의 사용내역 | `/api/pssrpCptalUseDtls.json` | `opendart get public-offering-fund-use` | `periodic` | `status`, `message`, `list[]` |
| company | 2020017 | 사모자금의 사용내역 | `/api/prvsrpCptalUseDtls.json` | `opendart get private-placement-fund-use` | `periodic` | `status`, `message`, `list[]` |
| company | 2026001 | 이사·감사의 개인별 보수현황(5억원 이상) (Ver 2.0) | `/api/hmvAuditIndvdlBySttusV2.json` | `opendart get director-auditor-pay-individual-v2` | `periodic` | `status`, `message`, `list[]` |
| company | 2026002 | 개인별 보수지급 금액(5억이상 상위5인) (Ver 2.0) | `/api/indvdlByPayV2.json` | `opendart get individual-pay-top5-v2` | `periodic` | `status`, `message`, `list[]` |
| financial | 2019016 | 단일회사 주요계정 | `/api/fnlttSinglAcnt.json` | `opendart get financial-statement` | `periodic` | `status`, `message`, `list[]` major account fields |
| financial | 2019017 | 다중회사 주요계정 | `/api/fnlttMultiAcnt.json` | `opendart get financial-statement-multi` | `periodic` | `status`, `message`, `list[]` major account fields |
| financial | 2019019 | 재무제표 원본파일(XBRL) | `/api/fnlttXbrl.xml` | `opendart download financial-xbrl` | `receipt-report` | file bytes or base64 envelope |
| financial | 2019020 | 단일회사 전체 재무제표 | `/api/fnlttSinglAcntAll.json` | `opendart get financial-statement-full` | `periodic` + `fs_div` | `status`, `message`, `list[]` account fields |
| financial | 2020001 | XBRL택사노미재무제표양식 | `/api/xbrlTaxonomy.json` | `opendart get xbrl-taxonomy` | `taxonomy` | `status`, `message`, `list[]` taxonomy fields |
| financial | 2022001 | 단일회사 주요 재무지표 | `/api/fnlttSinglIndx.json` | `opendart get financial-index` | `financial-index` | `status`, `message`, `list[]` index fields |
| financial | 2022002 | 다중회사 주요 재무지표 | `/api/fnlttCmpnyIndx.json` | `opendart get financial-index-multi` | `financial-index` | `status`, `message`, `list[]` index fields |
| ownership | 2019021 | 대량보유 상황보고 | `/api/majorstock.json` | `opendart get major-stock-ownership` | `corp` | `status`, `message`, `list[]` ownership fields |
| ownership | 2019022 | 임원ㆍ주요주주 소유보고 | `/api/elestock.json` | `opendart get executive-stock-ownership` | `corp` | `status`, `message`, `list[]` ownership fields |
| material | 2020018 | 자산양수도(기타), 풋백옵션 | `/api/astInhtrfEtcPtbkOpt.json` | `opendart get asset-transfer-putback-option` | `material` | `status`, `message`, `list[]` |
| material | 2020019 | 부도발생 | `/api/dfOcr.json` | `opendart get default-occurrence` | `material` | `status`, `message`, `list[]` |
| material | 2020020 | 영업정지 | `/api/bsnSp.json` | `opendart get business-suspension` | `material` | `status`, `message`, `list[]` |
| material | 2020021 | 회생절차 개시신청 | `/api/ctrcvsBgrq.json` | `opendart get rehabilitation-application` | `material` | `status`, `message`, `list[]` |
| material | 2020022 | 해산사유 발생 | `/api/dsRsOcr.json` | `opendart get dissolution-cause` | `material` | `status`, `message`, `list[]` |
| material | 2020023 | 유상증자 결정 | `/api/piicDecsn.json` | `opendart get paid-in-capital-increase` | `material` | `status`, `message`, `list[]` |
| material | 2020024 | 무상증자 결정 | `/api/fricDecsn.json` | `opendart get free-capital-increase` | `material` | `status`, `message`, `list[]` |
| material | 2020025 | 유무상증자 결정 | `/api/pifricDecsn.json` | `opendart get paid-and-free-capital-increase` | `material` | `status`, `message`, `list[]` |
| material | 2020026 | 감자 결정 | `/api/crDecsn.json` | `opendart get capital-reduction` | `material` | `status`, `message`, `list[]` |
| material | 2020027 | 채권은행 등의 관리절차 개시 | `/api/bnkMngtPcbg.json` | `opendart get bank-management-start` | `material` | `status`, `message`, `list[]` |
| material | 2020028 | 소송 등의 제기 | `/api/lwstLg.json` | `opendart get lawsuit` | `material` | `status`, `message`, `list[]` |
| material | 2020029 | 해외 증권시장 주권등 상장 결정 | `/api/ovLstDecsn.json` | `opendart get overseas-listing-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020030 | 해외 증권시장 주권등 상장폐지 결정 | `/api/ovDlstDecsn.json` | `opendart get overseas-delisting-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020031 | 해외 증권시장 주권등 상장 | `/api/ovLst.json` | `opendart get overseas-listing` | `material` | `status`, `message`, `list[]` |
| material | 2020032 | 해외 증권시장 주권등 상장폐지 | `/api/ovDlst.json` | `opendart get overseas-delisting` | `material` | `status`, `message`, `list[]` |
| material | 2020033 | 전환사채권 발행결정 | `/api/cvbdIsDecsn.json` | `opendart get convertible-bond-issuance` | `material` | `status`, `message`, `list[]` |
| material | 2020034 | 신주인수권부사채권 발행결정 | `/api/bdwtIsDecsn.json` | `opendart get bond-with-warrant-issuance` | `material` | `status`, `message`, `list[]` |
| material | 2020035 | 교환사채권 발행결정 | `/api/exbdIsDecsn.json` | `opendart get exchangeable-bond-issuance` | `material` | `status`, `message`, `list[]` |
| material | 2020036 | 채권은행 등의 관리절차 중단 | `/api/bnkMngtPcsp.json` | `opendart get bank-management-stop` | `material` | `status`, `message`, `list[]` |
| material | 2020037 | 상각형 조건부자본증권 발행결정 | `/api/wdCocobdIsDecsn.json` | `opendart get write-down-coco-bond-issuance` | `material` | `status`, `message`, `list[]` |
| material | 2020038 | 자기주식 취득 결정 | `/api/tsstkAqDecsn.json` | `opendart get treasury-stock-acquisition-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020039 | 자기주식 처분 결정 | `/api/tsstkDpDecsn.json` | `opendart get treasury-stock-disposal-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020040 | 자기주식취득 신탁계약 체결 결정 | `/api/tsstkAqTrctrCnsDecsn.json` | `opendart get treasury-stock-trust-contract` | `material` | `status`, `message`, `list[]` |
| material | 2020041 | 자기주식취득 신탁계약 해지 결정 | `/api/tsstkAqTrctrCcDecsn.json` | `opendart get treasury-stock-trust-cancellation` | `material` | `status`, `message`, `list[]` |
| material | 2020042 | 영업양수 결정 | `/api/bsnInhDecsn.json` | `opendart get business-acquisition` | `material` | `status`, `message`, `list[]` |
| material | 2020043 | 영업양도 결정 | `/api/bsnTrfDecsn.json` | `opendart get business-transfer` | `material` | `status`, `message`, `list[]` |
| material | 2020044 | 유형자산 양수 결정 | `/api/tgastInhDecsn.json` | `opendart get tangible-asset-acquisition` | `material` | `status`, `message`, `list[]` |
| material | 2020045 | 유형자산 양도 결정 | `/api/tgastTrfDecsn.json` | `opendart get tangible-asset-transfer` | `material` | `status`, `message`, `list[]` |
| material | 2020046 | 타법인 주식 및 출자증권 양수결정 | `/api/otcprStkInvscrInhDecsn.json` | `opendart get other-company-stock-acquisition` | `material` | `status`, `message`, `list[]` |
| material | 2020047 | 타법인 주식 및 출자증권 양도결정 | `/api/otcprStkInvscrTrfDecsn.json` | `opendart get other-company-stock-transfer` | `material` | `status`, `message`, `list[]` |
| material | 2020048 | 주권 관련 사채권 양수 결정 | `/api/stkrtbdInhDecsn.json` | `opendart get stock-related-bond-acquisition` | `material` | `status`, `message`, `list[]` |
| material | 2020049 | 주권 관련 사채권 양도 결정 | `/api/stkrtbdTrfDecsn.json` | `opendart get stock-related-bond-transfer` | `material` | `status`, `message`, `list[]` |
| material | 2020050 | 회사합병 결정 | `/api/cmpMgDecsn.json` | `opendart get merger-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020051 | 회사분할 결정 | `/api/cmpDvDecsn.json` | `opendart get split-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020052 | 회사분할합병 결정 | `/api/cmpDvmgDecsn.json` | `opendart get split-merger-decision` | `material` | `status`, `message`, `list[]` |
| material | 2020053 | 주식교환·이전 결정 | `/api/stkExtrDecsn.json` | `opendart get stock-exchange-transfer-decision` | `material` | `status`, `message`, `list[]` |
| registration | 2020054 | 지분증권 | `/api/estkRs.json` | `opendart get registration-equity` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020055 | 채무증권 | `/api/bdRs.json` | `opendart get registration-debt` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020056 | 증권예탁증권 | `/api/stkdpRs.json` | `opendart get registration-depositary-receipt` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020057 | 합병 | `/api/mgRs.json` | `opendart get registration-merger` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020058 | 주식의포괄적교환·이전 | `/api/extrRs.json` | `opendart get registration-share-exchange-transfer` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020059 | 분할 | `/api/dvRs.json` | `opendart get registration-division` | `material` | `status`, `message`, `group[]`, `list[]` |

## 비즈니스 view 리소스

비즈니스 view는 OpenDART raw API 하나를 그대로 노출하는 command가 아니라, `fnlttSinglAcntAll` 원문 row를 정규화해 사용자가 자주 보는 재무 리소스로 재구성한 CLI 표면이다. 따라서 `summarize`, `compare`, `inspect` 같은 동작 verb를 public top-level command로 쓰지 않고 `opendart get <resource>` 아래에 둔다.

공통 flag:

- `--corp-code`: 단일 OpenDART 고유번호.
- `--corp-codes`: 여러 OpenDART 고유번호를 쉼표로 구분한 값. 여러 회사를 비교할 때 별도 `compare` verb 대신 사용한다.
- `--year`: 사업연도.
- `--fs-div`: 재무제표 구분. 예: `CFS`, `OFS`.
- `--view`: 표현 깊이. `summary`, `detail`, `source`.
- `--output`: 출력 형식. view 리소스는 `json`, `table`, `csv`를 지원한다. 파일/XML API의 `--output raw` 계약은 `download` 계열에 유지한다.

| CLI command | Source API | Resource 의미 | 주요 옵션 |
| --- | --- | --- | --- |
| `opendart get quarter-performance` | `fnlttSinglAcntAll` | 1Q~4Q 손익 주요 지표. 4Q는 사업보고서 누적 `11011`에서 3분기 누적 `11014`를 빼서 계산한다. | `--corp-code` 또는 `--corp-codes`, `--year`, `--fs-div`, `--view`, `--output` |
| `opendart get annual-performance` | `fnlttSinglAcntAll` | 연간 매출, 매출총이익, 영업이익, 순이익. | `--corp-code` 또는 `--corp-codes`, `--year`, `--fs-div`, `--view`, `--output` |
| `opendart get financial-position` | `fnlttSinglAcntAll` | 자산, 유동/비유동자산, 부채, 유동/비유동부채, 자본. | `--corp-code` 또는 `--corp-codes`, `--year`, `--fs-div`, `--view`, `--output` |
| `opendart get cash-flow` | `fnlttSinglAcntAll` | 현금및현금성자산, 영업활동현금흐름. | `--corp-code` 또는 `--corp-codes`, `--year`, `--fs-div`, `--view`, `--output` |
| `opendart get financial-metric` | `fnlttSinglAcntAll` | `--metric revenue` 같은 단일 정규화 지표. `--view source`는 원천 row를 함께 출력한다. | `--corp-code` 또는 `--corp-codes`, `--year`, `--fs-div`, `--metric`, `--view`, `--output` |

## 추적 기준

- `Group`, `API ID`, `Endpoint`, `Params`는 공식 상세 페이지에서 확인했다.
- `CLI command`는 생성 파일인 `internal/cli/catalog_gen.go`와 대응한다. 사람이 조정하는 verb/resource 이름은 `docs/apis/cli-names.yaml`에만 둔다.
- 비즈니스 view 리소스는 공식 API catalog 생성물이 아니라 `internal/cli/views.go`의 수동 CLI 레이어이며, root SDK public package와 raw API command 계약을 깨지 않는다.
- SDK typed method 대응표는 `docs/apis/typed-sdk-checklist.md`에서 별도로 추적한다.
- 전체 요청/응답 필드 스키마는 `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`, 원본 표 덤프는 `docs/apis/opendart-api-metadata.json`에서 확인한다.
- `corp_code`는 OpenDART 고유번호이며 KRX 종목코드가 아니다.
