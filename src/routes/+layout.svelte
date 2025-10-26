<script lang="ts">
	import ArkToaster from '$lib/ui/ArkToaster.svelte';
	import '../app.css';
	import { pb } from '$lib/pb';
	import { goto } from '$app/navigation';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';

	const queryClient = new QueryClient();

	let { children } = $props();

	if (!pb.authStore.isValid) {
		goto('/login');
	}
</script>

<QueryClientProvider client={queryClient}>
	<ArkToaster />
	{@render children()}
</QueryClientProvider>
