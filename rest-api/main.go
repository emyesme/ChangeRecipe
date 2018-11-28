package main

import(
	"database/sql"
)

var db *sql.DB

func main() {
	//user, port, database name
	db = getDB("root", "26257", "recipes")
	defer db.Close()
	initializeRoutes()
}
