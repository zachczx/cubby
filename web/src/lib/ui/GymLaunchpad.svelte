<script lang="ts">
	import Icon from '@iconify/svelte';
	import dayjs from 'dayjs';
	import { router } from '$lib/routes';

	interface Props {
		activeWorkout: WorkoutDB | undefined;
		gymThisMonth: number;
		gymMonthlyTarget: number;
		onStartWorkout: () => Promise<void>;
	}

	let { activeWorkout, gymThisMonth, gymMonthlyTarget, onStartWorkout }: Props = $props();

	let starting = $state(false);
	let elapsedSeconds = $state(0);

	let exerciseCount = $derived.by(() => {
		if (!activeWorkout?.sets || activeWorkout.sets.length === 0) return 0;
		const unique = new Set(activeWorkout.sets.map((s) => s.exerciseId));
		return unique.size;
	});

	let elapsedDisplay = $derived.by(() => {
		const hours = Math.floor(elapsedSeconds / 3600);
		const minutes = Math.floor((elapsedSeconds % 3600) / 60);
		const seconds = elapsedSeconds % 60;
		if (hours > 0) {
			return `${hours}h ${minutes.toString().padStart(2, '0')}m`;
		}
		return `${minutes}:${seconds.toString().padStart(2, '0')}`;
	});

	$effect(() => {
		if (!activeWorkout) return;

		function tick() {
			elapsedSeconds = dayjs().diff(dayjs(activeWorkout!.startTime), 'second');
		}
		tick();
		const id = setInterval(tick, 1000);
		return () => clearInterval(id);
	});

	async function handleStart() {
		starting = true;
		try {
			await onStartWorkout();
		} finally {
			starting = false;
		}
	}
</script>

<section class="border-primary/20 bg-primary/5 grid gap-3 rounded-2xl border p-4">
	<a href={router.gymStats} class="flex items-center justify-between">
		<span class="text-base-content/60 text-sm font-medium">Gym This Month</span>
		<div class="flex items-center gap-1">
			{#each Array(gymMonthlyTarget) as _, i}
				<div
					class={['size-3.5 rounded-sm', i < gymThisMonth ? 'bg-neutral' : 'bg-neutral/15']}
				></div>
			{/each}
		</div>
	</a>

	{#if activeWorkout}
		<div class="flex items-center justify-between">
			<div class="grid gap-0.5">
				<div class="flex items-center gap-2">
					<Icon icon="material-symbols:local-fire-department-rounded" class="text-primary size-5" />
					<span class="text-primary text-sm font-bold">Workout in Progress</span>
				</div>
				<div class="text-base-content/60 flex items-center gap-2 ps-7 text-sm">
					<span class="tabular-nums">{elapsedDisplay}</span>
					{#if exerciseCount > 0}
						<span class="text-base-content/30">&#x2022;</span>
						<span>{exerciseCount} exercise{exerciseCount !== 1 ? 's' : ''}</span>
					{/if}
				</div>
			</div>
			<a href={router.gym(activeWorkout.id)} class="btn btn-primary btn-sm rounded-full">
				Continue
			</a>
		</div>
	{:else}
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			onclick={handleStart}
			disabled={starting}
		>
			{#if starting}
				<span class="loading loading-spinner loading-sm"></span>
			{:else}
				<Icon icon="material-symbols:play-arrow-rounded" class="size-6" />
			{/if}
			Start Workout
		</button>
	{/if}
</section>
