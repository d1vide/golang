package service

import "strings"

const (
	demoToken    = "demo-token"
	demoUsername = "student"
	demoPassword = "student"
)

type AuthService struct{}

func New() *AuthService { return &AuthService{} }

func (s *AuthService) Login(username, password string) (string, bool) {
	if username == demoUsername && password == demoPassword {
		return demoToken, true
	}
	return "", false
}

func (s *AuthService) Verify(authHeader string) (string, bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == demoToken {
		return demoUsername, true
	}
	return "", false
}
