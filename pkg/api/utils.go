package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any) {
	//Trying to serialise data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error during making JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonData)

}
