# 단일회사 전체 재무제표

- 공식 문서: https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS003&apiId=2019020
- 그룹: 정기보고서 재무정보
- SDK raw method: `Client.FnlttSinglAcntAll(ctx, params)`
- Friendly name history: `Client.FullFinancialStatement(ctx, query)` is preserved in `docs/apis/sdk-names.yaml`.

## 기본 정보

| 항목 | 값 |
| --- | --- |
| Method | `GET` |
| JSON URL | `https://opendart.fss.or.kr/api/fnlttSinglAcntAll.json` |
| XML URL | `https://opendart.fss.or.kr/api/fnlttSinglAcntAll.xml` |
| Encoding | UTF-8 |

## 요청 인자

| key | 명칭 | type | 필수 | 설명 |
| --- | --- | --- | --- | --- |
| `crtfc_key` | API 인증키 | `STRING(40)` | Y | 발급받은 인증키 |
| `corp_code` | 고유번호 | `STRING(8)` | Y | 공시대상회사의 고유번호 |
| `bsns_year` | 사업연도 | `STRING(4)` | Y | 2015년 이후 제공 |
| `reprt_code` | 보고서 코드 | `STRING(5)` | Y | `11013`, `11012`, `11014`, `11011` |
| `fs_div` | 개별/연결구분 | `STRING(3)` | Y | `OFS` 또는 `CFS` |

## 코드 값

| code | 의미 |
| --- | --- |
| `OFS` | 재무제표, 별도 |
| `CFS` | 연결재무제표 |
| `BS` | 재무상태표 |
| `IS` | 손익계산서 |
| `CIS` | 포괄손익계산서 |
| `CF` | 현금흐름표 |
| `SCE` | 자본변동표 |

## SDK 코드 상수

root package `opendart`는 공식 코드값을 public const로 제공한다.

```go
statements, err := client.FnlttSinglAcntAll(ctx, opendart.FnlttSinglAcntAllParams{
	CorpCode:  "00126380",
	BsnsYear:  "2025",
	ReprtCode: opendart.ReportCodeAnnual,
	FsDiv:     opendart.FinancialStatementDivisionConsolidated,
})
```

## 응답 필드

| key | 설명 |
| --- | --- |
| `status` | 에러 및 정보 코드 |
| `message` | 에러 및 정보 메시지 |
| `rcept_no` | 접수번호 |
| `reprt_code` | 보고서 코드 |
| `bsns_year` | 사업연도 |
| `corp_code` | 공시대상회사 고유번호 |
| `sj_div` | 재무제표구분 |
| `sj_nm` | 재무제표명 |
| `account_id` | XBRL 표준계정ID |
| `account_nm` | 계정명 |
| `account_detail` | 계정상세 |
| `thstrm_nm` | 당기명 |
| `thstrm_amount` | 당기금액 |
| `thstrm_add_amount` | 당기누적금액 |
| `frmtrm_nm` | 전기명 |
| `frmtrm_amount` | 전기금액 |
| `frmtrm_q_nm` | 전기명(분/반기) |
| `frmtrm_q_amount` | 전기금액(분/반기) |
| `frmtrm_add_amount` | 전기누적금액 |
| `bfefrmtrm_nm` | 전전기명 |
| `bfefrmtrm_amount` | 전전기금액 |
| `ord` | 계정과목 정렬순서 |
| `currency` | 통화 단위 |

## Normalized financial metrics

`single-account-all`은 회사별 계정 구조 차이가 크기 때문에 원문 row를 고정 스키마로 강제하지 않는다. SDK는 raw response를 그대로 보존하고, 공통 분석 지표가 필요할 때 별도 normalized layer를 선택적으로 사용한다.

```go
metrics, err := opendart.NormalizeFnlttSinglAcntAllResponse(raw, opendart.FnlttSinglAcntAllParams{
	CorpCode:  "00126380",
	BsnsYear:  "2024",
	ReprtCode: opendart.ReportCodeAnnual,
	FsDiv:     opendart.FinancialStatementDivisionConsolidated,
})
```

정규화 결과의 각 metric은 다음 추적 정보를 가진다.

| field | 설명 |
| --- | --- |
| `metric_code` | `revenue`, `operating_income`, `net_income`, `assets`, `liabilities`, `equity` 같은 공통 지표 코드 |
| `amount`, `currency` | 정규화된 당기 금액과 통화 |
| `corp_code`, `business_year`, `report_code`, `fs_div`, `statement_div` | 호출/row 기준 식별 정보 |
| `source_account_id`, `source_account_name`, `source_account_detail` | 원문 계정 추적 정보 |
| `source_row_index`, `source_row` | 원문 row로 돌아가기 위한 참조와 row 사본 |
| `match_method`, `confidence` | `account_id_exact`, `account_name_alias`, `override` 구분과 매핑 신뢰도 |

매핑 우선순위는 다음과 같다.

1. `WithFinancialMetricOverrideRules`로 넘긴 수동 override rule
2. 기본/추가 rule의 `account_id` exact match
3. 기본/추가 rule의 `account_nm` alias match

`-표준계정코드 미사용-`, 회사 고유 계정, 업종 특화 계정은 정규화 실패로 취급하지 않는다. 기본 rule에 맞지 않는 row는 `UnmappedRows`에 남고, 수동 승격이 필요하면 override rule로 별도 metric을 만들 수 있다.

## Analysis row schema

원문 SDK 응답 타입은 생성 코드 그대로 유지한다. 분석 코드에서 nullability를 더 명확히 다루고 싶으면 `AnalyzeFnlttSinglAcntAllResponse` 또는 `AnalyzeFnlttSinglAcntAllRows`를 사용한다.

```go
analysis, err := opendart.AnalyzeFnlttSinglAcntAllResponse(raw)
if err != nil {
	log.Fatal(err)
}
for _, row := range analysis.Rows {
	if row.CurrentAmount != nil {
		log.Println(row.CorpCode, row.AccountID, row.CurrentAmount.Value)
	}
}
```

2024년 사업보고서(`11011`), 연결(`CFS`) 기준 주요 200개 기업 audit에서 성공 조회된 188개 회사의 row 분포를 기준으로 field 정책을 나눴다.

| 정책 | fields | 근거 |
| --- | --- | --- |
| required string | `account_id`, `account_nm`, `account_detail`, `corp_code`, `bsns_year`, `reprt_code`, `sj_div`, `sj_nm`, `currency`, `ord`, `rcept_no`, `thstrm_nm` | 188개 회사, 46,625개 row에서 100% 채워짐 |
| nullable amount | `thstrm_amount`, `thstrm_add_amount`, `frmtrm_amount`, `frmtrm_q_amount`, `frmtrm_add_amount`, `bfefrmtrm_amount` | 금액 row 일부가 빈 값이며, 빈 값은 parse 실패가 아니라 `nil`로 다룸 |
| nullable period name | `frmtrm_nm`, `frmtrm_q_nm`, `bfefrmtrm_nm` | 전기/전전기명이 일부 회사 또는 일부 row에서 비어 있음 |

분석 스키마는 금액이 비어 있으면 `nil`로 두고, 값이 있는데 숫자로 파싱할 수 없을 때만 `Issues`에 기록한다. `SourceRow`는 원문 row 사본을 그대로 보존하므로 generated raw 응답이나 CLI JSON 출력 계약과 분리되어 있다.

## Live audit

실제 기업 표본에서 어떤 raw field와 metric이 공통으로 나타나는지는 `e2e` build tag 테스트로 측정한다. KOSPI 주요 200개 기업처럼 큰 표본은 target 파일로 넘기고, 기본 테스트 경로에는 포함하지 않는다.

```sh
OPENDART_API_KEY=... \
OPENDART_E2E_TARGETS_FILE=tests/e2e/kospi200.targets \
OPENDART_E2E_MIN_TARGETS=200 \
go test -tags=e2e -run TestE2EFinancialMetricAuditForMajorCompanies -count=1 -v .
```

이 테스트는 `FnlttSinglAcntAll` raw row를 받은 뒤 정규화 레이어를 통과시키고, 다음 분포를 JSON 로그로 남긴다.

| 분포 | 설명 |
| --- | --- |
| raw field coverage | `account_id`, `account_nm`, `account_detail`, `thstrm_amount` 같은 원문 field가 몇 개 회사/row에서 채워졌는지 |
| metric coverage | `revenue`, `operating_income`, `net_income`, `assets`, `liabilities`, `equity` 등 normalized metric이 몇 개 회사에서 추출됐는지 |
| match method | `account_id_exact`, `account_name_alias`, `override` 별 매칭 수 |
| account coverage | 표본 전체에서 자주 등장하는 `account_id`, `account_nm` |
| company gaps | 회사별 누락 metric과 `sj_div` 구성 |

2024년 사업보고서(`11011`), 연결(`CFS`) 조건의 주요 200개 기업 audit 결과는 다음 판단의 기준이다.

| 항목 | 결과 |
| --- | --- |
| target / fetched / no data | 200개 / 188개 / 12개 |
| raw rows / normalized metrics / unmapped rows / issues | 46,625 / 2,608 / 44,017 / 0 |
| 100% metric coverage | `assets`, `liabilities`, `equity`, `net_income`, `operating_cash_flow` |
| 90% 이상 metric coverage | `operating_income`, `cash_and_cash_equivalents`, `revenue` |
| `-표준계정코드 미사용-` | 183개 회사, 3,889 rows |

`status=013`과 `message=조회된 데이타가 없습니다.`는 target 표본의 데이터 부재로 audit report의 `no_data`에 남긴다. 네트워크 오류, 인증 오류, decode 오류, 그 외 OpenDART business error는 여전히 `failures`로 남기고 테스트 실패로 본다.
