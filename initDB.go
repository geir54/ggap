package main

import (
	"database/sql"
	"log"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

func initDB() *gorp.DbMap {
	db, err := sql.Open("postgres", "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	userTable := dbMap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	userTable.ColMap("Email").SetNotNull(true).SetUnique(true)
	userTable.ColMap("Username").SetNotNull(true).SetUnique(true)
	userTable.ColMap("Password").SetNotNull(true)
	userTable.ColMap("Salt").SetNotNull(true)

	err = dbMap.DropTablesIfExists() // TODO: Remove
	if err != nil {
		log.Fatal(err)
	}

	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	return dbMap
}
