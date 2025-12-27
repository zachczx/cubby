<script lang="ts">
	import MaterialSymbolsCalendarClock from '$lib/assets/svg/MaterialSymbolsCalendarClock.svelte';
	import MaterialSymbolsNestClockFarsightAnalogOutline from '$lib/assets/svg/MaterialSymbolsNestClockFarsightAnalogOutline.svelte';
	import type { Component, Snippet } from 'svelte';

	let {
		left,
		leftTitle,
		leftIcon,
		right,
		rightTitle,
		rightIcon
	}: {
		left: Snippet;
		leftTitle: string;
		leftIcon?: Component;
		right: Snippet;
		rightTitle: string;
		rightIcon?: Component;
	} = $props();
</script>

<div class="w-full">
	<div class="text-md grid grid-cols-2 content-center gap-4">
		<div
			class="border-base-content/5 bg-base-200/20 grid content-center justify-items-center rounded-xl border p-4 shadow"
		>
			<div class="flex items-center gap-2">
				{#if leftIcon}
					{@const Comp = leftIcon}
					<Comp class="size-[1.3em] opacity-50" />
				{:else}
					<MaterialSymbolsNestClockFarsightAnalogOutline class="size-[1.3em] opacity-50" />
				{/if}
				{leftTitle}
			</div>
			<div class="grid min-h-20 content-center justify-items-center text-center text-2xl font-bold">
				{@render left()}
			</div>
		</div>
		<div
			class="border-base-content/5 bg-base-200/20 grid justify-items-center rounded-xl border p-4 shadow"
		>
			<div class="flex items-center gap-2">
				{#if rightIcon}
					{@const Comp = rightIcon}
					<Comp class="size-[1.3em] opacity-50" />
				{:else}
					<MaterialSymbolsCalendarClock class="size-[1.3em] opacity-50" />
				{/if}
				{rightTitle}
			</div>
			<div class="grid min-h-20 content-center justify-items-center text-center text-2xl font-bold">
				{@render right()}
			</div>
		</div>
	</div>
</div>
