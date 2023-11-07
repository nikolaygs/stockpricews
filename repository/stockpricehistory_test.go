package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"stockpricews/entity"
	"time"

	"log"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestStockQuotesPerTimeSlice(t *testing.T) {
	db, mock := NewMock()
	repo := &DBRepository{db: db}

	from := time.Unix(1699356339, 0)
	to := time.Unix(2699356339, 0)
	rows := sqlmock.NewRows([]string{"id", "symbol", "price", "datapoint"}).
		AddRow("1", "UBER", "19.99", time.Unix(1999356339, 0))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM stock_quote WHERE symbol = ? AND datepoint > ? AND datepoint < ? ORDER BY datepoint ASC")).
		WithArgs("UBER", from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05")).WillReturnRows(rows)

	history, err := repo.StockQuotesPerTimeSlice(entity.StockQuoteRequest{Symbol: "UBER", Begin: from, End: to})
	assert.NotNil(t, history)
	assert.NoError(t, err)
	assert.True(t, len(history) == 1)
	assert.Equal(t, entity.StockQuote{ID: 1, Symbol: "UBER", Datepoint: time.Date(2033, time.May, 10, 19, 45, 39, 0, time.Local), Price: 19.99}, history[0])
}
