<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Logo from '$lib/assets/logo.webp?w=200&enhanced';
	import { pb } from '$lib/pb';
	import { onMount } from 'svelte';

	onMount(() => {
		if (pb.authStore.isValid) {
			goto(resolve('/app'));
		} else {
			goto(resolve('/login'));
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
