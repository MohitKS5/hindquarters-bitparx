package response_types

// Order define order info
type Order struct {
	Symbol           string `json:"symbol"`
	OrderID          int64  `json:"orderId"`
	ClientOrderID    string `json:"clientOrderId"`
	Price            string `json:"price"`
	OrigQuantity     string `json:"origQty"`
	ExecutedQuantity string `json:"executedQty"`
	Status           string `json:"status"`
	TimeInForce      string `json:"timeInForce"`
	Type             string `json:"type"`
	Side             string `json:"side"`
	StopPrice        string `json:"stopPrice"`
	IcebergQuantity  string `json:"icebergQty"`
	Time             int64  `json:"time"`
}
