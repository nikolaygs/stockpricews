package entity

import "time"

type StockQuoteRequest struct {
	Symbol string
	Begin  time.Time
	End    time.Time
}

type StockQuote struct {
	ID        int64
	Symbol    string
	Datepoint time.Time
	Price     float64
}

type TradePoint struct {
	Price float64   `json:"price"`
	Date  time.Time `json:"date"`
}

type MaxProfitPoints struct {
	BuyPoint  TradePoint `json:"buyPoint"`
	SellPoint TradePoint `json:"sellPoint"`
}
