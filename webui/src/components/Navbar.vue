<script>
import {computed} from "vue";
import {useRoute} from "vue-router";

export default {
	name: "Navbar",
	setup() {
		const route = useRoute();
		const isHome = computed(() => route.path === '/home');
		const isProfile = computed(() => route.path === `/user/${localStorage.getItem('userId')}`);

		return { isHome, isProfile };
	},
	data() {
		return {
			username: 'Utente'
		};
	},
	mounted() {
		this.username = localStorage.getItem('username');
	},
	methods: {
		logout() {
			localStorage.removeItem('userId');
			localStorage.removeItem('authToken');
			this.$router.push('/login');
		},
		showProfile() {
			this.$router.push(`/user/${localStorage.getItem('userId')}`);
		},
		goHome() {
			this.$router.push('/home');
		},
		goProfile() {
			this.$router.push(`/user/${localStorage.getItem('userId')}`);
		}
	}
}
</script>

<template>
	<nav class="navbar navbar-expand-lg navbar-dark sticky-top" style="background-color: rgba(45, 90, 110, 0.8); backdrop-filter: blur(10px);">
		<div class="container-fluid">
			<div class="navbar-brand" @click="goHome()">
				<img src="../assets/images/logo.png" alt="Wasaphoto Logo" height="45" class="d-inline-block align-text-top">
			</div>

			<div class="navbar-center-elements" style="position: absolute; left: 50%; transform: translateX(-50%);">
				<ul class="navbar-nav">
					<li class="nav-item">
						<div :class="['nav-link', isHome ? 'active' : '']" @click="goHome()">Stream</div>
					</li>
					<li class="nav-item">
						<div :class="['nav-link', isProfile ? 'active' : '']" @click="goProfile()">Profile</div>
					</li>
				</ul>
			</div>

			<ul class="navbar-nav ms-auto">
				<li class="nav-item dropdown card">
					<div class="nav-link dropdown-toggle justify-content-center align-items-center d-flex" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false" style="color: #1762cb">
						<i class="bi bi-person-circle"></i>
						<span class="ms-2">{{ username }}&nbsp;</span>
					</div>
					<ul class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
						<li class="dropdown-item" @click="showProfile()">Show profile</li>
						<li class="dropdown-item" @click="logout" style="color: red;">Logout</li>
					</ul>
				</li>
			</ul>
		</div>
	</nav>
</template>

<style>
.navbar-brand {
	cursor: pointer;
	padding: 0 !important;
	background-color: white !important;
	border-radius: 5px !important;
}
.nav-item {
	cursor: pointer !important;
}
.active {
	font-weight: bold;
	color: #fff;
}
</style>
