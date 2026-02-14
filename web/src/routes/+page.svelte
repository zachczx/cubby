<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Logo from '$lib/assets/logo.webp?w=200&enhanced';
	import { onMount } from 'svelte';

	onMount(() => {
		goto(resolve('/app'));
	});
</script>

<svelte:head>
	<title>Cubby</title>
</svelte:head>

<div class="grid min-h-dvh content-center justify-items-center">
	<enhanced:img src={Logo} alt="logo" loading="eager" />
	<p>Redirecting...</p>
</div>
