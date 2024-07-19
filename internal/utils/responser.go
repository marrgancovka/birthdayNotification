package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	resp, err := json.Marshal(MessageResponse{Message: err.Error()})
	if err != nil {
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(v)
	if err != nil {
		return
	}
	_, _ = w.Write(resp)
}

func ReadRequestData(r *http.Request, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}
	return nil
}
