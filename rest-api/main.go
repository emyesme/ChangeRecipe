package main

func main() {
	//user, password, database name
	db := getDB("postgres", "root", "pruebaReceta")
	defer db.Close()
	initializeRoutes(db)
}
