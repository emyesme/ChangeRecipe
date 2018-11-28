<template>
  <div class="update">
    <b-navbar toggleable="md" type="dark" variant="secondary" fixed="top">
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-navbar-brand to="/">Recetario</b-navbar-brand>
      <b-collapse is-nav id="nav_collapse">
        <b-navbar-nav class="ml-auto">
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <hr>
    <hr>
    <!-- fin navbar -->
    <b-container class="recipe">
      <b-card  bg-variant="dark" text-variant="white" :title="'Modificando Receta: '+ title">
        <b-container fluid>
          <b-form>
            <b-form-group id="title"
                          label="Titulo Receta:"
                          label-for="titleNew">
              <b-form-input id="titleNew"
                            type="text"
                            v-model=title>
              </b-form-input>
            </b-form-group>
            <b-form-group id="ingredients"
                          label="Ingredientes:">
              <b-row v-for="i in ingredients" v-bind:key="i.idIngredient" sm>
                <b-col sm>
                  <b-form-group id="name"
                                label="Nombre:"
                                label-for="nameNew">
                    <b-form-input id="nameNew"
                                type="text"
                                v-model=i.name>
                    </b-form-input>
                  </b-form-group>
                </b-col>
                <b-col sm>
                  <b-form-group id="unit"
                                label="Unidad:"
                                label-for="unitNew">
                    <b-form-input id="unitNew"
                                  type="text"
                                  v-model=i.unit>
                    </b-form-input>
                  </b-form-group>
                </b-col>
                <b-col sm>
                  <b-form-group id="quantity"
                                label="Cantidad:"
                                label-for="quantityNew">
                    <b-form-input id="quantityNew"
                                  type="text"
                                  v-model=i.quantity>
                    </b-form-input>
                  </b-form-group>
                </b-col>
                <b-col>
                  <b-form-group label="___">
                    <b-button variant="danger"
                              size="sm"
                              @click='deleteIngredient(i.idIngredient)'>
                    x
                    </b-button>
                  </b-form-group>
                </b-col>
              </b-row>
              <b-button @click="createIngredient=true" variant="success">Agregar Ingrediente</b-button>
            </b-form-group>
            <b-form-group id="description"
                          label="Descripcion:"
                          label-for="descriptionNew">
              <b-form-textarea id="descriptionNew"
                               v-model=description
                               :rows="3"
                               :max-rows="6">
              </b-form-textarea>
            </b-form-group>
            <b-alert :show="emptyInput" variant="danger" dismissible>Algunos elementos estan vacios</b-alert>
            <b-button @click="updateRecipe" variant="success">Guardar</b-button>
            <b-button to="/" variant="danger" @click="emptyInput=false">Cancelar</b-button>
          </b-form>
        </b-container>
      </b-card>
    </b-container>
    <!--  -->
    <b-modal ref="guardado" hide-footer>
      <div class="d-block text-center">
        <h3>Cambios guardados exitosamente.</h3>
      </div>
      <b-btn variant="primary" block @click="hide" to="/">Cerrar</b-btn>
    </b-modal>
    <!--  -->
    <b-modal v-model="createIngredient" title="Nuevo Ingrediente">
      <b-container fluid>
        <b-form v-if="createIngredient">
          <b-form-group id="name"
                        label="Nombre:"
                        label-for="nameNew">
            <b-form-input id="nameNew"
                          type="text"
                          placeholder="Nombre"
                          v-model=name>
            </b-form-input>
          </b-form-group>
          <b-form-group id="unit"
                        label="Unidad:"
                        label-for="unitNew">
            <b-form-input id="unitNew"
                          type="text"
                          placeholder="Unidad"
                          v-model=unit>
            </b-form-input>
          </b-form-group>
          <b-form-group id="quantity"
                        label="Cantidad:"
                        label-for="quantityNew">
            <b-form-input id="quantityNew"
                          type="text"
                          placeholder="Cantidad del Ingrediente"
                          v-model=quantity>
            </b-form-input>
          </b-form-group>
          <b-alert :show="emptyCreate" variant="danger" dismissible>
            Algunos elementos estan vacios.
          </b-alert>
          <b-button @click="addIngredient" variant="primary">Crear</b-button>
        </b-form>
      </b-container>
      <div slot="modal-footer" class="w-100">
        <b-btn size="sm" class="float-right" variant="danger" @click="closeCreate">Cancelar</b-btn>
      </div>
    </b-modal>
  </div>
</template>

<script>
import axios from "axios";
const api = "http://localhost:8083/api";

export default {
    name: 'update',
    data: () => ({
      emptyInput:false,
      emptyCreate: false,
      createIngredient: false,
      idRecipe : 0,
      title: null,
      description: null,
      name: null,
      unit: null,
      quantity: 0,
      idIngredient:0,
      ingredients: [],
    }),
    beforeMount(){
      this.infoRecipe()
      this.infoIngredients()
    },
    methods:{
      infoRecipe: function(){
        axios.get(api + '/recipes/' + this.$route.params.idRecipe)
        .then( response => {
          this.idRecipe = response.data.idRecipe;
          this.title = response.data.title;
          this.description = response.data.description;
        }).catch( function (error) { console.log(error)})
      },
      infoIngredients: function(){
        axios.get(api + '/recipes/' + this.$route.params.idRecipe + '/ingredients')
        .then( response => {
        this.ingredients = response.data;
        response.data.sort()
        }).catch( function (error) { console.log(error)})
      },
      addIngredient: function(){
        if(this.name == "" || this.unit == "" || this.quantity == 0){
          this.emptyCreate=true
          return
        }
        axios.post(api+"/recipe/"+this.$route.params.idRecipe,
        {idRecipe: this.idRecipe, name:this.name, quantity:Number(this.quantity), unit:this.unit})
        .catch( function (error) { console.log(error)})
        this.createIngredient=false
        this.name = ""
        this.unit = ""
        this.quantity = 0
        this.emptyCreate=false
        this.infoIngredients()
      },
      updateRecipe: function(){
        if(this.title == "" || this.description == ""){
          this.emptyInput = true
          return
        }
        axios.put(api +"/recipe/"+this.idRecipe ,
        { idRecipe: this.idRecipe, title: this.title, description:this.description})
        .catch( function (error) { console.log(error)})
        //update ingredients
        this.updateIngredient()
      },
      updateIngredient: function(){
        for(var i=0;i<this.ingredients.length;i++){
          if(this.ingredients[i].name == "" || this.ingredients[i].unit == "" || this.ingredients[i].quantity == 0){
            this.emptyInput=this.emptyInput || true
            return
          }
          axios.put(api+'/recipe/'+this.idRecipe+'/'+this.ingredients[i].idIngredient,
          { idIngredient:this.ingredients[i].idIngredient,
          idRecipe:this.idRecipe,
          name:this.ingredients[i].name,
          quantity:Number(this.ingredients[i].quantity),
          unit:this.ingredients[i].unit})
          .catch( function (error) { console.log(error)})
        }
        if(!this.emtpyInput){
          this.$refs.guardado.show()
          this.emptyInput=false
        }
      },
      hide:function(){
        this.$refs.guardado.hide()
      },
      deleteIngredient(id){
        axios.delete(api + '/recipe/'+this.$route.params.idRecipe+'/'+id)
        .catch( function (error) { console.log(error)})
        for( var i = 0; i < this.ingredients.length ;i++){
          if(this.ingredients[i].idIngredient == id){
            this.ingredients.splice(i,1)
          }
        }
      },
      closeCreate: function(){
        this.createIngredient=false
        this.emptyCreate=false
      }
    }
}
</script>
