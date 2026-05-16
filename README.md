# opendart

`opendart`는 OpenDART API를 위한 독립 Go 클라이언트 라이브러리와 CLI입니다. Go 사용자는 계속 root package `github.com/ev3rlit/opendart`를 import하고, CLI executable은 `cmd/opendart`에서 별도로 관리합니다.

## 설치

```sh
go get github.com/ev3rlit/opendart
```

## 사용 예시

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/ev3rlit/opendart"
)

func main() {
	client, err := opendart.New(opendart.Config{
		APIKey: os.Getenv("OPENDART_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	codes, err := client.CorpCode(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("corp code file bytes:", len(codes.Body))

	statements, err := client.FnlttSinglAcnt(ctx, opendart.FnlttSinglAcntParams{
		CorpCode:  "00126380",
		BsnsYear:  "2025",
		ReprtCode: "11011",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("financial statements:", len(statements.List))
}
```

## 인증키

OpenDART 인증키는 코드에 직접 쓰지 말고 환경변수로 주입합니다.

```sh
export OPENDART_API_KEY="발급받은_인증키"
```

CLI도 같은 환경변수를 기본으로 사용합니다. 필요하면 `--api-key` flag로 한 번만 넘길 수 있습니다.

## corp_code와 KRX symbol 차이

OpenDART의 `corp_code`는 DART 공시대상회사 고유번호입니다. KRX 종목코드 또는 단축코드와 다릅니다.

- `corp_code`: OpenDART API 호출에 쓰는 8자리 고유번호 예: `00126380`
- `stock_code`: 상장회사인 경우 제공되는 6자리 종목코드 예: `005930`

재무제표 API는 `stock_code`가 아니라 `corp_code`를 요구합니다.

## 지원 범위

- SDK generated raw API 지원: OpenAPI `operationId` 기반 root package method를 제공합니다.
  - 예: `Client.List(ctx, params)`, `Client.CompanyRaw(ctx, params)`, `Client.DocumentRaw(ctx, params)`
  - 예: `Client.IrdsSttus(ctx, params)`, `Client.CmpMgDecsn(ctx, params)`
  - 예: `Client.FnlttSinglAcnt(ctx, params)`, `Client.FnlttXbrl(ctx, params)`
- CLI API 지원: `internal/cli/catalog.go` 기준으로 기존 공식 API command를 제공합니다.
  - JSON API는 기본적으로 OpenDART 원문 JSON을 stdout에 출력합니다.
  - 파일/XML API는 기본 `json` 출력에서 base64 envelope를 쓰고, `--output raw`일 때 원문 bytes를 stdout에 씁니다.
- 2026-05-14 기준 공식 개발가이드에는 85개 API가 있으며, OpenAPI index는 `docs/apis/opendart.openapi.json`, 단일 파일 bundle은 `docs/apis/opendart.openapi.bundle.json`, API별 스키마는 `docs/apis/openapi/apis/*.json`에 있습니다.
- 과거 수작업 friendly method 이름은 `docs/apis/sdk-names.yaml`에 보존되어 있습니다.

## CLI

설치:

```sh
go install github.com/ev3rlit/opendart/cmd/opendart@latest
```

로컬 checkout에서 실행:

```sh
go run ./cmd/opendart --help
go run ./cmd/opendart corp-codes
go run ./cmd/opendart financial-statement \
  --corp-code 00126380 \
  --business-year 2025 \
  --report-code 11011
```

공통 flag:

- `--api-key`: OpenDART 인증키. 비어 있으면 `OPENDART_API_KEY`를 사용합니다.
- `--base-url`: fake server나 테스트용 base URL을 지정합니다.
- `--output`: 기본값은 `json`입니다. 파일/XML API는 `raw`를 지정하면 원문 bytes를 stdout에 씁니다.

API 그룹별 command:

| Group | 예시 |
| --- | --- |
| `disclosure` | `opendart disclosure list`, `opendart disclosure company`, `opendart disclosure document` |
| `company` | `opendart company irds-sttus`, `opendart company alot-matter` |
| `financial` | `opendart financial single-account`, `opendart financial single-account-all`, `opendart financial xbrl` |
| `ownership` | `opendart ownership major-stock`, `opendart ownership executive-stock` |
| `material` | `opendart material cmp-mg-decsn`, `opendart material stk-extr-decsn` |
| `registration` | `opendart registration equity`, `opendart registration debt` |

공식 API 전체 목록과 CLI 대응표는 `docs/apis/official-inventory.md`에 있고, 공식 응답 필드 기반 OpenAPI 추출물은 `docs/apis/openapi.md`에서 설명합니다.

## 오류 처리

다음 오류 타입으로 원인을 구분할 수 있습니다.

- `*opendart.HTTPError`: HTTP status가 2xx가 아닌 경우
- `*opendart.DecodeError`: JSON, XML, ZIP 파싱 실패
- `*opendart.APIError`: OpenDART business status가 `000`이 아닌 경우

## 개발

```sh
go mod tidy
go test ./...
go test -race ./...
git diff --check
```

기본 테스트는 live OpenDART 호출을 하지 않습니다.

### Live e2e smoke

실제 OpenDART 서버를 호출하는 e2e smoke는 opt-in으로만 실행합니다. 인증키는 파일에 저장하지 않고 `OPENDART_API_KEY` 환경변수로 주입합니다.

```sh
OPENDART_API_KEY=... scripts/e2e-newman.sh
```

스크립트는 `tests/e2e/postman/opendart-smoke.postman_collection.json`을 Newman으로 실행하고 JUnit report를 `test-results/newman.xml`에 씁니다. `newman` 실행 파일이 없으면 `npx --yes newman`을 사용합니다. CLI reporter는 요청 URL의 query string에 인증키를 표시할 수 있어 기본으로 사용하지 않습니다.

로컬 Postman 환경 예시는 `tests/e2e/postman/opendart.local.postman_environment.example.json`에 있으며, 스크립트의 기본 base URL은 `https://opendart.fss.or.kr`입니다. 다른 endpoint로 실행해야 하면 `OPENDART_BASE_URL`로 덮어씁니다.

직접 실행할 때는 아래처럼 환경변수만 넘깁니다.

```sh
newman run tests/e2e/postman/opendart-smoke.postman_collection.json \
  --env-var "base_url=https://opendart.fss.or.kr" \
  --env-var "crtfc_key=$OPENDART_API_KEY" \
  --reporters junit \
  --reporter-junit-export test-results/newman.xml
```
