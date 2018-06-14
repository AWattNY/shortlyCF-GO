package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
	fmt.Println(newURL)
	response := JSONResponse{Data: 1234}
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
	ExampleDB_Model()
	r := mux.NewRouter()
	r.HandleFunc("/", IndexEndPoint).Methods("GET")
	r.HandleFunc("/api/shorten", ShortenURLEndPoint).Methods("POST")
	r.HandleFunc("/slug", RedirectEndPoint).Methods("GET")
	r.HandleFunc("/stats/:slug/:statsParam", StatsEndPoint).Methods("GET")
	if err := http.ListenAndServe(":6060", r); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Id     int64
	Name   string
	Emails []string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

type Story struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *User
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Title, s.Author)
}

func ExampleDB_Model() {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	err = db.Insert(&User{
		Name:   "root",
		Emails: []string{"root1@root", "root2@root"},
	})
	if err != nil {
		panic(err)
	}

	story1 := &Story{
		Title:    "Cool story",
		AuthorId: user1.Id,
	}
	err = db.Insert(story1)
	if err != nil {
		panic(err)
	}

	// Select user by primary key.
	user := &User{Id: user1.Id}
	err = db.Select(user)
	if err != nil {
		panic(err)
	}

	// Select all users.
	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		panic(err)
	}

	// Select story and associated author in one query.
	story := new(Story)
	err = db.Model(story).
		Relation("Author").
		Where("story.id = ?", story1.Id).
		Select()
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
	fmt.Println(users)
	fmt.Println(story)
	// Output: User<1 admin [admin1@admin admin2@admin]>
	// [User<1 admin [admin1@admin admin2@admin]> User<2 root [root1@root root2@root]>]
	// Story<1 Cool story User<1 admin [admin1@admin admin2@admin]>>
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Story)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{})
		if err != nil {

			return err
		}
	}
	return nil
}
