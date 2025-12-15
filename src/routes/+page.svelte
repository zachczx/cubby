<script lang="ts">
	import { goto } from '$app/navigation';
	import Logo from '$lib/assets/logo.webp?w=200&enhanced';
	import { pb } from '$lib/pb';
	import { onMount } from 'svelte';

	onMount(() => {
		if (pb.authStore.isValid) {
			setTimeout(() => goto('/app'), 2000);
		}

		if (!pb.authStore.isValid) {
			setTimeout(() => goto('/login'), 2000);
		}
	});
</script>

<svelte:head>
	<title>Cubby</title>
</svelte:head>

<div class="grid min-h-dvh content-center justify-items-center">
	<enhanced:img src={Logo} alt="logo" loading="eager" />
	<p>Redirecting...</p>
</div>
