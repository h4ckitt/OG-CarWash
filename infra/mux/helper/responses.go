package helper

import (
	"car_wash/apperror"
	"encoding/json"
	"errors"
	"net/http"
)

func ReturnFailure(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var appErr *apperror.AppError
	switch {
	case errors.As(err, &appErr):
		code, msg := appErr.StatusAndMessage()
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]string{"message": msg})
	}
}

func ReturnSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if msg, ok := data.(map[string]string); ok {
		json.NewEncoder(w).Encode(msg)
	} else {
		resp := struct {
			Message string `json:"message"`
			Data    any    `json:"data"`
		}{"success", data}

		json.NewEncoder(w).Encode(resp)
	}
}
