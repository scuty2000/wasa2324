<script>
export default {
	name: "PhotoUploadModal",
	data() {
		return {
			imageFile: null,
			imagePreview: null,
			uploading: false,
			uploadComplete: false,
		};
	},
	methods: {
		triggerFileInput() {
			this.$refs.fileInput.click();
		},
		onFileChange(event) {
			const file = event.target.files[0];
			if (file && (file.type === "image/png" || file.type === "image/jpeg")) {
				if (file.size <= 65 * 1024 * 1024) {
					this.imageFile = file;
					this.imagePreview = URL.createObjectURL(file);
				} else {
					alert("File too large");
				}
			} else {
				alert("Invalid file type");
			}
		},
		onDrop(event) {
			const file = event.dataTransfer.files[0];
			this.onFileChange({ target: { files: [file] } });
		},
		async uploadImage() {
			if (!this.imageFile) return;

			this.uploading = true;
			setTimeout(async () => {
				try {
					const headers = {
						'Authorization': localStorage.getItem('authToken'),
						'X-Requesting-User-UUID': localStorage.getItem('userId'),
						'Content-Type': this.imageFile.type
					};
					await this.$axios.post('/photos', this.imageFile, { headers });

					this.uploading = false;
					this.uploadComplete = true;
				} catch (error) {
					console.error("Error uploading image:", error);
					this.errormsg = "Error uploading image: " + error.toString();
					this.uploading = false;
				}
			}, 2000);
		},
		closeModal() {
			this.$emit('close');
		},
	}
}
</script>

<template>
	<div class="photo-upload-modal">
		<div class="modal-overlay" @click="closeModal">
			<div class="modal-card" @click.stop>
				<div class="modal-header">
					<button class="close-button" @click="closeModal"><i class="bi bi-x-lg"></i></button>
				</div>

				<div class="image-drop-area" @drop.prevent="onDrop" @dragover.prevent @dragenter.prevent>
					<input type="file" ref="fileInput" accept="image/png, image/jpeg" @change="onFileChange" hidden>
					<div v-if="imagePreview">
						<img :src="imagePreview" class="preview-image" alt="Uploaded image preview">
					</div>
					<div v-else @click="triggerFileInput">Drag & drop here or click to select an image</div>
				</div>

				<button class="btn btn-success mt-3" :disabled="uploading || uploadComplete" @click="uploadImage">
					<span v-if="uploading">Uploading...</span>
					<span v-else-if="uploadComplete">Upload Complete</span>
					<span v-else>Upload Photo</span>
				</button>
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
	max-width: 600px;
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
	color: grey;
	font-size: 25px;
	cursor: pointer;
}

.image-drop-area {
	border: 2px dashed #ccc;
	padding: 20px;
	text-align: center;
	cursor: pointer;
}

.preview-image {
	max-width: 100%;
	max-height: 400px;
	object-fit: contain;
}

.photo-upload-modal {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: rgba(0, 0, 0, 0.5); /* Sfondo semi-trasparente */
	z-index: 1050; /* Maggiore di qualsiasi altro elemento, incluso la navbar */
}
</style>
