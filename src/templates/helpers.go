package templates

import (
	"context"
	"log"
	"net/http"
	"time"

	"../models"
	"../universe"
)

const authCookieName = "SessionID"

type contextUserNameType struct{}

var contextUserNameKey = &contextUserNameType{}

func isAuthenticated(r *http.Request) *models.Session {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		log.Println(err)
		return nil
	}
	if cookie == nil {
		return nil
	}

	sess, err := universe.Get().Auth.CheckSession(cookie.Value)
	if err != nil {
		log.Println(err)
		return nil
	}
	return sess
}

func addAuthCookie(w http.ResponseWriter, id string) {
	cookie := http.Cookie{
		Name:    authCookieName,
		Value:   id,
		Expires: time.Now().Add(time.Hour),
	}
	http.SetCookie(w, &cookie)
}

func unsetAuthCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    authCookieName,
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &cookie)
}

func loginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := isAuthenticated(r)
		if sess == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		ctx := context.WithValue(r.Context(), contextUserNameKey, sess.UserLogin)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func notLoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := isAuthenticated(r)
		// TODO check back redirect
		if sess != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func passUserName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := isAuthenticated(r)
		if sess != nil {
			ctx := context.WithValue(r.Context(), contextUserNameKey, sess.UserLogin)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
