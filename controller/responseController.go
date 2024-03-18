package controller

import (
	"encoding/json"
	"net/http"
	"uts/model"
)

func responseMessage(w http.ResponseWriter, status int, message string) {
	var response model.Response
	response.Status = status
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendSuccessResponseWithData(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var response model.ResponseWithData
	response.Status = http.StatusOK
	response.Data = value
	json.NewEncoder(w).Encode(response)
}
