<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { pb } from '$lib/pb';
	import { onMount } from 'svelte';

	onMount(async () => {
		await pb.authStore.clear();

		goto(resolve('/login'));
	});
</script>
