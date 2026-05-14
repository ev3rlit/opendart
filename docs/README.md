# OpenDART 문서 수집

이 디렉터리는 OpenDART 공식 개발가이드에서 확인한 API 정보를 repo-local 형태로 정리한다.

## 기준

- 확인 날짜: 2026-05-14
- 공식 개발가이드:
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS001
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS002
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS003
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS004
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS005
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS006

## 수집 범위

- 공식 API 전체 인벤토리 85개
- SDK typed method 대응표
- 공식 응답 필드 기준 OpenAPI 추출물
- 고유번호, 단일회사 주요계정, 단일회사 전체 재무제표 상세 필드 문서

현재 SDK는 기존 공식 인벤토리의 83개 API를 root package typed method로 제공한다. 공식 개발가이드 기준 85개 API의 endpoint, 요청 파라미터, 응답 필드는 `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`에서 추적하고, SDK method 대응은 `docs/apis/typed-sdk-checklist.md`에서 추적한다.

## 갱신 방법

1. `python3 scripts/scrape_opendart_openapi.py`로 공식 개발가이드 목록/상세 페이지를 다시 수집한다.
2. `docs/apis/opendart.openapi.json`, `docs/apis/opendart.openapi.bundle.json`, `docs/apis/openapi/apis/*.json`, `docs/apis/opendart-api-metadata.json` 변경을 확인한다.
3. 새 API나 필드가 public API에 영향을 주면 모델과 테스트 fixture를 함께 갱신한다.
