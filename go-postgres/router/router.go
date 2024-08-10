package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET")
	r.HandleFunc("/api/Stocks", middleware.GetAllStocks).Methods("GET" )
	r.HandleFunc("/api/NewStock", middleware.CreateStock).Methods("POST")
	r.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT" )
	r.HandleFunc("/api/deletestock", middleware.DeleteStocks).Methods("DELETE")

	return r

}
