# 고유번호

- 공식 문서: https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS001&apiId=2019018
- 그룹: 공시정보
- SDK raw method: `Client.CorpCode(ctx)`
- Friendly name history: `Client.CorpCodes(ctx)` is preserved in `docs/apis/sdk-names.yaml`.

## 기본 정보

| 항목 | 값 |
| --- | --- |
| Method | `GET` |
| URL | `https://opendart.fss.or.kr/api/corpCode.xml` |
| Encoding | UTF-8 |
| Output | ZIP FILE(binary) |

## 요청 인자

| key | 명칭 | type | 필수 | 설명 |
| --- | --- | --- | --- | --- |
| `crtfc_key` | API 인증키 | `STRING(40)` | Y | 발급받은 인증키 |

## 응답 구조

응답은 ZIP 파일이며 ZIP 안의 XML 파일에 회사 목록이 들어 있다.

| key | 명칭 | 설명 |
| --- | --- | --- |
| `status` | 에러 및 정보 코드 | 메시지 설명 참고 |
| `message` | 에러 및 정보 메시지 | 메시지 설명 참고 |
| `corp_code` | 고유번호 | 공시대상회사의 8자리 고유번호 |
| `corp_name` | 정식명칭 | 정식 회사명 |
| `corp_eng_name` | 영문 정식명칭 | 영문 정식 회사명 |
| `stock_code` | 종목코드 | 상장회사인 경우 6자리 종목코드 |
| `modify_date` | 최종변경일자 | 기업개황정보 최종변경일자, `YYYYMMDD` |

## 주의

`corp_code`는 OpenDART에서 쓰는 8자리 회사 고유번호이고, KRX 단축코드나 종목코드가 아니다. 재무제표 API 호출에는 `stock_code`가 아니라 `corp_code`를 사용한다.
