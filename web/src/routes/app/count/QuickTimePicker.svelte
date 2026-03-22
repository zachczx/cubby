<script lang="ts">
	import Icon from '@iconify/svelte';

	let {
		minutes = $bindable(5),
		seconds = $bindable(0),
		disabled = false,
		displayMinutes,
		displaySeconds
	}: {
		minutes: number;
		seconds: number;
		disabled?: boolean;
		displayMinutes?: string;
		displaySeconds?: string;
	} = $props();

	let shownMinutes = $derived(displayMinutes ?? String(minutes).padStart(2, '0'));
	let shownSeconds = $derived(displaySeconds ?? String(seconds).padStart(2, '0'));

	function changeMinutes(direction: 'increment' | 'decrement') {
		if (direction === 'increment') minutes += 1;
		if (direction === 'decrement' && minutes > 0) minutes -= 1;
	}

	function changeSeconds(direction: 'increment' | 'decrement') {
		if (direction === 'increment' && seconds < 59) seconds += 1;
		if (direction === 'decrement' && seconds > 0) seconds -= 1;
	}
</script>

<div class="flex w-full max-w-lg items-center justify-center gap-4">
	<div class="flex flex-col items-center text-center">
		<button
			class={[
				'btn btn-ghost transition-opacity duration-300',
				disabled ? 'text-neutral pointer-events-none opacity-30' : 'text-primary opacity-75'
			]}
			aria-label="+1 minute"
			onclick={() => changeMinutes('increment')}
		>
			<Icon icon="material-symbols:keyboard-arrow-up" class="size-8" />
		</button>
		<span class="text-7xl tabular-nums" aria-live="polite" aria-label={shownMinutes}>
			{shownMinutes}
		</span>
		<button
			class={[
				'btn btn-ghost transition-opacity duration-300',
				disabled ? 'text-neutral pointer-events-none opacity-30' : 'text-primary opacity-75'
			]}
			aria-label="-1 minute"
			onclick={() => changeMinutes('decrement')}
		>
			<Icon icon="material-symbols:keyboard-arrow-down" class="size-8" />
		</button>
	</div>
	<span class="text-base-content/30 text-7xl">:</span>
	<div class="flex flex-col items-center text-center">
		<button
			class={[
				'btn btn-ghost transition-opacity duration-300',
				disabled ? 'text-neutral pointer-events-none opacity-30' : 'text-primary opacity-75'
			]}
			aria-label="+1 second"
			onclick={() => changeSeconds('increment')}
		>
			<Icon icon="material-symbols:keyboard-arrow-up" class="size-8" />
		</button>
		<span class="text-7xl tabular-nums" aria-live="polite" aria-label={shownSeconds}>
			{shownSeconds}
		</span>
		<button
			class={[
				'btn btn-ghost transition-opacity duration-300',
				disabled ? 'text-neutral pointer-events-none opacity-30' : 'text-primary opacity-75'
			]}
			aria-label="-1 second"
			onclick={() => changeSeconds('decrement')}
		>
			<Icon icon="material-symbols:keyboard-arrow-down" class="size-8" />
		</button>
	</div>
</div>
