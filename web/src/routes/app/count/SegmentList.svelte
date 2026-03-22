<script lang="ts">
	import Icon from '@iconify/svelte';

	let {
		segments,
		durations = $bindable([])
	}: { segments: TimerSegmentDef[]; durations: number[] } = $props();

	function getMin(totalSeconds: number) {
		return Math.floor(totalSeconds / 60);
	}

	function getSec(totalSeconds: number) {
		return totalSeconds % 60;
	}

	function setMin(index: number, min: number) {
		if (min < 0) return;
		const sec = getSec(durations[index]);
		durations[index] = min * 60 + sec;
	}

	function setSec(index: number, sec: number) {
		if (sec < 0 || sec > 59) return;
		const min = getMin(durations[index]);
		durations[index] = min * 60 + sec;
	}
</script>

<div class="grid w-full max-w-lg gap-3">
	{#each segments as segment, i (i)}
		<div class="border-base-300/50 flex items-center gap-3 rounded-xl border bg-white/70 px-4 py-3">
			<span class="text-base-content/70 min-w-16 text-sm font-semibold">{segment.label}</span>
			<div class="flex items-center gap-1">
				<button
					class="btn btn-ghost btn-xs text-primary"
					onclick={() => setMin(i, getMin(durations[i]) - 1)}
					aria-label="decrease minutes"
				>
					<Icon icon="material-symbols:keyboard-arrow-down" width="16" height="16" />
				</button>
				<span class="min-w-10 text-center text-lg font-bold tabular-nums">
					{String(getMin(durations[i])).padStart(2, '0')}
				</span>
				<button
					class="btn btn-ghost btn-xs text-primary"
					onclick={() => setMin(i, getMin(durations[i]) + 1)}
					aria-label="increase minutes"
				>
					<Icon icon="material-symbols:keyboard-arrow-up" width="16" height="16" />
				</button>
			</div>
			<span class="text-base-content/30 text-lg">:</span>
			<div class="flex items-center gap-1">
				<button
					class="btn btn-ghost btn-xs text-primary"
					onclick={() => setSec(i, getSec(durations[i]) - 1)}
					aria-label="decrease seconds"
				>
					<Icon icon="material-symbols:keyboard-arrow-down" width="16" height="16" />
				</button>
				<span class="min-w-10 text-center text-lg font-bold tabular-nums">
					{String(getSec(durations[i])).padStart(2, '0')}
				</span>
				<button
					class="btn btn-ghost btn-xs text-primary"
					onclick={() => setSec(i, getSec(durations[i]) + 1)}
					aria-label="increase seconds"
				>
					<Icon icon="material-symbols:keyboard-arrow-up" width="16" height="16" />
				</button>
			</div>
		</div>
	{/each}
</div>
