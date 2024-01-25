<script>
import {RouterView, useRoute} from 'vue-router'
import Navbar from "./components/Navbar.vue";
import {computed} from "vue";

export default {
	components: {
		Navbar,
		RouterView,
	},
	setup() {
		const route = useRoute();
		const showNavbar = computed(() => route.path !== '/login');

		return { showNavbar };
	}
}
</script>

<template>
	<Navbar v-if="showNavbar"/>
	<div class="container-fluid">
		<RouterView />
	</div>
</template>

<style>
</style>
