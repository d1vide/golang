package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	Access  string `json:"Access"`
	Refresh string `json:"Refresh"`
}

type Validator interface {
	SignPair(userID int64, email, role string) (TokenPair, error)
	ParseAccess(token string) (map[string]any, error)
	ParseRefresh(token string) (map[string]any, error)
	RevokeRefresh(jti string, exp int64)
	IsRefreshRevoked(jti string) bool
}

type key struct {
	id  string
	prv *rsa.PrivateKey
	pub *rsa.PublicKey
}

type RS256 struct {
	keys       []key
	accessTTL  time.Duration
	refreshTTL time.Duration
	blacklist  map[string]int64 // jti -> exp
}

func NewRS256() *RS256 {
	gen := func(id string) key {
		k, _ := rsa.GenerateKey(rand.Reader, 2048)
		return key{id: id, prv: k, pub: &k.PublicKey}
	}
	return &RS256{
		keys:       []key{gen("k1"), gen("k2")},
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
		blacklist:  map[string]int64{},
	}
}

func (r *RS256) SignPair(userID int64, email, role string) (TokenPair, error) {
	now := time.Now()
	jti := jwt.NewNumericDate(now).String()

	sign := func(ttl time.Duration, extra map[string]any) (string, error) {
		claims := jwt.MapClaims{
			"sub":   userID,
			"email": email,
			"role":  role,
			"iat":   now.Unix(),
			"exp":   now.Add(ttl).Unix(),
			"iss":   "pz10-auth",
			"aud":   "pz10-clients",
		}
		for k, v := range extra {
			claims[k] = v
		}

		k := r.keys[0]
		t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		t.Header["kid"] = k.id
		return t.SignedString(k.prv)
	}

	access, err := sign(r.accessTTL, nil)
	if err != nil {
		return TokenPair{}, err
	}
	refresh, err := sign(r.refreshTTL, map[string]any{"jti": jti, "typ": "refresh"})
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{Access: access, Refresh: refresh}, nil
}

func (r *RS256) parse(tokenStr string, refresh bool) (map[string]any, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		for _, k := range r.keys {
			if k.id == kid {
				return k.pub, nil
			}
		}
		return nil, jwt.ErrTokenUnverifiable
	},
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithAudience("pz10-clients"),
		jwt.WithIssuer("pz10-auth"),
	)
	if err != nil || !t.Valid {
		return nil, err
	}
	claims := t.Claims.(jwt.MapClaims)
	if refresh && claims["typ"] != "refresh" {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return map[string]any(claims), nil
}

func (r *RS256) ParseAccess(token string) (map[string]any, error)  { return r.parse(token, false) }
func (r *RS256) ParseRefresh(token string) (map[string]any, error) { return r.parse(token, true) }

func (r *RS256) RevokeRefresh(jti string, exp int64) {
	r.blacklist[jti] = exp
}

func (r *RS256) IsRefreshRevoked(jti string) bool {
	exp, ok := r.blacklist[jti]
	return ok && time.Now().Unix() < exp
}
