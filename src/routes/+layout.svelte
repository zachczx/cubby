<script lang="ts">
	import ArkToaster from '$lib/ui/ArkToaster.svelte';
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { onNavigate } from '$app/navigation';
	import { topLevelRoutes } from '$lib/shell/nav';

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		let direction = '';

		const isNoAnimation = topLevelRoutes.noAnimation.some(
			(route) =>
				navigation.from?.url.pathname.startsWith(route.href) ||
				navigation.to?.url.pathname.startsWith(route.href)
		);

		if (!isNoAnimation) {
			const fromIndex = topLevelRoutes.animation.findIndex(
				(route) => route.href === navigation.from?.url.pathname
			);
			const toIndex = topLevelRoutes.animation.findIndex(
				(route) => route.href === navigation.to?.url.pathname
			);

			if (fromIndex !== -1 && toIndex !== -1) {
				// Top level route

				if (toIndex < fromIndex) {
					direction = 'back';
				} else {
					direction = 'forward';
				}
			} else {
				// 2nd level routes
				const fromId = navigation.from?.route.id;
				const toId = navigation.to?.route.id;

				const fromLevel = fromId?.split('/');
				const toLevel = toId?.split('/');

				if (fromLevel && toLevel) {
					if (fromLevel.length > toLevel.length) {
						direction = 'back';
					} else {
						direction = 'forward';
					}
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
