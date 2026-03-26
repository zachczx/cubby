<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { gymMusclesQueryOptions } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import Chart from 'chart.js/auto';
	import Icon from '@iconify/svelte';
	import dayjs from 'dayjs';
	import type { Action } from 'svelte/action';

	const exerciseMap = new Map(exercises.map((e) => [e.id, e]));

	const allMuscles = [
		'chest',
		'shoulders',
		'biceps',
		'triceps',
		'forearms',
		'lats',
		'middle back',
		'lower back',
		'abdominals',
		'quadriceps',
		'hamstrings',
		'glutes',
		'calves',
		'abductors',
		'adductors',
		'traps'
	];

	let selectedWeeks = $state(8);

	const musclesDb = createQuery(() => gymMusclesQueryOptions(selectedWeeks));

	interface MuscleAggregate {
		muscle: string;
		failureCount: number;
		totalSets: number;
		lastFailureDate: string | null;
	}

	function aggregateByMuscle(stats: ExerciseFailureStatsDB[]): MuscleAggregate[] {
		const muscleMap = new Map<string, MuscleAggregate>();

		for (const stat of stats) {
			const primaryMuscles = exerciseMap.get(stat.exerciseId)?.primaryMuscles ?? [];
			for (const muscle of primaryMuscles) {
				const existing = muscleMap.get(muscle);
				if (existing) {
					existing.failureCount += stat.failureCount;
					existing.totalSets += stat.totalSets;
					if (
						stat.lastFailureDate &&
						(!existing.lastFailureDate || stat.lastFailureDate > existing.lastFailureDate)
					) {
						existing.lastFailureDate = stat.lastFailureDate;
					}
				} else {
					muscleMap.set(muscle, {
						muscle,
						failureCount: stat.failureCount,
						totalSets: stat.totalSets,
						lastFailureDate: stat.lastFailureDate
					});
				}
			}
		}

		return allMuscles
			.map((m) => {
				const agg = muscleMap.get(m) || {
					muscle: m,
					failureCount: 0,
					totalSets: 0,
					lastFailureDate: null
				};
				return agg;
			})
			.sort((a, b) => b.failureCount - a.failureCount);
	}

	let muscleAggregates = $derived.by(() => {
		if (musclesDb.isSuccess && musclesDb.data) {
			return aggregateByMuscle(musclesDb.data);
		}

		return [];
	});

	let mostTrained = $derived(muscleAggregates.length > 0 ? muscleAggregates[0] : null);
	let leastTrained = $derived(
		muscleAggregates.length > 0
			? [...muscleAggregates].reverse().find((m: MuscleAggregate) => m.totalSets > 0) || null
			: null
	);

	const renderChart: Action<HTMLCanvasElement, MuscleAggregate[]> = (node, data) => {
		const chart = new Chart(node, {
			type: 'bar',
			data: {
				labels: data.map((m) => m.muscle),
				datasets: [
					{
						label: 'Failure Count',
						data: data.map((m) => m.failureCount),
						backgroundColor: '#3d6b5e',
						borderColor: 'rgba(59, 130, 246, 1)',
						borderWidth: 1
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: { legend: { display: false } },
				scales: {
					y: { beginAtZero: true, ticks: { stepSize: 1 } }
				}
			}
		});

		return {
			update(newData) {
				chart.data.labels = newData.map((m) => m.muscle);
				chart.data.datasets[0].data = newData.map((m) => m.failureCount);
				chart.update();
			},

			destroy() {
				chart.destroy();
			}
		};
	};

	let timeRangeOptions = [
		{ value: 4, label: '4 weeks' },
		{ value: 8, label: '8 weeks' },
		{ value: 12, label: '12 weeks' }
	];
</script>

<PageWrapper title="Muscle Failure Stats">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<div class="flex items-center justify-between">
					<h2 class="text-base-content/70 text-lg font-bold">Time Range</h2>
					<div class="join">
						{#each timeRangeOptions as option}
							<button
								class="join-item btn btn-sm {selectedWeeks === option.value
									? 'btn-active btn-primary'
									: 'btn-ghost'}"
								onclick={() => (selectedWeeks = option.value)}
							>
								{option.label}
							</button>
						{/each}
					</div>
				</div>
			</section>

			{#if musclesDb.isSuccess && musclesDb.data}
				{#if musclesDb.data.length === 0}
					<div class="grid justify-items-center gap-2 py-12">
						<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-16" />
						<p class="text-base-content/50 text-center">
							No failure sets in the last {selectedWeeks} weeks.<br />Push yourself to failure to
							see stats!
						</p>
					</div>
				{:else}
					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Summary</h2>
						<div class="grid grid-cols-2 gap-4">
							{#if mostTrained && mostTrained.failureCount > 0}
								<div
									class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
								>
									<div class="text-base-content/50 flex items-center gap-2 text-sm">
										<Icon icon="material-symbols:trending-up" class="size-4" />
										Most to Failure
									</div>
									<div class="text-lg font-bold capitalize">{mostTrained.muscle}</div>
									<div class="text-base-content/50 text-xs">{mostTrained.failureCount}x</div>
								</div>
							{/if}
							{#if leastTrained}
								<div
									class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
								>
									<div class="text-base-content/50 flex items-center gap-2 text-sm">
										<Icon icon="material-symbols:trending-down" class="size-4" />
										Least to Failure
									</div>
									<div class="text-lg font-bold capitalize">{leastTrained.muscle}</div>
									<div class="text-base-content/50 text-xs">{leastTrained.failureCount}x</div>
								</div>
							{/if}
						</div>
					</section>

					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Failure Count by Muscle</h2>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border p-4">
							<div class="h-64">
								<canvas use:renderChart={muscleAggregates}></canvas>
							</div>
						</div>
					</section>

					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Muscle Details</h2>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
							<table class="table-sm table">
								<thead>
									<tr>
										<th>Muscle</th>
										<th class="text-right">Failures</th>
										<th class="text-right">Total Sets</th>
										<th class="text-right">Rate</th>
										<th>Last</th>
									</tr>
								</thead>
								<tbody>
									{#each muscleAggregates.filter((m) => m.totalSets > 0) as muscle, i (muscle.muscle)}
										<tr
											class={i < muscleAggregates.filter((m) => m.totalSets > 0).length - 1
												? 'border-b-base-300/50 border-b'
												: ''}
										>
											<td class="capitalize">{muscle.muscle}</td>
											<td class="text-error/70 text-right font-semibold">{muscle.failureCount}x</td>
											<td class="text-right">{muscle.totalSets}</td>
											<td class="text-right"
												>{((muscle.failureCount / muscle.totalSets) * 100).toFixed(0)}%</td
											>
											<td class="text-base-content/50 text-xs">
												{muscle.lastFailureDate
													? dayjs(muscle.lastFailureDate).format('D MMM')
													: '-'}
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</section>
				{/if}
			{:else if musclesDb.isError}
				<div class="grid justify-items-center gap-2 py-12">
					<p class="text-base-content/50">Failed to load stats</p>
				</div>
			{:else}
				<section class="grid gap-4 py-2">
					<div class="skeleton h-8 w-32"></div>
					<div class="skeleton h-64 rounded-2xl"></div>
				</section>
				<section class="grid gap-4 py-2">
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
