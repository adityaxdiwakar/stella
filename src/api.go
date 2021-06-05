package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type StringResponse struct {
	Payload string `json:"payload"`
	Code    int    `json:"status_code"`
}

type PayloadResponse struct {
	Payload interface{} `json:"payload"`
	Code    int         `json:"status_code"`
}

func openAPI() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homePage)
	r.Get("/logs", logOutput)
	r.NotFound(notFoundHandler)

	http.ListenAndServe(fmt.Sprintf(":5862"), r)
}

func encode(data interface{}, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := PayloadResponse{
		Payload: data,
		Code:    statusCode,
	}
	str, _ := json.MarshalIndent(resp, "", "    ")
	str = bytes.Replace(str, []byte("\\u003c"), []byte("<"), -1)
	str = bytes.Replace(str, []byte("\\u003e"), []byte(">"), -1)
	str = bytes.Replace(str, []byte("\\u0026"), []byte("&"), -1)
	w.Write(str)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	encode("Welcome to the direct link to Stella", w, 200)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	encode("This is not a valid route, please refer to the documentation", w, 404)
}

func logOutput(w http.ResponseWriter, r *http.Request) {
	encode("Log output from Stella", w, 200)
}
