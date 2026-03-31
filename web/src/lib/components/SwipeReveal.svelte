<script lang="ts">
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	let {
		actions,
		children,
		hint = false
	}: {
		actions: Snippet;
		children: Snippet;
		hint?: boolean;
	} = $props();

	const threshold = 50;
	const actionsWidth = 144;

	const HINT_KEY = 'swipe-reveal-hinted';

	let el = $state<HTMLDivElement | null>(null);
	let startX = 0;
	let startY = 0;
	let currentX = $state(0);
	let swiping = $state(false);
	let open = $state(false);
	let bouncing = $state(false);
	let locked = false;

	function onTouchStart(e: TouchEvent) {
		startX = e.touches[0].clientX;
		startY = e.touches[0].clientY;
		locked = false;
		swiping = false;
	}

	function onTouchMove(e: TouchEvent) {
		const dx = e.touches[0].clientX - startX;
		const dy = e.touches[0].clientY - startY;

		if (!swiping && !locked) {
			if (Math.abs(dy) > Math.abs(dx)) {
				locked = true;
				return;
			}
			if (Math.abs(dx) > 8) {
				swiping = true;
			}
		}

		if (locked) return;
		if (!swiping) return;

		e.preventDefault();

		const base = open ? -actionsWidth : 0;
		const raw = base + dx;
		currentX = Math.max(-actionsWidth, Math.min(0, raw));
	}

	function onTouchEnd() {
		if (!swiping) return;
		swiping = false;

		if (open) {
			open = Math.abs(currentX) > actionsWidth - threshold;
		} else {
			open = Math.abs(currentX) > threshold;
		}

		currentX = open ? -actionsWidth : 0;
	}

	onMount(() => {
		el!.addEventListener('touchmove', onTouchMove, { passive: false });

		if (hint && !localStorage.getItem(HINT_KEY)) {
			localStorage.setItem(HINT_KEY, '1');
			const timer = setTimeout(() => {
				bouncing = true;
				currentX = -50;
				setTimeout(() => {
					currentX = 0;
					setTimeout(() => (bouncing = false), 300);
				}, 400);
			}, 500);
			return () => {
				clearTimeout(timer);
				el!.removeEventListener('touchmove', onTouchMove);
			};
		}

		return () => el!.removeEventListener('touchmove', onTouchMove);
	});

	export function close() {
		open = false;
		currentX = 0;
	}
</script>

<div
	role="listitem"
	bind:this={el}
	class="swipe-reveal"
	ontouchstart={onTouchStart}
	ontouchend={onTouchEnd}
>
	<div class="swipe-reveal-actions" style="width: {actionsWidth}px">
		{@render actions()}
	</div>
	<div
		class="swipe-reveal-content rounded-r-2xl"
		class:transitioning={!swiping || bouncing}
		style="transform: translateX({currentX}px)"
	>
		{@render children()}
		<div
			class="text-base-content/50 pointer-events-none absolute top-1/2 right-1.5 -translate-y-1/2 transition-opacity duration-200"
			class:opacity-0={open || swiping}
		>
			<Icon icon="material-symbols:chevron-left" class="size-5" />
		</div>
	</div>
</div>

<style>
	.swipe-reveal {
		position: relative;
		overflow: hidden;
	}

	.swipe-reveal-actions {
		position: absolute;
		top: 1px;
		right: 1px;
		bottom: 1px;
		display: flex;
		align-items: stretch;
	}

	.swipe-reveal-content {
		position: relative;
		z-index: 1;
		will-change: transform;
	}

	.swipe-reveal-content.transitioning {
		transition: transform 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
	}
</style>
