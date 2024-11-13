package http_writer

import (
	"encoding/json"
	"net/http"
)

// JSONError is a generic structure for error responses.
type JSONError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// JSONResponse is a generic structure for success responses.
type JSONResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// WriteError writes an error response in JSON format.
func WriteError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonError := JSONError{
		Error:   http.StatusText(statusCode),
		Message: err.Error(),
	}

	json.NewEncoder(w).Encode(jsonError)
}

// WriteJSON writes a success response with a given status code and data.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResponse := JSONResponse{
		Status: http.StatusText(statusCode),
		Data:   data,
	}

	json.NewEncoder(w).Encode(jsonResponse)
}

// WriteOK writes a 200 OK response with data.
func WriteOK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, http.StatusOK)
}

// WriteCreated writes a 201 Created response with data.
func WriteCreated(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, data, http.StatusCreated)
}

// WriteNoContent writes a 204 No Content response with no body.
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
