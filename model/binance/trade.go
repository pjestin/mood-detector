package binance

type Trade struct {
	Symbol          string
	Id              int
	OrderId         int
	OrderListId     int
	Price           string
	Qty             string
	QuoteQty        string
	Commission      string
	CommissionAsset string
	Time            int
	IsBuyer         bool
	IsBuyerMaker    bool
	IsBestMatch     bool
}
