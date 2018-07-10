package database

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Url ...
type Url struct {
	LongURL  string    `json:"longURL"`
	ShortURL string    `json:"shortURL"`
	Date     time.Time `json:"date"`
}

// CreateSchema add comment here
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Url)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {

			return err
		}
	}
	return nil
}

// ConnectDB Add comment here ...
func ConnectDB() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr: "db:5432",
		User: "postgres",
	})
}
