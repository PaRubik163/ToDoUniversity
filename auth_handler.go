package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	"homework-reminder-bot/internal/service"
)

type authHandler struct{ auth *service.AuthService }

func newAuthHandler(auth *service.AuthService) *authHandler { return &authHandler{auth: auth} }

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var in credentials
	if decodeCredentials(w, r, &in) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	p, err := h.auth.Register(r.Context(), in.Email, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRegistration):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrEmailExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Printf("register account: %v", err)
			http.Error(w, "could not create account", http.StatusInternalServerError)
		}
		return
	}
	h.respond(w, r, p)
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var in credentials
	if decodeCredentials(w, r, &in) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	p, err := h.auth.Login(r.Context(), in.Email, in.Password)
	if err != nil {
		http.Error(w, "invalid email or password", 401)
		return
	}
	h.respond(w, r, p)
}

func decodeCredentials(w http.ResponseWriter, r *http.Request, in *credentials) error {
	mediaType, _, mediaTypeErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaTypeErr != nil || mediaType != "application/json" {
		return errors.New("content type must be application/json")
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(in); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return errors.New("request must contain one JSON object")
	}
	return nil
}

func (h *authHandler) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	c, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token required", 401)
		return
	}
	p, err := h.auth.Refresh(r.Context(), c.Value)
	if err != nil {
		http.Error(w, "invalid refresh token", 401)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: p.RefreshToken, Path: "/api/auth", HttpOnly: true, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https", SameSite: http.SameSiteStrictMode, MaxAge: 30 * 24 * 3600})
	p.RefreshToken = ""
	writeJSON(w, p)
}

func (h *authHandler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	if c, err := r.Cookie("refresh_token"); err == nil {
		if err := h.auth.Logout(r.Context(), c.Value); err != nil {
			http.Error(w, "could not revoke token", 500)
			return
		}
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/api/auth", HttpOnly: true, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https", SameSite: http.SameSiteStrictMode, MaxAge: -1})
	w.WriteHeader(http.StatusNoContent)
}

func (h *authHandler) linkCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accountID, err := h.auth.ValidateAccess(r.Context(), strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	code, err := h.auth.IssueLinkCode(r.Context(), accountID)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyLinked) {
			http.Error(w, "Telegram account is already linked", http.StatusConflict)
			return
		}
		log.Printf("create account link code: %v", err)
		http.Error(w, "could not create link code", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, struct {
		Code      string `json:"code"`
		ExpiresIn int    `json:"expires_in"`
	}{Code: code, ExpiresIn: 600})
}

func (h *authHandler) respond(w http.ResponseWriter, r *http.Request, p service.TokenPair) {
	w.Header().Set("Cache-Control", "no-store")
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: p.RefreshToken, Path: "/api/auth", HttpOnly: true, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https", SameSite: http.SameSiteStrictMode, MaxAge: 30 * 24 * 3600})
	p.RefreshToken = ""
	writeJSON(w, p)
}
