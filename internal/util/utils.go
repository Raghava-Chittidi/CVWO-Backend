package util

import (
	"encoding/json"
	"net/http"
)

type ResponseJSON struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int, headers ...http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload = ResponseJSON {
		Error: true,
		Message: err.Error(),
	}

	WriteJSON(w, payload, statusCode)
}