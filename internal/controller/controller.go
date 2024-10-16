package controller

import (
	"auth_service/internal/models"
	"auth_service/internal/service"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func NewAuthManager(creater *service.UserManager, session *service.SessionManager, checker *service.EmailManager) *AuthManager {
	return &AuthManager{creater: creater, session: session, checker: checker}
}

func Signup(a *AuthManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("Signup", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			slog.Error("Signup", slog.String("place", "Decoder"), slog.String("error", err.Error()))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		err = a.creater.CreateUser(ctx, user)

		if err != nil {
			slog.Error("Signup", slog.String("place", "CreateSend"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Signin(a *AuthManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("Signin", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		var entryInfo models.EntryInfo

		err := json.NewDecoder(r.Body).Decode(&entryInfo)

		if err != nil {
			slog.Error("Signin", slog.String("place", "Decoder"), slog.String("error", err.Error()))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		token, err := a.session.CreateSession(ctx, entryInfo)

		if err != nil {
			slog.Error("Signin", slog.String("place", "CreateSession"), slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusCreated)
	}
}

func Logout(a *AuthManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("Logout", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		cookie, err := r.Cookie("session_token")

		if err != nil {
			slog.Error("Logout", slog.String("place", "Cookie read"), slog.String("error", err.Error()))
			http.Error(w, "Unread cookie", http.StatusBadRequest)
			return
		}

		err = a.session.DeleteSession(ctx, cookie.Name)

		if err != nil {
			slog.Error("Logout", slog.String("place", "DeleteSession"), slog.String("error", err.Error()))
			http.Error(w, "Сouldn't delete the session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)
	}
}

func VerifyEmail(a *AuthManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Error("VerifyEmail", slog.String("expected", http.MethodPost), slog.String("received", r.Method))
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		err := r.ParseForm()
		if err != nil {
			slog.Error("VerifyEmail", slog.String("place", "ParseForm"), slog.String("error", err.Error()))
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		token := r.FormValue("token")

		err = a.checker.CheckSend(ctx, token)

		if err != nil {
			slog.Error("VerifyEmail", slog.String("place", "CheckSend"), slog.String("error", err.Error()))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
