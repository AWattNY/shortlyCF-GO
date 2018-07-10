package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"

	"github.com/go-pg/pg"

	"github.com/AWattNY/shortlyCF-GO/database"
	"github.com/asaskevich/govalidator"
)

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
	body := r.Form
	fmt.Printf("body  = %v \n", body)
	urlToBeShortened := r.FormValue("url")
	validURL := govalidator.IsURL(urlToBeShortened)
	fmt.Printf("validURL  = %v \n", validURL)
	if !validURL {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL, Must Submit a Valid URL\n"))
		return
	}

	fmt.Printf("urlToBeShortened  = %v \n", urlToBeShortened)
	slug, _ := shortid.Generate()
	date := time.Now()
	fmt.Printf("longURL  = %v \n", urlToBeShortened)
	newURL := database.Url{
		LongURL:  urlToBeShortened,
		ShortURL: slug,
		Date:     date,
	}
	fmt.Println(newURL)
	db := pg.Connect(&pg.Options{
		Addr: "db:5432",
		User: "postgres",
	})
	defer db.Close()

	err := db.Insert(&newURL)
	if err != nil {
		panic(err)
	}
	response := JSONResponse{Data: newURL}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// RedirectEndPoint Add description here.
func RedirectEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	db := pg.Connect(&pg.Options{
		Addr: "db:5432",
		User: "postgres",
	})
	defer db.Close()
	urlToBeLoaded, err := getlongURL(db, slug)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
		return
	}
	fmt.Printf("urlToBeLoaded = %#v \n", urlToBeLoaded)
	LongURL := urlToBeLoaded.LongURL
	fmt.Printf("LongURL  = %#v \n", LongURL)

	http.Redirect(w, r, LongURL, 302)
}

// GetlongURL Add Description here
func getlongURL(db *pg.DB, slug string) (*database.Url, error) {
	var url database.Url
	err := db.Model(&url).
		Where("url.short_url = ?", slug).
		Select()
	fmt.Printf("slug  = %v \n", slug)
	fmt.Printf("&url  = %v \n", &url)
	fmt.Printf("err  = %v \n", err)
	return &url, err
}

// StatsEndPoint Add description here.
func StatsEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func main() {
	db := database.ConnectDB()
	defer db.Close()

	err := database.CreateSchema(db)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", IndexEndPoint).Methods("GET")
	r.HandleFunc("/api/shorten", ShortenURLEndPoint).Methods("POST")
	r.HandleFunc("/{slug}", RedirectEndPoint).Methods("GET")
	r.HandleFunc("/stats/:slug/:statsParam", StatsEndPoint).Methods("GET")
	if err := http.ListenAndServe(":6060", r); err != nil {
		log.Fatal(err)
	}
}
