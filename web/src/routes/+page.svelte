<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api';
	import loading from '$lib/assets/loading.webp?w=400&enhanced';
	import { onMount } from 'svelte';
	import { loadingQuips } from './loading-quips';

	onMount(async () => {
		const response = await api.get('check');

		const timeout = setTimeout(() => {
			if (response.status === 204) {
				goto(resolve('/app'));
			} else {
				goto(resolve('/login'));
			}

			clearTimeout(timeout);
		}, 500);
	});

	function randomIntFromInterval(min: number, max: number) {
		return Math.floor(Math.random() * (max - min + 1) + min);
	}

	const quip = loadingQuips[randomIntFromInterval(0, loadingQuips.length - 1)];
</script>

<svelte:head>
	<title>Cubby</title>
</svelte:head>

<div class="grid min-h-dvh content-center justify-items-center">
	<div class="w-full max-w-3xs"><enhanced:img src={loading} alt="loading" loading="eager" /></div>
	<p class="italic">{quip}...</p>
</div>
