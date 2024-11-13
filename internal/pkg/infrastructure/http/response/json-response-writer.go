package http_response

import (
	"encoding/json"
	"net/http"
)

type JsonResponseWriter struct {
}

func NewJsonResponseWriter() *JsonResponseWriter {
	return &JsonResponseWriter{}
}

func (jrw *JsonResponseWriter) WriteErrorResponse(w http.ResponseWriter, errors []error, httpStatus int, previousError error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	errorMessages := make([]string, len(errors))
	for i, err := range errors {
		errorMessages[i] = err.Error()
	}

	response := map[string]interface{}{
		"errors":        errorMessages,
		"previousError": previousError.Error(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode error response", http.StatusInternalServerError)
	}
}

func (jrw *JsonResponseWriter) WriteResponse(w http.ResponseWriter, payload interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
	}
}
