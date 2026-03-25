<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import updateLocale from 'dayjs/plugin/updateLocale';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { allWorkoutsQueryOptions, getAllWorkoutsQueryKey } from '$lib/queries';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { exercises } from '$lib/exercises';
	import CorgiGym from '$lib/assets/corgi_gym.webp?w=240&enhanced';
	import Icon from '@iconify/svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(updateLocale);

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

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

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

	const maxNumberOfExercisesToShow = 3;
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
					<a href="/app/gym/stats" class="btn btn-ghost btn-sm gap-1">
						<Icon icon="material-symbols:bar-chart" class="size-4" />
						Stats
					</a>
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

							<div class="border-base-300/50 bg-base-50 rounded-2xl border">
								<div class="flex items-center justify-between px-4 pt-3 pb-2">
									<a href="/app/gym/{workout.id}" class="grid grow gap-0">
										<div
											class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
										>
											{dayjs(workout.startTime).fromNow()}
										</div>
										<div class="text-base-content/65 flex items-center gap-2 text-sm font-light">
											<span>{dayjs(workout.startTime).format('ddd, D MMM')}</span>
											<span class="text-base-content/40 text-xs font-extralight">•</span>
											<span>{dayjs(workout.startTime).format('h:mma')}</span>
										</div>
									</a>
									<div class="px-1">
										<Icon
											icon="material-symbols:local-fire-department-rounded"
											class="text-accent size-5"
										/>
									</div>
								</div>

								{#if workout.notes}
									<p class="text-base-content/70 px-4 pb-2">{workout.notes}</p>
								{/if}

								{#if exerciseGroups.length > 0}
									<div class="grid gap-0 px-4 text-sm">
										{#each exerciseGroups as group, gi (group.exerciseId)}
											{#if gi < 3}
												<div class="flex items-center">
													<div
														class={[
															'grow py-3',
															gi < exerciseGroups.length - 1 && 'border-b-base-300/50 border-b'
														]}
													>
														<p class="text-base-content/80 font-medium">
															{group.exerciseName}
														</p>
													</div>
													<div class="text-base-content/65 px-2 text-sm font-medium">
														x{group.sets.length}
													</div>
												</div>
											{/if}
										{/each}
										{#if exerciseGroups.length > maxNumberOfExercisesToShow}
											<div class="text-base-content/60 py-3 text-sm italic">
												+{exerciseGroups.length - 3} more exercises
											</div>
										{/if}
									</div>
								{/if}

								<div class="border-t-base-300/50 border-t px-4 py-3">
									<a href="/app/gym/{workout.id}" class="btn btn-soft btn-primary w-full">
										View Workout
									</a>
								</div>
							</div>
						{/each}
					{/if}
				{:else}
					<SkeletonActionCard />
					<SkeletonActionCard />
					<SkeletonActionCard />
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>
