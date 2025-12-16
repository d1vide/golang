package core

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/pz10-auth/internal/http/middleware"
	"example.com/pz10-auth/internal/platform/jwt"
	"github.com/go-chi/chi/v5"
)

type jwtSigner interface {
	SignPair(userID int64, email, role string) (jwt.TokenPair, error)
	ParseRefresh(string) (map[string]any, error)
	RevokeRefresh(string, int64)
	IsRefreshRevoked(string) bool
}

type UserInfo struct {
	ID    int64
	Email string
	Role  string
}

type userRepo interface {
	CheckPassword(email, pass string) (UserInfo, error)
}

type Service struct {
	repo userRepo
	jwt  jwtSigner
}

func NewService(r userRepo, j jwtSigner) *Service { return &Service{repo: r, jwt: j} }

func (s *Service) MeHandler(w http.ResponseWriter, r *http.Request) {
	// клеймы положим в контекст в AuthN-мидлваре
	claims, _ := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	jsonOK(w, map[string]any{
		"id": claims["sub"], "email": claims["email"], "role": claims["role"],
	})
}

func (s *Service) AdminStats(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]any{"users": 2, "version": "1.0"})
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.Email == "" || in.Password == "" {
		httpError(w, 400, "invalid_credentials")
		return
	}
	u, err := s.repo.CheckPassword(in.Email, in.Password)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}
	pair, _ := s.jwt.SignPair(u.ID, u.Email, u.Role)
	jsonOK(w, pair)
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var in struct{ Refresh string }
	if json.NewDecoder(r.Body).Decode(&in) != nil || in.Refresh == "" {
		httpError(w, 400, "bad_request")
		return
	}
	claims, err := s.jwt.ParseRefresh(in.Refresh)
	if err != nil {
		httpError(w, 401, "unauthorized")
		return
	}
	jti := claims["jti"].(string)
	exp := int64(claims["exp"].(float64))
	if s.jwt.IsRefreshRevoked(jti) {
		httpError(w, 401, "revoked")
		return
	}
	s.jwt.RevokeRefresh(jti, exp)
	pair, _ := s.jwt.SignPair(
		int64(claims["sub"].(float64)),
		claims["email"].(string),
		claims["role"].(string),
	)
	jsonOK(w, pair)
}

func (s *Service) UserByID(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.CtxClaimsKey).(map[string]any)
	sub := int64(claims["sub"].(float64))
	role := claims["role"].(string)

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if role == "user" && id != sub {
		httpError(w, 403, "forbidden")
		return
	}
	jsonOK(w, map[string]any{"id": id})
}
