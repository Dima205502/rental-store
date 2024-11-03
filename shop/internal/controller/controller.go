package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"shop/internal/models"
	"shop/internal/service"
	"time"

	"github.com/go-playground/validator"
)

func NewShopManager(manager *service.ThingManagerService, getter *service.ThingGetterService) *ShopManager {
	return &ShopManager{manager: manager, getter: getter}
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

		err = s.manager.CreateThing(ctx, thing)

		if err != nil {
			slog.Error("AddThing", slog.String("place", "createThing"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func RemuveThing(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("RemuveThing", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		var rmvThng models.ThingId

		err := json.NewDecoder(r.Body).Decode(&rmvThng)

		if err != nil {
			slog.Error("RemuveThing", slog.String("place", "Decoder"), slog.String("error", err.Error()))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = validator.New().Struct(rmvThng)

		if err != nil {
			slog.Error("RemuveThing", slog.String("place", "validate"), slog.String("error", err.Error()))
			http.Error(w, "Error validate JSON", http.StatusBadRequest)
			return
		}

		nickCookie, err := r.Cookie("nickname")

		if err != nil {
			slog.Error("RemuveThing", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		err = s.manager.RemuveThing(ctx, nickCookie.Value, rmvThng.ThingId)

		if err != nil {
			slog.Error("RemuveThing", slog.String("place", "RemuveThing"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func BuyThing(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("BuyThing", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		var buyThng models.BuyThingRequest

		err := json.NewDecoder(r.Body).Decode(&buyThng)

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "Decoder"), slog.String("error", err.Error()))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = validator.New().Struct(buyThng)

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "validate"), slog.String("error", err.Error()))
			http.Error(w, "Error validate JSON", http.StatusBadRequest)
			return
		}

		nickCookie, err := r.Cookie("nickname")

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		emailCookie, err := r.Cookie("email")

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		err = s.manager.BuyThing(ctx, buyThng, nickCookie.Value, emailCookie.Value)

		if err != nil {
			slog.Error("BuyThing", slog.String("place", "BuyThing"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ShowAllThings(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("ShowAllThing", slog.String("expected", http.MethodGet), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		listThings, err := s.getter.GetAllThings(ctx)

		if err != nil {
			slog.Error("ShowAllThings", slog.String("place", "GetAllThings"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(listThings)
	}
}

func ShowRentalThings(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("ShowRentalThing", slog.String("expected", http.MethodGet), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		nickCookie, err := r.Cookie("nickname")

		if err != nil {
			slog.Error("ShowRentalThing", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		listThings, err := s.getter.GetRentalThings(ctx, nickCookie.Value)

		if err != nil {
			slog.Error("ShowRentalThings", slog.String("place", "GetRentalThings"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listThings)
	}
}

func ShowSaleThings(s *ShopManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("ShowSaleThing", slog.String("expected", http.MethodGet), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		nickCookie, err := r.Cookie("nickname")

		if err != nil {
			slog.Error("ShowSaleThing", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		listThings, err := s.getter.GetSaleThings(ctx, nickCookie.Value)

		if err != nil {
			slog.Error("ShowSaleThings", slog.String("place", "GetSaleThings"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listThings)
	}
}
