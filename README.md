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

	codes, err := client.CorpCodes(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("corp codes:", len(codes))

	statements, err := client.FinancialStatement(ctx, opendart.FinancialStatementQuery{
		CorpCode:     "00126380",
		BusinessYear: "2025",
		ReportCode:   opendart.ReportAnnual,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("financial statements:", len(statements))
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

## MVP 지원 API

- 고유번호 조회: `Client.CorpCodes(ctx)`
  - 공식 endpoint: `GET /api/corpCode.xml`
  - ZIP 안의 XML을 파싱해 `CorpCode` 목록을 반환합니다.
- 단일회사 주요계정: `Client.FinancialStatement(ctx, query)`
  - 공식 endpoint: `GET /api/fnlttSinglAcnt.json`
  - 최소 query: `corp_code`, `bsns_year`, `reprt_code`

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

공식 API 전체 목록과 CLI 대응표는 `docs/apis/official-inventory.md`에 있습니다.

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
