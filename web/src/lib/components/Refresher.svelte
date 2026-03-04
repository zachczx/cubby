<script lang="ts">
	import type { Snippet } from 'svelte';
	import Icon from '@iconify/svelte';

	let {
		resistance = 0.4,
		onRefresh = async () => {},
		children
	}: {
		resistance?: number;
		onRefresh?: () => Promise<void>;
		children: Snippet;
	} = $props();

	let startY = $state(0);
	let currentY = $state(0);
	let pulling = $state(false);
	let refreshing = $state(false);
	let rotateDeg = $state(0);
	let shouldRefresh = $state(false);
	let translateY = $state(0);

	const ontouchstart = (event: TouchEvent) => {
		if (refreshing) return;
		startY = event.touches[0].clientY;
	};

	const ontouchmove = (event: TouchEvent) => {
		if (refreshing) return;

		const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
		if (scrollTop > 0) return;

		currentY = event.touches[0].clientY;
		const delta = currentY - startY;

		if (delta > 20) {
			pulling = true;
			rotateDeg = delta;
			translateY = delta * resistance;
			shouldRefresh = delta > 160;
		} else {
			pulling = false;
		}
	};

	const ontouchend = async () => {
		if (refreshing) return;

		if (shouldRefresh) {
			refreshing = true;
			rotateDeg = 0;
			translateY = 60;

			await onRefresh();

			translateY = 0;
			pulling = false;
			shouldRefresh = false;
			refreshing = false;
		} else {
			translateY = 0;
			pulling = false;
			shouldRefresh = false;
		}
	};
</script>

<div role="status" {ontouchstart} {ontouchmove} {ontouchend} class="refresher bg-base-200">
	<!-- {#if pulling || refreshing} -->
	<div class="fixed top-3 right-0 left-0 w-full">
		<div class="grid w-full content-start justify-items-center">
			{#if shouldRefresh || refreshing}
				<div class="bg-base-200 flex items-center justify-center gap-3">
					<span class="loading loading-spinner loading-sm top-4"></span>
					<div class="text-base-content text-lg">Refreshing</div>
				</div>
			{:else}
				<div class="bg-base-200 flex items-center justify-center gap-3">
					<Icon icon="material-symbols:arrow-downward" class="size-6 opacity-75" />
					<!-- <Icon
							icon="material-symbols:refresh"
							class="size-6"
							style="transform: rotate({rotateDeg}deg);"
						/> -->
					<div class="text-base-content text-lg">Pull down to refresh</div>
				</div>
			{/if}
		</div>
	</div>
	<!-- {/if} -->

	<div class="content-wrapper" style="transform: translateY({translateY}px)">
		{@render children()}
	</div>
</div>

<style>
	.refresher {
		height: 100%;
		position: relative;
		overscroll-behavior-y: contain;
	}

	.content-wrapper {
		transition: transform 0.5s cubic-bezier(0.2, 0.8, 0.2, 1);
	}

	@keyframes spin {
		100% {
			transform: rotate(360deg);
		}
	}
</style>
