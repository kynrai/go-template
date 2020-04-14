package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var _ Auth0Repo = (*auth0repo)(nil)

const (
	AUTO0_API_IDENTIFIER = "https://domain/v1"
	AUTH0_ISSUSER        = "https://domain.eu.auth0.com/"
	AUTH0_JWKS           = "https://domain.eu.auth0.com/.well-known/jwks.json"
	AUTH0_DOMAIN         = "https://domain.eu.auth0.com"
)

var (
	ErrInvalidAudience = errors.New("invalid audience")
	ErrInvalidIssuer   = errors.New("invalid issuer")
	ErrNoKey           = errors.New("Uuable to find appropriate key")
)

type Jwks struct {
	Keys []struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

type Auth0Repo interface {
	Validator
}

type auth0repo struct {
	client     *http.Client
	middleware *jwtmiddleware.JWTMiddleware
}

func NewAuth0() Auth0Repo {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	r := &auth0repo{
		client: client,
	}
	r.middleware = jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: r.validationKeyGetter,
			SigningMethod:       jwt.SigningMethodRS256,
		})
	return r
}

func (r *auth0repo) Validate(h http.Handler) http.Handler {
	return r.middleware.Handler(h)
}

func (r *auth0repo) validationKeyGetter(token *jwt.Token) (interface{}, error) {
	// Based off the guide found https://auth0.com/docs/quickstart/backend/golang
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(AUTO0_API_IDENTIFIER, false)
	if !checkAud {
		return token, ErrInvalidAudience
	}
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(AUTH0_ISSUSER, false)
	if !checkIss {
		return token, ErrInvalidIssuer
	}
	cert, err := r.getPemCert(token)
	if err != nil {
		return token, err
	}
	return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
}

func (r *auth0repo) getPemCert(token *jwt.Token) (string, error) {
	resp, err := r.client.Get(AUTH0_JWKS)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	jwks := Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return "", err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			return fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwks.Keys[k].X5c[0]), nil
		}
	}
	return "", ErrNoKey
}
