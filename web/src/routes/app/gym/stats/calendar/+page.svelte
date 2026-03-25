<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { gymCalendarQueryOptions } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import Icon from '@iconify/svelte';
	import { tick } from 'svelte';

	dayjs.extend(utc);

	const calendarDb = createQuery(gymCalendarQueryOptions);

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	function getExerciseName(exerciseId: string): string {
		return exerciseMap.get(exerciseId)?.name ?? exerciseId;
	}

	function getMuscleGroups(exerciseIds: string[]): string[] {
		const muscles = new Set<string>();
		for (const id of exerciseIds) {
			for (const m of exerciseMap.get(id)?.primaryMuscles ?? []) {
				muscles.add(m);
			}
		}
		return [...muscles];
	}

	let modal = $state<HTMLDialogElement>();
	let selectedEntry = $state<WorkoutCalendarEntryDB | null>(null);

	let calendarOptions: Calendar.Options = $derived.by(() => {
		if (!calendarDb.isSuccess || !calendarDb.data) {
			return { view: 'dayGridMonth', events: [], firstDay: 1 };
		}

		const events = calendarDb.data.map((entry) => {
			const start = dayjs.utc(entry.startTime);
			return {
				start: start.toDate(),
				end: start.add(1, 'second').toDate(),
				title: `${entry.exerciseCount} exercises`,
				backgroundColor: 'var(--color-primary)',
				extendedProps: { workoutId: entry.workoutId }
			};
		});

		return {
			view: 'dayGridMonth',
			events,
			firstDay: 1,
			titleFormat: (date) => {
				const month = dayjs(date).get('month');
				if (month === 8) {
					return dayjs(date).format('MMMM YY');
				}
				return dayjs(date).format('MMMM YYYY');
			},
			dateClick: async (info) => {
				const dayEntries = calendarDb.data?.filter((e) =>
					dayjs.utc(e.startTime).isSame(info.date, 'day')
				);
				if (dayEntries && dayEntries.length > 0) {
					selectedEntry = dayEntries[0];
					await tick();
					modal?.showModal();
				}
			},
			eventClick: async (info) => {
				const entry = calendarDb.data?.find(
					(e) => e.workoutId === info.event.extendedProps.workoutId
				);
				if (entry) {
					selectedEntry = entry;
					await tick();
					modal?.showModal();
				}
			}
		};
	});
</script>

<PageWrapper title="Workout Calendar">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			{#if calendarDb.isSuccess && calendarDb.data}
				{#if calendarDb.data.length === 0}
					<div class="grid justify-items-center gap-2 py-12">
						<Icon
							icon="material-symbols:calendar-month-outline"
							class="text-base-content/20 size-16"
						/>
						<p class="text-base-content/50 text-center">
							No workouts yet.<br />Start one to see your calendar!
						</p>
						<a href="/app/gym" class="btn btn-primary btn-soft rounded-full">Go to Gym</a>
					</div>
				{:else}
					<section class="py-2">
						<div class="grid w-full grid-cols-[minmax(0,1fr)] overflow-hidden">
							{#key calendarDb.data}
								<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
							{/key}
						</div>
					</section>

					<section class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
						{#each calendarDb.data.slice(0, 10) as entry, i (entry.workoutId)}
							<a
								href="/app/gym/{entry.workoutId}"
								class={[
									'flex items-center px-4 py-3',
									i < Math.min(calendarDb.data.length, 10) - 1 &&
										'border-b-base-300/50 border-b'
								]}
							>
								<div class="grid grow gap-0.5">
									<p class="text-base-content/80 font-medium">
										{dayjs.utc(entry.startTime).format('ddd, D MMM')}
									</p>
									<p class="text-base-content/50 text-xs">
										{getMuscleGroups(entry.exerciseIds).join(', ')}
									</p>
								</div>
								<span class="text-base-content/60 text-sm font-medium">
									{entry.exerciseCount} ex · {entry.setCount} sets
								</span>
							</a>
						{/each}
					</section>
				{/if}
			{:else if calendarDb.isError}
				<div class="grid justify-items-center gap-2 py-12">
					<p class="text-base-content/50">Failed to load calendar</p>
				</div>
			{:else}
				<section class="py-2">
					<div class="skeleton h-80 rounded-2xl"></div>
				</section>
				<section class="grid gap-4 py-2">
					<div class="skeleton h-48 rounded-2xl"></div>
				</section>
			{/if}
		</div>
	</main>
</PageWrapper>

<dialog bind:this={modal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		{#if selectedEntry}
			<div class="flex items-center">
				<h3 class="grow text-lg font-bold">
					{dayjs.utc(selectedEntry.startTime).format('ddd, D MMM YYYY')}
				</h3>
				<form method="dialog">
					<button class="btn btn-ghost -me-2">
						<Icon icon="material-symbols:close" class="size-5" />
					</button>
				</form>
			</div>
			<div class="text-base-content/60 mb-4 text-sm">
				{selectedEntry.exerciseCount} exercises · {selectedEntry.setCount} sets
			</div>
			<div class="grid gap-2">
				{#each selectedEntry.exerciseIds as exerciseId (exerciseId)}
					<div class="border-base-300/50 bg-base-200/30 rounded-xl border px-3 py-2">
						<p class="font-medium">{getExerciseName(exerciseId)}</p>
						<p class="text-base-content/50 text-xs">
							{exerciseMap.get(exerciseId)?.primaryMuscles.join(', ')}
						</p>
					</div>
				{/each}
			</div>
			<div class="mt-4">
				<a href="/app/gym/{selectedEntry.workoutId}" class="btn btn-primary w-full rounded-full">
					View Workout
				</a>
			</div>
		{/if}
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
