<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		progress = 1,
		active = false,
		viewBox = 300,
		strokeWidth = 10,
		children
	}: {
		progress: number;
		active: boolean;
		viewBox?: number;
		strokeWidth?: number;
		children: Snippet;
	} = $props();

	let radius = $derived((viewBox - strokeWidth) / 2);
	let circumference = $derived(2 * Math.PI * radius);
	let offset = $derived(circumference * (1 - progress));
</script>

<div class="relative flex w-10/12 max-w-xs items-center justify-center">
	<svg
		class="absolute inset-0 h-full w-full"
		viewBox="0 0 {viewBox} {viewBox}"
		preserveAspectRatio="xMidYMid meet"
	>
		<circle
			cx={viewBox / 2}
			cy={viewBox / 2}
			r={radius}
			fill="none"
			stroke="currentColor"
			stroke-width={strokeWidth}
			class={active ? 'text-base-300/50' : 'text-transparent'}
			class:transition-colors={true}
			class:duration-300={true}
		/>

		{#if active}
			<circle
				cx={viewBox / 2}
				cy={viewBox / 2}
				r={radius}
				fill="none"
				stroke="currentColor"
				stroke-width={strokeWidth}
				stroke-linecap="round"
				stroke-dasharray={circumference}
				stroke-dashoffset={offset}
				class="text-primary transition-[stroke-dashoffset] duration-1000 ease-linear"
				transform="rotate(-90 {viewBox / 2} {viewBox / 2})"
			/>
		{/if}
	</svg>
	<div class="relative aspect-square w-full flex items-center justify-center">
		{@render children()}
	</div>
</div>
