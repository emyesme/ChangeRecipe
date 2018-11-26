package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, pack interface{}) {
	response, _ := json.Marshal(pack)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

func respondErrorJSON(w http.ResponseWriter, r *http.Request, code int, message string) {
	respondWithJSON(w, r, code, map[string]string{"error": message})
}

func createRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var auxrecipe recipe
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&auxrecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = createRecipeDB(auxrecipe, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusCreated, auxrecipe)
}

func createIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	decoder := json.NewDecoder(r.Body)
	var auxingredient ingredient
	err := decoder.Decode(&auxingredient)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = createIngredientDB(auxingredient, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusCreated, auxingredient)
}

func getRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r) //retorna variables de la ruta
	id, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id")
		return
	}
	auxrecipe, errget := getRecipeDB(id, db)
	if errget != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, errget.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxrecipe)
}

func getIngredients(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	idRecipe, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id ")
		return
	}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count < 1 || count > 10 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	ingredients, err := getIngredientsDB(start, count, idRecipe, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, ingredients)
}

func deleteFullRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	idRecipe, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id")
		return
	}

	err = deleteFullRecipeDB(idRecipe, db)

	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "Recipe and ingredients successfully deleted"})

}

func deleteIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	idIngredient, err := strconv.Atoi(vars["idIngredient"])
	idRecipe, errRecipe := strconv.Atoi(vars["idRecipe"])

	if err != nil && errRecipe != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid ids")
		return
	}
	err = deleteIngredientDB(idIngredient, idRecipe, db)

	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "Ingredient successfully deleted"})

}

func updateRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	idRecipe, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id")
		return
	}
	var auxrecipe recipe
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&auxrecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadGateway, err.Error())
		return
	}
	auxrecipe.IDRecipe = idRecipe
	err = updateRecipeDB(auxrecipe, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxrecipe)
}

func updateIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	idIngredient, err := strconv.Atoi(vars["idIngredient"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid ingredient id")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var auxingredient ingredient
	err = decoder.Decode(&auxingredient)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadGateway, "Invalid request payload")
		return
	}
	auxingredient.IDIngredient = idIngredient
	err = updateIngredientDB(auxingredient, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxingredient)
}

func getRecipes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count < 1 || count > 100 {
		count = 100
	}
	if start < 0 {
		start = 0
	}
	recipes, err := getRecipesDB(start, count, db)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, recipes)
}

func getIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idIngredient"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid ingredient id")
		return
	}
	auxIngredient, errget := getIngredientDB(id, db)
	if errget != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, errget.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxIngredient)
}

func searchRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB){
	vars := mux.Vars(r)
	pattern := vars["pattern"]
	answersRecipes, err := searchRecipeDB(pattern,db)
	if err != nil{
		respondErrorJSON(w,r,http.StatusInternalServerError, "invalid pattern")
		return
	}
	respondWithJSON(w,r,http.StatusOK, answersRecipes)
}

func initializeRoutes(db *sql.DB) {

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Request-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api/recipe", func(w http.ResponseWriter, r *http.Request) { createRecipe(w, r, db) }).Methods("POST")
	router.HandleFunc("/api/recipe/{idRecipe}", func(w http.ResponseWriter, r *http.Request) { createIngredient(w, r, db) }).Methods("POST")
	router.HandleFunc("/api/recipes/{idRecipe}", func(w http.ResponseWriter, r *http.Request) { getRecipe(w, r, db) }).Methods("GET")
	router.HandleFunc("/api/searchRecipe/{pattern}", func(w http.ResponseWriter, r *http.Request) { searchRecipe(w, r, db) }).Methods("GET")
	router.HandleFunc("/api/recipes/{idRecipe}/ingredients", func(w http.ResponseWriter, r *http.Request) { getIngredients(w, r, db) }).Methods("GET")
	router.HandleFunc("/api/recipes/{idRecipe}/{idIngredient}", func(w http.ResponseWriter, r *http.Request) { getIngredient(w, r, db) }).Methods("GET")
	router.HandleFunc("/api/recipes", func(w http.ResponseWriter, r *http.Request) { getRecipes(w, r, db) }).Methods("GET")
	router.HandleFunc("/api/recipe/{idRecipe}", func(w http.ResponseWriter, r *http.Request) { deleteFullRecipe(w, r, db) }).Methods("DELETE")
	router.HandleFunc("/api/recipe/{idRecipe}/{idIngredient}", func(w http.ResponseWriter, r *http.Request) { deleteIngredient(w, r, db) }).Methods("DELETE")
	router.HandleFunc("/api/recipe/{idRecipe}", func(w http.ResponseWriter, r *http.Request) { updateRecipe(w, r, db) }).Methods("PUT")
	router.HandleFunc("/api/recipe/{idRecipe}/{idIngredient}", func(w http.ResponseWriter, r *http.Request) { updateIngredient(w, r, db) }).Methods("PUT")

	fmt.Println("Starting...8083")
	log.Fatal(http.ListenAndServe(":8083", handlers.CORS(headers, methods, origins)(router)))

}
