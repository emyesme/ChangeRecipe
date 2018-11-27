import Vue from 'vue'
import Router from 'vue-router'
import BootstrapVue from "bootstrap-vue"
import Recipes from '@/components/recipes'
import Update from '@/components/update'
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap-vue/dist/bootstrap-vue.css"

Vue.use(Router)
Vue.use(BootstrapVue)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Recipes',
      component: Recipes
    },
    {
      path:'/update/:idRecipe',
      name: 'Update',
      component: Update
    }
  ]
})
