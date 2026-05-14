# API 목록

이 디렉터리는 OpenDART 공식 개발가이드의 전체 API 목록과 CLI command 대응 관계를 추적한다.

## SDK typed 구현

| API | Endpoint | SDK method | CLI |
| --- | --- | --- | --- |
| 고유번호 | `GET /api/corpCode.xml` | `Client.CorpCodes(ctx)` | `opendart corp-codes` |
| 단일회사 주요계정 | `GET /api/fnlttSinglAcnt.json` | `Client.FinancialStatement(ctx, query)` | `opendart financial-statement` |

## CLI 전체 API 구현 방식

- 공식 JSON API는 `internal/cli`의 catalog 기반 command로 제공하고, 기본적으로 OpenDART 원문 JSON을 stdout에 출력한다.
- `corp-codes`, `financial-statement`는 root SDK의 typed method를 호출한다.
- SDK typed method가 아직 없는 API는 CLI 내부 generic request helper로 호출한다. public SDK API를 불필요하게 넓히지 않기 위한 선택이다.
- `document.xml`, `corpCode.xml`, `fnlttXbrl.xml` 같은 파일/XML 계열은 기본 `json` 출력에서 base64 envelope로 감싸고, `--output raw`일 때 원문 bytes를 stdout에 쓴다.
- 에러와 진단은 stderr로 보내며, 인증키 값은 출력하지 않는다.

전체 command와 공식 문서 대응은 `docs/apis/official-inventory.md`에 남긴다.

## 공통 메시지 코드

| status | 의미 |
| --- | --- |
| `000` | 정상 |
| `010` | 등록되지 않은 키 |
| `011` | 사용할 수 없는 키 |
| `012` | 접근할 수 없는 IP |
| `013` | 조회된 데이터 없음 |
| `014` | 파일 없음 |
| `020` | 요청 제한 초과 |
| `021` | 조회 가능한 회사 개수 초과 |
| `100` | 필드 값 부적절 |
| `101` | 부적절한 접근 |
| `800` | 시스템 점검 |
| `900` | 정의되지 않은 오류 |
| `901` | 계정 개인정보 보유기간 만료 |
