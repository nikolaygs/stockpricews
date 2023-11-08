package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"stockpricews/controller"
	"stockpricews/entity"
	"strconv"
	"time"
)

type StockPriceHandler struct {
	Controller controller.Controller
}

// New initializes new StockPriceHandler that currently provides just one REST endpoint 'GET /maxprofit'
func New(controller controller.Controller, port int) (StockPriceHandler, error) {
	handerImpl := StockPriceHandler{Controller: controller}
	http.HandleFunc("/maxprofit", handerImpl.MaxProfitForPeriod)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return handerImpl, err
}

// MaxProfitForPeriod is HTTP handler that returns to client the maximum profit that could be realized within given time slice.
// Usage: curl GET /maxprofit?begin=<begin_time_in_seconds>&end=<end_time_in_seconds>&symbol=<STOCK_SYMBOL>
// Result status codes:
//  - 200 OK - when a profit can be realized within the given time slice. Body contains entity.MaxProfitPoints as json
//  - 400 Bad Request - if any of the query params is not passed or doesn't have a correct format (seconds). Body contains entity.ErrorMessage as json so the client can handle it accordingly
//  - 404 Not Found - if stock quote data can't be found for the given time slice or it's not possible to realize a profit. Body contains entity.ErrorMessage as json so the client can handle it accordingly
//  - 500 Intenal Server Error - if any expected error occur.
func (h StockPriceHandler) MaxProfitForPeriod(w http.ResponseWriter, r *http.Request) {
	// Access-Control-Allow-Origin is set as the client might run in a separate machine
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse request data and report BadRequest if any of the params can't be found/parsed
	timeSlice, err := parseRequestData(r)
	if err != nil {
		writeErrorToStatusCodeAndMessage(err, w)
		return
	}

	// Calculate max profit for the given time slice and report error if any
	maxProfitPrices, err := h.Controller.MaxProfitForPeriod(timeSlice)
	if err != nil {
		writeErrorToStatusCodeAndMessage(err, w)
		return
	}

	// Marshal the response to JSON and report successful execution
	maxProfitPricesJSON, _ := json.Marshal(maxProfitPrices) // this shouldn't happen for simplicity we won't handle the error explicitly

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(maxProfitPricesJSON))
}

func parseRequestData(r *http.Request) (entity.StockQuoteRequest, error) {
	if r == nil || r.URL == nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("failed to read request URL: %w", entity.ErrBadRequest)
	}

	if !(r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS") {
		return entity.StockQuoteRequest{}, fmt.Errorf("method %s not allowed: %w", r.Method, entity.ErrMethodNotAllowed)
	}

	if !r.URL.Query().Has("begin") {
		return entity.StockQuoteRequest{}, fmt.Errorf("begin param is missing: %w", entity.ErrBadRequest)
	}

	if !r.URL.Query().Has("end") {
		return entity.StockQuoteRequest{}, fmt.Errorf("end param is missing: %w", entity.ErrBadRequest)
	}

	if !r.URL.Query().Has("symbol") {
		return entity.StockQuoteRequest{}, fmt.Errorf("symbol param is missing: %w", entity.ErrBadRequest)
	}

	begin, err := strconv.ParseInt(r.URL.Query().Get("begin"), 10, 64)
	if err != nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("begin param can't be parsed as seconds: %w", entity.ErrBadRequest)
	}

	end, err := strconv.ParseInt(r.URL.Query().Get("end"), 10, 64)
	if err != nil {
		return entity.StockQuoteRequest{}, fmt.Errorf("end param can't be parsed as seconds: %w", entity.ErrBadRequest)
	}

	timeSlice := entity.StockQuoteRequest{Begin: time.Unix(begin, 0), End: time.Unix(end, 0)}
	if timeSlice.Begin.After(timeSlice.End) {
		return entity.StockQuoteRequest{}, fmt.Errorf("begin interval is after the end interval: %w", entity.ErrBadRequest)
	}

	symbol := r.URL.Query().Get("symbol")
	if len(symbol) < 1 || len(symbol) > 4 {
		return entity.StockQuoteRequest{}, fmt.Errorf("stock symbol must be between 1 and 4 chars long: %w", entity.ErrBadRequest)
	}
	timeSlice.Symbol = symbol

	return timeSlice, nil
}

func writeErrorToStatusCodeAndMessage(err error, w http.ResponseWriter) {
	// log the error at the server log for debug purposes
	fmt.Println(err)

	statusCode := 200
	body := errorAsJSON(err)
	if errors.Is(err, entity.ErrBadRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, entity.ErrMethodNotAllowed) {
		statusCode = http.StatusMethodNotAllowed
	} else {
		// we don't want to leak internal messages to the client
		statusCode = http.StatusInternalServerError
		body = "Internal server error"
	}

	w.WriteHeader(statusCode)
	io.WriteString(w, body)
}

func errorAsJSON(err error) string {
	errMsg := entity.ErrorMessage{Message: err.Error()}
	res, _ := json.Marshal(errMsg) // this shouldn't happen for simplicity we won't handle the error explicitly
	return string(res)
}
