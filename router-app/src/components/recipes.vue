<template>
  <div class="recipes">
    <b-navbar toggleable="md" type="dark" variant="secondary" fixed="top">
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-navbar-brand  @click="listRecipes('next',0)">Recetario</b-navbar-brand>
      <b-collapse is-nav id="nav_collapse">
        <b-navbar-nav>
          <b-nav-item @click="showAdd" variant="primary">Crear Receta</b-nav-item>
        </b-navbar-nav>
        <b-navbar-nav class="ml-auto">
          <b-nav-form>
            <b-form-input size="sm" class="mr-sm-2" type="text" v-model="pattern" placeholder="Buscar..."/>
            <b-button size="sm" class="my-2 my-sm-0" type="submit" @click="search">Buscar</b-button>
          </b-nav-form>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <hr>
    <hr>
    <!-- deberia haber propiedad para evitar los hr -->
    <b-modal v-model="showCreate" title="Crear Receta">
      <b-container fluid>
        <b-form v-if="showCreate">
          <b-form-group id="title"
                        label="Titulo Receta:"
                        label-for="titleNew">
            <b-form-input id="titleNew"
                          type="text"
                          placeholder="Titulo"
                          v-model=title>
            </b-form-input>
          </b-form-group>
          <b-form-group id="description"
                        label="Descripcion:"
                        label-for="descriptionNew">
            <b-form-textarea id="descriptionNew"
                     placeholder="Ingrese la descripcion de su receta"
                     v-model="description"
                     :rows="3"
                     :max-rows="6">
            </b-form-textarea>
          </b-form-group>
          <b-alert :show="emptyCreate" variant="danger" dismissible>
            Algunos elementos estan vacios.
          </b-alert>
          <b-button @click="addRecipe" variant="primary">Crear</b-button>
        </b-form>
      </b-container>
      <div slot="modal-footer" class="w-100">
        <b-btn size="sm" class="float-right" variant="danger" @click="closeCreate">Cancelar</b-btn>
      </div>
    </b-modal>
    <!-- lista de recetas -->
    <b-container class="recipes-list">
      <b-row class="text-center" v-for="recipe in recipes" v-bind:key="recipe.idRecipe" >
        <b-col sm>
          <b-card  bg-variant="dark" text-variant="white" :title="recipe.title">
            <p class="card-text">
              {{ recipe.description }}
            </p>
            <b-button @click='showRecipe(recipe.idRecipe)' variant="secondary">Mostrar</b-button>
            <b-button to="/" @click='deleteRecipe(recipe.idRecipe)' variant="danger">Eliminar</b-button>
          </b-card>
        </b-col>
      </b-row>
      <center>
        <b-pagination @input="changePage" align="center" size="md" :total-rows="total" v-model="page" :per-page="10">
        </b-pagination>
      </center>
    </b-container>
    <!-- modal receta -->
    <b-modal v-model="showInfo" title="Receta">
      <b-container fluid>
        <b-form v-if="showInfo">
          <h4>{{ title }}</h4>
          <b-list-group v-for="i in ingredients" v-bind:key="i.idIngredient" >
            <b-list-group-item>
              {{ i.quantity }}  ({{ i.unit }})  {{ i.name }}
            </b-list-group-item>
          </b-list-group>
          <hr>
          <p> <strong> Descripción: </strong></p>
          <p> {{ description }}</p>
          <b-button :to="'/update/'+idRecipe" variant="primary">
          Modificar
          </b-button>
        </b-form>
      </b-container>
      <div slot="modal-footer" class="w-100">
        <b-btn size="sm" class="float-right" variant="danger" @click="showInfo=false">Cerrar</b-btn>
      </div>
    </b-modal>
  </div>
</template>

<script>
import axios from 'axios';
const api='http://localhost:8083/api';
export default {

  name: 'recipes',
  data: () => ({
    total:0,
    previouspage:1,
    page: 1,
    emptyCreate: false,
    showCreate: false,
    showInfo: false,
    pattern:null,
    idRecipe: 0,
    title: null,
    description: null,
    idIngredient: 0,
    name: null,
    quantity: 0,
    unit: null,
    recipes: [],
    ingredients: [],
  }),
  beforeMount(){
    this.totalRecipes()
    this.listRecipes("next",0)
  },
  methods:{
    changePage: function(){
      if(this.page == this.previouspage){
        return
      }
      if( this.page > this.previouspage){
        this.page = this.previouspage + 1
        this.listRecipes("next",this.recipes[this.recipes.length-1].idRecipe)
      }
      else{
        this.page = 1
        this.listRecipes("previous",this.recipes[0].idRecipe)
      }
      this.previouspage = this.page
    },
    totalRecipes: function(){
      axios.get(api+'/total')
      .then( response => { this.total = response.data})
      .catch( function(error) { console.log(error)})
    },
    listRecipes: function(direction,id){
      axios.get(api + '/allRecipes/'+direction+'/'+id)
      .then(response => {
      this.recipes = response.data})
      .catch( function (error) { console.log(error)})
    },
    deleteRecipe: function(id){
      axios.delete(api + '/recipe/'+ id)
      .catch( function(error) { console.log(error)})
      for( var i = 0; i < this.recipes.length ;i++){
        if(this.recipes[i].idRecipe == id){
          this.recipes.splice(i,1)
        }
      }
      this.previouspage = 1
      this.page = 1
    },
    showRecipe: function(id){
      axios.get(api + '/recipes/' + id)
      .then( response => {
        this.idRecipe = response.data.idRecipe;
        this.title = response.data.title;
        this.description = response.data.description;
      }).catch( function (error) { console.log(error)})
      //get info ingredients
      axios.get(api + '/recipes/' + id + '/ingredients')
      .then( response => {
        this.ingredients = response.data;
      }).catch( function (error) { console.log(error)})
      this.showInfo=true
    },
    addRecipe: function(){
      //create recipe
      if ( this.title == "" || this.description == ""){
        this.emptyCreate=true
        return
      }
      axios.post(api + "/recipe" ,
      { title: this.title, description:this.description})
      .catch( function (error){ console.log(error)})
      this.emptyCreate=false
      this.showCreate=false
      this.totalRecipes()
      this.listRecipes("next",0)
      this.previouspage = 1
      this.page = 1
    },
    showAdd: function(){
      this.title=""
      this.description=""
      this.showCreate=true
    },
    closeCreate: function(){
      this.showInfo=false
      this.emptyCreate=false
    },
    search: function(){
      axios.get(api+'/searchRecipe/'+this.pattern)
      .then( response => {
        this.recipes = response.data;
      }).catch( function(error){ console.log(error)})
      this.pattern=""
    }
  }
}
</script>


