package binance

type Order struct {
	Symbol              string
	OrderId             int
	OrderListId         int
	ClientOrderId       string
	TransactTime        uint64
	Price               string
	OrigQty             string
	ExecutedQty         string
	CummulativeQuoteQty string
	Status              string
	TimeInForce         string
	Type                string
	Side                string
}
