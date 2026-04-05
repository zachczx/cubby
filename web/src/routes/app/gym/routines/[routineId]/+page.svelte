<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		routinesQueryOptions,
		getRoutinesQueryKey,
		favouriteExercisesQueryOptions,
		getFavouriteExercisesQueryKey,
		getAllWorkoutsQueryKey
	} from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { router } from '$lib/routes';
	import Icon from '@iconify/svelte';
	import BitsDialog from '$lib/ui/Dialog.svelte';
	import DeleteConfirmDialog from '$lib/ui/DeleteConfirmDialog.svelte';

	let { data } = $props();

	const queryClient = useQueryClient();
	const routinesDb = createQuery(routinesQueryOptions);
	const favouritesDb = createQuery(favouriteExercisesQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id.toLowerCase(), e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId.toLowerCase())?.name ?? exerciseId;
	}

	let currentRoutine = $derived.by(() => {
		if (!routinesDb.isSuccess || !routinesDb.data) return;
		return routinesDb.data.find((r) => r.id === data.routineId);
	});

	function updateRoutinesCache(updater: (routines: RoutineDB[]) => RoutineDB[]) {
		const current = queryClient.getQueryData<RoutineDB[]>(getRoutinesQueryKey());
		if (current) queryClient.setQueryData(getRoutinesQueryKey(), updater(current));
	}

	// Rename
	let renameDialogOpen = $state(false);
	let editName = $state('');
	let isSavingName = $state(false);

	function openRename() {
		editName = currentRoutine?.name ?? '';
		renameDialogOpen = true;
	}

	async function saveRename() {
		if (!currentRoutine || !editName.trim()) return;
		isSavingName = true;
		const response = await api.patch(`gym/routines/${currentRoutine.id}`, {
			json: { name: editName.trim() }
		});
		isSavingName = false;
		if (response.status === 204) {
			updateRoutinesCache((routines) =>
				routines.map((r) => (r.id === currentRoutine!.id ? { ...r, name: editName.trim() } : r))
			);
			renameDialogOpen = false;
		} else {
			addToast('error', 'Failed to rename routine');
		}
	}

	// Delete routine
	let deleteDialogOpen = $state(false);
	let isDeleting = $state(false);

	async function confirmDelete() {
		if (!currentRoutine) return;
		isDeleting = true;
		const response = await api.delete(`gym/routines/${currentRoutine.id}`);
		isDeleting = false;
		if (response.status === 204) {
			addToast('success', 'Routine deleted');
			updateRoutinesCache((routines) => routines.filter((r) => r.id !== currentRoutine!.id));
			goto(router.gymRoutines);
		} else {
			addToast('error', 'Failed to delete routine');
		}
	}

	// Add exercise
	let exerciseDialogOpen = $state(false);
	let exerciseSearch = $state('');
	let setsCount = $state(3);
	let isAddingExercise = $state(false);

	let filteredExercises = $derived(
		exerciseSearch
			? exercises.filter((e) => e.name.toLowerCase().includes(exerciseSearch.toLowerCase()))
			: exercises
	);

	let favouriteIds = $derived(favouritesDb.data?.exerciseIds ?? []);
	let favouriteExercises = $derived(exercises.filter((e) => favouriteIds.includes(e.id)));

	async function toggleFavourite(exerciseId: string) {
		const response = await api.post('gym/favourites', {
			json: { exerciseId }
		});
		if (response.ok) {
			const data: FavouriteExercisesDB = await response.json();
			queryClient.setQueryData(getFavouriteExercisesQueryKey(), data);
		}
	}

	async function addExercise(exerciseId: string) {
		if (!currentRoutine) return;
		isAddingExercise = true;
		const response = await api.post(`gym/routines/${currentRoutine.id}/exercises`, {
			json: { exerciseId, sets: setsCount }
		});
		isAddingExercise = false;
		if (response.status === 201) {
			const exercise: RoutineExerciseDB = await response.json();
			updateRoutinesCache((routines) =>
				routines.map((r) =>
					r.id === currentRoutine!.id
						? { ...r, exercises: [...r.exercises, exercise] }
						: r
				)
			);
			exerciseDialogOpen = false;
			exerciseSearch = '';
			setsCount = 3;
		} else {
			addToast('error', 'Failed to add exercise');
		}
	}

	// Edit sets count
	let editSetsDialogOpen = $state(false);
	let editingSetsExercise = $state<RoutineExerciseDB | null>(null);
	let editSetsCount = $state(3);
	let isSavingSets = $state(false);

	function openEditSets(ex: RoutineExerciseDB) {
		editingSetsExercise = ex;
		editSetsCount = ex.sets;
		editSetsDialogOpen = true;
	}

	async function saveEditSets() {
		if (!editingSetsExercise) return;
		isSavingSets = true;
		const response = await api.patch(
			`gym/routines/${editingSetsExercise.routineId}/exercises/${editingSetsExercise.id}`,
			{ json: { exerciseId: editingSetsExercise.exerciseId, sets: editSetsCount } }
		);
		isSavingSets = false;
		if (response.status === 204) {
			const updatedId = editingSetsExercise.id;
			updateRoutinesCache((routines) =>
				routines.map((r) => ({
					...r,
					exercises: r.exercises.map((e) =>
						e.id === updatedId ? { ...e, sets: editSetsCount } : e
					)
				}))
			);
			editSetsDialogOpen = false;
			editingSetsExercise = null;
		} else {
			addToast('error', 'Failed to update sets');
		}
	}

	// Remove exercise
	let removeDialogOpen = $state(false);
	let removingExercise = $state<RoutineExerciseDB | null>(null);
	let isRemoving = $state(false);

	function openRemove(ex: RoutineExerciseDB) {
		removingExercise = ex;
		removeDialogOpen = true;
	}

	async function confirmRemove() {
		if (!removingExercise) return;
		isRemoving = true;
		const response = await api.delete(
			`gym/routines/${removingExercise.routineId}/exercises/${removingExercise.id}`
		);
		isRemoving = false;
		if (response.status === 204) {
			const removedId = removingExercise.id;
			updateRoutinesCache((routines) =>
				routines.map((r) => ({
					...r,
					exercises: r.exercises.filter((e) => e.id !== removedId)
				}))
			);
			removeDialogOpen = false;
			removingExercise = null;
		} else {
			addToast('error', 'Failed to remove exercise');
		}
	}

	// Reorder
	let isReordering = $state(false);

	async function reorderExercise(exerciseId: string, direction: 'up' | 'down') {
		if (!currentRoutine) return;
		isReordering = true;

		// Optimistic update
		updateRoutinesCache((routines) =>
			routines.map((r) => {
				if (r.id !== currentRoutine!.id) return r;
				const sorted = [...r.exercises].sort((a, b) => a.position - b.position);
				const idx = sorted.findIndex((e) => e.id === exerciseId);
				const neighborIdx = direction === 'up' ? idx - 1 : idx + 1;
				if (neighborIdx < 0 || neighborIdx >= sorted.length) return r;
				const current = sorted[idx];
				const neighbor = sorted[neighborIdx];
				return {
					...r,
					exercises: r.exercises.map((e) => {
						if (e.id === current.id) return { ...e, position: neighbor.position };
						if (e.id === neighbor.id) return { ...e, position: current.position };
						return e;
					})
				};
			})
		);

		const response = await api.post(`gym/routines/${currentRoutine.id}/exercises/reorder`, {
			json: { exerciseId, direction }
		});
		isReordering = false;
		if (response.status !== 204) {
			addToast('error', 'Failed to reorder');
			queryClient.invalidateQueries({ queryKey: getRoutinesQueryKey() });
		}
	}

	// Start workout
	let isStarting = $state(false);

	async function startFromRoutine() {
		if (!currentRoutine) return;
		isStarting = true;
		const response = await api.post(`gym/routines/${currentRoutine.id}/start`);
		isStarting = false;
		if (response.status === 201) {
			const workout: WorkoutDB = await response.json();
			addToast('success', 'Workout started!');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
			goto(router.gym(workout.id));
		} else {
			addToast('error', 'Failed to start workout');
		}
	}

	let sortedExercises = $derived(
		currentRoutine ? [...currentRoutine.exercises].sort((a, b) => a.position - b.position) : []
	);
</script>

<PageWrapper title={currentRoutine?.name ?? 'Routine'}>
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			{#if routinesDb.isSuccess && currentRoutine}
				<section class="grid gap-4 py-2">
					<div class="flex items-center justify-between">
						<h2 class="text-base-content/70 text-lg font-bold">{currentRoutine.name}</h2>
						<button class="btn btn-ghost btn-sm" onclick={openRename}>
							<Icon icon="material-symbols:edit-outline" class="size-4" />
							Rename
						</button>
					</div>

					<button
						class="btn btn-primary btn-lg w-full rounded-full"
						disabled={isStarting || sortedExercises.length === 0}
						onclick={startFromRoutine}
					>
						{#if isStarting}<span class="loading loading-spinner loading-sm"></span>{/if}
						Start Workout
					</button>
				</section>

				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Exercises</h2>

					{#if sortedExercises.length === 0}
						<div class="grid justify-items-center gap-2 py-8">
							<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-12" />
							<p class="text-base-content/50 font-medium">No exercises yet</p>
						</div>
					{:else}
						<div class="border-base-300/50 bg-base-50 rounded-2xl border">
							{#each sortedExercises as ex, i (ex.id)}
								<div
									class={[
										'px-4 py-3',
										i < sortedExercises.length - 1 && 'border-b-base-300/50 border-b'
									]}
								>
									<div class="flex items-center gap-2">
										<div class="grow">
											<p class="text-base-content/80 font-semibold">
												{getExerciseName(ex.exerciseId)}
											</p>
											<button
												class="text-base-content/50 text-sm hover:underline"
												onclick={() => openEditSets(ex)}
											>
												{ex.sets} sets
											</button>
										</div>
										<div class="flex items-center gap-0.5">
											{#if i > 0}
												<button
													class="btn btn-ghost btn-xs btn-square"
													disabled={isReordering}
													onclick={() => reorderExercise(ex.id, 'up')}
												>
													<Icon
														icon="material-symbols:arrow-upward"
														class="text-base-content/50 size-3.5"
													/>
												</button>
											{/if}
											{#if i < sortedExercises.length - 1}
												<button
													class="btn btn-ghost btn-xs btn-square"
													disabled={isReordering}
													onclick={() => reorderExercise(ex.id, 'down')}
												>
													<Icon
														icon="material-symbols:arrow-downward"
														class="text-base-content/50 size-3.5"
													/>
												</button>
											{/if}
											<button
												class="btn btn-ghost btn-xs btn-square"
												onclick={() => openRemove(ex)}
											>
												<Icon
													icon="material-symbols:close"
													class="text-base-content/50 size-3.5"
												/>
											</button>
										</div>
									</div>
								</div>
							{/each}
						</div>
					{/if}

					<button
						class="btn btn-primary btn-soft w-full rounded-full"
						onclick={() => {
							exerciseSearch = '';
							setsCount = 3;
							exerciseDialogOpen = true;
						}}
					>
						<Icon icon="material-symbols:add" class="size-4" />
						Add Exercise
					</button>
				</section>

				<section class="py-4">
					<button
						class="btn btn-ghost btn-error flex w-full items-center gap-2 rounded-full"
						onclick={() => (deleteDialogOpen = true)}
					>
						<Icon icon="material-symbols:delete-outline" class="text-error size-5" />
						Delete Routine
					</button>
				</section>
			{:else if routinesDb.isSuccess && !currentRoutine}
				<div class="grid justify-items-center gap-4 py-16">
					<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-16" />
					<p class="text-base-content/50 text-lg font-medium">Routine not found</p>
					<a href={router.gymRoutines} class="btn btn-primary btn-soft rounded-full"
						>Back to Routines</a
					>
				</div>
			{:else}
				<div class="border-base-300/50 bg-base-50 grid h-28 gap-3 rounded-2xl border p-4">
					<div class="skeleton h-5 w-32 rounded-full"></div>
					<div class="skeleton h-4 w-48 rounded-full"></div>
				</div>
			{/if}
		</div>
	</main>
</PageWrapper>

<!-- Exercise picker -->
{#snippet exerciseListContent()}
	<div class="grid gap-3">
		<input
			type="text"
			class="input input-lg w-full"
			placeholder="Search exercises..."
			bind:value={exerciseSearch}
		/>

		<div class="mb-2 flex items-center gap-3">
			<span class="text-base-content/60 text-sm">Sets:</span>
			<div class="join">
				{#each [1, 2, 3, 4, 5] as n (n)}
					<input
						type="radio"
						name="sets-count"
						class="btn btn-sm join-item checked:bg-segmented checked:text-primary-content"
						aria-label={String(n)}
						checked={setsCount === n}
						onchange={() => (setsCount = n)}
					/>
				{/each}
			</div>
		</div>

		<div class="max-h-[60vh] overflow-y-auto">
			{#if favouriteExercises.length > 0}
				<p
					class="text-base-content/50 px-3 pt-2 pb-1 text-xs font-semibold tracking-wider uppercase"
				>
					Favourites
				</p>
				{#each favouriteExercises as ex (ex.id)}
					<div class="hover:bg-base-200 flex w-full items-center gap-1 rounded-lg px-3 py-2.5">
						<button
							class="flex-1 cursor-pointer text-left"
							disabled={isAddingExercise}
							onclick={() => addExercise(ex.id)}
						>
							{ex.name}
						</button>
						<button
							class="btn btn-ghost btn-xs btn-square"
							onclick={() => toggleFavourite(ex.id)}
						>
							<Icon icon="material-symbols:star" class="text-warning size-4" />
						</button>
					</div>
				{/each}
				<div class="divider my-0"></div>
			{/if}
			{#each filteredExercises as ex (ex.id)}
				<div class="hover:bg-base-200 flex w-full items-center gap-1 rounded-lg px-3 py-2.5">
					<button
						class="flex-1 cursor-pointer text-left"
						disabled={isAddingExercise}
						onclick={() => addExercise(ex.id)}
					>
						{ex.name}
					</button>
					<button
						class="btn btn-ghost btn-xs btn-square"
						onclick={() => toggleFavourite(ex.id)}
					>
						<Icon
							icon={favouriteIds.includes(ex.id)
								? 'material-symbols:star'
								: 'material-symbols:star-outline'}
							class={favouriteIds.includes(ex.id)
								? 'text-warning size-4'
								: 'text-base-content/30 size-4'}
						/>
					</button>
				</div>
			{/each}
		</div>
	</div>
{/snippet}

<BitsDialog bind:open={exerciseDialogOpen} title="Add Exercise">
	{@render exerciseListContent()}
</BitsDialog>

<!-- Rename dialog -->
<BitsDialog bind:open={renameDialogOpen} title="Rename Routine">
	<div class="grid gap-4">
		<input type="text" class="input input-lg w-full" bind:value={editName} />
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isSavingName || !editName.trim()}
			onclick={saveRename}
		>
			{#if isSavingName}<span class="loading loading-spinner loading-sm"></span>{/if}
			Save
		</button>
	</div>
</BitsDialog>

<!-- Edit sets count dialog -->
<BitsDialog bind:open={editSetsDialogOpen} title="Sets">
	<div class="grid gap-4">
		<p class="text-base-content/70">
			{editingSetsExercise ? getExerciseName(editingSetsExercise.exerciseId) : ''}
		</p>
		<div class="join w-full">
			{#each [1, 2, 3, 4, 5, 6] as n (n)}
				<input
					type="radio"
					name="edit-sets-count"
					class="btn join-item checked:bg-segmented checked:text-primary-content flex-1"
					aria-label={String(n)}
					checked={editSetsCount === n}
					onchange={() => (editSetsCount = n)}
				/>
			{/each}
		</div>
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isSavingSets}
			onclick={saveEditSets}
		>
			{#if isSavingSets}<span class="loading loading-spinner loading-sm"></span>{/if}
			Save
		</button>
	</div>
</BitsDialog>

<!-- Remove exercise dialog -->
<DeleteConfirmDialog bind:open={removeDialogOpen} title="Remove exercise?" confirmLabel="Remove" isLoading={isRemoving} onconfirm={confirmRemove}>
	{#if removingExercise}
		<p class="text-base-content/60">
			Remove <span class="text-base-content font-semibold"
				>{getExerciseName(removingExercise.exerciseId)}</span
			> from this routine?
		</p>
	{/if}
</DeleteConfirmDialog>

<!-- Delete routine dialog -->
<DeleteConfirmDialog bind:open={deleteDialogOpen} title="Delete routine?" isLoading={isDeleting} onconfirm={confirmDelete}>
	<p class="text-base-content/60">
		This will permanently delete <span class="text-base-content font-semibold"
			>{currentRoutine?.name}</span
		> and all its exercises.
	</p>
</DeleteConfirmDialog>
