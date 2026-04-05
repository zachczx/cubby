<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import updateLocale from 'dayjs/plugin/updateLocale';
	import isToday from 'dayjs/plugin/isToday';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { allWorkoutsQueryOptions, getAllWorkoutsQueryKey } from '$lib/queries';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { exercises } from '$lib/exercises';
	import Icon from '@iconify/svelte';
	import CorgiGym from '$lib/assets/corgi_gym.webp?w=240&enhanced';
	import { detectPr, type PrResult } from '$lib/pr';
	import StarBurst from '$lib/ui/StarBurst.svelte';
	import DeleteConfirmDialog from '$lib/ui/DeleteConfirmDialog.svelte';
	import ExercisePickerDialog from '$lib/ui/ExercisePickerDialog.svelte';
	import AddSetDialog from '$lib/ui/AddSetDialog.svelte';
	import EditSetDialog from '$lib/ui/EditSetDialog.svelte';
	import NotesDialog from '$lib/ui/NotesDialog.svelte';

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

	const exerciseMap = new Map(exercises.map((e) => [e.id.toLowerCase(), e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId.toLowerCase())?.name ?? exerciseId;
	}

	function updateWorkoutsCache(updater: (workouts: WorkoutDB[]) => WorkoutDB[]) {
		const current = queryClient.getQueryData<WorkoutDB[]>(getAllWorkoutsQueryKey());
		if (current) queryClient.setQueryData(getAllWorkoutsQueryKey(), updater(current));
	}

	function getAllSetsForExercise(exerciseId: string): SetDB[] {
		const workouts = queryClient.getQueryData<WorkoutDB[]>(getAllWorkoutsQueryKey());
		if (!workouts) return [];
		return workouts
			.sort((a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime())
			.flatMap((w) => w.sets.filter((s) => s.exerciseId === exerciseId));
	}

	const kgToLb = (kg: number) => Math.round(kg * 2.20462 * 10) / 10;

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

	// Map of setId → PrResult for sets in the current workout
	let prMap = $derived.by(() => {
		const map = new Map<string, PrResult>();
		if (!workoutsDb.isSuccess || !currentWorkout) return map;

		// All workouts sorted oldest first
		const sorted = [...workoutsDb.data].sort(
			(a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime()
		);

		for (const set of currentWorkout.sets) {
			// Collect all previous sets for this exercise across all workouts
			const historical: SetDB[] = [];
			for (const w of sorted) {
				for (const s of w.sets) {
					if (s.exerciseId !== set.exerciseId) continue;
					if (w.id === currentWorkout.id && s.id === set.id) break;
					if (w.id === currentWorkout.id) {
						historical.push(s);
						continue;
					}
					historical.push(s);
				}
				if (w.id === currentWorkout.id) break;
			}

			const pr = detectPr(set, historical);
			if (pr) map.set(set.id, pr);
		}
		return map;
	});

	// Shared state
	let weightUnit = $state<'kg' | 'lb'>('kg');
	let prStarSetId = $state<string | null>(null);

	function triggerPrCelebration(setId: string) {
		prStarSetId = setId;
		setTimeout(() => (prStarSetId = null), 2000);
	}

	// Exercise picker
	let exerciseDialogOpen = $state(false);

	function onPickExercise(exerciseId: string) {
		selectedExerciseId = exerciseId;
		const ex = exerciseMap.get(exerciseId.toLowerCase());
		exerciseDialogOpen = false;
		addSetRef?.prefill({
			weight: null,
			reps: null,
			setType: 'working',
			weightMode: ex?.equipment === 'dumbbell' ? 'each' : 'total'
		});
		addSetDialogOpen = true;
	}

	// Add set
	let addSetDialogOpen = $state(false);
	let addSetRef = $state<ReturnType<typeof AddSetDialog> | null>(null);
	let addingSetToWorkout = $state<string | null>(null);
	let selectedExerciseId = $state('');
	let isSubmittingSet = $state(false);

	function openAddSet(workoutId: string, group: { exerciseId: string; sets: SetDB[] }) {
		addingSetToWorkout = workoutId;
		selectedExerciseId = group.exerciseId;
		const ex = exerciseMap.get(group.exerciseId.toLowerCase());
		const lastSet = group.sets[group.sets.length - 1];
		addSetRef?.prefill({
			weight: lastSet?.weightKg != null
				? weightUnit === 'lb'
					? kgToLb(lastSet.weightKg)
					: lastSet.weightKg
				: null,
			reps: lastSet?.reps ?? null,
			setType: lastSet?.setType ?? 'working',
			weightMode: ex?.equipment === 'dumbbell' ? 'each' : 'total'
		});
		addSetDialogOpen = true;
	}

	async function handleAddSet(setData: { weightKg: number | null; reps: number | null; setType: string }) {
		if (!addingSetToWorkout || !selectedExerciseId) return;
		isSubmittingSet = true;

		const response = await api.post(`gym/workouts/${addingSetToWorkout}/sets`, {
			json: {
				exerciseId: selectedExerciseId,
				weightKg: setData.weightKg,
				reps: setData.reps,
				setType: setData.setType
			}
		});

		isSubmittingSet = false;

		if (response.status === 201) {
			const set: SetDB = await response.json();
			const historical = getAllSetsForExercise(set.exerciseId);
			const pr = detectPr(set, historical);

			updateWorkoutsCache((workouts) =>
				workouts.map((w) => (w.id === addingSetToWorkout ? { ...w, sets: [...w.sets, set] } : w))
			);
			if (pr) triggerPrCelebration(set.id);
			selectedExerciseId = '';
			addingSetToWorkout = null;
			addSetDialogOpen = false;
			addSetRef?.reset();
		} else {
			addToast('error', 'Failed to add set');
		}
	}

	async function duplicateSet(workoutId: string, set: SetDB) {
		const response = await api.post(`gym/workouts/${workoutId}/sets`, {
			json: {
				exerciseId: set.exerciseId,
				weightKg: set.weightKg,
				reps: set.reps,
				setType: set.setType
			}
		});

		if (response.status === 201) {
			const newSet: SetDB = await response.json();
			const historical = getAllSetsForExercise(newSet.exerciseId);
			const pr = detectPr(newSet, historical);

			updateWorkoutsCache((workouts) =>
				workouts.map((w) => (w.id === workoutId ? { ...w, sets: [...w.sets, newSet] } : w))
			);
			if (pr) triggerPrCelebration(newSet.id);
		} else {
			addToast('error', 'Failed to add set');
		}
	}

	// Reorder set
	let reorderingSetId = $state<string | null>(null);

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
			json: { setId, direction }
		});
		reorderingSetId = null;
		if (response.status !== 204) {
			addToast('error', 'Failed to reorder set');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		}
	}

	// Delete set
	let deletingSet = $state<SetDB | null>(null);
	let deleteSetDialogOpen = $state(false);
	let isDeletingSet = $state(false);

	function openDeleteSet(set: SetDB) {
		deletingSet = set;
		deleteSetDialogOpen = true;
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
			deleteSetDialogOpen = false;
			deletingSet = null;
		} else {
			addToast('error', 'Failed to delete set');
		}
	}

	// Edit set
	let editingSet = $state<SetDB | null>(null);
	let editDialogOpen = $state(false);
	let editSetRef = $state<ReturnType<typeof EditSetDialog> | null>(null);
	let isSavingEdit = $state(false);

	function openEditSet(set: SetDB) {
		editingSet = set;
		editSetRef?.prefill({
			weight:
				set.weightKg != null ? (weightUnit === 'lb' ? kgToLb(set.weightKg) : set.weightKg) : null,
			reps: set.reps,
			setType: set.setType
		});
		editDialogOpen = true;
	}

	async function handleEditSet(setData: { weightKg: number | null; reps: number | null; setType: string }) {
		if (!editingSet) return;
		isSavingEdit = true;
		const response = await api.patch(`gym/sets/${editingSet.id}`, {
			json: {
				exerciseId: editingSet.exerciseId,
				weightKg: setData.weightKg,
				reps: setData.reps,
				setType: setData.setType
			}
		});
		isSavingEdit = false;
		if (response.status === 204) {
			const updatedSet = { ...editingSet, weightKg: setData.weightKg, reps: setData.reps, setType: setData.setType };
			updateWorkoutsCache((workouts) =>
				workouts.map((w) => ({
					...w,
					sets: w.sets.map((s) => (s.id === updatedSet.id ? updatedSet : s))
				}))
			);
			editDialogOpen = false;
			editingSet = null;
		} else {
			addToast('error', 'Failed to update set');
		}
	}

	// Notes
	let notesDialogOpen = $state(false);
	let notesRef = $state<ReturnType<typeof NotesDialog> | null>(null);
	let isSavingNotes = $state(false);

	function openEditNotes() {
		notesRef?.prefill(currentWorkout?.notes ?? '');
		notesDialogOpen = true;
	}

	async function handleSaveNotes(notes: string) {
		if (!currentWorkout) return;
		isSavingNotes = true;
		const response = await api.patch(`gym/workouts/${currentWorkout.id}`, {
			json: {
				startTime: currentWorkout.startTime,
				notes: notes || null
			}
		});
		isSavingNotes = false;
		if (response.status === 204) {
			updateWorkoutsCache((workouts) =>
				workouts.map((w) =>
					w.id === currentWorkout.id ? { ...w, notes: notes || null } : w
				)
			);
			notesDialogOpen = false;
		} else {
			addToast('error', 'Failed to save notes');
		}
	}

	// Delete workout
	let deletingWorkoutId = $state<string | null>(null);
	let deleteWorkoutDialogOpen = $state(false);
	let isDeletingWorkout = $state(false);

	function openDeleteWorkout(workoutId: string) {
		deletingWorkoutId = workoutId;
		deleteWorkoutDialogOpen = true;
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
			deleteWorkoutDialogOpen = false;
			deletingWorkoutId = null;
			goto('/app/gym');
		} else {
			addToast('error', 'Failed to delete workout');
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
											href="/app/gym/stats/exercises/{encodeURIComponent(group.exerciseId)}"
											class="btn btn-ghost btn-xs"
										>
											<Icon icon="material-symbols:show-chart" class="size-4" />
											History
										</a>
									</div>
									<div class="grid gap-2">
										{#each group.sets as set, si (set.id)}
											<div class="relative flex items-center gap-3">
												<StarBurst show={prStarSetId === set.id} size="14px" count={8} />
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
														{#if prMap.get(set.id)}
															<span class="bg-warning/15 text-warning ml-1 rounded-full px-1.5 py-0.5 text-xs font-bold">
																{prMap.get(set.id)?.label}
															</span>
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
							exerciseDialogOpen = true;
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

<ExercisePickerDialog bind:open={exerciseDialogOpen} onpick={onPickExercise} />

<AddSetDialog
	bind:this={addSetRef}
	bind:open={addSetDialogOpen}
	exerciseName={selectedExerciseId ? getExerciseName(selectedExerciseId) : ''}
	bind:weightUnit
	isLoading={isSubmittingSet}
	onsubmit={handleAddSet}
/>

<EditSetDialog
	bind:this={editSetRef}
	bind:open={editDialogOpen}
	exerciseName={editingSet ? getExerciseName(editingSet.exerciseId) : ''}
	bind:weightUnit
	isLoading={isSavingEdit}
	onsave={handleEditSet}
/>

<NotesDialog
	bind:this={notesRef}
	bind:open={notesDialogOpen}
	isLoading={isSavingNotes}
	onsave={handleSaveNotes}
/>

<DeleteConfirmDialog bind:open={deleteSetDialogOpen} title="Delete this set?" isLoading={isDeletingSet} onconfirm={confirmDeleteSet}>
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
</DeleteConfirmDialog>

<DeleteConfirmDialog bind:open={deleteWorkoutDialogOpen} title="Delete workout?" isLoading={isDeletingWorkout} onconfirm={confirmDeleteWorkout}>
	<p class="text-base-content/60">All sets in this workout will be deleted.</p>
</DeleteConfirmDialog>
