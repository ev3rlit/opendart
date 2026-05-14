# API 목록

이 디렉터리는 OpenDART 공식 개발가이드의 전체 API 목록, 응답 필드, CLI command 대응 관계를 추적한다.

## SDK typed 구현

공식 개발가이드에는 2026-05-14 기준 85개 API가 있다. 현재 root package typed method는 기존 83개 API를 지원하며, 2026년에 추가된 Ver 2.0 2개 API는 `docs/apis/openapi.md`에서 별도로 추적한다.

- 전체 대응표: `docs/apis/typed-sdk-checklist.md`
- 기준 인벤토리: `docs/apis/official-inventory.md`
- OpenAPI index: `docs/apis/opendart.openapi.json`
- OpenAPI bundle: `docs/apis/opendart.openapi.bundle.json`
- API별 OpenAPI 파일: `docs/apis/openapi/apis/*.json`
- OpenAPI 생성 기준: `docs/apis/openapi.md`
- SDK package: `github.com/ev3rlit/opendart`

## CLI 구현 방식

- 현재 CLI 구현 범위의 JSON API는 `internal/cli`의 catalog 기반 command로 제공하고, 기본적으로 OpenDART 원문 JSON을 stdout에 출력한다.
- `corp-codes`, `financial-statement`와 파일 API 일부는 root SDK의 typed method를 호출한다.
- 나머지 JSON command는 원문 JSON stdout 계약을 유지하기 위해 CLI 내부 generic request helper로 호출한다.
- `document.xml`, `corpCode.xml`, `fnlttXbrl.xml` 같은 파일/XML 계열은 기본 `json` 출력에서 base64 envelope로 감싸고, `--output raw`일 때 원문 bytes를 stdout에 쓴다.
- 에러와 진단은 stderr로 보내며, 인증키 값은 출력하지 않는다.

command와 공식 문서 대응은 `docs/apis/official-inventory.md`에 남긴다.

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
