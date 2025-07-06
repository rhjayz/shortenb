package middleware

import (
	"context"
	"net/http"
	"time"

	"shortb/config"
	"shortb/models"
)

type contextKey string

const userKey = contextKey("user")

func UsertoContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")

		if err == nil {
			var session models.Session
			if dbErr := config.DB.Preload("User").
				Where("token = ?", cookie.Value).
				First(&session).Error; dbErr == nil && session.ExpiredAt.After(time.Now()) {
				ctx := context.WithValue(r.Context(), userKey, session.User)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(r *http.Request) (models.User, bool) {
	user, ok := r.Context().Value(userKey).(models.User)
	return user, ok
}

func BeforeLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if ok && user.ID != 0 {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AfterLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok || user.ID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
