# RULE.md

이 저장소는 `github.com/ev3rlit/opendart` 독립 Go SDK를 제공한다. `mwosa`나 다른 상위 애플리케이션에 의존하지 않고, OpenDART 공식 API를 얇고 안정적인 타입 메서드로 감싼다.

## Go 코드

- package name은 `opendart`를 사용한다.
- 공개 생성자는 `New(Config, ...Option)` 형태를 유지한다.
- 필수 값은 명시적인 `Config` struct에 둔다.
- 선택 값은 `With...` option으로 추가한다.
- 내부 HTTP 호출은 `resty`를 사용한다.
- I/O를 수행하는 공개 메서드는 첫 인자로 `context.Context`를 받는다.
- 공개 타입과 메서드는 Go doc comment를 작성한다.
- 불필요한 `any`, `interface{}`, broad catch, silent fallback을 피한다.

## Error Handling

- HTTP status error, JSON/XML/ZIP decode error, OpenDART business error를 구분한다.
- OpenDART 응답의 `status`와 `message`는 손실 없이 보존한다.
- decode 실패를 성공처럼 처리하지 않는다.
- secret, 인증키, 원문 credential은 error string, log, test fixture에 남기지 않는다.

## Tests

- 기본 테스트는 offline fake HTTP server 기반으로 작성한다.
- 테스트 도구는 `testing`과 `github.com/stretchr/testify`를 사용한다.
- `require`는 후속 검증이 의미 없어지는 조건에 사용하고, `assert`는 독립 필드 검증에 사용한다.
- live OpenDART 호출은 기본 `go test ./...`에 포함하지 않는다.
- live/e2e가 필요하면 `//go:build e2e`와 secret 기반 optional workflow로 분리한다.

## Documentation

- 공식 OpenDART 문서에서 확인한 요청 파라미터, 응답 필드, 메시지 코드를 repo-local 문서로 남긴다.
- 링크 목록만 남기지 않고 SDK 구현에 필요한 필드 의미와 제약을 같이 정리한다.
- 문서 갱신 시 기준 URL과 확인 날짜를 함께 남긴다.
- `corp_code`는 KRX 단축코드/종목코드와 다르다는 점을 README와 API 문서에 명확히 적는다.

## Secrets

- 인증키, 샘플 secret, 개인 계정 값은 커밋하지 않는다.
- README와 테스트는 `OPENDART_API_KEY` 같은 환경변수 이름만 예시로 사용한다.
- fixture에는 실제 OpenDART 인증키나 실사용 계정 정보를 넣지 않는다.

## Verification

기본 변경은 아래 명령을 통과해야 한다.

```sh
go mod tidy
go test ./...
git diff --check
```

가능하면 아래 검증도 함께 수행한다.

```sh
go test -race ./...
go vet ./...
```

