package opendart

import "context"

// MajorStockReport is a 대량보유 상황보고 row.
type MajorStockReport ReportFields

// ExecutiveStockReport is an 임원ㆍ주요주주 소유보고 row.
type ExecutiveStockReport ReportFields

// MajorStock returns 대량보유 상황보고 rows.
func (client *Client) MajorStock(ctx context.Context, query CorpCodeQuery) ([]MajorStockReport, error) {
	return getCorpList[MajorStockReport](ctx, client, "MajorStock", "/api/majorstock.json", query)
}

// ExecutiveStock returns 임원ㆍ주요주주 소유보고 rows.
func (client *Client) ExecutiveStock(ctx context.Context, query CorpCodeQuery) ([]ExecutiveStockReport, error) {
	return getCorpList[ExecutiveStockReport](ctx, client, "ExecutiveStock", "/api/elestock.json", query)
}
