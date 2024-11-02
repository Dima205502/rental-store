package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"shop/models"
	"time"

	"github.com/go-playground/validator"
)

func NewShopManager() *ShopManager {
	return *ShopManager{}
}

func AddThing(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("AddThing", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		var thing models.Thing
		err := json.NewDecoder(r.Body).Decode(&thing)

		if err != nil {
			slog.Error("AddThing", slog.String("place", "Decoder"), slog.String("error", err.Error()))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = validator.New().Struct(thing)

		if err != nil {
			slog.Error("AddThing", slog.String("place", "validate"), slog.String("error", err.Error()))
			http.Error(w, "Error validate JSON", http.StatusBadRequest)
			return
		}

		err = s.manager.createThing(ctx, thing)

		if err != nil {
			slog.Error("AddThing", slog.String("place", "createThing"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
