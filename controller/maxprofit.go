package controller

import (
	"fmt"
	"stockpricews/entity"
	"stockpricews/repository"
)

type MaxProfitController struct {
	Repository repository.Repository
}

// New initializes MaxProfitController that is used to calculate the maximum possible profit in a given historical time slice
func New(repository repository.Repository) MaxProfitController {
	return MaxProfitController{Repository: repository}
}

func (c MaxProfitController) MaxProfitForPeriod(req entity.StockQuoteRequest) (entity.MaxProfitPoints, error) {
	history, err := c.Repository.StockQuotesPerTimeSlice(req)
	if err != nil {
		return entity.MaxProfitPoints{}, err
	}

	return maxProfitForPeriod(history)
}

func maxProfitForPeriod(history []entity.StockQuote) (entity.MaxProfitPoints, error) {
	if len(history) == 0 {
		return entity.MaxProfitPoints{}, fmt.Errorf("no records found for the given period: %w", entity.ErrNotFound)
	}

	// holds the current max margin
	maxMargin := float64(0)
	// holds the current lowest price
	lowestPrice := history[0].Price

	// indexes of current low price and the prices that were used to compute the max margin
	var currLowIdx, lowIdx, highIdx int
	for i := 1; i < len(history); i++ {
		currentPrice := history[i].Price
		margin := currentPrice - lowestPrice
		if margin < 0 {
			// new lowest price - update the price and indexes
			lowestPrice = currentPrice
			currLowIdx = i
		} else if margin > maxMargin {
			// new max margin - update the margin and indexes
			maxMargin = margin
			highIdx = i
			lowIdx = currLowIdx
		}
	}

	if maxMargin == 0 {
		return entity.MaxProfitPoints{}, fmt.Errorf("it's not possible to realize a profit in the given period: %w", entity.ErrNotFound)
	}

	return entity.MaxProfitPoints{
		BuyPoint:  entity.TradePoint{Price: history[lowIdx].Price, Date: history[lowIdx].Datepoint},
		SellPoint: entity.TradePoint{Price: history[highIdx].Price, Date: history[highIdx].Datepoint},
	}, nil
}
