<script lang="ts">
	let { minutes, seconds }: { minutes: string; seconds: string } = $props();
</script>

<div class="flex w-full max-w-lg items-center justify-center gap-4">
	<span class="text-8xl tabular-nums" aria-live="polite" aria-label={minutes}>{minutes}</span>
	<span class="text-8xl text-base-content/30">:</span>
	<span class="text-8xl tabular-nums" aria-live="polite" aria-label={seconds}>{seconds}</span>
</div>
