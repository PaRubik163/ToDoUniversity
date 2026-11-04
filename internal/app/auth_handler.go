package app

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"homework-reminder-bot/internal/ratelimit"
	"homework-reminder-bot/internal/service"
)

type authHandler struct {
	auth              *service.AuthService
	limiter           ratelimit.Limiter
	trustProxyHeaders bool
}

func newAuthHandler(auth *service.AuthService, limiter ratelimit.Limiter, trustProxyHeaders bool) *authHandler {
	return &authHandler{auth: auth, limiter: limiter, trustProxyHeaders: trustProxyHeaders}
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	clientIP := requestIP(r, h.trustProxyHeaders)
	if !enforceRateLimits(h.limiter, w, r, rateLimitPolicy{key: "auth:register:ip:" + clientIP, limit: 10, window: time.Hour}) {
		return
	}
	var in credentials
	if decodeCredentials(w, r, &in) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	if !enforceRateLimits(h.limiter, w, r, rateLimitPolicy{key: "auth:register:email:" + strings.ToLower(strings.TrimSpace(in.Email)), limit: 3, window: time.Hour}) {
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
	clientIP := requestIP(r, h.trustProxyHeaders)
	if !enforceRateLimits(h.limiter, w, r, rateLimitPolicy{key: "auth:login:ip:" + clientIP, limit: 30, window: 15 * time.Minute}) {
		return
	}
	var in credentials
	if decodeCredentials(w, r, &in) != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	if !enforceRateLimits(h.limiter, w, r, rateLimitPolicy{key: "auth:login:email:" + strings.ToLower(strings.TrimSpace(in.Email)), limit: 8, window: 15 * time.Minute}) {
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
	if !validCookieRequestOrigin(r, h.trustProxyHeaders) {
		http.Error(w, "invalid request origin", http.StatusForbidden)
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
	setRefreshCookie(w, r, p.RefreshToken, 30*24*3600, h.trustProxyHeaders)
	p.RefreshToken = ""
	writeJSON(w, p)
}

func (h *authHandler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	if !validCookieRequestOrigin(r, h.trustProxyHeaders) {
		http.Error(w, "invalid request origin", http.StatusForbidden)
		return
	}
	if c, err := r.Cookie("refresh_token"); err == nil {
		if err := h.auth.Logout(r.Context(), c.Value); err != nil {
			http.Error(w, "could not revoke token", 500)
			return
		}
	}
	setRefreshCookie(w, r, "", -1, h.trustProxyHeaders)
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
	clientIP := requestIP(r, h.trustProxyHeaders)
	if !enforceRateLimits(h.limiter, w, r,
		rateLimitPolicy{key: "auth:link-code:ip:" + clientIP, limit: 20, window: 10 * time.Minute},
		rateLimitPolicy{key: "auth:link-code:account:" + strconv.FormatInt(accountID, 10), limit: 5, window: 10 * time.Minute},
	) {
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
	setRefreshCookie(w, r, p.RefreshToken, 30*24*3600, h.trustProxyHeaders)
	p.RefreshToken = ""
	writeJSON(w, p)
}

func setRefreshCookie(w http.ResponseWriter, r *http.Request, value string, maxAge int, trustProxyHeaders bool) {
	secure := r.TLS != nil || (trustProxyHeaders && strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https"))
	sameSite := http.SameSiteLaxMode
	if secure {
		// The Mini App runs in a Telegram WebView, which treats the site as a
		// third-party context. Strict cookies are omitted from refresh requests.
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token", Value: value, Path: "/api/auth", HttpOnly: true,
		Secure: secure, SameSite: sameSite, MaxAge: maxAge,
	})
}

func validCookieRequestOrigin(r *http.Request, trustProxyHeaders bool) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		if referer := strings.TrimSpace(r.Header.Get("Referer")); referer != "" {
			parsed, err := url.Parse(referer)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return false
			}
			origin = parsed.Scheme + "://" + parsed.Host
		} else {
			return false
		}
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil || parsed.Path != "" && parsed.Path != "/" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	scheme := "http"
	if r.TLS != nil || (trustProxyHeaders && strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https")) {
		scheme = "https"
	}
	return strings.EqualFold(parsed.Scheme, scheme) && strings.EqualFold(parsed.Host, r.Host)
}
