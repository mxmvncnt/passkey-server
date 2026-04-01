import HomeView from "./views/HomeView.vue";
import UserView from "./views/UserView.vue";
import {createRouter, createWebHistory} from "vue-router";
import {BASE_PATH} from "./config.ts";

const routes = [
    { path: '/', component: HomeView },
    { path: '/user/:userId', component: UserView },
]

const router = createRouter({
    history: createWebHistory(BASE_PATH || '/'),
    routes,
})

export default router