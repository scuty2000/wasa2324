<script>
import {useRoute} from "vue-router";
import {reactive} from "vue";
import instance from "../services/axios";
import PhotoCard from "../components/PhotoCard.vue";

export default {
	name: "UserProfileView",
	components: { PhotoCard },
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
				this.errormsg = error.toString();
				console.error("Error loading user data:", error);
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
			await this.$nextTick(() => {
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
			await this.$nextTick(() => {
				if (this.$refs.banButton) {
					this.$refs.banButton.blur();
				}
			});
		},
		handlePhotoDeleted() {
			console.log("Photo deleted, reloading user data");
			this.fetchUserData(this.userProfile.uuid);
		},
		async loadMorePhotos() {
			if (!this.hasMorePhotos) return;
			if (this.userProfile.uuid == null) return;

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
			} catch (error) {
				console.error("Error loading more photos:", error);
			}
		},

		handleIntersect(entries) {
			if (entries[0].isIntersecting) {
				this.loadMorePhotos();
			}
		}
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
		<h2 class="text-center mb-4">{{ userProfile.username }}</h2>

		<div class="card mx-auto mb-4" style="max-width: 20rem;">
			<div class="card-body">
				<div class="row text-center">
					<div class="col">
						<div class="counter-label">Photos</div>
						<div class="counter-number">{{ userProfile.photosCount }}</div>
					</div>

					<div class="col">
						<div class="counter-label">Followers</div>
						<div class="counter-number">{{ userProfile.followersCount }}</div>
					</div>

					<div class="col">
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

	</div>
</template>

<style scoped>
</style>
