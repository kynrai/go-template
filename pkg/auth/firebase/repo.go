package auth

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	Validator
}

type Validator interface {
	Validate(ctx context.Context, token string) (*auth.Token, error)
}

type repo struct {
	ac *auth.Client
}

func New() (Repo, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v\n", err)
	}
	r := &repo{
		ac: client,
	}
	return r, nil
}

func MustNew() Repo {
	client, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (r *repo) Validate(ctx context.Context, token string) (*auth.Token, error) {
	return r.ac.VerifyIDToken(ctx, token)
}
