<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import updateLocale from 'dayjs/plugin/updateLocale';
	import isToday from 'dayjs/plugin/isToday';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		allWorkoutsQueryOptions,
		getAllWorkoutsQueryKey,
		favouriteExercisesQueryOptions,
		getFavouriteExercisesQueryKey
	} from '$lib/queries';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { exercises } from '$lib/exercises';
	import Icon from '@iconify/svelte';
	import CorgiGym from '$lib/assets/corgi_gym.webp?w=240&enhanced';

	let { data } = $props();

	dayjs.extend(relativeTime);
	dayjs.extend(updateLocale);
	dayjs.extend(isToday);

	dayjs.updateLocale('en', {
		relativeTime: {
			future: 'in %s',
			past: '%s ago',
			s: 'a few seconds',
			m: '1 minute',
			mm: '%d minutes',
			h: '1 hour',
			hh: '%d hours',
			d: '1 day',
			dd: '%d days',
			M: '1 month',
			MM: '%d months',
			y: '1 year',
			yy: '%d years'
		}
	});

	const queryClient = useQueryClient();
	const workoutsDb = createQuery(allWorkoutsQueryOptions);
	const favouritesDb = createQuery(favouriteExercisesQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

	function updateWorkoutsCache(updater: (workouts: WorkoutDB[]) => WorkoutDB[]) {
		const current = queryClient.getQueryData<WorkoutDB[]>(getAllWorkoutsQueryKey());
		if (current) queryClient.setQueryData(getAllWorkoutsQueryKey(), updater(current));
	}

	let exerciseSearch = $state('');
	let addSetDialog = $state<HTMLDialogElement>();
	let exerciseDialog = $state<HTMLDialogElement>();
	let filteredExercises = $derived(
		exerciseSearch
			? exercises.filter((e) => e.name.toLowerCase().includes(exerciseSearch.toLowerCase()))
			: exercises
	);

	let favouriteIds = $derived(favouritesDb.data?.exerciseIds ?? []);
	let favouriteExercises = $derived(exercises.filter((e) => favouriteIds.includes(e.id)));

	async function toggleFavourite(exerciseId: string) {
		const response = await api.post('gym/favourites', {
			body: JSON.stringify({ exerciseId })
		});
		if (response.ok) {
			const data: FavouriteExercisesDB = await response.json();
			queryClient.setQueryData(getFavouriteExercisesQueryKey(), data);
		}
	}

	function pickExercise(id: string) {
		selectedExerciseId = id;
		const ex = exerciseMap.get(id);
		weightMode = ex?.equipment === 'dumbbell' ? 'each' : 'total';
		setWeight = null;
		setReps = null;
		setType = 'working';
		exerciseDialog?.close();
		addSetDialog?.showModal();
	}

	function openAddSet(workoutId: string, group: { exerciseId: string; sets: SetDB[] }) {
		addingSetToWorkout = workoutId;
		selectedExerciseId = group.exerciseId;
		const ex = exerciseMap.get(group.exerciseId);
		weightMode = ex?.equipment === 'dumbbell' ? 'each' : 'total';
		const lastSet = group.sets[group.sets.length - 1];
		if (lastSet) {
			setWeight =
				lastSet.weightKg != null
					? weightUnit === 'lb'
						? kgToLb(lastSet.weightKg)
						: lastSet.weightKg
					: null;
			setReps = lastSet.reps;
			setType = lastSet.setType;
		} else {
			setWeight = null;
			setReps = null;
			setType = 'working';
		}
		addSetDialog?.showModal();
	}

	let addingSetToWorkout = $state<string | null>(null);
	let selectedExerciseId = $state('');
	let setWeight = $state<number | null>(null);
	let setReps = $state<number | null>(null);
	let setType = $state('working');
	let weightUnit = $state<'kg' | 'lb'>('kg');
	let weightMode = $state<'total' | 'each'>('total');
	let isSubmittingSet = $state(false);
	let reorderingSetId = $state<string | null>(null);
	let isDeletingSet = $state(false);
	let isDeletingWorkout = $state(false);
	let isSavingEdit = $state(false);
	let isSavingNotes = $state(false);
	let notesDialog = $state<HTMLDialogElement>();
	let editingNotes = $state('');

	function openEditNotes() {
		editingNotes = currentWorkout?.notes ?? '';
		notesDialog?.showModal();
	}

	async function saveNotes() {
		if (!currentWorkout) return;
		isSavingNotes = true;
		const response = await api.patch(`gym/workouts/${currentWorkout.id}`, {
			body: JSON.stringify({
				startTime: currentWorkout.startTime,
				notes: editingNotes || null
			})
		});
		isSavingNotes = false;
		if (response.status === 204) {
			updateWorkoutsCache((workouts) =>
				workouts.map((w) =>
					w.id === currentWorkout.id ? { ...w, notes: editingNotes || null } : w
				)
			);
			notesDialog?.close();
		} else {
			addToast('error', 'Failed to save notes');
		}
	}

	const kgToLb = (kg: number) => Math.round(kg * 2.20462 * 10) / 10;
	const lbToKg = (lb: number) => Math.round((lb / 2.20462) * 10) / 10;

	function groupSetsByExercise(sets: SetDB[]) {
		const groups: { exerciseId: string; exerciseName: string; sets: SetDB[] }[] = [];
		const map = new Map<string, SetDB[]>();
		const order: string[] = [];

		for (const set of sets) {
			if (!map.has(set.exerciseId)) {
				map.set(set.exerciseId, []);
				order.push(set.exerciseId);
			}
			map.get(set.exerciseId)!.push(set);
		}

		for (const exerciseId of order) {
			groups.push({
				exerciseId,
				exerciseName: getExerciseName(exerciseId),
				sets: map.get(exerciseId)!
			});
		}

		return groups;
	}

	let currentWorkout = $derived.by(() => {
		if (workoutsDb.isPending || !workoutsDb.isSuccess) return;
		return workoutsDb.data.find((w) => w.id === data.workoutId);
	});

	let deletingWorkoutId = $state<string | null>(null);
	let deleteWorkoutDialog = $state<HTMLDialogElement>();

	function openDeleteWorkout(workoutId: string) {
		deletingWorkoutId = workoutId;
		deleteWorkoutDialog?.showModal();
	}

	async function confirmDeleteWorkout() {
		if (!deletingWorkoutId) return;
		isDeletingWorkout = true;
		const response = await api.delete(`gym/workouts/${deletingWorkoutId}`);
		isDeletingWorkout = false;
		if (response.status === 204) {
			addToast('success', 'Workout deleted');
			const deletedId = deletingWorkoutId;
			updateWorkoutsCache((workouts) => workouts.filter((w) => w.id !== deletedId));
			deleteWorkoutDialog?.close();
			deletingWorkoutId = null;
			goto('/app/gym');
		} else {
			addToast('error', 'Failed to delete workout');
		}
	}

	async function addSet(workoutId: string) {
		if (!selectedExerciseId) return;
		isSubmittingSet = true;

		const displayWeight = setWeight != null && weightMode === 'each' ? setWeight * 2 : setWeight;
		const weightKg =
			displayWeight != null && weightUnit === 'lb' ? lbToKg(displayWeight) : displayWeight;

		const response = await api.post(`gym/workouts/${workoutId}/sets`, {
			body: JSON.stringify({
				exerciseId: selectedExerciseId,
				weightKg,
				reps: setReps,
				setType
			})
		});

		isSubmittingSet = false;

		if (response.status === 201) {
			const set: SetDB = await response.json();
			updateWorkoutsCache((workouts) =>
				workouts.map((w) => (w.id === workoutId ? { ...w, sets: [...w.sets, set] } : w))
			);
			selectedExerciseId = '';
			setWeight = null;
			setReps = null;
			setType = 'working';
			weightMode = 'total';
			addingSetToWorkout = null;
			addSetDialog?.close();
		} else {
			addToast('error', 'Failed to add set');
		}
	}

	async function duplicateSet(workoutId: string, set: SetDB) {
		const response = await api.post(`gym/workouts/${workoutId}/sets`, {
			body: JSON.stringify({
				exerciseId: set.exerciseId,
				weightKg: set.weightKg,
				reps: set.reps,
				setType: set.setType
			})
		});

		if (response.status === 201) {
			const newSet: SetDB = await response.json();
			updateWorkoutsCache((workouts) =>
				workouts.map((w) => (w.id === workoutId ? { ...w, sets: [...w.sets, newSet] } : w))
			);
		} else {
			addToast('error', 'Failed to add set');
		}
	}

	async function reorderSet(setId: string, direction: 'up' | 'down') {
		reorderingSetId = setId;

		updateWorkoutsCache((workouts) =>
			workouts.map((w) => {
				const current = w.sets.find((s) => s.id === setId);
				if (!current) return w;
				const sameEx = w.sets
					.filter((s) => s.exerciseId === current.exerciseId)
					.sort((a, b) => a.position - b.position);
				const idx = sameEx.findIndex((s) => s.id === setId);
				const neighborIdx = direction === 'up' ? idx - 1 : idx + 1;
				if (neighborIdx < 0 || neighborIdx >= sameEx.length) return w;
				const neighbor = sameEx[neighborIdx];
				return {
					...w,
					sets: w.sets.map((s) => {
						if (s.id === current.id) return { ...s, position: neighbor.position };
						if (s.id === neighbor.id) return { ...s, position: current.position };
						return s;
					})
				};
			})
		);

		const response = await api.post('gym/sets/reorder', {
			body: JSON.stringify({ setId, direction })
		});
		reorderingSetId = null;
		if (response.status !== 204) {
			addToast('error', 'Failed to reorder set');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		}
	}

	let deletingSet = $state<SetDB | null>(null);
	let deleteSetDialog = $state<HTMLDialogElement>();

	function openDeleteSet(set: SetDB) {
		deletingSet = set;
		deleteSetDialog?.showModal();
	}

	async function confirmDeleteSet() {
		if (!deletingSet) return;
		isDeletingSet = true;
		const response = await api.delete(`gym/sets/${deletingSet.id}`);
		isDeletingSet = false;
		if (response.status === 204) {
			const deletedId = deletingSet.id;
			updateWorkoutsCache((workouts) =>
				workouts.map((w) => ({ ...w, sets: w.sets.filter((s) => s.id !== deletedId) }))
			);
			deleteSetDialog?.close();
			deletingSet = null;
		} else {
			addToast('error', 'Failed to delete set');
		}
	}

	let editingSet = $state<SetDB | null>(null);
	let editWeight = $state<number | null>(null);
	let editReps = $state<number | null>(null);
	let editDialog = $state<HTMLDialogElement>();

	function openEditSet(set: SetDB) {
		editingSet = set;
		editWeight =
			set.weightKg != null ? (weightUnit === 'lb' ? kgToLb(set.weightKg) : set.weightKg) : null;
		editReps = set.reps;
		editDialog?.showModal();
	}

	async function saveEditSet() {
		if (!editingSet) return;
		isSavingEdit = true;
		const weightKg = editWeight != null && weightUnit === 'lb' ? lbToKg(editWeight) : editWeight;
		const response = await api.patch(`gym/sets/${editingSet.id}`, {
			body: JSON.stringify({
				exerciseId: editingSet.exerciseId,
				weightKg,
				reps: editReps,
				setType: editingSet.setType
			})
		});
		isSavingEdit = false;
		if (response.status === 204) {
			const updatedSet = { ...editingSet, weightKg, reps: editReps };
			updateWorkoutsCache((workouts) =>
				workouts.map((w) => ({
					...w,
					sets: w.sets.map((s) => (s.id === updatedSet.id ? updatedSet : s))
				}))
			);
			editDialog?.close();
			editingSet = null;
		} else {
			addToast('error', 'Failed to update set');
		}
	}
</script>

<PageWrapper title="Gym">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				{#if workoutsDb.isSuccess && currentWorkout}
					{@const exerciseGroups = groupSetsByExercise(currentWorkout.sets)}
					<div class="flex items-center justify-between">
						<div class="flex items-center justify-between px-4 pt-3 pb-2">
							<div class="grid grow gap-0">
								<div
									class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
								>
									{dayjs(currentWorkout.startTime).fromNow()}
									{#if dayjs(currentWorkout.startTime).isToday()}
										<div
											class="bg-primary/10 text-primary rounded-full px-2 py-0.5 text-xs font-bold"
										>
											Today
										</div>
									{/if}
								</div>
								<div class="text-base-content/65 flex items-center gap-2 text-sm font-light">
									<span>{dayjs(currentWorkout.startTime).format('ddd, D MMM')}</span>
									<span class="text-base-content/40 text-xs font-extralight">•</span>
									<span>{dayjs(currentWorkout.startTime).format('h:mma')}</span>
								</div>
							</div>
						</div>
						<div
							class={[
								'join text-sm',
								!workoutsDb.data || workoutsDb.data.length === 0 ? 'hidden' : undefined
							]}
						>
							<input
								type="radio"
								name="weight-unit"
								class="btn btn-sm join-item checked:bg-segmented checked:text-primary-content"
								aria-label="kg"
								checked={weightUnit === 'kg'}
								onchange={() => (weightUnit = 'kg')}
							/>
							<input
								type="radio"
								name="weight-unit"
								class="btn btn-sm join-item checked:bg-segmented checked:text-primary-content"
								aria-label="lb"
								checked={weightUnit === 'lb'}
								onchange={() => (weightUnit = 'lb')}
							/>
						</div>
					</div>

					<div class="px-4">
						{#if currentWorkout.notes}
							<div class="flex items-center gap-2 rounded-lg py-2 italic">
								<p class="text-base-content/70 text-sm whitespace-pre-wrap">
									{currentWorkout.notes}
								</p>
								<button class="btn btn-ghost btn-xs btn-square shrink-0" onclick={openEditNotes}>
									<Icon
										icon="material-symbols:edit-outline"
										class="text-base-content/50 size-3.5"
									/>
								</button>
							</div>
						{:else}
							<button
								class="btn btn-ghost btn-sm text-base-content/50 w-full gap-1 rounded-lg"
								onclick={openEditNotes}
							>
								<Icon icon="material-symbols:note-add-outline" class="size-4" />
								Add notes
							</button>
						{/if}
					</div>

					{#if exerciseGroups.length > 0}
						{#each exerciseGroups as group, gi (group.exerciseId)}
							<div class="border-base-300/50 bg-base-50 rounded-2xl border px-4">
								<div
									class={[
										'py-3',
										gi < exerciseGroups.length - 1 && 'border-b-base-300/50 border-b'
									]}
								>
									<div class="mb-2 flex items-center justify-between">
										<p class="text-base-content/80 font-semibold tracking-wide">
											{group.exerciseName}
										</p>
										<a
											href="/app/gym/stats/exercises/{group.exerciseId}"
											class="btn btn-ghost btn-xs"
										>
											<Icon icon="material-symbols:show-chart" class="size-4" />
											History
										</a>
									</div>
									<div class="grid gap-2">
										{#each group.sets as set, si (set.id)}
											<div class="flex items-center gap-3">
												<span class="text-base-content/40 w-5 text-right text-sm font-medium"
													>{si + 1}</span
												>
												<div
													class="bg-base-200/60 flex grow items-center justify-between rounded-lg px-3 py-2"
												>
													<span>
														{#if set.weightKg !== null}
															<span class="font-semibold"
																>{weightUnit === 'lb'
																	? kgToLb(set.weightKg)
																	: set.weightKg}{weightUnit}</span
															>
														{/if}
														{#if set.reps !== null}
															<span class="text-base-content/60">× {set.reps}</span>
														{/if}
														{#if set.setType !== 'working'}
															<span class="text-base-content/30 ml-1">({set.setType})</span>
														{/if}
													</span>
													<div class="flex items-center gap-1">
														{#if si > 0}
															<button
																class="btn btn-ghost btn-xs btn-square"
																disabled={reorderingSetId === set.id}
																onclick={() => reorderSet(set.id, 'up')}
															>
																<Icon
																	icon="material-symbols:arrow-upward"
																	class="text-base-content/50 size-3.5"
																/>
															</button>
														{/if}
														{#if si < group.sets.length - 1}
															<button
																class="btn btn-ghost btn-xs btn-square"
																disabled={reorderingSetId === set.id}
																onclick={() => reorderSet(set.id, 'down')}
															>
																<Icon
																	icon="material-symbols:arrow-downward"
																	class="text-base-content/50 size-3.5"
																/>
															</button>
														{/if}
														<button
															class="btn btn-ghost btn-xs btn-square"
															onclick={() => openEditSet(set)}
														>
															<Icon
																icon="material-symbols:edit-outline"
																class="text-base-content/50 size-3.5"
															/>
														</button>
														<button
															class="btn btn-ghost btn-xs btn-square"
															onclick={() => duplicateSet(currentWorkout.id, set)}
														>
															<Icon
																icon="material-symbols:replay"
																class="text-base-content/50 size-3.5"
															/>
														</button>
														<button
															class="btn btn-ghost btn-xs btn-square"
															onclick={() => openDeleteSet(set)}
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
								</div>
								<div class="border-t-base-300/50 border-t px-4 py-3">
									<button
										class="btn btn-ghost btn-sm text-primary w-full"
										onclick={() => openAddSet(currentWorkout.id, group)}
									>
										<Icon icon="material-symbols:add" class="size-4" />
										Add Set
									</button>
								</div>
							</div>
						{/each}
					{/if}

					<button
						class="btn btn-primary btn-soft w-full rounded-full"
						onclick={() => {
							addingSetToWorkout = currentWorkout.id;
							exerciseSearch = '';
							exerciseDialog?.showModal();
						}}
					>
						<Icon icon="material-symbols:add" class="size-4" />
						Add Exercise
					</button>
				{:else if workoutsDb.isSuccess && !currentWorkout}
					<div class="grid justify-items-center gap-4 py-16">
						<enhanced:img src={CorgiGym} alt="Corgi Gymming" />
						<p class="text-base-content/50 text-lg font-medium">Workout not found</p>
						<a href="/app/gym" class="btn btn-primary btn-soft rounded-full">Back to Gym</a>
					</div>
				{:else}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{/if}

				{#if currentWorkout}
					<button
						class="btn btn-ghost btn-error mt-2 flex w-full items-center gap-2 rounded-full"
						onclick={() => openDeleteWorkout(currentWorkout.id)}
					>
						<Icon icon="material-symbols:delete-outline" class="text-error size-5" />Delete Workout
					</button>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>

<dialog bind:this={addSetDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-4">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-bold">
				{selectedExerciseId ? getExerciseName(selectedExerciseId) : 'Add Set'}
			</h3>
			<form method="dialog">
				<button class="btn btn-ghost btn-sm btn-square">
					<Icon icon="material-symbols:close" class="size-5" />
				</button>
			</form>
		</div>

		<div class="flex items-center gap-2">
			<input
				type="number"
				class="input input-bordered w-full rounded-lg"
				placeholder="Reps"
				bind:value={setReps}
			/>
			<span class="text-base-content/40 text-lg font-bold">×</span>
			<div class="flex grow">
				<input
					type="number"
					class="input input-bordered w-full rounded-lg"
					placeholder={weightUnit}
					bind:value={setWeight}
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
		{#if weightMode === 'each' && setWeight != null}
			<p class="text-base-content/50 -mt-2 text-sm">
				{setWeight} × 2 = {setWeight * 2}
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
			disabled={isSubmittingSet}
			onclick={() => addingSetToWorkout && addSet(addingSetToWorkout)}
		>
			{#if isSubmittingSet}<span class="loading loading-spinner loading-sm"></span>{/if}
			Add Set
		</button>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>

<dialog bind:this={exerciseDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box flex max-h-[80vh] flex-col gap-3">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-bold">Select Exercise</h3>
			<form method="dialog">
				<button class="btn btn-ghost btn-sm btn-square">
					<Icon icon="material-symbols:close" class="size-5" />
				</button>
			</form>
		</div>
		<input
			type="text"
			class="input input-lg w-full"
			placeholder="Search exercises..."
			bind:value={exerciseSearch}
		/>
		<div class="flex-1 overflow-y-auto">
			{#if favouriteExercises.length > 0}
				<p
					class="text-base-content/50 px-3 pt-2 pb-1 text-xs font-semibold tracking-wider uppercase"
				>
					Favourites
				</p>
				{#each favouriteExercises as ex (ex.id)}
					<div class="hover:bg-base-200 flex w-full items-center gap-1 px-3 py-2.5">
						<button class="flex-1 cursor-pointer text-left" onclick={() => pickExercise(ex.id)}>
							{ex.name}
						</button>
						<button class="btn btn-ghost btn-xs btn-square" onclick={() => toggleFavourite(ex.id)}>
							<Icon icon="material-symbols:star" class="text-warning size-4" />
						</button>
					</div>
				{/each}
				<div class="divider my-0"></div>
			{/if}
			{#each filteredExercises as ex (ex.id)}
				<div class="hover:bg-base-200 flex w-full items-center gap-1 px-3 py-2.5">
					<button class="flex-1 cursor-pointer text-left" onclick={() => pickExercise(ex.id)}>
						{ex.name}
					</button>
					<button class="btn btn-ghost btn-xs btn-square" onclick={() => toggleFavourite(ex.id)}>
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
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>

<dialog bind:this={editDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-4">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-bold">
				{editingSet ? getExerciseName(editingSet.exerciseId) : ''}
			</h3>
			<form method="dialog">
				<button class="btn btn-ghost btn-sm btn-square">
					<Icon icon="material-symbols:close" class="size-5" />
				</button>
			</form>
		</div>

		<div class="flex items-center gap-2">
			<input
				type="number"
				class="input input-bordered w-full rounded-lg"
				placeholder="Reps"
				bind:value={editReps}
			/>
			<span class="text-base-content/40 text-lg font-bold">×</span>
			<div class="flex grow">
				<input
					type="number"
					class="input input-bordered w-full rounded-lg"
					placeholder={weightUnit}
					bind:value={editWeight}
				/>
				<button
					class="btn btn-ghost"
					onclick={() => {
						if (editWeight != null) {
							editWeight = weightUnit === 'kg' ? kgToLb(editWeight) : lbToKg(editWeight);
						}
						weightUnit = weightUnit === 'kg' ? 'lb' : 'kg';
					}}
				>
					{weightUnit}
				</button>
			</div>
		</div>

		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isSavingEdit}
			onclick={saveEditSet}
		>
			{#if isSavingEdit}<span class="loading loading-spinner loading-sm"></span>{/if}
			Save
		</button>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>

<dialog bind:this={notesDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-4">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-bold">Notes</h3>
			<form method="dialog">
				<button class="btn btn-ghost btn-sm btn-square">
					<Icon icon="material-symbols:close" class="size-5" />
				</button>
			</form>
		</div>
		<textarea
			class="textarea textarea-bordered w-full resize-none"
			rows="4"
			placeholder="Add notes about this workout..."
			bind:value={editingNotes}
		></textarea>
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isSavingNotes}
			onclick={saveNotes}
		>
			{#if isSavingNotes}<span class="loading loading-spinner loading-sm"></span>{/if}
			Save
		</button>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>

<dialog bind:this={deleteSetDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-8">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<div class="grid gap-4">
			<div
				class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center rounded-full"
			>
				<Icon icon="material-symbols:delete-outline" class="size-10" />
			</div>
			<h2 class="text-2xl font-bold">Delete this set?</h2>
			{#if deletingSet}
				<p class="text-base-content/60">
					This will permanently remove the
					<span class="text-base-content font-semibold">
						{#if deletingSet.weightKg != null}{weightUnit === 'lb'
								? kgToLb(deletingSet.weightKg)
								: deletingSet.weightKg}{weightUnit}{/if}{#if deletingSet.weightKg != null && deletingSet.reps != null}
							×
						{/if}{#if deletingSet.reps != null}{deletingSet.reps} reps{/if}
					</span>
					set from {getExerciseName(deletingSet.exerciseId)}.
				</p>
			{/if}
		</div>
		<div class="grid gap-4">
			<button class="btn btn-error btn-lg" disabled={isDeletingSet} onclick={confirmDeleteSet}>
				{#if isDeletingSet}<span class="loading loading-spinner loading-sm"></span>{/if}
				Delete
			</button>
			<button
				class="btn btn-outline btn-neutral btn-lg w-full"
				onclick={() => deleteSetDialog?.close()}>Cancel</button
			>
		</div>
	</div>
</dialog>

<dialog bind:this={deleteWorkoutDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-8">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<div class="grid gap-4">
			<div
				class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center rounded-full"
			>
				<Icon icon="material-symbols:delete-outline" class="size-10" />
			</div>
			<h2 class="text-2xl font-bold">Delete workout?</h2>
			<p class="text-base-content/60">All sets in this workout will be deleted.</p>
		</div>
		<div class="grid gap-4">
			<button
				class="btn btn-error btn-lg"
				disabled={isDeletingWorkout}
				onclick={confirmDeleteWorkout}
			>
				{#if isDeletingWorkout}<span class="loading loading-spinner loading-sm"></span>{/if}
				Delete
			</button>
			<button
				class="btn btn-outline btn-neutral btn-lg w-full"
				onclick={() => deleteWorkoutDialog?.close()}>Cancel</button
			>
		</div>
	</div>
</dialog>
