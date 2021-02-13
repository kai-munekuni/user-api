package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func success(w http.ResponseWriter, response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		internalServerError(w, "marshal error")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		log.Println(err)
		internalServerError(w, "io error")
		return
	}
}

func badRequest(w http.ResponseWriter, message, cause string) {
	data, err := json.Marshal(badRequestResponse{
		Message: message,
		Cause:   cause,
	})
	if err != nil {
		log.Println(err)
		internalServerError(w, "marshal error")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write(data); err != nil {
		log.Println(err)
		internalServerError(w, "io error")
		return
	}
}

type badRequestResponse struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

func unauthorized(w http.ResponseWriter) {
	httpError(w, http.StatusUnauthorized, "Authentication Faild")
}

func forbidden(w http.ResponseWriter) {
	httpError(w, http.StatusForbidden, "No Permission for Update")
}

func notFound(w http.ResponseWriter) {
	httpError(w, http.StatusNotFound, "No User Found")
}

func internalServerError(w http.ResponseWriter, message string) {
	httpError(w, http.StatusInternalServerError, message)
}

func httpError(w http.ResponseWriter, code int, message string) {
	data, err := json.Marshal(errorResponse{
		Message: message,
	})
	if err != nil {
		log.Println(err)
		internalServerError(w, "marshal error")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Println(err)
		internalServerError(w, "io error")
		return
	}
}

type errorResponse struct {
	Message string `json:"message"`
}
