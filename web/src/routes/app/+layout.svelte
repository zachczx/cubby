<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { pb } from '$lib/pb';
	import { type Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	beforeNavigate(({ cancel }) => {
		if (!pb.authStore.isValid) {
			cancel();
			goto(resolve('/login'));
		}
	});
</script>

{@render children()}
