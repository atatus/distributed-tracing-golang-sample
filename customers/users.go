package main

import (
	"fmt"
	"net/http"

	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
)

type customer struct {
	ID           int64  `json:"id" validate:"-"`
	CustomerName string `json:"customer_name" validate:"required"`
	Account      string `json:"account" validate:"required"`
	Amount       int
}

type paymentData struct {
	Amount int `json:"amount" validate:"required"`
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var u customer
	if err := utils.ReadBody(w, r, &u); err != nil {
		return
	}

	ctx, span := tracer.Start(r.Context(), "create customer")
	defer span.End()
	id, err := db.InsertOne(ctx, datastore.InsertParams{
		Query: `INSERT INTO CUSTOMERS(CUSTOMER_NAME, ACCOUNT) VALUES (?, ?)`,
		Vars:  []interface{}{u.CustomerName, u.Account},
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("create customer error: %w", err))
		return
	}

	u.ID = id
	utils.WriteResponse(w, http.StatusCreated, u)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerID"]
	var u customer

	ctx, span := tracer.Start(r.Context(), "get customer")
	defer span.End()
	if err := db.SelectOne(ctx, datastore.SelectParams{
		Query:   `select ID, CUSTOMER_NAME, ACCOUNT, AMOUNT from CUSTOMERS where ID = ?`,
		Filters: []interface{}{customerID},
		Result:  []interface{}{&u.ID, &u.CustomerName, &u.Account, &u.Amount},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("get customer error: %w", err))
		return
	}

	utils.WriteResponse(w, http.StatusOK, u)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerID"]
	var data paymentData
	if err := utils.ReadBody(w, r, &data); err != nil {
		return
	}

	ctx, span := tracer.Start(r.Context(), "update customer amount")
	defer span.End()
	if err := db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update CUSTOMERS set AMOUNT = AMOUNT + ? where ID = ?`,
		Vars:  []interface{}{data.Amount, customerID},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("get customer error: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
