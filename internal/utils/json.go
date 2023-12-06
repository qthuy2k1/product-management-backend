package utils

import (
	"encoding/json"
	"net/http"
)

// RenderJson marshals the given data and write the json to the connection
func RenderJson(w http.ResponseWriter, data interface{}, code int) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}
