package repository

import (
	"stockpricews/entity"
)

// Repository an interface for loading stock quotes for given time period
type Repository interface {
	StockQuotesPerTimeSlice(timeSlice entity.StockQuoteRequest) ([]entity.StockQuote, error)
}
