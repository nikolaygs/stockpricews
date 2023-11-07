package controller

import (
	"stockpricews/entity"
)

type Controller interface {
	MaxProfitForPeriod(timeSlice entity.StockQuoteRequest) (entity.MaxProfitPoints, error)
}
