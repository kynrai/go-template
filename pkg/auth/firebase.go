package auth

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

const defaultID = "project-id"

// ProjectID is the id of the Google Cloud project you are connected to.
// Useful when acquiring clients for google managed infrastructure (eg. Datastore, BigQuery), because they require it.
//
// ProjectID will be set to the Google Cloud Project ID with the following order:
// 1) Look for the GOOGLE_PROJECT_ID envar
// 2) Use the metadata API to get ID, this will only work in Google Cloud
// 3) Any failure or timeout (3s) will presume that the code is running outside the cloud
// in which case a default project ID is returned.
var ProjectID string

func init() {
	setProjectID()
}

func setProjectID() {
	if id := os.Getenv("GOOGLE_CLOUD_PROJECT"); id != "" {
		ProjectID = id
		return
	}

	if id := os.Getenv("GOOGLE_PROJECT_ID"); id != "" {
		ProjectID = id
		return
	}

	if id := os.Getenv("GCP_PROJECT"); id != "" {
		ProjectID = id
		return
	}

	if id, err := metadata.ProjectID(); err != nil {
		ProjectID = defaultID
	} else if id == "" {
		ProjectID = defaultID
	} else {
		ProjectID = id
	}
}

var _ FirebaseRepo = (*firebaseRepo)(nil)

type FirebaseRepo interface {
	Validator
}

type firebaseRepo struct {
	a *auth.Client
}

func NewFirebase() FirebaseRepo {

	config := &firebase.Config{ProjectID: ProjectID}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_CREDENTIALS")))

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	r := &firebaseRepo{
		a: client,
	}
	return r
}

func (f *firebaseRepo) Validate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract the token
		token := ""
		t, err := f.a.VerifyIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		log.Println(t.UID)
		h.ServeHTTP(w, r)
	})
}
