<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { gymExerciseQueryOptions } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import Chart from 'chart.js/auto';
	import Icon from '@iconify/svelte';
	import dayjs from 'dayjs';
	import type { Action } from 'svelte/action';
	import type { PageData } from './$types';
	import { detectAllPrs } from '$lib/pr';

	let { data }: { data: PageData } = $props();
	let exerciseId = $derived(data.exerciseId);

	let exerciseName = $derived.by(() => {
		const lower = exerciseId.toLowerCase();
		const exercise = exercises.find((e) => e.id.toLowerCase() === lower);
		return exercise?.name ?? exerciseId;
	});

	const exerciseDb = createQuery(() => gymExerciseQueryOptions(exerciseId));

	type ChartDataPoint = { date: string; weight: number; reps: number; setType: string };
	type VolumeDataPoint = { date: string; volume: number };

	const renderLineChart: Action<HTMLCanvasElement, ChartDataPoint[]> = (node, dataPoints) => {
		const chart = new Chart(node, {
			type: 'line',
			data: {
				labels: dataPoints.map((d) => dayjs(d.date).format('D MMM')),
				datasets: [
					{
						label: 'Weight (kg)',
						data: dataPoints.map((d) => d.weight),
						borderColor: '#3d6b5e',
						backgroundColor: 'rgba(61, 107, 94, 0.1)',
						borderWidth: 2,
						tension: 0.3,
						fill: true,
						pointRadius: 4,
						pointHoverRadius: 6
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: { display: false }
				},
				scales: {
					y: {
						beginAtZero: false,
						ticks: { stepSize: 2.5 }
					}
				}
			}
		});

		return {
			update(newData) {
				chart.data.labels = newData.map((d) => dayjs(d.date).format('D MMM'));
				chart.data.datasets[0].data = newData.map((d) => d.weight);
				chart.update();
			},
			destroy() {
				chart.destroy();
			}
		};
	};

	const renderVolumeChart: Action<HTMLCanvasElement, VolumeDataPoint[]> = (node, dataPoints) => {
		const chart = new Chart(node, {
			type: 'line',
			data: {
				labels: dataPoints.map((d) => dayjs(d.date).format('D MMM')),
				datasets: [
					{
						label: 'Volume (kg)',
						data: dataPoints.map((d) => d.volume),
						borderColor: 'rgba(34, 197, 94, 1)',
						backgroundColor: 'rgba(34, 197, 94, 0.1)',
						borderWidth: 2,
						tension: 0.3,
						fill: true,
						pointRadius: 4,
						pointHoverRadius: 6
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: { display: false }
				},
				scales: {
					y: {
						beginAtZero: true,
						ticks: { stepSize: 100 }
					}
				}
			}
		});

		return {
			update(newData) {
				chart.data.labels = newData.map((d) => dayjs(d.date).format('D MMM'));
				chart.data.datasets[0].data = newData.map((d) => d.volume);
				chart.update();
			},
			destroy() {
				chart.destroy();
			}
		};
	};

	type ViewMode = 'weight' | 'volume';
	let viewMode = $state<ViewMode>('weight');

	let chartData = $derived.by(() => {
		if (!exerciseDb.isSuccess || !exerciseDb.data) return [];
		return exerciseDb.data
			.filter((s) => s.weightKg != null)
			.map((s) => ({
				date: s.date,
				weight: s.weightKg!,
				reps: s.reps ?? 0,
				setType: s.setType
			}));
	});

	let volumeData = $derived.by(() => {
		if (!exerciseDb.isSuccess || !exerciseDb.data) return [];
		const grouped = new Map<string, number>();
		for (const set of exerciseDb.data) {
			const vol = (set.weightKg ?? 0) * (set.reps ?? 0);
			grouped.set(set.date, (grouped.get(set.date) ?? 0) + vol);
		}
		return Array.from(grouped.entries())
			.map(([date, volume]) => ({ date, volume }))
			.sort((a, b) => a.date.localeCompare(b.date));
	});

	let maxWeight = $derived(chartData.length > 0 ? Math.max(...chartData.map((d) => d.weight)) : 0);
	let totalVolume = $derived(
		volumeData.length > 0 ? volumeData.reduce((sum, d) => sum + d.volume, 0) : 0
	);
	let totalSets = $derived(exerciseDb.data?.length ?? 0);

	let prSet = $derived.by(() => {
		if (!exerciseDb.isSuccess || !exerciseDb.data) return new Map();
		return detectAllPrs(exerciseDb.data);
	});

</script>

<PageWrapper title={exerciseName}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			{#if exerciseDb.isSuccess && exerciseDb.data}
				{#if exerciseDb.data.length === 0}
					<div class="grid justify-items-center gap-2 py-12">
						<Icon icon="material-symbols:fitness-center" class="text-base-content/20 size-16" />
						<p class="text-base-content/50 text-center">
							No history for this exercise yet.<br />Start a workout to log sets!
						</p>
					</div>
				{:else}
					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Summary</h2>
						<div class="grid grid-cols-3 gap-4">
							<div
								class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
							>
								<div class="text-base-content/50 flex items-center gap-2 text-sm">
									<Icon icon="material-symbols:weight" class="size-4" />
									Max
								</div>
								<div class="text-xl font-bold">{maxWeight.toFixed(1)}kg</div>
							</div>
							<div
								class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
							>
								<div class="text-base-content/50 flex items-center gap-2 text-sm">
									<Icon icon="material-symbols:format-list-numbered" class="size-4" />
									Sets
								</div>
								<div class="text-xl font-bold">{totalSets}</div>
							</div>
							<div
								class="border-base-300/50 bg-base-50 grid content-center justify-items-center rounded-2xl border p-4"
							>
								<div class="text-base-content/50 flex items-center gap-2 text-sm">
									<Icon icon="material-symbols:calculate" class="size-4" />
									Volume
								</div>
								<div class="text-xl font-bold">{Math.round(totalVolume)}kg</div>
							</div>
						</div>
					</section>

					<section class="grid gap-4 py-2">
						<div class="flex items-center justify-between">
							<h2 class="text-base-content/70 text-lg font-bold">Progress</h2>
							<div class="join">
								<button
									class="join-item btn btn-sm {viewMode === 'weight'
										? 'btn-active btn-primary'
										: 'btn-ghost'}"
									onclick={() => (viewMode = 'weight')}
								>
									Weight
								</button>
								<button
									class="join-item btn btn-sm {viewMode === 'volume'
										? 'btn-active btn-primary'
										: 'btn-ghost'}"
									onclick={() => (viewMode = 'volume')}
								>
									Volume
								</button>
							</div>
						</div>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border p-4">
							<div class="h-64">
								{#if viewMode === 'weight'}
									<canvas use:renderLineChart={chartData}></canvas>
								{:else}
									<canvas use:renderVolumeChart={volumeData}></canvas>
								{/if}
							</div>
						</div>
					</section>

					<section class="grid gap-4 py-2">
						<h2 class="text-base-content/70 text-lg font-bold">Set History</h2>
						<div class="border-base-300/50 bg-base-50 overflow-hidden rounded-2xl border">
							<table class="table-sm table">
								<thead>
									<tr>
										<th>Date</th>
										<th class="text-right">Weight</th>
										<th class="text-right">Reps</th>
										<th class="text-right">Type</th>
									</tr>
								</thead>
								<tbody>
									{#each exerciseDb.data as set, i (i)}
										<tr
											class={i < exerciseDb.data.length - 1 ? 'border-b-base-300/50 border-b' : ''}
										>
											<td class="text-base-content/60 text-sm">
												{dayjs(set.date).format('D MMM YYYY')}
											</td>
											<td class="text-right font-semibold">
												{set.weightKg != null ? `${set.weightKg}kg` : '–'}
												{#if prSet.get(i)}
													<span
														class="bg-warning/15 text-warning ml-1 rounded-full px-1.5 py-0.5 text-xs font-bold"
													>
														{prSet.get(i)?.label}
													</span>
												{/if}
											</td>
											<td class="text-right">{set.reps ?? '–'}</td>
											<td class="text-right">
												{#if set.setType === 'failure'}
													<span class="text-error/70">failure</span>
												{:else if set.setType === 'dropset'}
													<span class="text-warning/70">dropset</span>
												{:else}
													<span class="text-base-content/50">working</span>
												{/if}
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</section>
				{/if}
			{:else if exerciseDb.isError}
				<div class="grid justify-items-center gap-2 py-12">
					<p class="text-base-content/50">Failed to load stats</p>
				</div>
			{:else}
				<section class="grid gap-4 py-2">
					<div class="skeleton h-24 rounded-2xl"></div>
				</section>
				<section class="grid gap-4 py-2">
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
						href="/app/gym/stats/exercises"
						class="border-b-base-300/50 flex items-center gap-3 border-b px-4 py-3"
					>
						<Icon icon="material-symbols:chevron-left" class="text-primary size-5" />
						<span class="grow font-medium">Back to Exercises</span>
					</a>
					<a href="/app/gym/stats/muscles" class="flex items-center gap-3 px-4 py-3">
						<Icon icon="material-symbols:chevron-left" class="text-primary size-5" />
						<span class="grow font-medium">Muscle Stats</span>
					</a>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>
