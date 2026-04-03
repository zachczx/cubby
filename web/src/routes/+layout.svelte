<script lang="ts">
	import ArkToaster from '$lib/ui/ArkToaster.svelte';
	import '../app.css';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { onNavigate } from '$app/navigation';
	import { topLevelRoutes } from '$lib/shell/nav';
	import type { OnNavigate } from '@sveltejs/kit';
	import { queryClient } from '$lib/queries';
	import { Capacitor } from '@capacitor/core';
	import { App } from '@capacitor/app';
	import { Keyboard } from '@capacitor/keyboard';
	import { onMount } from 'svelte';

	onMount(() => {
		if (Capacitor.isNativePlatform()) {
			App.addListener('backButton', ({ canGoBack }) => {
				if (canGoBack) {
					window.history.back();
				} else {
					App.exitApp();
				}
			});

			Keyboard.addListener('keyboardWillShow', (info) => {
				document.documentElement.style.setProperty(
					'--keyboard-height',
					`${info.keyboardHeight}px`
				);
				document.documentElement.classList.add('keyboard-open');
			});

			Keyboard.addListener('keyboardWillHide', () => {
				document.documentElement.style.setProperty('--keyboard-height', '0px');
				document.documentElement.classList.remove('keyboard-open');
			});
		}

	});

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		document.documentElement.dataset.direction = getAnimationStatus(topLevelRoutes, navigation);

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});

	let { children } = $props();

	type TopLevelRoutes = typeof topLevelRoutes;

	function getAnimationStatus(topLevelRoutes: TopLevelRoutes, navigation: OnNavigate) {
		let direction = '';

		const isSamePage = navigation.from?.url.pathname === navigation.to?.url.pathname;

		const isNoAnimation = topLevelRoutes.noAnimation.some(
			(route) =>
				navigation.from?.url.pathname.startsWith(route.href) ||
				navigation.to?.url.pathname.startsWith(route.href)
		);

		if (!isSamePage && !isNoAnimation) {
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
				// Profile and other topNavAnimation would end up here (always -1)
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

		return direction;
	}
</script>

<QueryClientProvider client={queryClient}>
	<ArkToaster />
	{@render children()}
</QueryClientProvider>
