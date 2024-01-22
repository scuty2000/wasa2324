<script>
import ErrorMsg from "./ErrorMsg.vue";

export default {
	name: "PhotoCard",
	components: {ErrorMsg},
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
			showDeleteButton: false,
			showDeleteConfirmation: false,
		};
	},
	emits: ['photo-deleted'],
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
				this.showDeleteButton = this.photoDetails.author === localStorage.getItem('userId');
				console.log(this.showDeleteButton);
			} catch (error) {
				console.error("Error loading photo data:", error);
				this.errormsg = "Error loading photo data: " + error.toString();
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
				this.errormsg = "Error toggling like status: " + error.toString();
			}
		},
		showComments() {
			this.showCommentsSection = !this.showCommentsSection;
		},
		postComment() {
			// Aggiungi qui la logica per inviare il nuovo commento
		},
		toggleDeleteConfirmation() {
			this.showDeleteConfirmation = !this.showDeleteConfirmation;
		},
		async deletePhoto() {
			try {
				await this.$axios.delete(`/photos/${this.photoUUID}`, {
					headers: {
						'X-Requesting-User-UUID': localStorage.getItem('userId'),
						'Authorization': localStorage.getItem('authToken')
					}
				});
				this.showDeleteConfirmation = false;

				this.$emit('photo-deleted', this.photoUUID);
			} catch (error) {
				console.error("Error deleting photo:", error);
				this.errormsg = "Error deleting photo: " + error.toString();
			}
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
			<div v-if="errormsg == null" class="d-flex justify-content-between align-items-center">
				<div class="left-icons d-flex align-items-center">
					<span @click="toggleLike" class="icon-btn d-flex align-items-center">
						<i :class="['bi', photoDetails.liked ? 'bi-heart-fill' : 'bi-heart']"></i>
						<span class="counter">{{ photoDetails.likesCount }}</span>
					</span>
					<span @click="showComments" class="icon-btn d-flex align-items-center">
						<i :class="['bi', showCommentsSection ? 'bi-chat-fill' : 'bi-chat']"></i>
						<span class="counter">{{ photoDetails.commentsCount }}</span>
					</span>
				</div>
				<span v-if="showDeleteButton && !showDeleteConfirmation" @click="toggleDeleteConfirmation" class="icon-btn" style="color: red;">
					<i class="bi bi-trash-fill"></i>
				</span>
			</div>
			<div class="photo-info d-flex justify-content-end">
				<span class="photo-date">{{ formatDate(photoDetails.date) }}</span>
			</div>
			<ErrorMsg v-if="errormsg != null" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
	<div class="comments-section" v-if="showCommentsSection && !showDeleteConfirmation && errormsg == null">
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
	<div class="delete-confirmation" v-if="showDeleteConfirmation && errormsg == null">
		<div class="confirmation-message">
			Are you sure you want to delete this photo?
		</div>
		<div class="confirmation-buttons">
			<button @click="deletePhoto" class="btn btn-danger">
				<i class="bi bi-trash-fill"></i> Delete
			</button>
			<button @click="toggleDeleteConfirmation" class="btn btn-secondary">
				<i class="bi bi-x-circle-fill"></i> Cancel
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

.left-icons {
	display: flex;
	align-items: center;
}

.icon-btn {
	cursor: pointer;
	display: flex;
	align-items: center;
	margin-right: 20px;
}

.icon-btn:last-child {
	margin-right: 0;
}

.counter {
	margin-left: 5px;
}

.bi-chat, .bi-heart, .bi-heart-fill, .bi-trash-fill {
	vertical-align: middle;
}

.counter {
	margin-left: 5px;
}

.photo-info {
	margin-top: 10px;
}

.photo-date {
	font-size: 0.8rem;
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
	margin-right: 0.5rem;
}

.new-comment button {
	background-color: #82868a;
	color: #ffffff;
	border: none;
}

.new-comment button:hover {
	background-color: #94979c;
}

.confirmation-buttons .btn + .btn {
	margin-left: 10px;
}

.delete-confirmation {
	background-color: #dadce0;
	padding: 25px 1rem 1rem;
	border-radius: 0 0 0.5rem 0.5rem;
	width: 100%;
	max-width: 500px;
	margin-top: -15px;
}

.confirmation-message {
	display: flex;
	margin-bottom: 1rem;
	justify-content: center;
}

.confirmation-buttons {
	display: flex;
	justify-content: center;
}
</style>
