package main

import (
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

func createRecipe(w http.ResponseWriter, r *http.Request) {
	var auxrecipe recipe
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&auxrecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = createRecipeDB(auxrecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusCreated, auxrecipe)
}

func createIngredient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var auxingredient ingredient
	err := decoder.Decode(&auxingredient)
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = createIngredientDB(auxingredient)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusCreated, auxingredient)
}

func getRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //retorna variables de la ruta
	id, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id")
		return
	}
	auxrecipe, errget := getRecipeDB(id)
	if errget != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, errget.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxrecipe)
}

func getIngredients(w http.ResponseWriter, r *http.Request) {
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
	ingredients, err := getIngredientsDB(start, count, idRecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, ingredients)
}

func deleteFullRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRecipe, err := strconv.Atoi(vars["idRecipe"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid recipe id")
		return
	}

	err = deleteFullRecipeDB(idRecipe)

	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "Recipe and ingredients successfully deleted"})

}

func deleteIngredient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idIngredient, err := strconv.Atoi(vars["idIngredient"])
	idRecipe, errRecipe := strconv.Atoi(vars["idRecipe"])

	if err != nil && errRecipe != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid ids")
		return
	}
	err = deleteIngredientDB(idIngredient, idRecipe)

	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "Ingredient successfully deleted"})

}

func updateRecipe(w http.ResponseWriter, r *http.Request) {
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
	err = updateRecipeDB(auxrecipe)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxrecipe)
}

func updateIngredient(w http.ResponseWriter, r *http.Request) {
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
	err = updateIngredientDB(auxingredient)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxingredient)
}

func getRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("bueno si envio")
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if page < 0{ 
		page = 0
	}
	if (err != nil){
		respondErrorJSON(w,r,http.StatusBadRequest, err.Error())
		return 
	}
	recipes, err := getRecipesDB(page)
	if err != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, recipes)
}

func getIngredient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idIngredient"])
	if err != nil {
		respondErrorJSON(w, r, http.StatusBadRequest, "Invalid ingredient id")
		return
	}
	auxIngredient, errget := getIngredientDB(id)
	if errget != nil {
		respondErrorJSON(w, r, http.StatusInternalServerError, errget.Error())
		return
	}
	respondWithJSON(w, r, http.StatusOK, auxIngredient)
}

func searchRecipe(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pattern := vars["pattern"]
	answersRecipes, err := searchRecipeDB(pattern)
	if err != nil{
		respondErrorJSON(w,r,http.StatusInternalServerError, "invalid pattern")
		return
	}
	respondWithJSON(w,r,http.StatusOK, answersRecipes)
}

func initializeRoutes() {

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Request-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api/recipe", createRecipe).Methods("POST")
	router.HandleFunc("/api/allRecipes/{page}", getRecipes).Methods("GET")
	router.HandleFunc("/api/recipes/{idRecipe}", getRecipe).Methods("GET")
	router.HandleFunc("/api/searchRecipe/{pattern}", searchRecipe).Methods("GET")
	router.HandleFunc("/api/recipe/{idRecipe}", deleteFullRecipe).Methods("DELETE")	
	router.HandleFunc("/api/recipe/{idRecipe}", updateRecipe).Methods("PUT")

	router.HandleFunc("/api/recipe/{idRecipe}", createIngredient).Methods("POST")
	router.HandleFunc("/api/recipes/{idRecipe}/ingredients", getIngredients).Methods("GET")
	router.HandleFunc("/api/recipes/{idRecipe}/{idIngredient}", getIngredient).Methods("GET")
	router.HandleFunc("/api/recipe/{idRecipe}/{idIngredient}", deleteIngredient).Methods("DELETE")
	router.HandleFunc("/api/recipe/{idRecipe}/{idIngredient}", updateIngredient).Methods("PUT")

	fmt.Println("Starting...8083")
	log.Fatal(http.ListenAndServe(":8083", handlers.CORS(headers, methods, origins)(router)))

}
