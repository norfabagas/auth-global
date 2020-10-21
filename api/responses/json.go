package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.WriteHeader(statusCode)
	resp := Response{
		Success: success,
		Message: message,
		Data:    data,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, false, http.StatusText(statusCode), struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, false, http.StatusText(http.StatusBadRequest), nil)
}
