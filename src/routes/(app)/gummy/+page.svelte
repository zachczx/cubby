<script lang="ts">
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import timezone from 'dayjs/plugin/timezone';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import Chart from 'chart.js/auto';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createGummyQueryOptions,
		createGummyRefetchOptions,
		createUserQueryOptions,
		createVacationQueryOptions
	} from '$lib/queries';
	import { getCalendarEntries } from '$lib/calendar';
	import CustomDateModal from '$lib/ui/CustomDateModal.svelte';
	import TwoColumnCard from '$lib/ui/TwoColumnCard.svelte';
	import StatusDescriptions from '$lib/ui/StatusDescriptions.svelte';
	import StatusHeroImage from '$lib/ui/StatusHeroImage.svelte';
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import SingleDayModal from '$lib/ui/SingleDayModal.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let singleDay: GummyDB[] | undefined = $state([]);
	let modal = $state<HTMLDialogElement>();

	const gummies = createQuery(createGummyQueryOptions);
	const user = createQuery(createUserQueryOptions);
	const vacations = createQuery(createVacationQueryOptions);
	const tanstackClient = useQueryClient();

	const query = async () =>
		await pb.collection('gummy').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: daysToNext
		});
	const refetch = async () => await tanstackClient.refetchQueries(createGummyRefetchOptions());

	let times = $derived.by(() => getCalendarEntries(gummies, 'Gummy'));
	let vacationTimes = $derived.by(() => getCalendarEntries(vacations, 'Vacation', '✈️'));

	let calendarOptions: Calendar.Options = $derived.by(() => {
		return {
			view: 'dayGridMonth',
			events: [...times, ...vacationTimes],
			selectBackgroundColor: 'red',
			eventBackgroundColor: '#4a4a7d',
			firstDay: 1,
			titleFormat: (date) => {
				const month = dayjs(date).get('month');
				if (month === 8) {
					return dayjs(date).format('MMMM YY');
				}
				return dayjs(date).format('MMMM YYYY');
			},
			dateClick: async (info) => {
				if (gummies.isSuccess) {
					singleDay = gummies.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.date).get('date');
					});
					modal?.showModal();
				}
			},
			eventClick: async (info) => {
				if (gummies.isSuccess) {
					singleDay = gummies.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
					modal?.showModal();
				}
			}
		};
	});

	function getStatusColorFromValue(val: number): string {
		if (!daysToNext) return '';

		if (val === 0) return 'empty';

		const day = 24;

		if (val > 0 && val <= daysToNext * day) return 'green';

		if (val > 2 * day && val <= 4 * day) return 'yellow';

		if (val > 4 * day && val <= 6 * day) return 'orange';

		if (val > daysToNext * day && val <= 999 * day) return 'red';

		return '';
	}

	let lastGummy: number | undefined = $derived.by(() => {
		let lastGummy;
		if (gummies.isSuccess) {
			if (gummies.data && gummies.data.length > 0) {
				const lastGummyDate = dayjs(gummies.data[0].time);
				const today = dayjs();
				lastGummy = today.diff(lastGummyDate, 'hours', true);
				return lastGummy;
			}
		}
	});

	let status = $derived.by(() => {
		if (!lastGummy) {
			return 'empty';
		}

		return getStatusColorFromValue(lastGummy);
	});

	let currentTab = $state('overview');

	let daysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.gummyInterval;
	});

	let gummyRecords = $derived.by(() => {
		if (gummies.isSuccess && gummies.data) {
			return gummies.data.map((record, i, allRecords) => {
				const nextRecord = allRecords[i + 1];
				const gap = nextRecord ? dayjs(record.time).diff(nextRecord.time, 'day', true) : 0;
				return { ...record, gap };
			});
		}
		return [];
	});

	let longestGap: GummyRecord | undefined = $derived.by(() => {
		if (gummyRecords.length <= 1) return;
		const avoidMutatingOriginalArray = [...gummyRecords];
		const sorted = avoidMutatingOriginalArray.sort((a, b) => b.gap - a.gap);
		return sorted[0];
	});
	let averageGap: number = $derived.by(() => {
		const total = gummyRecords.reduce((accumulator, entry) => {
			return (accumulator += entry.gap);
		}, 0);
		const average = total / gummyRecords.length;
		return average;
	});

	let gaps = $derived.by(() => {
		const gaps: number[] = [];
		const numberOfRecords = 10;

		for (let i = 0; i < numberOfRecords; i++) {
			if (gummyRecords[i]) {
				gaps.unshift(gummyRecords[i].gap);
			}
		}
		return gaps;
	});

	let gapsDates = $derived.by(() => {
		const gapsDates: string[] = [];
		const numberOfRecords = 10;

		for (let i = 0; i < numberOfRecords; i++) {
			if (gummyRecords[i]) {
				gapsDates.unshift(dayjs(gummyRecords[i].time).format('D/M'));
			}
		}
		return gapsDates;
	});

	let lineChartEl = $state<HTMLCanvasElement>();
	let lineChart: Chart | undefined = $state();

	$effect(() => {
		if (lineChartEl && gummies.isSuccess && currentTab === 'overview' && !lineChart) {
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

		if (lineChart && gapsDates.length > 0 && gaps.length > 0) {
			lineChart.data.labels = [...gapsDates];
			lineChart.data.datasets[0].data = gaps;
			lineChart.update();
		}
	});
</script>

<PageWrapper title="Gummy" {pb}>
	<main class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			{#if gummies.isSuccess}
				<StatusHeroImage {status} />
			{/if}

			<ActionButton {query} {refetch} text="Ate" />

			<div class="flex justify-start">
				<CustomDateModal collectionName="gummy" {tanstackClient} {daysToNext} />
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
				<TwoColumnCard leftTitle="Status" rightTitle="Last Ate">
					{#snippet left()}
						{#if gummies.isSuccess}
							{#if status}
								{@const descriptions = {
									green: 'Dosed',
									yellow: 'Good',
									orange: 'Due',
									red: 'Overdue'
								}}
								<StatusDescriptions {status} {descriptions} />
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
					{/snippet}

					{#snippet right()}
						{#if gummies.isSuccess}
							{#if gummies.data.length > 0}
								{@const formatted = dayjs(gummies.data[0].time).fromNow()}
								{formatted}
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
						{#if gummies.isPending}
							<div class="custom-loader"></div>
						{/if}
					{/snippet}
				</TwoColumnCard>

				<div class="border-base-content/5 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Intervals</h2>
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
				{#key gummies.data}
					<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
				{/key}
			</div>
		</div>
	</main>
</PageWrapper>

<SingleDayModal {modal} {singleDay} />
