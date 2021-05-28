package httpapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
	"net/http"
	"strings"
	"time"
)

func extractTokenFromAuthHeader(r *http.Request) (usecases.AuthToken, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != bearer {
		return "", errors.New("incorrect authorization header format")
	}
	return usecases.AuthToken(parts[1]), nil
}

func (a *HttpApi) extractIdFromHeader(r *http.Request) (model.UserId, error) {
	token, err := extractTokenFromAuthHeader(r)
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", nil
	}
	userId, err := a.usecases.Decode(token)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (a *HttpApi) extractAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := a.extractIdFromHeader(r)
		if err != nil || userId == "" {
			handler(w, r)
			return
			//TODO handle errors?
		}
		handler(w, r.WithContext(context.WithValue(r.Context(), contextKeyUserId, userId)))
	}
}

type responseWriterObserver struct {
	http.ResponseWriter
	status int
	wroteHeader bool
}

func (o *responseWriterObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseWriterObserver) StatusCode() int {
	if !o.wroteHeader {
		return http.StatusOK
	}
	return o.status
}

func (a *HttpApi) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		o := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		fmt.Printf("method: %s; url: %s; status-code: %d; remote-addr: %s; duration: %v;\n",
			r.Method, r.URL.String(), o.StatusCode(), r.RemoteAddr, time.Since(start))
	})
}