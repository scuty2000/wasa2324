<script>
import {useRoute} from "vue-router";
import {reactive} from "vue";
import instance from "../services/axios";
import PhotoCard from "../components/PhotoCard.vue";
import PhotoUploadModal from "./PhotoUploadModal.vue";
import FollowListModal from "./FollowListModal.vue";

export default {
	name: "UserProfileView",
	components: {FollowListModal, PhotoUploadModal, PhotoCard },
	data: function() {
		return {
			errormsg: null,
			userProfile: reactive({
				uuid: null,
				username: null,
				photosCount: null,
				followersCount: null,
				followingCount: null,
				isFollowed: null,
				isBanned: null,
			}),
			photos: [],
			showActionButtons: false,
			photoPaginationIndex: 0,
			hasMorePhotos: true,
			isUploadModalOpen: false,
			showUsernameEdit: false,
			editingUsername: false,
			newUsername: '',
			isFollowListModalOpen: false,
			followListType: '',
			loading: false,
		}
	},
	mounted() {
		const route = useRoute();
		this.fetchUserData(route.params.userID);
		this.loadMorePhotos();

		const options = {
			root: null,
			rootMargin: '0px',
			threshold: 1.0
		};
		this.observer = new IntersectionObserver(this.handleIntersect, options);
		this.observer.observe(this.$refs.loadMorePhotosTrigger);
	},
	methods: {
		async fetchUserData(userID) {
			this.errormsg = null;
			try {
				const requestingUserID = localStorage.getItem('userId');
				const authToken = localStorage.getItem('authToken');
				this.showActionButtons = requestingUserID !== userID;
				this.showUsernameEdit = requestingUserID === userID;
				const response = await instance.get(`/users/${userID}`, {
					headers: {
						'X-Requesting-User-UUID': requestingUserID,
						'Authorization': authToken
					}
				});
				this.userProfile = response.data;
				this.photoPaginationIndex = 0;
				this.hasMorePhotos = true;
				this.photos = [];
				await this.loadMorePhotos();
			} catch (error) {
				console.error("Error loading user data:", error);
				if(error.response.status === 404) {
					this.errormsg = "User not found.";
					return;
				} else if (error.response.status === 403) {
					this.errormsg = "This user has banned you.";
					return;
				} else if (error.response.status === 401) {
					this.errormsg = "You must be logged in to view this user.";
					return;
				} else if (error.response.status === 400) {
					this.errormsg = "Invalid user ID.";
					return;
				} else if (error.response.status === 500) {
					this.errormsg = "Internal server error.";
					return;
				}
				this.errormsg = error.toString();
			}
		},
		async toggleUserFollow() {
			if (this.userProfile.isFollowed) {
				try {
					const requestingUserID = localStorage.getItem('userId');
					await instance.delete(`/users/${requestingUserID}/following/${this.userProfile.uuid}`, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				} catch (error) {
					this.errormsg = error.toString();
					console.error("Error unfollowing user:", error);
				}
			} else {
				try {
					const requestingUserID = localStorage.getItem('userId');
					await instance.put(`/users/${requestingUserID}/following/${this.userProfile.uuid}`, null, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				} catch (error) {
					this.errormsg = error.toString();
					console.error("Error following user:", error);
				}
			}
			await this.fetchUserData(this.userProfile.uuid);
			this.$nextTick(() => {
				if (this.$refs.followButton) {
					this.$refs.followButton.blur();
				}
			});
		},
		async toggleUserBan() {
			const requestingUserID = localStorage.getItem('userId');
			if (this.userProfile.isBanned) {
				try {
					await instance.delete(`/users/${requestingUserID}/banned/${this.userProfile.uuid}`, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				} catch (error) {
					this.errormsg = error.toString();
					console.error("Error unbanning user:", error);
				}
			} else {
				try {
					await instance.put(`/users/${requestingUserID}/banned/${this.userProfile.uuid}`, null, {
						headers: {
							'Authorization': localStorage.getItem('authToken')
						}
					});
				} catch (error) {
					this.errormsg = error.toString();
					console.error("Error banning user:", error);
				}
			}
			await this.fetchUserData(this.userProfile.uuid);
			this.$nextTick(() => {
				if (this.$refs.banButton) {
					this.$refs.banButton.blur();
				}
			});
		},
		handlePhotoDeleted() {
			this.fetchUserData(this.userProfile.uuid);
		},
		async loadMorePhotos() {
			if (!this.hasMorePhotos) return;
			if (this.userProfile.uuid == null) return;
			if (this.loading) return;
			this.loading = true;

			try {
				const requestingUserID = localStorage.getItem('userId');
				const authToken = localStorage.getItem('authToken');
				const response = await this.$axios.get('/photos', {
					headers: {
						'X-Requesting-User-UUID': requestingUserID,
						'Authorization': authToken
					},
					params: {
						'userID': this.userProfile.uuid,
						'paginationIndex': this.photoPaginationIndex,
					}
				});
				this.photos.push(...response.data["user-photos"]);
				this.photoPaginationIndex++;
				this.hasMorePhotos = response.data["paginationLimit"] >= this.photoPaginationIndex;
				this.loading = false;
			} catch (error) {
				console.error("Error loading more photos:", error);
			}
		},

		handleIntersect(entries) {
			if (entries[0].isIntersecting) {
				this.loadMorePhotos();
			}
		},
		openUploadModal() {
			this.isUploadModalOpen = true;
		},
		async closeUploadModal() {
			this.isUploadModalOpen = false;
			await this.fetchUserData(this.userProfile.uuid);
		},
		editUsername() {
			this.editingUsername = true;
			this.newUsername = this.userProfile.username;
		},
		async saveUsername() {
			try {
				const userID = this.userProfile.uuid;
				const authToken = localStorage.getItem('authToken');
				await instance.put(`/users/${userID}/username`, { username: this.newUsername }, {
					headers: {
						'Authorization': authToken
					}
				});
				this.userProfile.username = this.newUsername;
				localStorage.setItem('username', this.newUsername);
				this.editingUsername = false;
				window.location.reload();
			} catch (error) {
				console.error("Error updating username:", error);
				this.errormsg = error.toString();
			}
		},
		openFollowListModal(type) {
			this.followListType = type;
			if(type === 'followers' && this.userProfile.followersCount > 0)
				this.isFollowListModalOpen = true;
			else if(type === 'following' && this.userProfile.followingCount > 0)
				this.isFollowListModalOpen = true;
		},
		closeFollowListModal() {
			this.isFollowListModalOpen = false;
		},
	},
	watch: {
		'$route.params.userID'(newUserID) {
			this.fetchUserData(newUserID);
		}
	},
	unmounted() {
		if (this.observer) {
			this.observer.disconnect();
		}
	}
}
</script>

<template>
	<div v-if="errormsg != null">
		<br>
		<ErrorMsg :msg="errormsg"></ErrorMsg>
	</div>
	<div class="user-profile container py-4" v-if="errormsg == null">
		<div v-if="!editingUsername">
			<h2 class="text-center mb-4">{{ userProfile.username }}
				<i v-if="showUsernameEdit" class="bi bi-pencil-square" @click="editUsername"></i>
			</h2>
		</div>
		<div v-else class="username-edit-section">
			<input type="text" v-model="newUsername" class="username-edit-input" pattern="^[a-zA-Z0-9_]*$" maxlength="16" minlength="3">
			<button class="username-save-button" @click="saveUsername">Save</button>
		</div>

		<div class="card mx-auto mb-4" style="max-width: 20rem;">
			<div class="card-body">
				<div class="row text-center">
					<div class="col">
						<div class="counter-label">Photos</div>
						<div class="counter-number">{{ userProfile.photosCount }}</div>
					</div>

					<div class="col clickable" @click="openFollowListModal('followers')">
						<div class="counter-label">Followers</div>
						<div class="counter-number">{{ userProfile.followersCount }}</div>
					</div>

					<div class="col clickable" @click="openFollowListModal('following')">
						<div class="counter-label">Following</div>
						<div class="counter-number">{{ userProfile.followingCount }}</div>
					</div>
				</div>

				<div v-if="showActionButtons" class="row mt-3">
					<div class="col-6">
						<button ref="followButton" @click="toggleUserFollow()" :class="['btn', userProfile.isFollowed ? 'btn-outline-primary' : 'btn-primary', 'btn-sm', 'w-100', 'd-flex', 'justify-content-center', 'align-items-center']">
							<i :class="userProfile.isFollowed ? 'bi bi-person-fill-dash' : 'bi bi-person-fill-add'"></i>
							<span class="ms-2">{{ userProfile.isFollowed ? 'Unfollow' : 'Follow' }}</span>
						</button>
					</div>
					<div class="col-6">
						<button ref="banButton" @click="toggleUserBan()" :class="['btn', userProfile.isBanned ? 'btn-outline-danger' : 'btn-danger', 'btn-sm', 'w-100', 'd-flex', 'justify-content-center', 'align-items-center']">
							<i :class="userProfile.isBanned ? 'bi bi-check-circle-fill' : 'bi bi-slash-circle-fill'"></i>
							<span class="ms-2">{{ userProfile.isBanned ? 'Unban' : 'Ban' }}</span>
						</button>
					</div>
				</div>
			</div>
		</div>

		<div class="container" v-if="errormsg == null">
			<div class="row justify-content-center">
				<div class="d-flex flex-column align-items-center" v-for="photo in photos" :key="photo.uuid">
					<PhotoCard :photo-u-u-i-d="photo" @photo-deleted="handlePhotoDeleted"></PhotoCard><br>
				</div>
				<div ref="loadMorePhotosTrigger"></div>
			</div>
		</div>
		<button class="upload-button" @click="openUploadModal">
			<i class="bi bi-upload"></i> Upload a Photo
		</button>
		<PhotoUploadModal v-if="isUploadModalOpen" @close="closeUploadModal" />
		<FollowListModal
			v-if="isFollowListModalOpen"
			:userID="userProfile.uuid"
			:listType="followListType"
			:user-name="userProfile.username"
			@close="closeFollowListModal"
		/>
	</div>
</template>

<style scoped>
.upload-button {
	position: fixed;
	bottom: 20px;
	left: 20px;
	background-color: rgba(45, 90, 110, 0.8);
	color: white;
	border: none;
	padding: 10px 20px;
	border-radius: 5px;
	font-weight: bold;
	cursor: pointer;
	display: flex;
	align-items: center;
	z-index: 1000;
}

.upload-button i {
	margin-right: 10px;
}

.upload-button {
	position: fixed;
	bottom: 20px;
	left: 20px;
	background-color: rgba(45, 90, 110, 0.8);
	color: white;
	border: none;
	padding: 10px 20px;
	border-radius: 5px;
	font-weight: bold;
	cursor: pointer;
	display: flex;
	align-items: center;
	z-index: 1000;
}

.upload-button i {
	margin-right: 10px;
}

.user-profile {
	margin-top: 20px;
}

.username-edit-section {
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 10px;
}

.username-edit-input {
	flex-grow: 1;
	padding: 10px;
	border: 1px solid #ced4da;
	border-radius: 0.25rem;
	margin-right: 10px;
	font-size: 1rem;
	line-height: 1.5;
	color: #495057;
	background-color: #fff;
	background-clip: padding-box;
	transition: border-color .15s ease-in-out,box-shadow .15s ease-in-out;
	max-width: 200px;
}

.username-edit-input:focus {
	color: #495057;
	background-color: #fff;
	border-color: #80bdff;
	outline: 0;
	box-shadow: 0 0 0 0.2rem rgba(0,123,255,.25);
}

.username-save-button {
	background-color: #4CAF50;
	color: white;
	border: none;
	border-radius: 4px;
	padding: 5px 15px;
	cursor: pointer;
	font-weight: bold;
}

.bi-pencil-square {
	cursor: pointer;
	margin-left: 10px;
	color: #1762cb;
}

.username-edit-input:focus + .username-save-button {
	margin-left: 5px;
}

.username-save-button:disabled {
	background-color: #ccc;
	cursor: not-allowed;
}

.col {
	cursor: default;
}

.clickable {
	cursor: pointer;
}
</style>
