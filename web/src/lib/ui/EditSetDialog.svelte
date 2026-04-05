<script lang="ts">
	import BitsDialog from '$lib/ui/Dialog.svelte';

	let {
		open = $bindable(false),
		exerciseName = '',
		weightUnit = $bindable('kg'),
		isLoading = false,
		onsave
	}: {
		open?: boolean;
		exerciseName?: string;
		weightUnit?: 'kg' | 'lb';
		isLoading?: boolean;
		onsave: (data: {
			weightKg: number | null;
			reps: number | null;
			setType: string;
		}) => void;
	} = $props();

	let weight = $state<number | null>(null);
	let reps = $state<number | null>(null);
	let setType = $state('working');

	const kgToLb = (kg: number) => Math.round(kg * 2.20462 * 10) / 10;
	const lbToKg = (lb: number) => Math.round((lb / 2.20462) * 10) / 10;

	export function prefill(opts: {
		weight: number | null;
		reps: number | null;
		setType: string;
	}) {
		weight = opts.weight;
		reps = opts.reps;
		setType = opts.setType;
	}

	function submit() {
		const weightKg = weight != null && weightUnit === 'lb' ? lbToKg(weight) : weight;
		onsave({ weightKg, reps, setType });
	}
</script>

<BitsDialog bind:open title={exerciseName}>
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
					onclick={() => {
						if (weight != null) {
							weight = weightUnit === 'kg' ? kgToLb(weight) : lbToKg(weight);
						}
						weightUnit = weightUnit === 'kg' ? 'lb' : 'kg';
					}}
				>
					{weightUnit}
				</button>
			</div>
		</div>

		<fieldset class="join w-full text-sm">
			<legend class="sr-only">Set Type</legend>
			{#each [{ value: 'working', label: 'Working' }, { value: 'dropset', label: 'Drop' }, { value: 'failure', label: 'Failure' }] as opt (opt.value)}
				<input
					type="radio"
					name="edit-set-type"
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
			Save
		</button>
	</div>
</BitsDialog>
