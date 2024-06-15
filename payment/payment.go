package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
)

type paymentData struct {
	Amount int `json:"amount" validate:"required"`
}

func transferAmount(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerID"]
	var data paymentData
	if err := utils.ReadBody(w, r, &data); err != nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// send the request to customer service
	url := fmt.Sprintf("http://%s/customers/%s", customerUrl, customerID)
	resp, err := utils.SendRequest(r.Context(), http.MethodPut, url, payload)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("payment failed. got response: %s", b))
		return
	}

	utils.WriteResponse(w, http.StatusOK, data)
}
