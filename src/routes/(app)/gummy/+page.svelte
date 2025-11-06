<script lang="ts">
	import Red from '$lib/images/nosey_red.webp?w=150&enhanced';
	import Orange from '$lib/images/nosey_orange.webp?w=150&enhanced';
	import Yellow from '$lib/images/nosey_olive.webp?w=150&enhanced';
	import Green from '$lib/images/nosey_green.webp?w=150&enhanced';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import timezone from 'dayjs/plugin/timezone';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import Chart from 'chart.js/auto';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { onMount } from 'svelte';
	import {
		createGummyQueryOptions,
		createGummyRefetchOptions,
		createSprayRefetchOptions,
		createUserQueryOptions,
		createVacationQueryOptions
	} from '$lib/queries';
	import { getCalendarEntries } from '$lib/calendar';
	import EmptyState from '$lib/assets/svg/EmptyState.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let singleDay: GummyDB[] | undefined = $state([]);
	let singleDayModal = $state() as HTMLDialogElement;

	let gummyButtonStatus: 'default' | 'loading' | 'success' = $state('default');

	const gummies = createQuery(createGummyQueryOptions);
	const user = createQuery(createUserQueryOptions);
	const vacations = createQuery(createVacationQueryOptions);
	const tanstackClient = useQueryClient();

	async function addSprayHandler() {
		gummyButtonStatus = 'loading';

		const result = await pb.collection('gummy').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: daysToNext
		});

		if (result.id) {
			addToast('success', 'Added successfully!');

			gummyButtonStatus = 'success';

			setTimeout(() => {
				gummyButtonStatus = 'default';
			}, 3000);
		}

		await tanstackClient.refetchQueries(createGummyRefetchOptions());
	}

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
					singleDayModal.showModal();
				}
			},
			eventClick: async (info) => {
				if (gummies.isSuccess) {
					singleDay = gummies.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
					singleDayModal.showModal();
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

		return user.data?.defaultGummyInterval;
	});

	// For Statuses

	let gummyRecords: GummyRecord[] | undefined = $derived.by(() => {
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
			if (!gummyRecords[i]) {
				gaps.unshift(0);
			} else {
				gaps.unshift(gummyRecords[i].gap);
			}
		}
		return gaps;
	});

	let gapsDates = $derived.by(() => {
		if (!gummyRecords || gummyRecords.length === 0) return [];

		const gapsDates: string[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (!gummyRecords[i]) {
				break;
			} else {
				gapsDates.unshift(dayjs(gummyRecords[i].time).format('D/M'));
			}
		}
		return gapsDates;
	});

	let lineChartEl = $state() as HTMLCanvasElement;
	let lineChart: Chart | undefined = $state();

	onMount(() => {
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
						borderColor: 'oklch(56.86% 0.255 257.57)',
						tension: 0.1,
						showLine: true
					}
				]
			}
		});
	});

	$effect(() => {
		if (lineChart && gapsDates && gaps) {
			lineChart.data.labels = [...gapsDates];
			lineChart.data.datasets[0].data = gaps;
			lineChart.update();
		}
	});
</script>

<svelte:head>
	<title>Spray</title>
</svelte:head>

<PageWrapper title="Spray" back={true} {pb}>
	<main class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			{#if gummies.isSuccess}
				{#if status === 'green'}
					<div class="avatar">
						<div class="w-32 rounded-full bg-[#dbf0be]">
							<enhanced:img src={Green} alt="Still Good" class="" />
						</div>
					</div>
				{:else if status === 'yellow'}
					<div class="avatar">
						<div class="w-32 rounded-full bg-yellow-200/70">
							<enhanced:img src={Yellow} alt="Still Good" class="rounded-2xl" />
						</div>
					</div>
				{:else if status === 'orange'}
					<div class="avatar">
						<div class="w-32 rounded-full bg-orange-200/70">
							<enhanced:img src={Orange} alt="Ripening" class="rounded-2xl" />
						</div>
					</div>
				{:else if status === 'red'}
					<div class="avatar">
						<div class="w-32 rounded-full bg-red-200/70">
							<enhanced:img src={Red} alt="Time to Wash!" class="rounded-2xl" />
						</div>
					</div>
				{:else if status === 'empty'}
					<EmptyState class="text-primary/30 my-4" />
				{:else}
					<div class="skeleton h-[600px] w-[600px] max-lg:h-[342px] max-lg:w-[342px]"></div>
				{/if}
			{/if}

			<button
				class={[
					'btn btn-lg flex w-full items-center gap-2 rounded-full',
					gummyButtonStatus === 'default' && 'btn-primary',
					gummyButtonStatus === 'loading' && 'btn-primary',
					gummyButtonStatus === 'success' && 'btn-success'
				]}
				onclick={addSprayHandler}
			>
				{#if gummyButtonStatus === 'success'}
					<MaterialSymbolsCheck class="size-6" />Added!
				{:else if gummyButtonStatus === 'loading'}
					<span class="loading loading-spinner loading-md"></span>
				{:else}
					Just Ate Gummy
				{/if}
			</button>
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
				<div
					class="text-md border-base-content/5 bg-base-200/50 grid grid-cols-2 content-center gap-4 rounded-lg border shadow"
				>
					<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
						<div>Status</div>
						{#if gummies.isSuccess}
							{#if status}
								<div
									class="flex min-h-14 items-center justify-center gap-4 text-center text-2xl font-bold lg:min-h-20"
								>
									{#if status === 'green'}
										<div class="hidden h-6 w-6 rounded-full bg-lime-500 lg:flex"></div>
										<span class="text-lime-500">Dosed</span>
									{:else if status === 'yellow'}
										<div class="hidden h-6 w-6 rounded-full bg-yellow-500 lg:flex"></div>
										<span class="text-yellow-500">Good</span>
									{:else if status === 'orange'}
										<div class="hidden h-6 w-6 rounded-full bg-orange-400 lg:flex"></div>
										<span class="text-orange-400">Due</span>
									{:else if status === 'red'}
										<div class="hidden h-5 w-5 rounded-full bg-red-700 lg:flex"></div>
										<span class="text-red-700">Overdue</span>
									{:else}
										<span>Nil</span>
									{/if}
								</div>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
					</div>
					<div class="grid justify-items-center p-4">
						<div>Last Gummy</div>
						<div class="text-center text-2xl font-bold">
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
						</div>
					</div>
				</div>

				<div class="border-base-content/5 bg-base-200/50 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Intervals</h2>
					<div>
						<canvas bind:this={lineChartEl}></canvas>
					</div>
				</div>

				<div
					class="border-base-content/5 bg-base-200/50 grid w-full grid-cols-2 content-center gap-4 rounded-lg border shadow"
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
<span class="ec ec-time-grid ec-title hidden"></span>

<dialog bind:this={singleDayModal} class="modal">
	<div class="modal-box">
		{#if singleDay && singleDay.length > 0}
			{@const theDay = dayjs(singleDay[0].time).format('D MMM YYYY')}
			<h3 class="text-lg font-bold">{theDay}</h3>
			<div class="py-4">
				<ul class="list-disc">
					{#each singleDay as day}
						{@const formatted = dayjs(day.time).format('h:mma')}
						<li class="flex items-center gap-2">
							<MaterialSymbolsCheck class="text-success size-[1.3em]" />
							{formatted}
						</li>
					{/each}
				</ul>
			</div>
		{:else}
			<h3 class="text-lg font-bold">Nothing here!</h3>
		{/if}
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">Close</button>
			</form>
		</div>
	</div>
</dialog>

<style>
	:global {
		.ec.ec-day-grid {
			/* max-width: 90vw; */
			width: 100%;
		}

		.ec-title {
			font-size: 1.25rem;
			font-weight: 600 !important;
		}

		.ec-today {
			background-color: var(--color-primary-content) !important;
		}

		.ec-day {
			&:hover {
				background-color: rgba(242, 189, 205, 0.3);
				.ec-day-head {
					font-weight: 600 !important;
				}
				color: var(--color-primary);
			}
		}

		.ec-bg-events {
			&:hover {
				cursor: pointer;
			}
		}

		.ec-header {
			border-top-left-radius: 10px !important;
			border-top-right-radius: 10px !important;
		}
	}
</style>
