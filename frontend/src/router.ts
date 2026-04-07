import HomeView from "./views/HomeView.vue";
import UserView from "./views/UserView.vue";
import {createRouter, createWebHistory} from "vue-router";

const routes = [
    { path: '/', component: HomeView },
    { path: '/user/:userId', component: UserView },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router