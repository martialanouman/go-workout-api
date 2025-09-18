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
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadIdParam(r *http.Request) (int64, error) {
	paramId := chi.URLParam(r, "id")
	if paramId == "" {
		return 0, errors.New("invalid id parameter")
	}

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func ReadPaginationParams(r *http.Request) (int32, int32, error) {
	qs := r.URL.Query()
	paramTake, paramSkip := qs.Get("take"), qs.Get("skip")

	if paramTake == "" && paramSkip == "" {
		return 5, 0, nil
	}

	take, err := strconv.ParseInt(paramTake, 10, 32)
	if err != nil {
		return 0, 0, errors.New("invalid take parameter")
	}

	skip, err := strconv.ParseInt(paramSkip, 10, 32)
	if err != nil {
		return 0, 0, errors.New("invalid skip parameter")
	}

	return int32(take), int32(skip), nil
}
