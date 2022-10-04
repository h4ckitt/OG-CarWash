package authenticator

import (
	"car_wash/apperror"
	"car_wash/config"
	"car_wash/infra/mux/helper"
	"car_wash/repository"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
)

type Authenticator struct {
	firebaseAuth *auth.Client
	repo         repository.Repo
}

// Authentication MiddleWare

func (auth Authenticator) APIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")

		id, err := auth.repo.VerifyAPIKey(apiKey)

		if errors.Is(err, &apperror.NotFound) {
			err = &apperror.Unauthorized
		}

		if err != nil {
			helper.ReturnFailure(w, err)
			return
		}

		ctx := context.WithValue(context.Background(), "ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (auth Authenticator) BearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
			return
		}

		if strings.HasPrefix(bearerToken, "Bearer") {
			bearerToken = strings.TrimSpace(strings.Replace(bearerToken, "Bearer", "", 1))

			uid, err := auth.VerifyBearer(bearerToken)

			if err != nil {
				log.Println(err)
				helper.ReturnFailure(w, &apperror.Unauthorized)
				return
			}
			ctx := context.WithValue(context.Background(), "UID", uid)
			//ctx := context.WithValue(context.Background(), "UID", "4KzLpODPIfPIWL1BGZ23hiEUk9W2")
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "unauthorized"})
			return
		}
	})
}

func (auth Authenticator) VerifyBearer(token string) (string, error) {
	tk, err := auth.firebaseAuth.VerifyIDToken(context.Background(), token)

	if err != nil {
		return "", err
	}

	return tk.UID, nil
}

// Should Reimplement This Because VerifyAPIKey Is Not A Firebase Specific Function

func NewAuthenticator(repo repository.Repo) (*Authenticator, error) {
	opt := option.WithCredentialsFile(config.GetConfig().FirebaseConfig.ServiceFileName)

	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		return nil, err
	}

	authApp, err := app.Auth(context.Background())

	if err != nil {
		return nil, err
	}

	return &Authenticator{firebaseAuth: authApp, repo: repo}, nil
}
