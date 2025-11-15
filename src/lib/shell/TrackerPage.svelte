<script lang="ts">
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import timezone from 'dayjs/plugin/timezone';
	import calendar from 'dayjs/plugin/calendar';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Chart from 'chart.js/auto';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { createUserQueryOptions, createVacationQueryOptions } from '$lib/queries';
	import { getCalendarEntries } from '$lib/calendar';
	import CustomDateModal from '$lib/ui/CustomDateModal.svelte';
	import StatusDescriptions from '$lib/ui/StatusDescriptions.svelte';
	import TwoColumnCard from '$lib/ui/TwoColumnCard.svelte';
	import StatusHeroImage from '$lib/ui/StatusHeroImage.svelte';
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import SingleDayModal from '$lib/ui/SingleDayModal.svelte';
	import { defaultNotificationStatus, getNotificationStatus } from '$lib/notification';

	dayjs.extend(calendar);
	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let { options }: { options: TrackerPageOptions } = $props();

	let singleDay: TrackerDB[] | undefined = $state([]);
	let modal = $state<HTMLDialogElement>();

	const dbRecords = createQuery<TrackerDB[]>(options.queries.query);
	const user = createQuery(createUserQueryOptions);
	const vacations = createQuery(createVacationQueryOptions);
	const tanstackClient = useQueryClient();

	const query = async () =>
		await pb.collection(options.collectionName).create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: daysToNext
		});

	const refetch = async () => await tanstackClient.refetchQueries(options.queries.refetch());

	let times = $derived.by(() => getCalendarEntries(dbRecords, options.labels.ctaButtonText));
	let vacationTimes = $derived.by(() => getCalendarEntries(vacations, 'Vacation', '✈️'));

	let calendarOptions: Calendar.Options = $derived.by(() => {
		return {
			view: 'dayGridMonth',
			events: [...times, ...vacationTimes],
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
				if (dbRecords.isSuccess) {
					singleDay = dbRecords.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.date).get('date');
					});
					modal?.showModal();
				}
			},
			eventClick: async (info) => {
				if (dbRecords.isSuccess) {
					singleDay = dbRecords.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
					modal?.showModal();
				}
			}
		};
	});

	let last: number | undefined = $derived.by(() => {
		let last;
		if (dbRecords.isSuccess) {
			if (dbRecords.data && dbRecords.data.length > 0) {
				const lastRecord = dayjs(dbRecords.data[0].time);
				const today = dayjs();
				last = today.diff(lastRecord, 'hours', true);
				return last;
			}
		}
	});

	let notification = $derived.by(() => getNotificationStatus(dbRecords));

	type TabPages = 'overview' | 'stats' | 'calendar';

	let currentTab = $state<TabPages>('overview');

	let daysToNext = $derived.by(() => (user.isSuccess ? user.data?.sprayIntervalDays : undefined));

	/**
	 * Using $state + $effect instead of $derived due to TanStack Query store
	 * not properly triggering Svelte 5's fine-grained reactivity on async data changes
	 */

	let records: TrackerRecord[] = $derived.by(() => {
		if (dbRecords.isSuccess && dbRecords.data) {
			if (options.calculateGaps) {
				return options.calculateGaps(dbRecords.data, vacations.data ?? []);
			} else {
				return dbRecords.data.map((record, i, allRecords) => {
					const nextRecord = allRecords[i + 1];
					const gap = nextRecord ? dayjs(record.time).diff(nextRecord.time, 'day', true) : 0;
					return { ...record, gap };
				});
			}
		}

		return [];
	});

	let longestGap: TrackerRecord | undefined = $derived.by(() => {
		if (records.length <= 1) return;
		const avoidMutatingOriginalArray = [...records];
		const sorted = avoidMutatingOriginalArray.sort((a, b) => b.gap - a.gap);
		return sorted[0];
	});
	let averageGap: number = $derived.by(() => {
		const total = records.reduce((accumulator, entry) => {
			return (accumulator += entry.gap);
		}, 0);
		const average = total / records.length;
		return average;
	});

	let gaps = $derived.by(() => {
		const gaps: number[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (records[i]) {
				gaps.unshift(records[i].gap);
			}
		}
		return gaps;
	});

	let gapsDates = $derived.by(() => {
		const gapsDates: string[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (!records[i]) {
				gapsDates.unshift('-');
			} else {
				gapsDates.unshift(dayjs(records[i].time).format('D/M'));
			}
		}
		return gapsDates;
	});

	let lineChartEl = $state<HTMLCanvasElement>();
	let lineChart: Chart | undefined = $state();

	$effect(() => {
		if (lineChartEl && dbRecords.isSuccess && currentTab === 'stats' && !lineChart) {
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

		if (lineChart && gapsDates && gaps) {
			lineChart.data.labels = [...gapsDates];
			lineChart.data.datasets[0].data = gaps;
			lineChart.update();
		}
	});

	const dayjsCalendarOptions = {
		sameDay: '[Today], h:mma', // The same day ( Today at 2:30 AM )
		nextDay: '[Tomorrow], h:mma', // The next day ( Tomorrow at 2:30 AM )
		nextWeek: 'dddd, h:mma', // The next week ( Sunday at 2:30 AM )
		lastDay: '[Yesterday], h:mma', // The day before ( Yesterday at 2:30 AM )
		lastWeek: '[Last] dddd, h:mma', // Last week ( Last Monday at 2:30 AM )
		sameElse: 'DD/MM/YYYY' // Everything else ( 17/10/2011 )
	};
</script>

<PageWrapper title={options.labels.pageTitle} {pb}>
	<main class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			{#if dbRecords.isSuccess}
				<StatusHeroImage {notification} />
			{/if}

			<ActionButton {query} {refetch} text={options.labels.ctaButtonText} />

			<div class="flex justify-start">
				<CustomDateModal collectionName={options.collectionName} {tanstackClient} {daysToNext} />
			</div>
		</div>

		<div class="grid w-full content-start gap-8 pt-4 pb-8">
			<ul class="tracker-tabs grid w-full grid-cols-3 content-center justify-items-center">
				<li class="w-full" aria-current={currentTab === 'overview' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'overview')}
					>
						Overview</button
					>
				</li>

				<li class="w-full" aria-current={currentTab === 'stats' ? 'page' : undefined}>
					<button class="w-full cursor-pointer text-center" onclick={() => (currentTab = 'stats')}>
						Stats</button
					>
				</li>

				<li class="w-full" aria-current={currentTab === 'calendar' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'calendar')}
					>
						Calendar</button
					>
				</li>
			</ul>

			<div class={['grid w-full gap-8 px-4', currentTab === 'overview' ? undefined : 'hidden']}>
				<div class="flex items-center justify-center text-2xl font-bold">
					{#if dbRecords.isSuccess}
						{#if notification}
							<StatusDescriptions {notification} />
						{:else}
							<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
						{/if}
					{/if}
				</div>

				<div
					class="border-base-content/5 bg-primary/10 grid justify-items-center rounded-xl border p-4 shadow"
				>
					<h2 class="text-md text-center">Next</h2>
					<div class="grid min-h-20 content-center justify-items-center">
						{#if dbRecords.isSuccess}
							{#if notification}
								{@const semantic = dayjs(notification.next).calendar(dayjs(), dayjsCalendarOptions)}
								{@const formatted = dayjs(notification.next).format('D MMM YYYY')}
								<p class="text-primary text-2xl font-bold">{semantic}</p>
								<p>{formatted}</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
					</div>
				</div>

				<TwoColumnCard leftTitle="Frequency" rightTitle="Last">
					{#snippet left()}
						{#if user.isSuccess}
							{@const collectionName = dbRecords.data?.[0].collectionName}
							<p>
								{#if collectionName === 'spray'}
									{user.data.sprayIntervalDays} days
								{:else if collectionName === 'towel'}
									{user.data.towelIntervalDays} days
								{:else if collectionName === 'gummy'}
									{user.data.gummyIntervalDays} days
								{:else if collectionName === 'doggoBath'}
									{user.data.doggoBathIntervalDays} days
								{:else if collectionName === 'doggoChewable'}
									{user.data.doggoChewableIntervalMonths}
									{#if user.data.doggoChewableIntervalMonths === 1}month{:else}months{/if}
								{/if}
							</p>
							<p class="text-base-content/70 text-base font-normal">once</p>
						{/if}
					{/snippet}

					{#snippet right()}
						{#if dbRecords.isSuccess}
							{#if dbRecords.data.length > 0}
								{@const formatted = dayjs(dbRecords.data[0].time).fromNow(true)}
								<p>{formatted}</p>
								<p class="text-base-content/70 text-base font-normal">ago</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
						{#if dbRecords.isPending}
							<div class="custom-loader"></div>
						{/if}
					{/snippet}
				</TwoColumnCard>
			</div>

			<div class={['grid w-full gap-8 px-4', currentTab === 'stats' ? undefined : 'hidden']}>
				<div class="border-base-content/5 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Trend</h2>
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

			<div
				class={[
					'grid w-full grid-cols-[minmax(0,1fr)] overflow-hidden px-4',
					currentTab === 'calendar' ? undefined : 'hidden'
				]}
			>
				{#key dbRecords.data}
					<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
				{/key}
			</div>
		</div>
	</main>
</PageWrapper>

<SingleDayModal {modal} {singleDay} />

<style>
	.tracker-tabs {
		li {
			position: relative;
			padding-bottom: 0.25rem;

			&::after {
				position: absolute;
				content: '';
				bottom: 0;
				left: 50%;
				width: 100%;
				transform: translateX(-50%);
				height: 0.125rem;

				background-color: var(--color-base-300);
			}

			&[aria-current='page'] {
				font-weight: bold;
			}

			&[aria-current='page']::after {
				/* view-transition-name: active-tab; */
				position: absolute;
				content: '';
				bottom: 0px;
				left: 50%;
				width: 100%;
				transform: translateX(-50%);
				height: 0.125rem;

				background-color: var(--color-primary);
			}
		}
	}
</style>
