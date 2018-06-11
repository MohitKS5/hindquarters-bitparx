package response_types

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Symbols []Symbol `json:"symbols"`
}

// Symbol market symbol
type Symbol struct {
	Symbol             string              `json:"symbol"`
	Status             string              `json:"status"`
	BaseAsset          string              `json:"baseAsset"`
	BaseAssetPrecision int                 `json:"baseAssetPrecision"`
	QuoteAsset         string              `json:"quoteAsset"`
	QuotePrecision     int                 `json:"quotePrecision"`
	OrderTypes         []string            `json:"orderTypes"`
	IcebergAllowed     bool                `json:"icebergAllowed"`
	Filters            []map[string]string `json:"filters"`
}
