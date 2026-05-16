# 단일회사 주요계정

- 공식 문서: https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS003&apiId=2019016
- 그룹: 정기보고서 재무정보
- SDK raw method: `Client.FnlttSinglAcnt(ctx, params)`
- Friendly name history: `Client.FinancialStatement(ctx, query)` is preserved in `docs/apis/sdk-names.yaml`.

## 기본 정보

| 항목 | 값 |
| --- | --- |
| Method | `GET` |
| JSON URL | `https://opendart.fss.or.kr/api/fnlttSinglAcnt.json` |
| XML URL | `https://opendart.fss.or.kr/api/fnlttSinglAcnt.xml` |
| Encoding | UTF-8 |

## 요청 인자

| key | 명칭 | type | 필수 | 설명 |
| --- | --- | --- | --- | --- |
| `crtfc_key` | API 인증키 | `STRING(40)` | Y | 발급받은 인증키 |
| `corp_code` | 고유번호 | `STRING(8)` | Y | 공시대상회사의 고유번호 |
| `bsns_year` | 사업연도 | `STRING(4)` | Y | 2015년 이후 제공 |
| `reprt_code` | 보고서 코드 | `STRING(5)` | Y | `11013`, `11012`, `11014`, `11011` |

## 보고서 코드

| code | 의미 |
| --- | --- |
| `11013` | 1분기보고서 |
| `11012` | 반기보고서 |
| `11014` | 3분기보고서 |
| `11011` | 사업보고서 |

## 응답 필드

| key | 설명 |
| --- | --- |
| `status` | 에러 및 정보 코드 |
| `message` | 에러 및 정보 메시지 |
| `rcept_no` | 접수번호 |
| `bsns_year` | 사업연도 |
| `stock_code` | 상장회사 종목코드 |
| `reprt_code` | 보고서 코드 |
| `account_nm` | 계정명 |
| `fs_div` | 개별/연결구분, `OFS` 또는 `CFS` |
| `fs_nm` | 개별/연결명 |
| `sj_div` | 재무제표구분, `BS` 또는 `IS` |
| `sj_nm` | 재무제표명 |
| `thstrm_nm` | 당기명 |
| `thstrm_dt` | 당기일자 |
| `thstrm_amount` | 당기금액 |
| `thstrm_add_amount` | 당기누적금액 |
| `frmtrm_nm` | 전기명 |
| `frmtrm_dt` | 전기일자 |
| `frmtrm_amount` | 전기금액 |
| `frmtrm_add_amount` | 전기누적금액 |
| `bfefrmtrm_nm` | 전전기명 |
| `bfefrmtrm_dt` | 전전기일자 |
| `bfefrmtrm_amount` | 전전기금액 |
| `ord` | 계정과목 정렬순서 |
| `currency` | 통화 단위 |

## MVP 결정

MVP 재무제표 API는 단일회사 주요계정으로 구현한다. 공식 문서 기준 필수 query가 `corp_code`, `bsns_year`, `reprt_code`로 단순하고, 사용 예시 목표와도 맞다.
