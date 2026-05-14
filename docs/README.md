# OpenDART 문서 수집

이 디렉터리는 OpenDART 공식 개발가이드에서 확인한 API 정보를 repo-local 형태로 정리한다.

## 기준

- 확인 날짜: 2026-05-14
- 공식 개발가이드:
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS001
  - https://opendart.fss.or.kr/guide/main.do?apiGrpCd=DS003
  - https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS001&apiId=2019018
  - https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS003&apiId=2019016
  - https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=DS003&apiId=2019020

## 수집 범위

- 공식 API 전체 인벤토리 83개
- SDK typed method 전체 대응표
- 고유번호, 단일회사 주요계정, 단일회사 전체 재무제표 상세 필드 문서

현재 SDK는 공식 인벤토리의 83개 API를 root package typed method로 제공한다. API별 endpoint, CLI command, SDK method 대응은 `docs/apis/typed-sdk-checklist.md`에서 추적한다.

## 갱신 방법

1. 공식 개발가이드 목록 페이지에서 API 그룹별 목록을 확인한다.
2. 구현 대상 API의 상세 페이지를 열고 요청 URL, 필수 파라미터, 응답 필드, 메시지 코드를 확인한다.
3. `docs/apis/` 아래 문서를 갱신한다.
4. 변경된 필드가 public API에 영향을 주면 모델과 테스트 fixture를 함께 갱신한다.
