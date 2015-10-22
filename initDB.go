package main

import (
	"database/sql"
	"log"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

func populateLocations(locationsMap *gorp.DbMap) {
	tuftsResQuad := initLocation("Tufts Res Quad", 42.408565, -71.121765)
	hodgkinsPark := initLocation("Hodgkins Park", 42.399566, -71.124595)
	danehyPark := initLocation("Danehy Park", 42.388306, -71.137507)
	ledermanPark := initLocation("Lederman Park", 42.363649, -71.071774)

	err := locationsMap.Insert(tuftsResQuad, hodgkinsPark,
		danehyPark, ledermanPark)

	if err != nil {
		log.Fatal(err)
	}
}

func initDB() *gorp.DbMap {
	db, err := sql.Open("postgres", "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	locationsTable :=
		dbMap.AddTableWithName(Location{}, "locations").SetKeys(true, "Id")

	locationsTable.ColMap("Name").SetNotNull(true).SetUnique(true)
	locationsTable.ColMap("Lat").SetNotNull(true)
	locationsTable.ColMap("Lng").SetNotNull(true)

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

	populateLocations(dbMap)

	return dbMap
}
