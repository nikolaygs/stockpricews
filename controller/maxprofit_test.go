package controller

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"stockpricews/entity"
	"testing"
	"time"
)

func TestMaxProfitForPeriod(t *testing.T) {
	initialTime := time.Now()
	times := make([]time.Time, 10)
	for i := 0; i < 10; i++ {
		times[i] = initialTime.Add(time.Second * time.Duration(i))
	}

	testCases := []struct {
		name        string
		history     []entity.StockQuote
		expected    entity.MaxProfitPoints
		expectedErr error
	}{
		{
			name:        "No history records found - error not found",
			history:     []entity.StockQuote{},
			expectedErr: entity.ErrNotFound,
		},
		{
			name: "Prices in ascending order",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 1.0},
				{Datepoint: times[1], Price: 2.0},
				{Datepoint: times[2], Price: 3.0},
				{Datepoint: times[3], Price: 4.0},
			},
			expected: entity.MaxProfitPoints{
				BuyPoint:  entity.TradePoint{1.0, times[0]},
				SellPoint: entity.TradePoint{4.0, times[3]},
			},
		},
		{
			name: "Max profit is not at the max price",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 4.0},
				{Datepoint: times[1], Price: 7.0},
				{Datepoint: times[2], Price: 1.0},
				{Datepoint: times[3], Price: 5.0},
			},
			expected: entity.MaxProfitPoints{
				BuyPoint:  entity.TradePoint{1.0, times[2]},
				SellPoint: entity.TradePoint{5.0, times[3]},
			},
		},
		{
			name: "Two variants for max profit - take the earliest",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 4.0},
				{Datepoint: times[1], Price: 7.0},
				{Datepoint: times[2], Price: 2.0},
				{Datepoint: times[3], Price: 6.0},
				{Datepoint: times[4], Price: 1.0},
				{Datepoint: times[5], Price: 5.0},
			},
			expected: entity.MaxProfitPoints{
				BuyPoint:  entity.TradePoint{2.0, times[2]},
				SellPoint: entity.TradePoint{6.0, times[3]},
			},
		},
		{
			name: "Three variants for max profit - take the shortest",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 1.0},
				{Datepoint: times[1], Price: 4.0},
				{Datepoint: times[2], Price: 4.0},
				{Datepoint: times[3], Price: 4.0},
			},
			expected: entity.MaxProfitPoints{
				BuyPoint:  entity.TradePoint{1.0, times[0]},
				SellPoint: entity.TradePoint{4.0, times[1]},
			},
		},
		{
			name: "Two variants for max profit - take the earliest over shortest ",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 1.0},
				{Datepoint: times[1], Price: 1.0},
				{Datepoint: times[2], Price: 1.0},
				{Datepoint: times[3], Price: 2.0},
				{Datepoint: times[4], Price: 1.0},
				{Datepoint: times[5], Price: 2.0},
			},
			expected: entity.MaxProfitPoints{
				BuyPoint:  entity.TradePoint{1.0, times[0]},
				SellPoint: entity.TradePoint{2.0, times[3]},
			},
		},
		{
			name: "Prices in ascending order - we can't realize profits",
			history: []entity.StockQuote{
				{Datepoint: times[0], Price: 4.0},
				{Datepoint: times[1], Price: 3.0},
				{Datepoint: times[2], Price: 2.0},
				{Datepoint: times[3], Price: 1.0},
			},
			expectedErr: entity.ErrNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := maxProfitForPeriod(tt.history)
			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
