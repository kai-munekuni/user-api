package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kai-munekuni/user-api/internal/domain/repository"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		log.Printf("AuthHeader: %s\n", r.Header.Get("Authorization"))

		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			log.Print("bodyErr ", bodyErr.Error())
			internalServerError(w, "marshalerror")
			return
		}
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = rdr2
		log.Printf("BODY: %q", rdr1)
		next.ServeHTTP(w, r)
	})
}

type authMiddleware struct {
	userRepo repository.User
}

func (m authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth) != 2 || auth[0] != "Basic" {
			log.Println("invalid header")
			unauthorized(w)
			return
		}

		dec, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			log.Printf("%+v", err)
			unauthorized(w)
			return
		}
		tokens := strings.Split(string(dec), ":")

		if len(tokens) != 2 {
			log.Println("decoded token is invalid")
			unauthorized(w)
			return
		}
		userID := tokens[0]
		password := tokens[1]
		if err := m.userRepo.Authorize(ctx, userID, password); err != nil {
			log.Printf("%+v", err)
			log.Printf("fail to authorize userID=%s password=%s", userID, password)
			unauthorized(w)
			return
		}
		ctx = setUserID(ctx, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
