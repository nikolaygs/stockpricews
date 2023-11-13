package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	_ "golang.org/x/time/rate"
	"net/http"
	"stockpricews/controller"
	"stockpricews/entity"
	"strconv"
	"time"
)

const (
	begin  = "begin"
	end    = "end"
	symbol = "symbol"
)

type StockPriceHandler struct {
	Controller controller.Controller
}

// New initializes new StockPriceHandler that currently provides just one REST endpoint 'GET /maxprofit'
func New(controller controller.Controller, port int) (StockPriceHandler, error) {
	handerImpl := StockPriceHandler{Controller: controller}
	http.Handle("/maxprofit", rateLimiter(handerImpl.MaxProfitForPeriod))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return handerImpl, err
}

// MaxProfitForPeriod is HTTP handler that returns to client the maximum profit that could be realized within given time slice.
// Usage: curl GET /maxprofit?begin=<begin_time_in_seconds>&end=<end_time_in_seconds>&symbol=<STOCK_SYMBOL>
// Result status codes:
//  - 200 OK - when a profit can be realized within the given time slice. Body contains entity.MaxProfitPoints as json
//  - 400 Bad Request - if any of the query params is not passed or doesn't have a correct format (seconds). Body contains entity.ErrorMessage as json so the client can handle it accordingly
//  - 404 Not Found - if stock quote data can't be found for the given time slice or it's not possible to realize a profit. Body contains entity.ErrorMessage as json so the client can handle it accordingly
//  - 429 Too Many Requests if the client got rate limited.
//  - 500 Intenal Server Error - if any expected error occur.
func (h StockPriceHandler) MaxProfitForPeriod(w http.ResponseWriter, r *http.Request) {
	// Access-Control-Allow-Origin is set as the client might run in a separate machine
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse request data and report BadRequest if any of the params can't be found/parsed
	timeSlice, err := parseRequestData(r)
	if err != nil {
		respondWithError(err, w)
		return
	}

	// Calculate max profit for the given time slice and report error if any
	maxProfitPrices, err := h.Controller.MaxProfitForPeriod(timeSlice)
	if err != nil {
		respondWithError(err, w)
		return
	}

	// Marshal the response to JSON and report successful execution
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(maxProfitPrices)
}

// Simple rate limiting using Token Bucket
func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			errMsg := entity.ErrorMessage{Message: "The API is at capacity, try again later."}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			next(w, r)
		}
	})
}

func parseRequestData(r *http.Request) (entity.StockQuoteRequest, error) {
	if r == nil || r.URL == nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("failed to read request URL: %w", entity.ErrBadRequest)
	}

	if !(r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions) {
		return entity.StockQuoteRequest{}, fmt.Errorf("method %s not allowed: %w", r.Method, entity.ErrMethodNotAllowed)
	}

	if !r.URL.Query().Has(begin) {
		return entity.StockQuoteRequest{}, fmt.Errorf("%s param is missing: %w", begin, entity.ErrBadRequest)
	}

	if !r.URL.Query().Has(end) {
		return entity.StockQuoteRequest{}, fmt.Errorf("%s param is missing: %w", end, entity.ErrBadRequest)
	}

	if !r.URL.Query().Has(symbol) {
		return entity.StockQuoteRequest{}, fmt.Errorf("%s param is missing: %w", symbol, entity.ErrBadRequest)
	}

	beginSecs, err := strconv.ParseInt(r.URL.Query().Get(begin), 10, 64)
	if err != nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("%s param can't be parsed as seconds: %w", begin, entity.ErrBadRequest)
	}

	endSecs, err := strconv.ParseInt(r.URL.Query().Get("end"), 10, 64)
	if err != nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("%s param can't be parsed as seconds: %w", end, entity.ErrBadRequest)
	}

	timeSlice := entity.StockQuoteRequest{Begin: time.Unix(beginSecs, 0), End: time.Unix(endSecs, 0)}
	if timeSlice.Begin.After(timeSlice.End) {
		return entity.StockQuoteRequest{}, fmt.Errorf("begin period is after the end period: %w", entity.ErrBadRequest)
	}

	stockSymbol := r.URL.Query().Get(symbol)
	if len(stockSymbol) < 1 || len(stockSymbol) > 4 {
		return entity.StockQuoteRequest{}, fmt.Errorf("stock symbol must be between 1 and 4 chars long: %w", entity.ErrBadRequest)
	}
	timeSlice.Symbol = stockSymbol

	return timeSlice, nil
}

func respondWithError(err error, w http.ResponseWriter) {
	// log the error at the server log for debug purposes
	fmt.Println(err)

	statusCode := 200
	errMsg := err.Error()
	if errors.Is(err, entity.ErrBadRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, entity.ErrMethodNotAllowed) {
		statusCode = http.StatusMethodNotAllowed
	} else {
		// we don't want to leak internal messages to the client
		statusCode = http.StatusInternalServerError
		errMsg = "Internal server error"
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(entity.ErrorMessage{Message: errMsg})
}
