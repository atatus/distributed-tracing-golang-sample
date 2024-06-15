package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
)

type orderData struct {
	ID          int64  `json:"id"`
	CustomerID  int    `json:"customer_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}

type customer struct {
	ID           int64  `json:"id"`
	CustomerName string `json:"customer_name"`
	Account      string `json:"account"`
	Amount       int
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var request orderData
	if err := utils.ReadBody(w, r, &request); err != nil {
		return
	}

	// get customer details from customer service
	url := fmt.Sprintf("http://%s/customers/%d", customerUrl, request.CustomerID)
	customerResponse, err := utils.SendRequest(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Printf("%v", err)
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	b, err := ioutil.ReadAll(customerResponse.Body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer customerResponse.Body.Close()

	if customerResponse.StatusCode != http.StatusOK {
		utils.WriteErrorResponse(w, customerResponse.StatusCode, fmt.Errorf("payment failed. got response: %s", b))
		return
	}

	var customer customer
	if err := json.Unmarshal(b, &customer); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// basic check for the customer balance
	if customer.Amount < request.Price {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("insufficient balance. add %d more amount to account", request.Price-customer.Amount))
		return
	}

	// insert the order into order table
	ctx, insertSpan := tracer.Start(r.Context(), "insert order")
	id, err := db.InsertOne(ctx, datastore.InsertParams{
		Query: `insert into ORDERS(ACCOUNT, PRODUCT_NAME, PRICE, ORDER_STATUS) VALUES (?,?,?, ?)`,
		Vars:  []interface{}{customer.Account, request.ProductName, request.Price, "SUCCESS"},
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		insertSpan.End()
		return
	}
	insertSpan.End()

	// update the pending amount in customer table
	ctx, updateSpan := tracer.Start(r.Context(), "update customer amount")
	if err := db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update CUSTOMERS set AMOUNT = AMOUNT - ? where ID = ?`,
		Vars:  []interface{}{request.Price, customer.ID},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		updateSpan.End()
		return
	}
	updateSpan.End()

	// send response
	response := request
	response.ID = id
	utils.WriteResponse(w, http.StatusCreated, response)
}
