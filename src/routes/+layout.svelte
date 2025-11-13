<script lang="ts">
	import ArkToaster from '$lib/ui/ArkToaster.svelte';
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { onNavigate } from '$app/navigation';

	const topLevelRoutes = ['/', '/personal', '/household', '/pet'];

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		let direction = 'forward';

		// Bottom navbar
		if (
			topLevelRoutes.includes(navigation.from?.url.pathname ?? '') &&
			topLevelRoutes.includes(navigation.to?.url.pathname ?? '')
		) {
			const fromIndex = topLevelRoutes.indexOf(navigation.from?.url.pathname ?? '');
			const toIndex = topLevelRoutes.indexOf(navigation.to?.url.pathname ?? '');

			if (toIndex < fromIndex) {
				direction = 'back';
			}
		} else {
			const from = navigation.from?.route.id;
			const to = navigation.to?.route.id;

			const fromLevel = from?.split('/');
			const toLevel = to?.split('/');

			if (fromLevel && toLevel) {
				if (fromLevel.length > toLevel.length) {
					direction = 'back';
				}
			}
		}

		document.documentElement.dataset.direction = direction;

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
