<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import Logo from '$lib/assets/logo.webp?w=200&enhanced';
	import { onMount } from 'svelte';

	onMount(async () => {
		const response = await api.get('check');
		if (response.status === 204) {
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
