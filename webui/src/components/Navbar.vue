<script>
export default {
	name: "Navbar",
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
			this.$router.push('/user/' + localStorage.getItem('userId'));
		},
		goHome() {
			this.$router.push('/home');
		}
	}
}
</script>

<template>
	<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
		<div class="container-fluid">
			<div class="navbar-brand" @click="goHome()">WasaPhoto</div>
			<ul class="navbar-nav">
				<li class="nav-item dropdown card">
					<div class="nav-link dropdown-toggle justify-content-center align-items-center" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false" style="color: black">
						<i class="bi bi-person-circle"></i>
						<span class="ms-2">{{ username }}&nbsp;</span>
					</div>
					<ul class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
						<li class="dropdown-item" @click="showProfile()">Show profile</li>
						<li class="dropdown-item" @click="logout">Logout</li>
					</ul>
				</li>
			</ul>
		</div>
	</nav>
</template>

<style>
</style>
