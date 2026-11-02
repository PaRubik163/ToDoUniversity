package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"homework-reminder-bot/internal/repository"
)

var ErrCredentials = errors.New("invalid email or password")
var ErrRefresh = errors.New("invalid refresh token")
var ErrInvalidRegistration = errors.New("invalid email or password")
var ErrEmailExists = errors.New("email already registered")
var ErrAlreadyLinked = errors.New("Telegram account already linked")

type TokenPair struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	TokenType      string `json:"token_type"`
	ExpiresIn      int64  `json:"expires_in"`
	TelegramLinked bool   `json:"telegram_linked"`
}
type AuthService struct {
	repo   repository.AuthRepository
	secret []byte
	now    func() time.Time
}

func NewAuthService(repo repository.AuthRepository, secret string) *AuthService {
	return &AuthService{repo: repo, secret: []byte(secret), now: time.Now}
}
func (s *AuthService) Register(ctx context.Context, email, password string) (TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	parsedEmail, emailErr := mail.ParseAddress(email)
	if len(email) > 254 || emailErr != nil || parsedEmail.Address != email || len(password) < 8 || len(password) > 72 {
		return TokenPair{}, ErrInvalidRegistration
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, err
	}
	id, err := s.repo.Create(ctx, email, string(hash))
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return TokenPair{}, ErrEmailExists
		}
		return TokenPair{}, err
	}
	return s.issue(ctx, id)
}
func (s *AuthService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	a, err := s.repo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) != nil {
		return TokenPair{}, ErrCredentials
	}
	return s.issue(ctx, a.ID)
}
func (s *AuthService) Refresh(ctx context.Context, token string) (TokenPair, error) {
	if token == "" {
		return TokenPair{}, ErrRefresh
	}
	next, err := randomToken()
	if err != nil {
		return TokenPair{}, err
	}
	exp := s.now().Add(30 * 24 * time.Hour)
	id, err := s.repo.RotateRefresh(ctx, tokenHash(token), tokenHash(next), exp)
	if err != nil {
		return TokenPair{}, ErrRefresh
	}
	access, err := s.access(id)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: access, RefreshToken: next, TokenType: "Bearer", ExpiresIn: 900, TelegramLinked: id > 0}, nil
}
func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.RevokeRefresh(ctx, tokenHash(token))
}

func (s *AuthService) IssueLinkCode(ctx context.Context, accountID int64) (string, error) {
	if accountID > 0 {
		return "", ErrAlreadyLinked
	}
	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	code := strings.ToUpper(hex.EncodeToString(raw))
	if err := s.repo.CreateLinkCode(ctx, tokenHash(code), accountID, s.now().Add(10*time.Minute)); err != nil {
		return "", err
	}
	return code, nil
}

func (s *AuthService) LinkTelegramAccount(ctx context.Context, code string, telegramID int64) error {
	if telegramID <= 0 {
		return errors.New("invalid Telegram user")
	}
	return s.repo.LinkTelegramAccount(ctx, tokenHash(strings.ToUpper(strings.TrimSpace(code))), telegramID)
}
func (s *AuthService) issue(ctx context.Context, id int64) (TokenPair, error) {
	refresh, err := randomToken()
	if err != nil {
		return TokenPair{}, err
	}
	if err = s.repo.SaveRefresh(ctx, tokenHash(refresh), id, s.now().Add(30*24*time.Hour)); err != nil {
		return TokenPair{}, err
	}
	access, err := s.access(id)
	return TokenPair{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", ExpiresIn: 900, TelegramLinked: id > 0}, err
}
func (s *AuthService) access(id int64) (string, error) {
	enc := base64.RawURLEncoding
	head := enc.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, err := json.Marshal(map[string]any{"sub": fmt.Sprint(id), "iat": s.now().Unix(), "exp": s.now().Add(15 * time.Minute).Unix()})
	if err != nil {
		return "", err
	}
	body := head + "." + enc.EncodeToString(payload)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(body))
	return body + "." + enc.EncodeToString(mac.Sum(nil)), nil
}
func (s *AuthService) ValidateAccess(ctx context.Context, token string) (int64, error) {
	if len(token) > 4096 {
		return 0, ErrCredentials
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrCredentials
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, ErrCredentials
	}
	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if json.Unmarshal(headerBytes, &header) != nil || header.Algorithm != "HS256" || header.Type != "JWT" {
		return 0, ErrCredentials
	}
	body := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(body))
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
		return 0, ErrCredentials
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrCredentials
	}
	var claims struct {
		Sub string `json:"sub"`
		Iat int64  `json:"iat"`
		Exp int64  `json:"exp"`
	}
	now := s.now().Unix()
	if json.Unmarshal(b, &claims) != nil || claims.Exp <= now || claims.Iat > now+30 || claims.Exp <= claims.Iat || claims.Exp-claims.Iat > int64((15*time.Minute)/time.Second) {
		return 0, ErrCredentials
	}
	id, err := strconv.ParseInt(claims.Sub, 10, 64)
	if err != nil || id == 0 {
		return 0, ErrCredentials
	}
	exists, err := s.repo.AccountExists(ctx, id)
	if err != nil || !exists {
		return 0, ErrCredentials
	}
	return id, nil
}
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func tokenHash(v string) []byte { sum := sha256.Sum256([]byte(v)); return sum[:] }
