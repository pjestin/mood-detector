package binance

type Kline struct {
	OpenTime         uint64
	Open             string
	High             string
	Low              string
	Close            string
	Volume           string
	CloseTime        uint64
	QuoteAssetVolume string
	NumberOfTrades   uint64
}
