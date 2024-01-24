<script>
import ErrorMsg from "./ErrorMsg.vue";

export default {
	name: "PhotoCard",
	components: {ErrorMsg},
	props: {
		photoUUID: String,
		showAuthor: {
			type: Boolean,
			default: false
		}
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
			currentUser: null,
			userInfos: {},
			showCommentsSection: false,
			newCommentText: '',
			showDeleteButton: false,
			showDeleteConfirmation: false,
			comments: [],
			paginationIndex: 0,
			paginationLimit: 0,
			hasMoreComments: false,
			apiUrl: __API_URL__,
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
				if (!this.userInfos[this.photoDetails.author]) {
					const userResponse = await this.$axios.get(`/users/${this.photoDetails.author}`, {
						headers: {
							'X-Requesting-User-UUID': localStorage.getItem('userId'),
							'Authorization': localStorage.getItem('authToken')
						}
					});
					this.userInfos[this.photoDetails.author] = userResponse.data.username;
				}
				this.showDeleteButton = this.photoDetails.author === localStorage.getItem('userId');
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
		async postComment()
		{
			if (!this.newCommentText.trim()) {
				console.error("Comment text is empty");
				return;
			}

			try {
				const issuerUUID = localStorage.getItem('userId');
				await this.$axios.post(`/photos/${this.photoUUID}/comments`, {
					text: this.newCommentText,
					issuer: issuerUUID
				}, {
					headers: {
						'Authorization': localStorage.getItem('authToken')
					}
				});
				await this.fetchPhotoDetails()
				this.comments = []
				this.paginationIndex = 0
				await this.fetchComments()

				this.newCommentText = '';
			} catch (error) {
				console.error("Error posting comment:", error);
				this.errormsg = "Error posting comment: " + error.toString();
			}
		}
,
		toggleDeleteConfirmation() {
			this.showDeleteConfirmation = !this.showDeleteConfirmation;
			this.showCommentsSection = false;
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
		async fetchComments() {
			try {
				const response = await this.$axios.get(`/photos/${this.photoUUID}/comments?paginationIndex=${this.paginationIndex}`, {
					headers: {
						'X-Requesting-User-UUID': localStorage.getItem('userId'),
						'Authorization': localStorage.getItem('authToken')
					}
				});

				this.comments.push(...response.data.comments);
				this.paginationLimit = response.data.paginationLimit;
				this.hasMoreComments = this.paginationIndex < this.paginationLimit;

				await Promise.all(this.comments.map(async (comment) => {
					if (!this.userInfos[comment.issuer]) {
						const userResponse = await this.$axios.get(`/users/${comment.issuer}`, {
							headers: {
								'X-Requesting-User-UUID': localStorage.getItem('userId'),
								'Authorization': localStorage.getItem('authToken')
							}
						});
						this.userInfos[comment.issuer] = userResponse.data.username;
					}
				}));
			} catch (error) {
				console.error("Error fetching comments:", error);
				this.errormsg = "Error fetching comments: " + error.toString();
			}
		},
		loadMoreComments() {
			this.paginationIndex += 1;
			this.fetchComments();
		},
		async deleteComment(commentId) {
			try {
				await this.$axios.delete(`/photos/${this.photoUUID}/comments/${commentId}`, {
					headers: {
						'X-Requesting-User-UUID': localStorage.getItem('userId'),
						'Authorization': localStorage.getItem('authToken')
					}
				});
				await this.fetchPhotoDetails()
				this.comments = this.comments.filter(comment => comment.id !== commentId);
			} catch (error) {
				console.error("Error deleting comment:", error);
				this.errormsg = "Error deleting comment: " + error.toString();
			}
		},
	},
	mounted() {
		this.currentUser = localStorage.getItem('userId');
		this.fetchPhotoDetails();
		this.fetchComments();
	}
}
</script>

<template>
	<div class="photo-card card shadow-sm">
		<div v-if="showAuthor && photoDetails.author" class="photo-author">
			<router-link :to="`/user/${photoDetails.author}`" class="author-link">
				<i class="bi bi-person-badge"></i>
				{{ userInfos[photoDetails.author] || 'Unknown' }}
			</router-link>
		</div>
		<img :src="`${ apiUrl }/uploads/${photoDetails.author}/${photoUUID}.${photoDetails.extension}`" alt="Photo" class="card-img-top" />

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
		<hr class="my-2">
		<div v-for="comment in comments" :key="comment.id" class="comment">
			<div class="comment-header d-flex justify-content-between">
				<span class="comment-author">{{ userInfos[comment.issuer] || 'Unknown user' }}</span>
				<span class="comment-date">{{ formatDate(comment.date) }}</span>
			</div>
			<div class="comment-text">
				{{ comment.text }}
			</div>
			<div class="comment-header d-flex justify-content-end">
				<span v-if="comment.issuer === currentUser" @click="deleteComment(comment.id)" class="icon-btn">
					<i class="bi bi-trash" style="color: red;"></i>
				</span>
			</div>
		</div>
		<button v-if="hasMoreComments" @click="loadMoreComments" class="btn btn-secondary">Show more comments</button>

		<div class="new-comment">
	  		<textarea class="form-control rounded"
				v-model="newCommentText"
				maxlength="250"
				placeholder="Leave a comment..."
				rows="2"></textarea>
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

.new-comment textarea {
	flex-grow: 1;
	margin-right: 0.5rem;
	resize: vertical;
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

.comment {
	padding: 0.5rem;
	background-color: #f8f9fa;
	border-radius: 0.25rem;
	margin-bottom: 0.5rem;
}

.comment-header {
	font-size: 0.9rem;
	color: #6c757d;
}

.comment-author {
	font-weight: bold;
}

.comment-date {
	font-style: italic;
}

.comment-text {
	margin-top: 0.5rem;
	color: #495057;
}

hr.my-2 {
	border-top: 1px solid #e9ecef;
}

.photo-author {
	background-color: #f8f9fa;
	padding: 10px;
	border-top-left-radius: 0.5rem;
	border-top-right-radius: 0.5rem;
	text-align: left;
}

.author-link {
	color: #495057;
	font-weight: bold;
	display: flex;
	align-items: center;
	text-decoration: none;
}

.author-link i {
	font-size: 1.2rem;
	margin-right: 5px;
}

</style>
