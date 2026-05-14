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
- 파일/XML API인 `document.xml`, `corpCode.xml`, `fnlttXbrl.xml`은 OpenDART 원문 bytes를 반환한다. CLI 기본 JSON 출력에서는 `api_id`, `api_name`, `endpoint`, `content_type`, `content_base64` envelope로 감싼다.
- business error는 공통 메시지 코드(`status`, `message`)를 따른다. 기본 테스트는 live 호출 없이 fake server로 검증한다.
- 공식 응답 필드 전체는 `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`에 OpenAPI 3.1 형태로 추출했다.

## 파라미터 프로필

| Profile | 요청 파라미터 |
| --- | --- |
| `corp` | `corp_code` |
| `periodic` | `corp_code`, `bsns_year`, `reprt_code` |
| `material` | `corp_code`, `bgn_de`, `end_de` |
| `receipt-report` | `rcept_no`, `reprt_code` |
| `taxonomy` | `sj_div` |
| `financial-index` | `corp_code`, `bsns_year`, `reprt_code`, `idx_cl_code` |

## CLI Command 인벤토리

| Group | API ID | API | Endpoint | CLI command | Params | Response fields |
| --- | --- | --- | --- | --- | --- | --- |
| disclosure | 2019001 | 공시검색 | `/api/list.json` | `opendart disclosure list` | optional `corp_code`, `bgn_de`, `end_de`, `last_reprt_at`, `pblntf_ty`, `pblntf_detail_ty`, `corp_cls`, `sort`, `sort_mth`, `page_no`, `page_count` | `status`, `message`, paging fields, `list[]` disclosure fields |
| disclosure | 2019002 | 기업개황 | `/api/company.json` | `opendart disclosure company` | `corp` | `status`, `message`, company profile fields |
| disclosure | 2019003 | 공시서류원본파일 | `/api/document.xml` | `opendart disclosure document` | `rcept_no` | file bytes or base64 envelope |
| disclosure | 2019018 | 고유번호 | `/api/corpCode.xml` | `opendart disclosure corp-code`, `opendart corp-codes` | none beyond auth | `list[]` corp code fields |
| company | 2019004 | 증자(감자) 현황 | `/api/irdsSttus.json` | `opendart company irds-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019005 | 배당에 관한 사항 | `/api/alotMatter.json` | `opendart company alot-matter` | `periodic` | `status`, `message`, `list[]` |
| company | 2019006 | 자기주식 취득 및 처분 현황 | `/api/tesstkAcqsDspsSttus.json` | `opendart company tesstk-acqs-dsps-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019007 | 최대주주 현황 | `/api/hyslrSttus.json` | `opendart company hyslr-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019008 | 최대주주 변동현황 | `/api/hyslrChgSttus.json` | `opendart company hyslr-chg-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019009 | 소액주주 현황 | `/api/mrhlSttus.json` | `opendart company mrhl-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019010 | 임원 현황 | `/api/exctvSttus.json` | `opendart company exctv-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019011 | 직원 현황 | `/api/empSttus.json` | `opendart company emp-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019012 | 이사·감사의 개인별 보수현황(5억원 이상) | `/api/hmvAuditIndvdlBySttus.json` | `opendart company hmv-audit-indvdl-by-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019013 | 이사·감사 전체의 보수현황(보수지급금액 - 이사·감사 전체) | `/api/hmvAuditAllSttus.json` | `opendart company hmv-audit-all-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2019014 | 개인별 보수지급 금액(5억이상 상위5인) | `/api/indvdlByPay.json` | `opendart company indvdl-by-pay` | `periodic` | `status`, `message`, `list[]` |
| company | 2019015 | 타법인 출자현황 | `/api/otrCprInvstmntSttus.json` | `opendart company otr-cpr-invstmnt-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020002 | 주식의 총수 현황 | `/api/stockTotqySttus.json` | `opendart company stock-totqy-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020003 | 채무증권 발행실적 | `/api/detScritsIsuAcmslt.json` | `opendart company det-scrits-isu-acmslt` | `periodic` | `status`, `message`, `list[]` |
| company | 2020004 | 기업어음증권 미상환 잔액 | `/api/entrprsBilScritsNrdmpBlce.json` | `opendart company entrprs-bil-scrits-nrdmp-blce` | `periodic` | `status`, `message`, `list[]` |
| company | 2020005 | 단기사채 미상환 잔액 | `/api/srtpdPsndbtNrdmpBlce.json` | `opendart company srtpd-psndbt-nrdmp-blce` | `periodic` | `status`, `message`, `list[]` |
| company | 2020006 | 회사채 미상환 잔액 | `/api/cprndNrdmpBlce.json` | `opendart company cprnd-nrdmp-blce` | `periodic` | `status`, `message`, `list[]` |
| company | 2020007 | 신종자본증권 미상환 잔액 | `/api/newCaplScritsNrdmpBlce.json` | `opendart company new-capl-scrits-nrdmp-blce` | `periodic` | `status`, `message`, `list[]` |
| company | 2020008 | 조건부 자본증권 미상환 잔액 | `/api/cndlCaplScritsNrdmpBlce.json` | `opendart company cndl-capl-scrits-nrdmp-blce` | `periodic` | `status`, `message`, `list[]` |
| company | 2020009 | 회계감사인의 명칭 및 감사의견 | `/api/accnutAdtorNmNdAdtOpinion.json` | `opendart company accnut-adtor-nm-nd-adt-opinion` | `periodic` | `status`, `message`, `list[]` |
| company | 2020010 | 감사용역체결현황 | `/api/adtServcCnclsSttus.json` | `opendart company adt-servc-cncls-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020011 | 회계감사인과의 비감사용역 계약체결 현황 | `/api/accnutAdtorNonAdtServcCnclsSttus.json` | `opendart company accnut-adtor-non-adt-servc-cncls-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020012 | 사외이사 및 그 변동현황 | `/api/outcmpnyDrctrNdChangeSttus.json` | `opendart company outcmpny-drctr-nd-change-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020013 | 미등기임원 보수현황 | `/api/unrstExctvMendngSttus.json` | `opendart company unrst-exctv-mendng-sttus` | `periodic` | `status`, `message`, `list[]` |
| company | 2020014 | 이사·감사 전체의 보수현황(주주총회 승인금액) | `/api/drctrAdtAllMendngSttusGmtsckConfmAmount.json` | `opendart company drctr-adt-all-mendng-sttus-gmtsck-confm-amount` | `periodic` | `status`, `message`, `list[]` |
| company | 2020015 | 이사·감사 전체의 보수현황(보수지급금액 - 유형별) | `/api/drctrAdtAllMendngSttusMendngPymntamtTyCl.json` | `opendart company drctr-adt-all-mendng-sttus-mendng-pymntamt-ty-cl` | `periodic` | `status`, `message`, `list[]` |
| company | 2020016 | 공모자금의 사용내역 | `/api/pssrpCptalUseDtls.json` | `opendart company pssrp-cptal-use-dtls` | `periodic` | `status`, `message`, `list[]` |
| company | 2020017 | 사모자금의 사용내역 | `/api/prvsrpCptalUseDtls.json` | `opendart company prvsrp-cptal-use-dtls` | `periodic` | `status`, `message`, `list[]` |
| company | 2026001 | 이사·감사의 개인별 보수현황(5억원 이상) (Ver 2.0) | `/api/hmvAuditIndvdlBySttusV2.json` | 미구현 | `periodic` | `status`, `message`, `list[]` |
| company | 2026002 | 개인별 보수지급 금액(5억이상 상위5인) (Ver 2.0) | `/api/indvdlByPayV2.json` | 미구현 | `periodic` | `status`, `message`, `list[]` |
| financial | 2019016 | 단일회사 주요계정 | `/api/fnlttSinglAcnt.json` | `opendart financial single-account`, `opendart financial-statement` | `periodic` | `status`, `message`, `list[]` major account fields |
| financial | 2019017 | 다중회사 주요계정 | `/api/fnlttMultiAcnt.json` | `opendart financial multi-account` | `periodic` | `status`, `message`, `list[]` major account fields |
| financial | 2019019 | 재무제표 원본파일(XBRL) | `/api/fnlttXbrl.xml` | `opendart financial xbrl` | `receipt-report` | file bytes or base64 envelope |
| financial | 2019020 | 단일회사 전체 재무제표 | `/api/fnlttSinglAcntAll.json` | `opendart financial single-account-all` | `periodic` + `fs_div` | `status`, `message`, `list[]` account fields |
| financial | 2020001 | XBRL택사노미재무제표양식 | `/api/xbrlTaxonomy.json` | `opendart financial xbrl-taxonomy` | `taxonomy` | `status`, `message`, `list[]` taxonomy fields |
| financial | 2022001 | 단일회사 주요 재무지표 | `/api/fnlttSinglIndx.json` | `opendart financial single-index` | `financial-index` | `status`, `message`, `list[]` index fields |
| financial | 2022002 | 다중회사 주요 재무지표 | `/api/fnlttCmpnyIndx.json` | `opendart financial company-index` | `financial-index` | `status`, `message`, `list[]` index fields |
| ownership | 2019021 | 대량보유 상황보고 | `/api/majorstock.json` | `opendart ownership major-stock` | `corp` | `status`, `message`, `list[]` ownership fields |
| ownership | 2019022 | 임원ㆍ주요주주 소유보고 | `/api/elestock.json` | `opendart ownership executive-stock` | `corp` | `status`, `message`, `list[]` ownership fields |
| material | 2020018 | 자산양수도(기타), 풋백옵션 | `/api/astInhtrfEtcPtbkOpt.json` | `opendart material ast-inhtrf-etc-ptbk-opt` | `material` | `status`, `message`, `list[]` |
| material | 2020019 | 부도발생 | `/api/dfOcr.json` | `opendart material df-ocr` | `material` | `status`, `message`, `list[]` |
| material | 2020020 | 영업정지 | `/api/bsnSp.json` | `opendart material bsn-sp` | `material` | `status`, `message`, `list[]` |
| material | 2020021 | 회생절차 개시신청 | `/api/ctrcvsBgrq.json` | `opendart material ctrcvs-bgrq` | `material` | `status`, `message`, `list[]` |
| material | 2020022 | 해산사유 발생 | `/api/dsRsOcr.json` | `opendart material ds-rs-ocr` | `material` | `status`, `message`, `list[]` |
| material | 2020023 | 유상증자 결정 | `/api/piicDecsn.json` | `opendart material piic-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020024 | 무상증자 결정 | `/api/fricDecsn.json` | `opendart material fric-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020025 | 유무상증자 결정 | `/api/pifricDecsn.json` | `opendart material pifric-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020026 | 감자 결정 | `/api/crDecsn.json` | `opendart material cr-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020027 | 채권은행 등의 관리절차 개시 | `/api/bnkMngtPcbg.json` | `opendart material bnk-mngt-pcbg` | `material` | `status`, `message`, `list[]` |
| material | 2020028 | 소송 등의 제기 | `/api/lwstLg.json` | `opendart material lwst-lg` | `material` | `status`, `message`, `list[]` |
| material | 2020029 | 해외 증권시장 주권등 상장 결정 | `/api/ovLstDecsn.json` | `opendart material ov-lst-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020030 | 해외 증권시장 주권등 상장폐지 결정 | `/api/ovDlstDecsn.json` | `opendart material ov-dlst-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020031 | 해외 증권시장 주권등 상장 | `/api/ovLst.json` | `opendart material ov-lst` | `material` | `status`, `message`, `list[]` |
| material | 2020032 | 해외 증권시장 주권등 상장폐지 | `/api/ovDlst.json` | `opendart material ov-dlst` | `material` | `status`, `message`, `list[]` |
| material | 2020033 | 전환사채권 발행결정 | `/api/cvbdIsDecsn.json` | `opendart material cvbd-is-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020034 | 신주인수권부사채권 발행결정 | `/api/bdwtIsDecsn.json` | `opendart material bdwt-is-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020035 | 교환사채권 발행결정 | `/api/exbdIsDecsn.json` | `opendart material exbd-is-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020036 | 채권은행 등의 관리절차 중단 | `/api/bnkMngtPcsp.json` | `opendart material bnk-mngt-pcsp` | `material` | `status`, `message`, `list[]` |
| material | 2020037 | 상각형 조건부자본증권 발행결정 | `/api/wdCocobdIsDecsn.json` | `opendart material wd-cocobd-is-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020038 | 자기주식 취득 결정 | `/api/tsstkAqDecsn.json` | `opendart material tsstk-aq-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020039 | 자기주식 처분 결정 | `/api/tsstkDpDecsn.json` | `opendart material tsstk-dp-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020040 | 자기주식취득 신탁계약 체결 결정 | `/api/tsstkAqTrctrCnsDecsn.json` | `opendart material tsstk-aq-trctr-cns-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020041 | 자기주식취득 신탁계약 해지 결정 | `/api/tsstkAqTrctrCcDecsn.json` | `opendart material tsstk-aq-trctr-cc-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020042 | 영업양수 결정 | `/api/bsnInhDecsn.json` | `opendart material bsn-inh-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020043 | 영업양도 결정 | `/api/bsnTrfDecsn.json` | `opendart material bsn-trf-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020044 | 유형자산 양수 결정 | `/api/tgastInhDecsn.json` | `opendart material tgast-inh-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020045 | 유형자산 양도 결정 | `/api/tgastTrfDecsn.json` | `opendart material tgast-trf-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020046 | 타법인 주식 및 출자증권 양수결정 | `/api/otcprStkInvscrInhDecsn.json` | `opendart material otcpr-stk-invscr-inh-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020047 | 타법인 주식 및 출자증권 양도결정 | `/api/otcprStkInvscrTrfDecsn.json` | `opendart material otcpr-stk-invscr-trf-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020048 | 주권 관련 사채권 양수 결정 | `/api/stkrtbdInhDecsn.json` | `opendart material stkrtbd-inh-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020049 | 주권 관련 사채권 양도 결정 | `/api/stkrtbdTrfDecsn.json` | `opendart material stkrtbd-trf-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020050 | 회사합병 결정 | `/api/cmpMgDecsn.json` | `opendart material cmp-mg-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020051 | 회사분할 결정 | `/api/cmpDvDecsn.json` | `opendart material cmp-dv-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020052 | 회사분할합병 결정 | `/api/cmpDvmgDecsn.json` | `opendart material cmp-dvmg-decsn` | `material` | `status`, `message`, `list[]` |
| material | 2020053 | 주식교환·이전 결정 | `/api/stkExtrDecsn.json` | `opendart material stk-extr-decsn` | `material` | `status`, `message`, `list[]` |
| registration | 2020054 | 지분증권 | `/api/estkRs.json` | `opendart registration equity` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020055 | 채무증권 | `/api/bdRs.json` | `opendart registration debt` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020056 | 증권예탁증권 | `/api/stkdpRs.json` | `opendart registration depositary-receipt` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020057 | 합병 | `/api/mgRs.json` | `opendart registration merger` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020058 | 주식의포괄적교환·이전 | `/api/extrRs.json` | `opendart registration share-exchange-transfer` | `material` | `status`, `message`, `group[]`, `list[]` |
| registration | 2020059 | 분할 | `/api/dvRs.json` | `opendart registration division` | `material` | `status`, `message`, `group[]`, `list[]` |

## 추적 기준

- `Group`, `API ID`, `Endpoint`, `Params`는 공식 상세 페이지에서 확인했다.
- `CLI command`는 `internal/cli/catalog.go`와 대응한다. `미구현` 항목은 공식 개발가이드에 새로 추가되었지만 현재 CLI catalog에는 아직 없다.
- SDK typed method 대응표는 `docs/apis/typed-sdk-checklist.md`에서 별도로 추적한다.
- 전체 요청/응답 필드 스키마는 `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`, 원본 표 덤프는 `docs/apis/opendart-api-metadata.json`에서 확인한다.
- `corp_code`는 OpenDART 고유번호이며 KRX 종목코드가 아니다.
