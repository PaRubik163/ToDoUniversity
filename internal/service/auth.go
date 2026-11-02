package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"homework-reminder-bot/internal/repository"
)

var ErrCredentials = errors.New("invalid email or password")
var ErrRefresh = errors.New("invalid refresh token")

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
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
	if len(email) > 254 || !strings.Contains(email, "@") || len(password) < 8 || len(password) > 72 {
		return TokenPair{}, errors.New("email is invalid or password must contain 8 to 72 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return TokenPair{}, err
	}
	id, err := s.repo.Create(ctx, email, string(hash))
	if err != nil {
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
	return TokenPair{access, next, "Bearer", 900}, nil
}
func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.repo.RevokeRefresh(ctx, tokenHash(token))
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
	return TokenPair{access, refresh, "Bearer", 900}, err
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
func (s *AuthService) ValidateAccess(token string) (int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
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
		Exp int64  `json:"exp"`
	}
	if json.Unmarshal(b, &claims) != nil || claims.Exp <= s.now().Unix() {
		return 0, ErrCredentials
	}
	var id int64
	if _, err = fmt.Sscan(claims.Sub, &id); err != nil || id >= 0 {
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
