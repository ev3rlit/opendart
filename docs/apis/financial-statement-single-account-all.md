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
