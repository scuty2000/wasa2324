<script>
export default {
	name: "Navbar",
	data() {
		return {
			// Esempio di nome utente, sostituiscilo con la logica appropriata
			username: 'Utente'
		};
	},
	methods: {
		logout() {
			localStorage.removeItem('userId');
			localStorage.removeItem('authToken');
			this.$router.push('/login');
		}
	}
}
</script>

<template>
	<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
		<div class="container-fluid">
			<a class="navbar-brand" href="#">WasaPhoto</a>
			<ul class="navbar-nav">
				<li class="nav-item dropdown card">
					<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false" style="color: black">
						{{ username }}
					</a>
					<ul class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
						<li><a class="dropdown-item" href="#" @click="showProfile()">Show profile</a></li>
						<li><a class="dropdown-item" href="#" @click="logout">Logout</a></li>
					</ul>
				</li>
			</ul>
		</div>
	</nav>
</template>

<style>
</style>
