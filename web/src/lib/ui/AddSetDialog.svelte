<script lang="ts">
	import Icon from '@iconify/svelte';
	import Dialog from '$lib/ui/Dialog.svelte';

	let {
		open = $bindable(false),
		exerciseName,
		weightUnit = $bindable('kg'),
		isLoading = false,
		onsubmit
	}: {
		open?: boolean;
		exerciseName: string;
		weightUnit?: 'kg' | 'lb';
		isLoading?: boolean;
		onsubmit: (data: {
			weightKg: number | null;
			reps: number | null;
			setType: string;
		}) => void;
	} = $props();

	let weight = $state<number | null>(null);
	let reps = $state<number | null>(null);
	let setType = $state('working');
	let weightMode = $state<'total' | 'each'>('total');

	const lbToKg = (lb: number) => Math.round((lb / 2.20462) * 10) / 10;

	export function prefill(opts: {
		weight: number | null;
		reps: number | null;
		setType: string;
		weightMode: 'total' | 'each';
	}) {
		weight = opts.weight;
		reps = opts.reps;
		setType = opts.setType;
		weightMode = opts.weightMode;
	}

	export function reset() {
		weight = null;
		reps = null;
		setType = 'working';
		weightMode = 'total';
	}

	function submit() {
		const displayWeight = weight != null && weightMode === 'each' ? weight * 2 : weight;
		const weightKg =
			displayWeight != null && weightUnit === 'lb' ? lbToKg(displayWeight) : displayWeight;
		onsubmit({ weightKg, reps, setType });
	}
</script>

<Dialog bind:open title={exerciseName || 'Add Set'}>
	<div class="grid gap-4">
		<div class="flex items-center gap-2">
			<input
				type="number"
				class="input input-bordered w-full rounded-lg"
				placeholder="Reps"
				bind:value={reps}
			/>
			<span class="text-base-content/40 text-lg font-bold">&times;</span>
			<div class="flex grow">
				<input
					type="number"
					class="input input-bordered w-full rounded-lg"
					placeholder={weightUnit}
					bind:value={weight}
				/>
				<button
					class="btn btn-ghost"
					onclick={() => (weightUnit = weightUnit === 'kg' ? 'lb' : 'kg')}
				>
					{weightUnit}
				</button>
			</div>
			<button
				class="btn btn-ghost border-base-300 min-w-18 shrink-0 rounded-lg border"
				onclick={() => (weightMode = weightMode === 'total' ? 'each' : 'total')}
				title={weightMode === 'each' ? 'Per hand' : 'Both hands'}
			>
				{#if weightMode === 'each'}
					<Icon icon="mdi:hand-back-left" class="size-5" />
				{:else}
					<Icon icon="mdi:hand-back-left" class="-mr-2 size-5" />
					<Icon icon="mdi:hand-back-right" class="size-5" />
				{/if}
			</button>
		</div>
		{#if weightMode === 'each' && weight != null}
			<p class="text-base-content/50 -mt-2 text-sm">
				{weight} &times; 2 = {weight * 2}
				{weightUnit} total
			</p>
		{/if}

		<fieldset class="join w-full text-sm">
			<legend class="sr-only">Set Type</legend>
			{#each [{ value: 'working', label: 'Working' }, { value: 'dropset', label: 'Drop' }, { value: 'failure', label: 'Failure' }] as opt (opt.value)}
				<input
					type="radio"
					name="set-type"
					class="btn join-item checked:bg-segmented checked:text-primary-content flex-1"
					aria-label={opt.label}
					checked={setType === opt.value}
					onchange={() => (setType = opt.value)}
				/>
			{/each}
		</fieldset>

		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isLoading}
			onclick={submit}
		>
			{#if isLoading}<span class="loading loading-spinner loading-sm"></span>{/if}
			Add Set
		</button>
	</div>
</Dialog>
