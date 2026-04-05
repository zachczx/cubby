<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { routinesQueryOptions, getRoutinesQueryKey, getAllWorkoutsQueryKey } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { goto } from '$app/navigation';
	import { router } from '$lib/routes';
	import Icon from '@iconify/svelte';
	import Dialog from '$lib/ui/Dialog.svelte';

	const queryClient = useQueryClient();
	const routinesDb = createQuery(routinesQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id.toLowerCase(), e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId.toLowerCase())?.name ?? exerciseId;
	}

	let createDialogOpen = $state(false);
	let newRoutineName = $state('');
	let isCreating = $state(false);

	async function createRoutine() {
		if (!newRoutineName.trim()) return;
		isCreating = true;
		const response = await api.post('gym/routines', {
			json: { name: newRoutineName.trim() }
		});
		isCreating = false;
		if (response.status === 201) {
			const routine: RoutineDB = await response.json();
			addToast('success', 'Routine created!');
			queryClient.invalidateQueries({ queryKey: getRoutinesQueryKey() });
			createDialogOpen = false;
			newRoutineName = '';
			goto(router.gymRoutine(routine.id));
		} else {
			addToast('error', 'Failed to create routine');
		}
	}

	let sortedRoutines = $derived(
		routinesDb.data ? [...routinesDb.data].sort((a, b) => a.position - b.position) : []
	);

	function updateRoutinesCache(updater: (routines: RoutineDB[]) => RoutineDB[]) {
		const current = queryClient.getQueryData<RoutineDB[]>(getRoutinesQueryKey());
		if (current) queryClient.setQueryData(getRoutinesQueryKey(), updater(current));
	}

	let isReordering = $state(false);

	async function reorderRoutine(routineId: string, direction: 'up' | 'down') {
		isReordering = true;

		updateRoutinesCache((routines) => {
			const sorted = [...routines].sort((a, b) => a.position - b.position);
			const idx = sorted.findIndex((r) => r.id === routineId);
			const neighborIdx = direction === 'up' ? idx - 1 : idx + 1;
			if (neighborIdx < 0 || neighborIdx >= sorted.length) return routines;
			const current = sorted[idx];
			const neighbor = sorted[neighborIdx];
			return routines.map((r) => {
				if (r.id === current.id) return { ...r, position: neighbor.position };
				if (r.id === neighbor.id) return { ...r, position: current.position };
				return r;
			});
		});

		const response = await api.post('gym/routines/reorder', {
			json: { routineId, direction }
		});
		isReordering = false;
		if (response.status !== 204) {
			addToast('error', 'Failed to reorder');
			queryClient.invalidateQueries({ queryKey: getRoutinesQueryKey() });
		}
	}

	let isStarting = $state<string | null>(null);

	async function startFromRoutine(routineId: string) {
		isStarting = routineId;
		const response = await api.post(`gym/routines/${routineId}/start`);
		isStarting = null;
		if (response.status === 201) {
			const workout: WorkoutDB = await response.json();
			addToast('success', 'Workout started!');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
			goto(router.gym(workout.id));
		} else {
			addToast('error', 'Failed to start workout');
		}
	}
</script>

<PageWrapper title="Routines">
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<button
					class="btn btn-primary btn-lg w-full rounded-full"
					onclick={() => {
						newRoutineName = '';
						createDialogOpen = true;
					}}
				>
					<Icon icon="material-symbols:add" class="size-5" />
					Create Routine
				</button>
			</section>

			<section class="grid gap-4 py-2">
				{#if routinesDb.isSuccess}
					{#if sortedRoutines.length === 0}
						<div class="grid justify-items-center gap-2 py-8">
							<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-16" />
							<p class="text-base-content/50 font-medium">No routines yet</p>
							<p class="text-base-content/40 text-sm">
								Create a routine to quickly start workouts
							</p>
						</div>
					{:else}
						{#each sortedRoutines as routine, ri (routine.id)}
							<div class="border-base-300/50 bg-base-50 rounded-2xl border">
								<div class="flex items-center gap-0 px-4 pt-3 pb-2">
									<a href={router.gymRoutine(routine.id)} class="block grow">
										<div class="text-base-content/90 text-lg font-bold">
											{routine.name}
										</div>
										<div class="text-base-content/50 text-sm">
											{routine.exercises.length} exercise{routine.exercises.length !== 1 ? 's' : ''}
										</div>
									</a>
									{#if sortedRoutines.length > 1}
										<div class="flex items-center gap-0.5">
											{#if ri > 0}
												<button
													class="btn btn-ghost btn-xs btn-square"
													disabled={isReordering}
													onclick={() => reorderRoutine(routine.id, 'up')}
												>
													<Icon
														icon="material-symbols:arrow-upward"
														class="text-base-content/50 size-3.5"
													/>
												</button>
											{/if}
											{#if ri < sortedRoutines.length - 1}
												<button
													class="btn btn-ghost btn-xs btn-square"
													disabled={isReordering}
													onclick={() => reorderRoutine(routine.id, 'down')}
												>
													<Icon
														icon="material-symbols:arrow-downward"
														class="text-base-content/50 size-3.5"
													/>
												</button>
											{/if}
										</div>
									{/if}
								</div>

								{#if routine.exercises.length > 0}
									<div class="grid gap-0 px-4 text-sm">
										{#each routine.exercises.slice(0, 3) as ex, i (ex.id)}
											<div class="flex items-center">
												<div
													class={[
														'grow py-2.5',
														i < Math.min(routine.exercises.length, 3) - 1 &&
															'border-b-base-300/50 border-b'
													]}
												>
													<p class="text-base-content/80 font-medium">
														{getExerciseName(ex.exerciseId)}
													</p>
												</div>
												<div class="text-base-content/65 px-2 text-sm font-medium">
													{ex.sets} sets
												</div>
											</div>
										{/each}
										{#if routine.exercises.length > 3}
											<div class="text-base-content/60 py-2.5 text-sm italic">
												+{routine.exercises.length - 3} more
											</div>
										{/if}
									</div>
								{/if}

								<div class="border-t-base-300/50 border-t px-4 py-3">
									<button
										class="btn btn-soft btn-primary w-full"
										disabled={isStarting === routine.id}
										onclick={() => startFromRoutine(routine.id)}
									>
										{#if isStarting === routine.id}
											<span class="loading loading-spinner loading-sm"></span>
										{/if}
										Start Workout
									</button>
								</div>
							</div>
						{/each}
					{/if}
				{:else}
					<div class="border-base-300/50 bg-base-50 grid h-28 gap-3 rounded-2xl border p-4">
						<div class="skeleton h-5 w-32 rounded-full"></div>
						<div class="skeleton h-4 w-48 rounded-full"></div>
					</div>
					<div class="border-base-300/50 bg-base-50 grid h-28 gap-3 rounded-2xl border p-4">
						<div class="skeleton h-5 w-32 rounded-full"></div>
						<div class="skeleton h-4 w-48 rounded-full"></div>
					</div>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>

<Dialog bind:open={createDialogOpen} title="New Routine">
	<div class="grid gap-4">
		<input
			type="text"
			class="input input-lg w-full"
			placeholder="Routine name (e.g. Push Day)"
			bind:value={newRoutineName}
		/>
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isCreating || !newRoutineName.trim()}
			onclick={createRoutine}
		>
			{#if isCreating}<span class="loading loading-spinner loading-sm"></span>{/if}
			Create
		</button>
	</div>
</Dialog>
