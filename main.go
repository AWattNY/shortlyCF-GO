package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
)

type url struct {
	LongURL  string `json:"longURL"`
	ShortURL string `json:"shortURL"`
}

// JSONResponse with meta property
type JSONResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// IndexEndPoint Add description here.
func IndexEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome To Shortly !")
}

// ShortenURLEndPoint Add description here.
func ShortenURLEndPoint(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	urlToBeShortened := r.Form["url"][0]
	slug, _ := shortid.Generate()

	newURL := url{
		LongURL:  urlToBeShortened,
		ShortURL: slug,
	}
	// fmt.Fprintln(w, newURL)
	response := JSONResponse{Data: newURL}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// RedirectEndPoint Add description here.
func RedirectEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

// StatsEndPoint Add description here.
func StatsEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexEndPoint).Methods("GET")
	r.HandleFunc("/api/shorten", ShortenURLEndPoint).Methods("POST")
	r.HandleFunc("/slug", RedirectEndPoint).Methods("GET")
	r.HandleFunc("/stats/:slug/:statsParam", StatsEndPoint).Methods("GET")
	if err := http.ListenAndServe(":6060", r); err != nil {
		log.Fatal(err)
	}
}
