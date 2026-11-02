package main

import (
	"encoding/json"
	"net/http"

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
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	p, err := h.auth.Register(r.Context(), in.Email, in.Password)
	if err != nil {
		http.Error(w, err.Error(), 400)
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
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
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

func (h *authHandler) respond(w http.ResponseWriter, r *http.Request, p service.TokenPair) {
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: p.RefreshToken, Path: "/api/auth", HttpOnly: true, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https", SameSite: http.SameSiteStrictMode, MaxAge: 30 * 24 * 3600})
	p.RefreshToken = ""
	writeJSON(w, p)
}
