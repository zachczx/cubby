<script lang="ts">
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import Chart from 'chart.js/auto';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import { calculateVacationOverlap } from '$lib/overlap';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createTowelQueryOptions,
		createTowelRefetchOptions,
		createVacationQueryOptions
	} from '$lib/queries';
	import { getCalendarEntries } from '$lib/calendar';
	import { getStatusColorFromValue } from '$lib/towels';
	import CustomDateModal from '$lib/ui/CustomDateModal.svelte';
	import TwoColumnCard from '$lib/ui/TwoColumnCard.svelte';
	import StatusDescriptions from '$lib/ui/StatusDescriptions.svelte';
	import StatusHeroImage from '$lib/ui/StatusHeroImage.svelte';
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import SingleDayModal from '$lib/ui/SingleDayModal.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const vacations = createQuery(createVacationQueryOptions);
	const towels = createQuery(createTowelQueryOptions);
	const tanstackClient = useQueryClient();

	/**
	 * Using $state + $effect instead of $derived due to TanStack Query store
	 * not properly triggering Svelte 5's fine-grained reactivity on async data changes
	 */

	let towelRecords: TowelRecord[] = $state([]);
	$effect(() => {
		if (towels.isSuccess) {
			towelRecords = towels.data?.map((record, i, allRecords) => {
				const nextRecord = allRecords[i + 1];
				if (!nextRecord) {
					return { ...record, gap: 0 };
				}

				// 1. Calculate the raw gap
				const rawGap = dayjs(record.time).diff(nextRecord.time, 'day', true);

				// 2. Calculate vacation days within that gap
				const vacationDays = calculateVacationOverlap(
					record.time,
					nextRecord.time,
					vacations.data ?? [],
					'day'
				);

				// 3. Subtract vacation time from the raw gap
				const finalGap = rawGap - vacationDays;

				// Ensure gap is never negative
				return { ...record, gap: Math.max(0, finalGap) };
			});
		}
	});

	let longestGap: TowelRecord | undefined = $derived.by(() => {
		if (towelRecords.length <= 1) return;
		const avoidMutatingOriginalArray = [...towelRecords];
		const sorted = avoidMutatingOriginalArray.sort((a, b) => b.gap - a.gap);
		return sorted[0];
	});
	let averageGap: number = $derived.by(() => {
		const total = towelRecords.reduce((accumulator, entry) => {
			return (accumulator += entry.gap);
		}, 0);
		const average = total / towelRecords.length;
		return average;
	});

	let gaps = $derived.by(() => {
		const gaps: number[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (towelRecords[i]) {
				gaps.unshift(towelRecords[i].gap);
			}
		}
		return gaps;
	});

	let gapsDates = $derived.by(() => {
		const gapsDates: string[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (!towelRecords[i]) {
				gapsDates.unshift('-');
			} else {
				gapsDates.unshift(dayjs(towelRecords[i].time).format('D/M'));
			}
		}
		return gapsDates;
	});

	let towelTimes = $derived.by(() => getCalendarEntries(towels, 'Washed'));

	let vacationTimes = $derived.by(() => getCalendarEntries(vacations, 'Vacation', '✈️'));

	let singleDay: TowelDB[] | undefined = $state([]);
	let modal = $state<HTMLDialogElement>();

	let calendarOptions: Calendar.Options = $derived.by(() => {
		return {
			view: 'dayGridMonth',
			events: [...towelTimes, ...vacationTimes],
			selectBackgroundColor: 'red',
			eventBackgroundColor: 'var(--color-primary)',
			firstDay: 1,
			titleFormat: (date) => {
				const month = dayjs(date).get('month');
				if (month === 8) {
					return dayjs(date).format('MMMM YY');
				}
				return dayjs(date).format('MMMM YYYY');
			},
			dateClick: async (info) => {
				if (towels.isSuccess) {
					singleDay = towels.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.date).get('date');
					});
				}
				modal?.showModal();
			},
			eventClick: async (info) => {
				if (towels.isSuccess) {
					singleDay = towels.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
				}
				modal?.showModal();
			}
		};
	});

	let towelDirty: number | undefined = $derived.by(() => {
		let towelDirty;

		if (towelRecords && towelRecords.length > 0) {
			if (vacations.isSuccess) {
				const lastWashDate = dayjs(towelRecords[0].time);
				const now = dayjs();

				// 1. Calculate the raw difference in hours
				const rawDiffHours = now.diff(lastWashDate, 'hours', true);

				// 2. Calculate vacation hours between last wash and now
				const vacationHours = calculateVacationOverlap(
					now,
					lastWashDate,
					vacations.data ?? [],
					'hours'
				);

				// 3. Subtract vacation time
				towelDirty = rawDiffHours - vacationHours;

				// Ensure it's not negative
				return Math.max(0, towelDirty);
			}

			const lastWashDate = dayjs(towelRecords[0].time);
			const today = dayjs();
			/**
			 * https://stackoverflow.com/questions/36560806/the-left-hand-side-of-an-arithmetic-operation-must-be-of-type-any-number-or
			 */
			towelDirty = today.diff(lastWashDate, 'hours', true);
			return towelDirty;
		}
	});

	let lineChartEl = $state<HTMLCanvasElement>();
	let lineChart: Chart | undefined = $state();

	let currentTab = $state('overview');

	$effect(() => {
		if (lineChartEl && towels.isSuccess && currentTab === 'overview' && !lineChart) {
			lineChart = new Chart(lineChartEl, {
				type: 'line',
				options: { plugins: { legend: { display: false } } },
				data: {
					labels: [...gapsDates],
					datasets: [
						{
							label: '',
							data: gaps,
							fill: false,
							borderColor: '#3d6b5e',
							tension: 0.1,
							showLine: true
						}
					]
				}
			});
		}

		if (lineChart && towelRecords) {
			lineChart.data.labels = [...gapsDates];
			lineChart.data.datasets[0].data = gaps;
			lineChart.update();
		}
	});

	let status = $derived.by(() => {
		if (!towelDirty) {
			return 'empty';
		}
		return getStatusColorFromValue(towelDirty ?? 0);
	});

	const query = async () =>
		await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

	const refetch = async () => await tanstackClient.refetchQueries(createTowelRefetchOptions());
</script>

<PageWrapper title="Towel" {pb}>
	<div class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			{#key towelRecords}
				<StatusHeroImage {status} />
			{/key}

			<ActionButton {query} {refetch} text="Washed" />

			<div class="flex justify-start">
				<CustomDateModal collectionName="towel" {tanstackClient} />
			</div>
		</div>

		<div class="grid w-full content-start gap-8 pt-4 pb-8">
			<ul class="grid w-full grid-cols-2 content-center justify-items-center">
				<li class="w-full">
					<button
						class={[
							currentTab === 'overview' ? 'border-b-primary font-semibold' : 'border-b-base-300',
							'w-full cursor-pointer border-b-2 text-center'
						]}
						onclick={() => {
							currentTab = 'overview';
						}}
					>
						Overview</button
					>
				</li>
				<li class="w-full">
					<button
						class={[
							currentTab === 'calendar' ? 'border-b-primary font-semibold' : 'border-b-base-300',
							'w-full cursor-pointer border-b-2 text-center'
						]}
						onclick={() => {
							currentTab = 'calendar';
						}}
					>
						Calendar</button
					>
				</li>
			</ul>

			<div class="{currentTab === 'overview' ? 'grid' : 'hidden'} w-full gap-8 px-4">
				<TwoColumnCard leftTitle="Status" rightTitle="Last Washed">
					{#snippet left()}
						{#key towelRecords}
							{#if status}
								{@const descriptions = {
									green: 'Clean',
									yellow: 'Alright',
									orange: 'Funky',
									red: 'Wash now!'
								}}
								<StatusDescriptions {status} {descriptions} />
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/key}
					{/snippet}

					{#snippet right()}
						{#key towelRecords}
							{#if towelRecords && towelRecords.length > 0}
								{@const formatted = dayjs(towelRecords[0].time).fromNow()}
								{formatted}
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/key}
					{/snippet}
				</TwoColumnCard>

				<div class="border-base-content/5 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Gaps for Past Washes</h2>
					<div>
						<canvas bind:this={lineChartEl}></canvas>
					</div>
				</div>

				<div
					class="border-base-content/5 grid w-full grid-cols-2 content-center gap-4 rounded-lg border shadow"
				>
					<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
						<h2 class="text-md">Longest Gap</h2>
						<div class="text-center text-2xl font-bold">
							{#if longestGap}
								{longestGap.gap.toFixed(0)} days
							{:else}
								N/A
							{/if}
						</div>
					</div>

					<div class="grid justify-items-center p-4">
						<h2 class="text-md">Average Gap</h2>
						<div class="text-center text-2xl font-bold">
							{#if longestGap}
								{averageGap.toFixed(1)} days
							{:else}
								N/A
							{/if}
						</div>
					</div>
				</div>
			</div>

			<!-- The minmax(0, 1fr) or min-w-0 approaches explicitly tell the browser that the calendar can shrink below its content size, which is often what's needed for responsive calendar components. -->
			<div
				class="{currentTab === 'calendar'
					? 'grid'
					: 'hidden'} w-full grid-cols-[minmax(0,1fr)] overflow-hidden px-4"
			>
				{#key towels.data}
					<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
				{/key}
			</div>
		</div>
	</div>
</PageWrapper>

<SingleDayModal {modal} {singleDay} />
