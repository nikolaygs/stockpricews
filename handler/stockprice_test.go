package handler

import (
	"errors"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/url"
	"stockpricews/entity"

	"testing"
)

func TestParseRequestData(t *testing.T) {
	testCases := []struct {
		name        string
		req         *http.Request
		expected    entity.StockQuoteRequest
		expectedErr error
	}{
		{
			name:        "request is nil",
			req:         nil,
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "URL is nil",
			req:         &http.Request{URL: &url.URL{RawQuery: ""}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Begin param is missing",
			req:         &http.Request{URL: &url.URL{RawQuery: "end=1699228800"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "End param is missing",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=1699228800"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Begin param can't be parsed",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=asd"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Symbol is param missing",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=2699228800&end=1699228800"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "End param can't be parsed",
			req:         &http.Request{URL: &url.URL{RawQuery: "end=asd"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Begin param after End param",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=2699228800&end=1699228800"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Symbol param with zero length",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=2699228800&end=1699228800&symbol="}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:        "Symbol param with five length",
			req:         &http.Request{URL: &url.URL{RawQuery: "begin=2699228800&end=1699228800&symbol=TESLA"}, Method: "GET"},
			expectedErr: entity.ErrBadRequest,
		},
		{
			name:     "URL successfully parsed",
			req:      &http.Request{URL: &url.URL{RawQuery: "begin=1699228800&end=2699228800&symbol=UBER"}, Method: "GET"},
			expected: entity.StockQuoteRequest{Symbol: "UBER", Begin: time.Unix(1699228800, 0), End: time.Unix(2699228800, 0)},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRequestData(tt.req)
			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

type MockController struct {
	err error
}

func (c MockController) MaxProfitForPeriod(req entity.StockQuoteRequest) (entity.MaxProfitPoints, error) {
	return entity.MaxProfitPoints{}, c.err
}

func TestMaxProfitForPeriod_StatusCodes(t *testing.T) {
	testCases := []struct {
		name               string
		handler            Handler
		method             string
		url                string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Successfull GET request",
			method:             "GET",
			url:                "maxprofit?begin=1699228800&end=2699228800&symbol=UBER",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "{\"buyPoint\":{\"price\":0,\"date\":\"0001-01-01T00:00:00Z\"},\"sellPoint\":{\"price\":0,\"date\":\"0001-01-01T00:00:00Z\"}}\n",
		},
		{
			name:               "Bad request",
			method:             "GET",
			url:                "maxprofit",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "{\"message\":\"begin param is missing: bad request\"}\n",
		},
		{
			name:               "Non GET request",
			method:             "POST",
			url:                "maxprofit?begin=1699228800&end=2699228800&symbol=UBER",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedBody:       "{\"message\":\"method POST not allowed: method not allowed\"}\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.handler
			if handler == nil {
				handler = StockPriceHandler{
					Controller: MockController{},
				}
			}

			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handlerF := http.HandlerFunc(handler.MaxProfitForPeriod)
			handlerF.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.True(t, rr.Header().Get("Content-Type") != "")
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
