<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { gymUserExercisesQueryOptions } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import Icon from '@iconify/svelte';

	const userExercisesDb = createQuery(gymUserExercisesQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

	function getMuscleGroups(exerciseId: string): string[] {
		return exerciseMap.get(exerciseId)?.primaryMuscles ?? [];
	}

	let exerciseSearch = $state('');

	let userExercises = $derived(
		userExercisesDb.isSuccess && userExercisesDb.data
			? userExercisesDb.data
			: []
	);

	let filteredExercises = $derived(
		exerciseSearch
			? userExercises.filter((ue) => {
					const name = getExerciseName(ue.exerciseId).toLowerCase();
					return name.includes(exerciseSearch.toLowerCase());
				})
			: userExercises
	);
</script>

<PageWrapper title="Exercise Progress">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<input
					type="text"
					placeholder="Search exercises..."
					class="input input-bordered w-full"
					bind:value={exerciseSearch}
				/>
			</section>

			{#if userExercisesDb.isSuccess && userExercisesDb.data}
				{#if userExercisesDb.data.length === 0}
					<div class="grid justify-items-center gap-2 py-12">
						<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-16" />
						<p class="text-base-content/50 text-center">
							No exercises logged yet.<br />Start a workout to track your progress!
						</p>
					</div>
				{:else}
					<section class="grid gap-4 py-2">
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
							{#each filteredExercises as ue, i (ue.exerciseId)}
								<a
									href="/app/gym/stats/exercises/{encodeURIComponent(ue.exerciseId)}"
									class="flex items-center gap-3 px-4 py-3 {i < filteredExercises.length - 1 ? 'border-b-base-300/50 border-b' : ''}"
								>
									<div class="grow">
										<p class="text-base-content/80 font-medium">
											{getExerciseName(ue.exerciseId)}
										</p>
										<p class="text-base-content/50 text-xs">
											{getMuscleGroups(ue.exerciseId).join(', ')}
										</p>
									</div>
									<span class="text-base-content/50 text-sm">{ue.setCount} sets</span>
									<Icon icon="material-symbols:chevron-right" class="text-base-content/30 size-5" />
								</a>
							{/each}
						</div>
					</section>
				{/if}
			{:else if userExercisesDb.isError}
				<div class="grid justify-items-center gap-2 py-12">
					<p class="text-base-content/50">Failed to load exercises</p>
				</div>
			{:else}
				<section class="grid gap-4 py-2">
					<div class="skeleton h-12 rounded-xl"></div>
					<div class="skeleton h-48 rounded-2xl"></div>
				</section>
			{/if}

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Explore</h2>
				<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
					<a
						href="/app/gym/stats"
						class="border-b-base-300/50 flex items-center gap-3 border-b px-4 py-3"
					>
						<Icon icon="material-symbols:chevron-left" class="text-primary size-5" />
						<span class="grow font-medium">Back to Stats</span>
					</a>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>