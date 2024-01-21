<script>
export default {
	name: "LoginView",
	data() {
		return {
			username: ''
		};
	},
	methods: {
		async handleLogin() {
			try {
				const response = await this.$axios.post('/session', { name: this.username });
				const { identifier, token } = response.data;

				localStorage.setItem('userId', identifier);
				localStorage.setItem('authToken', token);

				this.$router.push('/home');
			} catch (error) {
				console.error(error);
				// Gestisci gli errori (es. mostrando un messaggio all'utente)
			}
		}
	}
};
</script>

<template>
	<div class="login-container d-flex align-items-center justify-content-center">
		<div class="login-form">
			<h2 class="text-center mb-4">Accedi a WasaPhoto</h2>
			<form @submit.prevent="handleLogin">
				<div class="row">
					<div class="col-12 mb-3">
						<input
							type="text"
							class="form-control"
							id="username"
							v-model="username"
							placeholder="Username"
							required
							pattern="^[a-zA-Z0-9_]{3,16}$"
							title="Username should be 3-16 characters long and can only contain letters, numbers, and underscores."
						>
					</div>
					<div class="col-12">
						<button type="submit" class="btn btn-primary btn-block">Login</button>
					</div>
				</div>
			</form>
		</div>
	</div>
</template>

<style scoped>
.login-container {
	height: 100vh;
	background-color: #f8f9fa;
}

.login-form {
	width: 100%;
	max-width: 400px;
	padding: 15px;
	margin: auto;
	background-color: white;
	border-radius: 10px;
	box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.login-form h2 {
	color: #333;
}

.form-control {
	border-radius: 20px;
	border: 1px solid #ced4da;
}

.btn-primary {
	border-radius: 20px;
	background-color: #007bff;
	border-color: #007bff;
}

.btn-primary:hover {
	background-color: #0069d9;
	border-color: #0062cc;
}
.btn-block {
	width: 100%;
}
</style>
