package main

import (
	"errors"
)

type recipe struct {
	IDRecipe    int    `json:"idRecipe,string"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ingredient struct {
	IDIngredient int     `json:"idIngredient,string"`
	IDRecipe     int     `json:"idRecipe,string"`
	Name         string  `json:"name"`
	Quantity     float32 `json:"quantity"`
	Unit         string  `json:"unit"`
}

func createIngredientDB(i ingredient) error {
	instruction := "INSERT INTO ingredients (IDRecipe,Name,Quantity,Unit) values($1,$2,$3,$4);"
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

func createRecipeDB(r recipe) error {
	instruction := "INSERT INTO recipe (Title,Description) values($1,$2);"
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

func deleteRecipe(id int) error {
	instruction := "DELETE FROM recipe WHERE idRecipe=$1;"
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

func deleteIngredientDB(id int, idrecipe int) error {
	instruction := "DELETE FROM ingredients WHERE idRecipe=$1 AND idIngredient=$2;"
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

func deleteIngredients(id int) error {
	instruction := "DELETE FROM ingredients WHERE idRecipe=$1;"
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

func getRecipeDB(idRecipe int) (*recipe, error) {
	instruction := "SELECT * FROM recipe WHERE idRecipe=$1;"
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

func getIngredientDB(idIngredient int) (*ingredient, error) {
	instruction := "SELECT * FROM ingredients WHERE idIngredient=$1;"
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

func getIngredientsDB(idRecipe int) ([]ingredient, error) {
	instruction := "SELECT idIngredient, name, quantity, unit FROM ingredients WHERE IDRecipe=$1;"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var idIngredient int
	var name string
	var quantity float32
	var unit string
	rows, errq := stmt.Query(idRecipe)
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

func getRecipesDB(direction string, id int) ([]recipe, error) {
	var instruction string
	if (direction != "next" && direction != "previous"){
		return nil, errors.New("error with the direction")
	}
	if direction == "next" {
		instruction = "SELECT idRecipe, title, description FROM recipe WHERE idRecipe > $1 LIMIT 10;"
	}else{
		instruction = "SELECT idRecipe, title, description FROM recipe WHERE idRecipe < $1 LIMIT 10;"
	}
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	var idRecipe int
	var title string
	var description string
	rows, errq := stmt.Query(id)
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

func updateRecipeDB(r recipe) error {
	instruction := "UPDATE recipe SET title=$1,description=$2 WHERE idRecipe=$3;"
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

func updateIngredientDB(i ingredient) error {
	instruction := "UPDATE ingredients SET name=$1,quantity=$2,unit=$3 WHERE idIngredient=$4;"
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

func deleteFullRecipeDB(idRecipe int) error {
	err := deleteIngredients(idRecipe)
	if err != nil {
		return err
	}
	err = deleteRecipe(idRecipe)
	if err != nil {
		return err
	}
	return nil
}

func searchRecipeDB(pattern string) ([]recipe, error) {
	answerR, err := searchRTitle(pattern)
	if err != nil{
		return nil, err
	}
	answerRDescription, errd := searchRDescription(pattern)
	if errd != nil{
		return answerR, errd
	}
	answerR = append(answerR, answerRDescription...)
	answerR = uniqueResults(answerR)
	return answerR, nil
}

func uniqueResults( items []recipe) []recipe{
	keys := make(map[recipe]bool)
	list := []recipe{}
	for _,entry := range items {
		if _,value := keys[entry]; !value{
			keys[entry] = true
			list = append(list,entry)
		}
	}
	return list
}

func searchRTitle(pattern string) ([]recipe, error){
	instruction := "SELECT * FROM recipe WHERE title LIKE $1;"
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

func searchRDescription(pattern string) ([]recipe, error){
	instruction := "SELECT * FROM recipe WHERE description LIKE $1;"
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

func countRecipesDB() (int,error){
	instruction := "SELECT count(*) FROM recipe;"
	stmt, err := db.Prepare(instruction)
	defer stmt.Close()
	if err != nil{
		return 0,err
	}
	var total int
	err = stmt.QueryRow().Scan(&total)
	if err != nil{
		return 0,err
	}
	return total, nil
}