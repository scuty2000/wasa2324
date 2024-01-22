<script>
import Navbar from "../components/Navbar.vue";
import PhotoCard from "../components/PhotoCard.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import instance from "../services/axios";
import PhotoUploadModal from "./PhotoUploadModal.vue";

export default {
	name: "HomeView",
	components: { PhotoCard, ErrorMsg, Navbar, PhotoUploadModal },
	data() {
		return {
			photos: [],
			errormsg: null,
			photoPaginationIndex: 0,
			hasMorePhotos: true,
			isUploadModalOpen: false,
		};
	},
	mounted() {
		this.loadMorePhotos();

		this.$nextTick(() => {
			if (this.$refs.loadMorePhotosTrigger) {
				const options = {
					root: null,
					rootMargin: '0px',
					threshold: 1.0
				};
				this.observer = new IntersectionObserver(this.handleIntersect, options);
				this.observer.observe(this.$refs.loadMorePhotosTrigger);
			}
		});

	},
	methods: {
		async loadMorePhotos() {
			if (!this.hasMorePhotos) return;
			this.hasMorePhotos = false;

			try {
				const userID = localStorage.getItem('userId');
				const authToken = localStorage.getItem('authToken');
				const response = await instance.get(`/users/${userID}/stream`, {
					headers: { 'Authorization': authToken },
					params: { 'paginationIndex': this.photoPaginationIndex }
				});

				this.photos.push(...response.data.photos);
				this.photoPaginationIndex++;
				this.hasMorePhotos = response.data["paginationLimit"] >= this.photoPaginationIndex;
			} catch (error) {
				console.error("Error loading more photos:", error);
				this.errormsg = error.toString();
			}
		},
		handleIntersect(entries) {
			if (entries[0].isIntersecting && this.hasMorePhotos) {
				this.loadMorePhotos();
			}
		},
		openUploadModal() {
			this.isUploadModalOpen = true;
		},
		closeUploadModal() {
			this.isUploadModalOpen = false;
		},
	},
	unmounted() {
		if (this.observer) {
			this.observer.disconnect();
		}
	}
};
</script>

<template>
	<div class="user-stream">
		<div class="container">
			<div v-if="errormsg">
				<ErrorMsg :msg="errormsg" />
			</div>
			<div v-else-if="photos != null && photos.length === 0" class="d-flex justify-content-center">
				<div class="no-photos-card card text-center">
					<div class="card-body d-flex align-items-center justify-content-center">
						<span role="img" aria-label="Thinking face" class="emoji">🤔</span>
						<p class="card-text">There's nothing to see here :( <br> Follow someone to add their photos to your stream!</p>
					</div>
				</div>
			</div>
			<div class="row justify-content-center" v-else>
				<div v-for="photo in photos" :key="photo.id" class="d-flex flex-column align-items-center">
					<PhotoCard :photo-u-u-i-d="photo.id" :show-author="true" />
					&nbsp;
				</div>
			</div>
		</div>
		<button class="upload-button" @click="openUploadModal">
			<i class="bi bi-upload"></i> Upload a Photo
		</button>
		<PhotoUploadModal v-if="isUploadModalOpen" @close="closeUploadModal" />
	</div>
	<div ref="loadMorePhotosTrigger" class="load-more-trigger"></div>
</template>

<style>
.user-stream {
	margin-top: 20px;
	position: relative;
}

.no-photos-card {
	background-color: #f2f2f2;
	border: none;
	max-width: 500px;
}

.no-photos-card .card-body {
	padding: 20px;
}

.no-photos-card .emoji {
	font-size: 4rem;
	margin-right: 20px;
}

.no-photos-card .card-text {
	text-align: left;
	margin: 0;
}

.load-more-trigger {
	min-height: 10px;
	width: 100%;
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
</style>
