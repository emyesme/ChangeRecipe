package main

import (
	"database/sql"
	//"fmt"
	"log"

	_ "github.com/lib/pq"
)

func getDB(user string, pass string, database string) *sql.DB {
	//connection := fmt.Sprintf("postgres://%s:%s@127.0.0.1:5432/%s?sslmode=disable", user, pass, database)
	connection := "postgresql://root@localhost:26257/pruebarecetas?sslmode=disable"


	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	initDB(db)
	return db
}

func initDB(db *sql.DB) {
	instruction := `CREATE TABLE IF NOT EXISTS recipe(
						idRecipe SERIAL NOT NULL,
						title TEXT,
						description TEXT,
						PRIMARY KEY(idRecipe));
					CREATE TABLE IF NOT EXISTS ingredients(
						idIngredient SERIAL NOT NULL,
						idRecipe int NOT NULL,
						name TEXT,
						quantity FLOAT,
						unit TEXT,
						PRIMARY KEY(idIngredient),
						FOREIGN KEY (idRecipe) REFERENCES recipe(idRecipe));`

	_, err := db.Exec(instruction)
	if err != nil {
		log.Fatal(err)
	}
}
