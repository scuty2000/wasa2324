<script>
export default {
	name: "PhotoCard",
	props: {
		photoUUID: String,
	},
	data() {
		return {
			errormsg: null,
			photoDetails: {
				id: null,
				author: null,
				date: null,
				extension: null,
				likesCount: 0,
				commentsCount: 0,
				liked: false
			},
			showCommentsSection: false,
			comments: [],
			newCommentText: '',
		};
	},
	methods: {
		async fetchPhotoDetails() {
			try {
				const response = await this.$axios.get(`/photos/${this.photoUUID}`, {
					headers: {
						'X-Requesting-User-UUID': localStorage.getItem('userId'),
						'Authorization': localStorage.getItem('authToken'),
					},
				});
				this.photoDetails = response.data;
			} catch (error) {
				console.error("Error loading photo data:", error);
			}
		},
		formatDate(dateString) {
			const date = new Date(dateString);
			const day = date.getDate().toString().padStart(2, '0');
			const month = (date.getMonth() + 1).toString().padStart(2, '0');
			const year = date.getFullYear();
			const hours = date.getHours().toString().padStart(2, '0');
			const minutes = date.getMinutes().toString().padStart(2, '0');

			return `${day}/${month}/${year} - ${hours}:${minutes}`;
		},
		async toggleLike() {
			try {
				const userID = localStorage.getItem('userId');
				if (this.photoDetails.liked) {
					await this.$axios.delete(`/photos/${this.photoUUID}/likes/${userID}`, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				} else {
					await this.$axios.put(`/photos/${this.photoUUID}/likes/${userID}`, null, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				}
				await this.fetchPhotoDetails();
			} catch (error) {
				console.error("Error toggling like status:", error);
				// Gestisci l'errore qui
			}
		},
		showComments() {
			this.showCommentsSection = !this.showCommentsSection;
		},
		postComment() {
			// Aggiungi qui la logica per inviare il nuovo commento
		},
	},
	mounted() {
		this.fetchPhotoDetails();
	}
}
</script>

<template>
	<div class="photo-card card shadow-sm">
		<img :src="`http://localhost:8080/uploads/${photoDetails.author}/${photoUUID}.${photoDetails.extension}`" alt="Photo" class="card-img-top" />

		<div class="card-body">
			<div class="likes-comments d-flex justify-content-between align-items-center">
				<span @click="toggleLike" class="icon-btn">
				  <i :class="['bi', photoDetails.liked ? 'bi-heart-fill' : 'bi-heart']"></i>
				  {{ photoDetails.likesCount }}
				</span>
				<span @click="showComments" class="icon-btn">
				  <i class="bi bi-chat"></i>
				  {{ photoDetails.commentsCount }}
				</span>
				<span class="photo-date">{{ formatDate(photoDetails.date) }}</span>
			</div>
		</div>
	</div>
	<div class="comments-section" v-if="showCommentsSection">
		<div v-for="comment in comments" :key="comment.id" class="comment">
			<span>{{ comment.author }}: {{ comment.text }}</span>
		</div>

		<div class="new-comment">
			<input type="text" class="form-control rounded" v-model="newCommentText" maxlength="250" placeholder="Leave a comment...">
			<button class="btn btn-light" @click="postComment">
				<i class="bi bi-send-fill"></i>
			</button>
		</div>
	</div>

</template>

<style scoped>
.photo-card {
	border-radius: 0.5rem;
	max-width: 500px;
	width: 100%;
	background-color: #f8f9fa;
}

.icon-btn {
	cursor: pointer;
	display: inline-flex;
	align-items: center;
}

.icon-btn i {
	margin-right: 0.5rem;
}

.photo-date {
	font-size: 0.9rem;
	color: #6c757d;
}

.bi-heart-fill {
	color: red;
}

.comments-section {
	background-color: #dadce0;
	color: #3d4249;
	padding: 1rem;
	border-radius: 0 0 0.5rem 0.5rem;
	width: 100%;
	max-width: 500px;
	margin-top: -15px;
}

.comment {
	margin-bottom: 0.5rem;
}

.new-comment {
	margin-top: 1rem;
	display: flex;
	align-items: center;
}

.new-comment input {
	flex-grow: 1;
	margin-right: 0.5rem; /* Spazio tra il campo di input e il pulsante */
}

.new-comment button {
	background-color: #82868a; /* Colore grigio scuro per il bottone */
	color: #ffffff; /* Colore del testo e dell'icona bianco */
	border: none; /* Rimuove il bordo per un look più pulito */
}

.new-comment button:hover {
	background-color: #94979c; /* Un grigio leggermente più chiaro per lo hover */
}
</style>
