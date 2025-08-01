package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	if data == nil {
		w.WriteHeader(status)
		return nil
	}

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return err
	}
	jsonData = append(jsonData, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
	return nil
}

func ReadIdParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("id parameter is required")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
