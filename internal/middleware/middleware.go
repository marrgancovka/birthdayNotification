package middleware

import (
	"birthdayNotification/internal/utils"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// CookieName represents the name of the JWT cookie.
const CookieName = "jwt-birthday-service"

// Auth is a middleware function that handles JWT authentication.
func Auth(next http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := cookie.Value

		claims, err := utils.ParseToken(token)
		if err != nil {
			log.Error("Error parsing token: ", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			log.Error("invalid token: ", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}
		if timeExp.Before(time.Now()) {
			log.Error("token is expired")
			utils.WriteError(w, http.StatusUnauthorized, errors.New("token is expired"))
			return
		}

		id, err := utils.ParseClaims(claims)
		if err != nil {
			log.Error("Error parsing claims: ", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, errors.New("error in parse token: "+err.Error()))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), CookieName, id))
		next.ServeHTTP(w, r)
	})
}
