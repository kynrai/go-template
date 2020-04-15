package auth

import (
	"context"

	"firebase.google.com/go/auth"
)

type (
	tKey int
)

var (
	tokenKey tKey
)

// NewContextWithToken returns a new ctx that carries a Token
func NewContextWithToken(ctx context.Context, token *auth.Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// TokenFromContext returns the Token stored in ctx, if any.
func TokenFromContext(ctx context.Context) (*auth.Token, bool) {
	t, ok := ctx.Value(tokenKey).(*auth.Token)
	return t, ok && t != nil
}
