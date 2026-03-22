<script lang="ts">
	import Icon from '@iconify/svelte';

	let {
		segments,
		currentIndex,
		isRunning
	}: { segments: TimerSegmentDef[]; currentIndex: number; isRunning: boolean } = $props();
</script>

<div class="flex w-full max-w-lg items-center justify-center gap-2">
	{#each segments as segment, i (i)}
		{#if i > 0}
			<div
				class={['h-0.5 w-4', i <= currentIndex ? 'bg-primary/50' : 'bg-base-300']}
			></div>
		{/if}
		<div
			class={[
				'flex items-center gap-1.5 rounded-full px-3 py-1.5 text-sm font-medium transition-colors',
				i === currentIndex && isRunning
					? 'bg-primary text-primary-content'
					: i < currentIndex
						? 'bg-primary/20 text-primary'
						: 'bg-base-200 text-base-content/50'
			]}
		>
			{#if i < currentIndex}
				<Icon icon="material-symbols:check" width="14" height="14" />
			{/if}
			{segment.label}
		</div>
	{/each}
</div>

<div class="text-base-content/50 text-center text-sm">
	{currentIndex + 1} / {segments.length}
</div>
