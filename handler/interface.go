package handler

import "net/http"

type Handler interface {
	MaxProfitForPeriod(w http.ResponseWriter, r *http.Request)
}
