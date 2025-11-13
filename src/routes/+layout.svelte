<script lang="ts">
	import ArkToaster from '$lib/ui/ArkToaster.svelte';
	import '../app.css';
	import { pb } from '$lib/pb';
	import { goto } from '$app/navigation';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { onNavigate } from '$app/navigation';

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});
	const queryClient = new QueryClient();

	let { children } = $props();
</script>

<QueryClientProvider client={queryClient}>
	<ArkToaster />
	{@render children()}
</QueryClientProvider>
