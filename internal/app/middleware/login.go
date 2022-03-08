package middleware

import (
	"context"
	"encoding/hex"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"net/http"
)

const authCookieName = "SessionID"

type UserCtx string

const UserCtxValue UserCtx = "userctx"

func Login(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user auth.User
			session, err := r.Cookie(authCookieName)

			if err != nil {
				user = auth.NewUser()
				cookie, err := createAuthCookie(user, secretKey)
				if err != nil {
					http.Error(w, "can't set cookie", http.StatusBadRequest)
				}
				http.SetCookie(w, cookie)
			} else {
				sessionID, err := hex.DecodeString(session.Value)
				if err != nil {
					http.Error(w, "can't decode cookie", http.StatusBadRequest)
				}

				user = auth.User{}
				if err = user.DecryptUserID([]byte(secretKey), sessionID); err != nil {
					http.Error(w, "can't decode cookie", http.StatusBadRequest)
				}
			}
			ctx := context.WithValue(r.Context(), UserCtxValue, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func createAuthCookie(user auth.User, secretKey string) (*http.Cookie, error) {
	sessionID, err := user.EncryptUserID([]byte(secretKey))

	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		HttpOnly: true,
		Name:     authCookieName,
		Value:    hex.EncodeToString(sessionID),
	}

	return cookie, nil
}
