# OpenAPI 추출물

- 확인 날짜: 2026-05-14
- 루트 산출물: `docs/apis/opendart.openapi.json`
- 번들 산출물: `docs/apis/opendart.openapi.bundle.json`
- API별 산출물: `docs/apis/openapi/apis/*.json`
- 원본 표 덤프: `docs/apis/opendart-api-metadata.json`
- 재생성 명령: `python3 scripts/scrape_opendart_openapi.py`

## 수집 범위

공식 개발가이드의 6개 그룹을 모두 순회해 상세 페이지 85개를 파싱한다.

| Group | 공식 API 수 |
| --- | ---: |
| 공시정보 | 4 |
| 정기보고서 주요정보 | 30 |
| 정기보고서 재무정보 | 7 |
| 지분공시 종합정보 | 2 |
| 주요사항보고서 주요정보 | 36 |
| 증권신고서 주요정보 | 6 |
| 합계 | 85 |

## 스키마화 규칙

- 루트 OpenAPI 파일은 85개 path를 API별 파일로 연결하는 `$ref` index로 둔다.
- 번들 OpenAPI 파일은 외부 `$ref` 없이 85개 API의 path와 schema를 한 파일에 모두 담는다.
- 각 API 파일은 `path.get` operation과 해당 API에서만 쓰는 `components.schemas`를 함께 가진다.
- `crtfc_key`는 모든 API의 공통 query 인증키라서 OpenAPI `apiKey` security scheme으로 정의한다.
- 공식 `STRING(n)` 요청 타입은 query parameter의 `type: string`, `maxLength: n`으로 변환한다.
- JSON 응답의 `status`, `message`, paging field, `list[]` item field를 공식 응답 표 기준으로 분리한다.
- 증권신고서 주요정보처럼 `group/title/list`가 반복되는 응답은 그룹별 item schema를 만들고, 그룹 제목은 `const`로 남긴다.
- `document.xml`, `corpCode.xml`, `fnlttXbrl.xml` 파일/XML 계열은 `application/octet-stream` 응답으로 둔다.

## 파일 배치

| 파일 | 용도 |
| --- | --- |
| `docs/apis/opendart.openapi.json` | 전체 API path index. 파일 크기를 줄이기 위해 API별 파일을 외부 `$ref`로 참조한다. |
| `docs/apis/opendart.openapi.bundle.json` | 코드 생성기나 단일 파일 소비 도구를 위한 전체 번들 OpenAPI 문서. |
| `docs/apis/openapi/apis/{apiId}-{endpoint}.json` | API 1개 단위의 operation, parameter, response schema. |
| `docs/apis/opendart-api-metadata.json` | 공식 HTML 표를 손실 없이 보존한 원본 파싱 결과. |

예를 들어 공시검색은 `docs/apis/openapi/apis/2019001-list.json`, 기업개황은 `docs/apis/openapi/apis/2019002-company.json`에 있다.

## 구현 참고

현재 공식 개발가이드에는 기존 SDK 구현 범위보다 2개 많은 85개 API가 있다.

| API ID | API | Endpoint | 비고 |
| --- | --- | --- | --- |
| `2026001` | 이사·감사의 개인별 보수현황(5억원 이상) (Ver 2.0) | `/api/hmvAuditIndvdlBySttusV2.json` | 2026년 5월 이후 제출 보고서 대상 |
| `2026002` | 개인별 보수지급 금액(5억이상 상위5인) (Ver 2.0) | `/api/indvdlByPayV2.json` | 2026년 5월 이후 제출 보고서 대상 |
