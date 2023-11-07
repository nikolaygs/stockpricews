package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"stockpricews/entity"
	"time"
)

type DBRepository struct {
	db *sql.DB
}

const getStockQuotesPerTimeSlice = "SELECT * FROM stock_quote WHERE symbol = ? AND datepoint > ? AND datepoint < ? ORDER BY datepoint ASC"

// New initializes a new DB repository that connects to MySQL database
func New(user, pass string, port int) (DBRepository, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(localhost:%d)/stockquotedb?charset=utf8mb4,utf8&parseTime=true", user, pass, port)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return DBRepository{}, err
	}

	// Connection pool options
	// Connection pooling is internally provided thus we don't need to explicitly handle it
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return DBRepository{db: db}, nil
}

func (r DBRepository) StockQuotesPerTimeSlice(req entity.StockQuoteRequest) ([]entity.StockQuote, error) {
	// db.Query uses prepared statement under the hook for a performance optimization and SQL injection protection
	rows, err := r.db.Query(getStockQuotesPerTimeSlice, req.Symbol,
		req.Begin.Format("2006-01-02 15:04:05"), req.End.Format("2006-01-02 15:04:05"))
	if err != nil {
		return []entity.StockQuote{}, err
	}
	// essentially not needed as the sql.DB will close it internally as soon as rows iteration is over
	defer rows.Close()

	var history []entity.StockQuote
	for rows.Next() {
		quote := entity.StockQuote{}
		if err = rows.Scan(&quote.ID, &quote.Symbol, &quote.Price, &quote.Datepoint); err != nil {
			return history, err
		}
		history = append(history, quote)
	}

	return history, nil
}
