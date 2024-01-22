import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue";
import UserProfileView from "../views/UserProfileView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/home', component: HomeView },
		{ path: '/login', component: LoginView },
		{ path: '/user/:userID', component: UserProfileView },
	]
});

router.beforeEach((to, from, next) => {
	const isAuthenticated = !!localStorage.getItem('authToken');

	if (to.path !== '/login' && !isAuthenticated) {
		next('/login');
	} else if (to.path === '/login' && isAuthenticated) {
		next('/home');
	} else {
		next();
	}
});

export default router;
