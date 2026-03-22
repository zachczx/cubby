<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { allWorkoutsQueryOptions, getAllWorkoutsQueryKey } from '$lib/queries';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { exercises } from '$lib/exercises';
	import Icon from '@iconify/svelte';
	import { Combobox, useListCollection } from '@ark-ui/svelte/combobox';
	import CorgiGym from '$lib/assets/corgi_gym.webp?w=240&enhanced';

	dayjs.extend(relativeTime);

	const queryClient = useQueryClient();
	const workoutsDb = createQuery(allWorkoutsQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

	const exerciseItems = exercises.map((e) => ({ label: e.name, value: e.id }));
	const exerciseCollection = useListCollection({
		initialItems: exerciseItems,
		filter: (itemString, filterText) => {
			return itemString.toLowerCase().includes(filterText.toLowerCase());
		}
	});

	let addingSetToWorkout = $state<string | null>(null);
	let selectedExerciseId = $state('');
	let setWeight = $state<number | null>(null);
	let setReps = $state<number | null>(null);
	let setType = $state('working');
	let weightUnit = $state<'kg' | 'lb'>('kg');

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

	async function startWorkout() {
		const response = await api.post('gym/workouts');
		if (response.status === 201) {
			addToast('success', 'Workout started!');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		} else {
			addToast('error', 'Failed to start workout');
		}
	}

	async function deleteWorkout(workoutId: string) {
		const response = await api.delete(`gym/workouts/${workoutId}`);
		if (response.status === 204) {
			addToast('success', 'Workout deleted');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		} else {
			addToast('error', 'Failed to delete workout');
		}
	}

	async function addSet(workoutId: string) {
		if (!selectedExerciseId) return;

		const weightKg = setWeight != null && weightUnit === 'lb' ? lbToKg(setWeight) : setWeight;

		const response = await api.post(`gym/workouts/${workoutId}/sets`, {
			body: JSON.stringify({
				exerciseId: selectedExerciseId,
				weightKg,
				reps: setReps,
				setType
			})
		});

		if (response.status === 201) {
			addToast('success', 'Set added!');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
			selectedExerciseId = '';
			setWeight = null;
			setReps = null;
			setType = 'working';
			addingSetToWorkout = null;
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
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		} else {
			addToast('error', 'Failed to add set');
		}
	}

	async function deleteSet(setId: string) {
		const response = await api.delete(`gym/sets/${setId}`);
		if (response.status === 204) {
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
		} else {
			addToast('error', 'Failed to delete set');
		}
	}
</script>

<PageWrapper title="Gym">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<button class="btn btn-primary btn-lg w-full rounded-full" onclick={startWorkout}>
					Start Workout
				</button>
			</section>

			<section class="grid gap-4 py-2">
				<div class="flex items-center justify-between">
					<h2 class="text-base-content/70 text-lg font-bold">Workouts</h2>
					<div
						class={[
							'border-base-300 flex overflow-hidden rounded-lg border text-sm',
							!workoutsDb.data || workoutsDb.data.length === 0 ? 'hidden' : undefined
						]}
					>
						<button
							class={[
								'px-3 py-1 font-medium transition-colors',
								weightUnit === 'kg'
									? 'bg-primary text-primary-content'
									: 'text-base-content/50 hover:bg-base-200'
							]}
							onclick={() => (weightUnit = 'kg')}
						>
							kg
						</button>
						<button
							class={[
								'px-3 py-1 font-medium transition-colors',
								weightUnit === 'lb'
									? 'bg-primary text-primary-content'
									: 'text-base-content/50 hover:bg-base-200'
							]}
							onclick={() => (weightUnit = 'lb')}
						>
							lb
						</button>
					</div>
				</div>

				{#if workoutsDb.isSuccess}
					{#if !workoutsDb.data || workoutsDb.data.length === 0}
						<div class="grid justify-items-center gap-2 py-8">
							<enhanced:img src={CorgiGym} alt="No workout" />
							<button class="text-base-content/70 font-medium" onclick={startWorkout}
								>Add your first workout!</button
							>
						</div>
					{:else}
						{#each workoutsDb.data as workout (workout.id)}
							{@const exerciseGroups = groupSetsByExercise(workout.sets)}
							<div class="border-base-300/50 rounded-2xl border bg-base-50">
								<div class="flex items-center justify-between px-4 pt-3 pb-2">
									<a href="/app/gym/{workout.id}" class="grid grow gap-0">
										<div
											class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
										>
											{dayjs(workout.startTime).format('ddd, D MMM')}
											{#if !workout.endTime}
												<div
													class="bg-primary/10 text-primary rounded-full px-2 py-0.5 text-xs font-bold"
												>
													Active
												</div>
											{/if}
										</div>
										<span class="text-base-content/50">
											{dayjs(workout.startTime).format('h:mm A')}
											{#if workout.endTime}
												— {dayjs(workout.endTime).format('h:mm A')}
											{/if}
										</span>
									</a>
									<button
										class="btn btn-ghost btn-sm btn-square"
										onclick={() => deleteWorkout(workout.id)}
									>
										<Icon icon="material-symbols:delete-outline" class="text-error size-5" />
									</button>
								</div>

								{#if workout.notes}
									<p class="text-base-content/70 px-4 pb-2">{workout.notes}</p>
								{/if}

								{#if exerciseGroups.length > 0}
									<div class="grid gap-0 px-4 pb-3">
										{#each exerciseGroups as group, gi (group.exerciseId)}
											<div
												class={[
													'py-3',
													gi < exerciseGroups.length - 1 && 'border-b-base-300/50 border-b'
												]}
											>
												<p class="text-base-content/80 mb-2 font-semibold tracking-wide">
													{group.exerciseName}
												</p>
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
																	<button
																		class="btn btn-ghost btn-xs btn-square"
																		onclick={() => duplicateSet(workout.id, set)}
																	>
																		<Icon
																			icon="material-symbols:replay"
																			class="text-base-content/50 size-3.5"
																		/>
																	</button>
																	<button
																		class="btn btn-ghost btn-xs btn-square"
																		onclick={() => deleteSet(set.id)}
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
										{/each}
									</div>
								{/if}

								{#if addingSetToWorkout === workout.id}
									<div class="border-t-base-300/50 grid gap-3 border-t px-4 py-3">
										<Combobox.Root
											collection={exerciseCollection.collection()}
											inputBehavior="autohighlight"
											openOnClick
											onValueChange={(e) => {
												selectedExerciseId = e.value[0] ?? '';
											}}
											onInputValueChange={(e) => {
												exerciseCollection.filter(e.inputValue);
											}}
										>
											<Combobox.Control class="flex gap-1">
												<Combobox.Input
													class="input input-bordered w-full rounded-lg"
													placeholder="Search exercises..."
												/>
												<Combobox.Trigger class="btn btn-ghost btn-square">
													<Icon icon="material-symbols:expand-more" class="size-5" />
												</Combobox.Trigger>
											</Combobox.Control>
											<Combobox.Positioner>
												<Combobox.Content
													class="bg-base-100 border-base-300 z-50 max-h-60 overflow-y-auto rounded-xl border shadow-lg"
												>
													{#each exerciseCollection.collection().items as item (item.value)}
														<Combobox.Item
															{item}
															class="hover:bg-base-200 data-highlighted:bg-base-200 cursor-pointer px-3 py-2"
														>
															{item.label}
														</Combobox.Item>
													{/each}
												</Combobox.Content>
											</Combobox.Positioner>
										</Combobox.Root>

										<div class="flex gap-2">
											<div class="flex w-full">
												<input
													type="number"
													class="input input-bordered w-full rounded-l-lg rounded-r-none"
													placeholder={weightUnit}
													bind:value={setWeight}
												/>
												<button
													class="btn btn-ghost border-base-300 rounded-l-none rounded-r-lg border"
													onclick={() => (weightUnit = weightUnit === 'kg' ? 'lb' : 'kg')}
												>
													{weightUnit}
												</button>
											</div>
											<input
												type="number"
												class="input input-bordered w-full rounded-lg"
												placeholder="Reps"
												bind:value={setReps}
											/>
											<select class="select select-bordered rounded-lg" bind:value={setType}>
												<option value="working">Working</option>
												<option value="warmup">Warmup</option>
												<option value="dropset">Drop</option>
												<option value="failure">Failure</option>
											</select>
										</div>

										<div class="flex gap-2">
											<button
												class="btn btn-primary grow rounded-lg"
												onclick={() => addSet(workout.id)}
											>
												Add Set
											</button>
											<button
												class="btn btn-ghost rounded-lg"
												onclick={() => {
													addingSetToWorkout = null;
													selectedExerciseId = '';
												}}
											>
												Cancel
											</button>
										</div>
									</div>
								{:else}
									<div class="border-t-base-300/50 border-t px-4 py-3">
										<button
											class="btn btn-ghost btn-sm text-primary w-full"
											onclick={() => (addingSetToWorkout = workout.id)}
										>
											<Icon icon="material-symbols:add" class="size-4" />
											Add Set
										</button>
									</div>
								{/if}
							</div>
						{/each}
					{/if}
				{:else}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>
