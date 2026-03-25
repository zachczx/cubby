<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { gymSummaryQueryOptions } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import Icon from '@iconify/svelte';

	const summaryDb = createQuery(gymSummaryQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

	function formatVolume(kg: number): string {
		if (kg >= 1000) return `${(kg / 1000).toFixed(1)}t`;
		return `${Math.round(kg)}kg`;
	}

	function getMuscleGroups(exerciseId: string): string[] {
		return exerciseMap.get(exerciseId)?.primaryMuscles ?? [];
	}
</script>

<PageWrapper title="Gym Stats">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			{#if summaryDb.isSuccess && summaryDb.data}
				{@const s = summaryDb.data}

				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">This Month</h2>
					<div class="grid grid-cols-2 gap-4">
						<div
							class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
						>
							<div class="text-base-content/50 flex items-center gap-2 text-sm">
								<Icon icon="material-symbols:exercise" class="size-4" />
								Workouts
							</div>
							<div class="text-2xl font-bold">{s.totalWorkoutsThisMonth}</div>
						</div>
						<div
							class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
						>
							<div class="text-base-content/50 flex items-center gap-2 text-sm">
								<Icon icon="material-symbols:weight" class="size-4" />
								Volume
							</div>
							<div class="text-2xl font-bold">{formatVolume(s.totalVolumeThisMonth)}</div>
						</div>
					</div>
				</section>

				{#if s.topExercises.length > 0}
					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Top Exercises</h2>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
							{#each s.topExercises as ex, i (ex.exerciseId)}
								<div
									class={[
										'flex items-center px-4 py-3',
										i < s.topExercises.length - 1 && 'border-b-base-300/50 border-b'
									]}
								>
									<div class="grid grow gap-0.5">
										<p class="text-base-content/80 font-medium">{getExerciseName(ex.exerciseId)}</p>
										<p class="text-base-content/50 text-xs">
											{getMuscleGroups(ex.exerciseId).join(', ')}
										</p>
									</div>
									<span class="text-base-content/60 text-sm font-semibold">{ex.count} sets</span>
								</div>
							{/each}
						</div>
					</section>
				{/if}

				{#if s.failureExercises.length > 0}
					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Failure Sets</h2>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
							{#each s.failureExercises as ex, i (ex.exerciseId)}
								<div
									class={[
										'flex items-center px-4 py-3',
										i < s.failureExercises.length - 1 && 'border-b-base-300/50 border-b'
									]}
								>
									<div class="grid grow gap-0.5">
										<p class="text-base-content/80 font-medium">{getExerciseName(ex.exerciseId)}</p>
										<p class="text-base-content/50 text-xs">
											{getMuscleGroups(ex.exerciseId).join(', ')}
										</p>
									</div>
									<span class="text-error/70 text-sm font-semibold">{ex.count}x failure</span>
								</div>
							{/each}
						</div>
					</section>
				{/if}

				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Explore</h2>
					<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
						<a
							href="/app/gym/stats/calendar"
							class="border-b-base-300/50 flex items-center gap-3 border-b px-4 py-3"
						>
							<Icon icon="material-symbols:calendar-month" class="text-primary size-5" />
							<span class="grow font-medium">Workout Calendar</span>
							<Icon icon="material-symbols:chevron-right" class="text-base-content/30 size-5" />
						</a>
					</div>
				</section>

				{#if s.totalWorkoutsThisMonth === 0}
					<div class="grid justify-items-center gap-2 py-12">
						<Icon icon="material-symbols:bar-chart-outline" class="text-base-content/20 size-16" />
						<p class="text-base-content/50 text-center">
							No workouts this month yet.<br />Start one to see your stats!
						</p>
						<a href="/app/gym" class="btn btn-primary btn-soft rounded-full">Go to Gym</a>
					</div>
				{/if}
			{:else if summaryDb.isError}
				<div class="grid justify-items-center gap-2 py-12">
					<p class="text-base-content/50">Failed to load stats</p>
				</div>
			{:else}
				<section class="grid gap-4 py-2">
					<div class="skeleton h-8 w-32"></div>
					<div class="grid grid-cols-2 gap-4">
						<div class="skeleton h-24 rounded-xl"></div>
						<div class="skeleton h-24 rounded-xl"></div>
					</div>
				</section>
				<section class="grid gap-4 py-2">
					<div class="skeleton h-8 w-32"></div>
					<div class="skeleton h-48 rounded-2xl"></div>
				</section>
			{/if}
		</div>
	</main>
</PageWrapper>
