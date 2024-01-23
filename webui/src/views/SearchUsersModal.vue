<script>
export default {
	name: "SearchUsersModal",
	data() {
		return {
			searchQuery: '',
			searchResults: [],
			errormsg: null,
		};
	},
	watch: {
		searchQuery(newQuery) {
			if (newQuery.length >= 3 && newQuery.length <= 16 && /^[a-zA-Z0-9_]*$/.test(newQuery)) {
				this.performSearch();
			} else {
				this.searchResults = [];
			}
		},
	},
	methods: {
		async performSearch() {
			try {
				const response = await this.$axios.get(`/users?searchQuery=${this.searchQuery}`, {
					headers: {
						'Authorization': localStorage.getItem('authToken'),
						'X-Requesting-User-UUID': localStorage.getItem('userId')
					}
				});
				this.searchResults = response.data.users;
				this.errormsg = null;
			} catch (error) {
				console.error("Error searching users:", error);
				this.errormsg = "Error searching users: " + error.toString();
			}
		},
		closeModal() {
			this.$emit('close');
		},
	}
}
</script>

<template>
	<div class="user-search-modal">
		<div class="modal-overlay" @click="closeModal">
			<div class="modal-card" @click.stop>
				<div class="modal-header">
					<button class="close-button" @click="closeModal"><i class="bi bi-x-lg"></i></button>
				</div>

				<input type="text" v-model="searchQuery" class="form-control" placeholder="Search users..." maxlength="16" />

				<div class="search-results mt-3">
					<div v-if="errormsg">{{ errormsg }}</div>
					<div v-else>
						<ul style="padding-left: 0;">
							<li v-for="user in searchResults" :key="user.uuid" class="search-result-card">
								<router-link @click="closeModal" :to="`/user/${user.uuid}`" class="user-link">
									<i class="bi bi-person-fill"></i>
									<b>{{ user.username }}</b>
								</router-link>
							</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
	backdrop-filter: blur(10px);
	background-color: rgba(0, 0, 0, 0.5);
}

.modal-card {
	background: white;
	padding: 0 20px 20px;
	border-radius: 10px;
	max-width: 400px;
	width: 100%;
	position: relative;
}

.modal-header {
	display: flex;
	justify-content: flex-end;
}

.close-button {
	border: none;
	background: none;
	color: #777777;
	font-size: 25px;
	cursor: pointer;
}

.user-search-modal {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: rgba(0, 0, 0, 0.5);
	z-index: 1050;
}

.search-result-card {
	list-style-type: none;
	background-color: #c7d3dc;
	padding: 10px 15px;
	margin-bottom: 10px;
	border-radius: 5px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.user-link {
	text-decoration: none;
	color: #000;
	display: flex;
	align-items: center;
}

.bi-person-fill {
	margin-right: 10px;
	font-size: 1.5rem;
	color: #2d5a6e;
}
</style>
