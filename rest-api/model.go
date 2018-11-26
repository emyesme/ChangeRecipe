package main

import (
	"database/sql"
	"errors"
)

type recipe struct {
	IDRecipe    int    `json:"idRecipe,string"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ingredient struct {
	IDIngredient int     `json:"idIngredient"`
	IDRecipe     int     `json:"idRecipe,string"`
	Name         string  `json:"name"`
	Quantity     float32 `json:"quantity"`
	Unit         string  `json:"unit"`
}

func createIngredientDB(i ingredient, db *sql.DB) error {
	instruction := "INSERT INTO ingredients (IDRecipe,Name,Quantity,Unit) values($1,$2,$3,$4)"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(i.IDRecipe, i.Name, i.Quantity, i.Unit)
	if err != nil {
		return err
	}
	return nil
}

func createRecipeDB(r recipe, db *sql.DB) error {
	instruction := "INSERT INTO recipe (Title,Description) values($1,$2)"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	result, err := stmt.Exec(r.Title, r.Description)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("More that i row affected")
	}
	return nil
}

func deleteRecipe(id int, db *sql.DB) error {
	instruction := "DELETE FROM recipe WHERE idRecipe=$1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("Do not exist the recipe")
	}
	return nil
}

func deleteIngredientDB(id int, idrecipe int, db *sql.DB) error {
	instruction := "DELETE FROM ingredients WHERE idRecipe=$1 AND idIngredient=$2"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	result, err := stmt.Exec(idrecipe, id)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return errors.New("Do not exist the ingredient")
	}
	return nil
}

func deleteIngredients(id int, db *sql.DB) error {
	instruction := "DELETE FROM ingredients WHERE idRecipe=$1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func getRecipeDB(idRecipe int, db *sql.DB) (*recipe, error) {
	instruction := "SELECT * FROM recipe WHERE idRecipe=$1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	var id int
	var titleGet string
	var description string
	if err != nil {
		return nil, err

	}
	err = stmt.QueryRow(idRecipe).Scan(&id, &titleGet, &description)
	if err != nil {
		return nil, err

	}
	auxrecipe := recipe{
		id,
		titleGet,
		description,
	}
	return &auxrecipe, nil
}

func getIngredientDB(idIngredient int, db *sql.DB) (*ingredient, error) {
	instruction := "SELECT * FROM ingredients WHERE idIngredient=$1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var id int
	var idRecipe int
	var name string
	var quantity float32
	var unit string
	err = stmt.QueryRow(idIngredient).Scan(&id, &idRecipe, &name, &quantity, &unit)
	if err != nil {
		return nil, err
	}
	auxingredient := ingredient{
		idIngredient,
		idRecipe,
		name,
		quantity,
		unit,
	}
	return &auxingredient, nil
}

func getIngredientsDB(start int, count int, idRecipe int, db *sql.DB) ([]ingredient, error) {
	instruction := "SELECT idIngredient, name, quantity, unit FROM ingredients WHERE IDRecipe=$1 LIMIT $2 offset $3"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var idIngredient int
	var name string
	var quantity float32
	var unit string
	rows, errq := stmt.Query(idRecipe, count, start)
	if errq != nil {
		return nil, err
	}
	ingredients := []ingredient{}
	for rows.Next() {
		var auxingredient ingredient
		err = rows.Scan(&idIngredient, &name, &quantity, &unit)
		if err != nil {
			return nil, err
		}
		auxingredient = ingredient{
			idIngredient,
			idRecipe,
			name,
			quantity,
			unit,
		}
		ingredients = append(ingredients, auxingredient)
	}
	return ingredients, nil
}

func getRecipesDB(start int, count int, db *sql.DB) ([]recipe, error) {
	instruction := "SELECT idRecipe,title,description FROM recipe LIMIT $1 offset $2"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var idRecipe int
	var title string
	var description string
	rows, errq := stmt.Query(count, start)
	if errq != nil {
		return nil, err
	}
	recipes := []recipe{}
	for rows.Next() {
		var auxrecipe recipe
		err = rows.Scan(&idRecipe, &title, &description)
		if err != nil {
			return nil, err
		}
		auxrecipe = recipe{
			idRecipe,
			title,
			description,
		}
		recipes = append(recipes, auxrecipe)
	}
	return recipes, nil
}

func updateRecipeDB(r recipe, db *sql.DB) error {
	instruction := "UPDATE recipe SET title=$1,description=$2 WHERE idRecipe=$3"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(r.Title, r.Description, r.IDRecipe)
	if err != nil {
		return err
	}
	return nil
}

func updateIngredientDB(i ingredient, db *sql.DB) error {
	instruction := "UPDATE ingredients SET name=$1,quantity=$2,unit=$3 WHERE idIngredient=$4"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(i.Name, i.Quantity, i.Unit, i.IDIngredient)
	if err != nil {
		return err
	}
	return nil
}

func deleteFullRecipeDB(idRecipe int, db *sql.DB) error {
	err := deleteIngredients(idRecipe, db)
	if err != nil {
		return err
	}
	err = deleteRecipe(idRecipe, db)
	if err != nil {
		return err
	}
	return nil
}

func searchRecipeDB(pattern string, db *sql.DB) ([]recipe, error) {
	answerR, err := searchRTitle(pattern,db)
	if err != nil{
		return nil, err
	}
	answerRDescription, errd := searchRDescription(pattern, db)
	if errd != nil{
		return answerR, errd
	}
	answerR = append(answerR, answerRDescription...)
	return answerR, nil
}

func searchRTitle(pattern string, db *sql.DB) ([]recipe, error){
	instruction := "SELECT * FROM recipe WHERE title LIKE $1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil{
		return nil,err
	}
	var idRecipe int
	var title string
	var description string
	rows, errq := stmt.Query(pattern)
	if errq != nil {
		return nil, err
	}
	recipes := []recipe{}
	for rows.Next() {
		var auxrecipe recipe
		err = rows.Scan(&idRecipe, &title, &description)
		if err != nil {
			return nil, err
		}
		auxrecipe = recipe{
			idRecipe,
			title,
			description,
		}
		recipes = append(recipes, auxrecipe)
	}
	return recipes, nil
}

func searchRDescription(pattern string, db *sql.DB) ([]recipe, error){
	instruction := "SELECT * FROM recipe WHERE description LIKE $1"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil{
		return nil,err
	}
	var idRecipe int
	var title string
	var description string
	rows, errq := stmt.Query(pattern)
	if errq != nil {
		return nil, err
	}
	recipes := []recipe{}
	for rows.Next() {
		var auxrecipe recipe
		err = rows.Scan(&idRecipe, &title, &description)
		if err != nil {
			return nil, err
		}
		auxrecipe = recipe{
			idRecipe,
			title,
			description,
		}
		recipes = append(recipes, auxrecipe)
	}
	return recipes, nil
}