package templates

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func createDataOnContext(ctx context.Context) map[string]interface{} {
	data := make(map[string]interface{})
	if userName, ok := ctx.Value(contextUserNameKey).(string); ok {
		data["UserName"] = userName
	}

	return data
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

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered: ", r, err)
				http.Error(w, "Server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type OutputElem struct {
	ID   int
	Data string
}

func jsonDataFormatter(w http.ResponseWriter, lts []models.LabeledTask) {
	data := make([]OutputElem, len(lts))

	for i, val := range lts {
		data[i].ID = val.OriginID
		data[i].Data = val.AnswerJSON
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func csvDataFormatter(w http.ResponseWriter, lts []models.LabeledTask) {
	writer := csv.NewWriter(w)
	writer.Write([]string{"ID", "Data"})

	for _, val := range lts {
		writer.Write([]string{strconv.Itoa(val.OriginID), val.AnswerJSON})
	}

	writer.Flush()
	return
}
