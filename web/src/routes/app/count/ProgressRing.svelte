<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		progress = 1,
		active = false,
		size = 380,
		strokeWidth = 10,
		children
	}: {
		progress: number;
		active: boolean;
		size?: number;
		strokeWidth?: number;
		children: Snippet;
	} = $props();

	let radius = $derived((size - strokeWidth) / 2);
	let circumference = $derived(2 * Math.PI * radius);
	let offset = $derived(circumference * (1 - progress));
</script>

<div class="relative flex items-center justify-center" style="width: {size}px; height: {size}px;">
	<svg class="absolute inset-0" width={size} height={size} viewBox="0 0 {size} {size}">
		<circle
			cx={size / 2}
			cy={size / 2}
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
				cx={size / 2}
				cy={size / 2}
				r={radius}
				fill="none"
				stroke="currentColor"
				stroke-width={strokeWidth}
				stroke-linecap="round"
				stroke-dasharray={circumference}
				stroke-dashoffset={offset}
				class="text-primary transition-[stroke-dashoffset] duration-1000 ease-linear"
				transform="rotate(-90 {size / 2} {size / 2})"
			/>
		{/if}
	</svg>
	<div class="relative">
		{@render children()}
	</div>
</div>
