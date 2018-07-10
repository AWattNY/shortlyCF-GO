package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// CreateSchema add comment here
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*url)(nil)} {
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
