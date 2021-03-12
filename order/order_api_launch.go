package main

import (
	"net/http"
	"github.com/gorilla/mux"
    "log"
    "github.com/basic-e-commerce/db"
    "github.com/basic-e-commerce/controllers/order"
)

func main(){

	if err := db.ConnectOracle("admin/admin@localhost:1521/mydb"); err != nil {
		log.Println("Error initializing db connection:", err.Error())
		return
	}
	defer db.CloseOracle()

	InitializeServerParams()
	router := mux.NewRouter()

	router.HandleFunc("/users/{id}/orders", order.AllOrders)
	router.HandleFunc("/users/{id}/orders/{order_id}", order.Order)
	router.HandleFunc("/users/{id}/orders/create", order.CreateOrder).Methods("POST")
	router.HandleFunc("/users/{id}/orders/{order_id}/update", order.UpdateOrder).Methods("POST")
	router.HandleFunc("/users/{id}/orders/{order_id}/delete", order.DeleteOrder)
	router.HandleFunc("/users/{id}/delete/orders", order.DeleteAllOrders)

	// We can check if the port is already in use and generate an Alert for the developers or dev-ops engineers
	log.Fatal(http.ListenAndServe("9002", router))
}