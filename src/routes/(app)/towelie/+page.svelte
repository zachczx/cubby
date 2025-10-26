<script lang="ts">
	import Red from '$lib/images/red.webp?w=600&enhanced';
	import Orange from '$lib/images/orange.webp?w=600&enhanced';
	import Yellow from '$lib/images/yellow.webp?w=600&enhanced';
	import Green from '$lib/images/green.webp?w=600&enhanced';
	import { onMount } from 'svelte';
	import { pb } from '$lib/pb';
	import { goto } from '$app/navigation';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import Chart from 'chart.js/auto';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import { calculateVacationOverlap } from '$lib/overlap';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const vacations = createQuery<VacationDB[]>(() => ({
		queryKey: ['vacations'],
		queryFn: async () => await pb.collection('vacation').getFullList({ sort: '-time' })
	}));

	const towels = createQuery<TowelDB[]>(() => ({
		queryKey: ['towels'],
		queryFn: async () => await pb.collection('towel').getFullList({ sort: '-time' })
	}));

	const tanstackClient = useQueryClient();

	let towelRecords: TowelRecord[] | undefined = $derived.by(() => {
		if (towels.isSuccess) {
			return towels.data.map((record, i, allRecords) => {
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

		return [];
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
			if (!towelRecords[i]) {
				gaps.unshift(0);
			} else {
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

	let times = $derived.by(() => {
		let times: Calendar.EventInput[] = [];
		if (towels.isSuccess) {
			for (const r of towels.data) {
				//Adding the timezone here creates a problem for mobile devices, not sure why => .tz('Asia/Singapore');
				const n = dayjs.utc(r.time);
				times.push({ start: n.toDate(), end: n.toDate(), title: `— Washed` });
			}
		}
		return times;
	});

	let singleDay: TowelDB[] | undefined = $state([]);
	let singleDayModal = $state() as HTMLDialogElement;

	let calendarOptions: Calendar.Options = $derived.by(() => {
		return {
			view: 'dayGridMonth',
			events: [...times],
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
				if (towels.isSuccess) {
					singleDay = towels.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.date).get('date');
					});
				}
				singleDayModal.showModal();
			},
			eventClick: async (info) => {
				if (towels.isSuccess) {
					singleDay = towels.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
				}
				singleDayModal.showModal();
			}
		};
	});

	function getStatusColorFromValue(val: number): string {
		const day = 24;

		if (val > 0 && val <= 2 * day) return 'green';

		if (val > 2 * day && val <= 4 * day) return 'yellow';

		if (val > 4 * day && val <= 6 * day) return 'orange';

		if (val > 6 * day && val <= 999 * day) return 'red';

		return '';
	}

	let towelDirty: number | undefined = $derived.by(() => {
		let towelDirty;

		if (towelRecords && towelRecords.length > 0) {
			if (!vacations.isSuccess) {
				const lastWashDate = dayjs(towelRecords[0].time);
				const today = dayjs();
				/**
				 * https://stackoverflow.com/questions/36560806/the-left-hand-side-of-an-arithmetic-operation-must-be-of-type-any-number-or
				 */
				towelDirty = today.diff(lastWashDate, 'hours', true);
				return towelDirty;
			}

			const lastWashDate = dayjs(towelRecords[0].time);
			const now = dayjs();

			// 1. Calculate the raw difference in hours
			const rawDiffHours = now.diff(lastWashDate, 'hours', true);

			// 2. Calculate vacation hours between last wash and now
			const vacationHours = calculateVacationOverlap(now, lastWashDate, vacations.data, 'hours');

			// 3. Subtract vacation time
			towelDirty = rawDiffHours - vacationHours;

			// Ensure it's not negative
			return Math.max(0, towelDirty);
		}
	});
	let spinner = $state(false);

	let lineChartEl = $state() as HTMLCanvasElement;
	let lineChart: Chart | undefined = $state();

	let currentTab = $state('overview');

	onMount(async () => {
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
		if (lineChart && towelRecords) {
			lineChart.data.labels = [...gapsDates];
			lineChart.data.datasets[0].data = gaps;
			lineChart.update();
		}
	});

	let status = $derived.by(() => {
		if (!towelDirty) {
			return '';
		}
		return getStatusColorFromValue(towelDirty);
	});

	let modal = $state() as HTMLDialogElement;

	async function addHandler() {
		const result = await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			spinner = false;
		}

		await tanstackClient.refetchQueries({
			queryKey: ['towels'],
			type: 'active',
			exact: true
		});
	}

	function stripNegative(input: number | undefined): string {
		if (!input) return '';

		if (input < 0) {
			input = input * -1;
		}

		let clean = '';

		if (input === 0) {
			clean = input.toFixed(0).toString() + ' day';
		} else if (input <= 1) {
			clean = input.toFixed(1).toString() + ' day';
		} else {
			clean = input.toFixed(1).toString() + ' days';
		}

		return clean;
	}

	function getColorClass(gapDays: number | undefined): string {
		if (!gapDays) return '';

		if (gapDays > 0 && gapDays <= 2) return 'lime-500';

		if (gapDays > 2 && gapDays <= 4) return 'yellow-500';

		if (gapDays > 4 && gapDays <= 6) return 'orange-400';

		if (gapDays > 6 && gapDays <= 999) return 'red-700';

		return '';
	}
</script>

<svelte:head>
	<title>Towelie</title>
</svelte:head>

<PageWrapper title="Towelie" {pb} back={true} noPadding={true}>
	<div
		class="grid w-full max-w-[1200px] content-start justify-items-center gap-4 justify-self-center lg:grid-cols-2 lg:pt-8"
	>
		<div class="grid content-start justify-items-center gap-4 px-4 pt-8">
			{#key towelRecords}
				{#if status === 'green'}
					<enhanced:img src={Green} alt="Clean" class="rounded-3xl" />
				{:else if status === 'yellow'}
					<enhanced:img src={Yellow} alt="Still Good" class="rounded-3xl" />
				{:else if status === 'orange'}
					<enhanced:img src={Orange} alt="Ripening" class="rounded-3xl" />
				{:else if status === 'red'}
					<enhanced:img src={Red} alt="Time to Wash!" class="rounded-3xl" />
				{:else}
					<div class="skeleton h-[600px] w-[600px] max-lg:h-[342px] max-lg:w-[342px]"></div>
				{/if}
			{/key}
			<form class="w-full">
				<button
					class="btn btn-xl btn-primary flex w-full items-center gap-2 rounded-2xl"
					onclick={addHandler}>Just Washed My Towel!</button
				>
			</form>
		</div>

		<div class="grid w-full content-start gap-8 py-4">
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
				<div class="w-full">
					<div
						class="text-md border-base-content/5 bg-base-200/50 grid grid-cols-2 content-center gap-4 rounded-lg border shadow"
					>
						<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
							<div>Status</div>
							{#key towelRecords}
								{#if status}
									<div
										class="flex min-h-20 items-center justify-center gap-4 text-center text-2xl font-bold"
									>
										{#if status === 'green'}
											<div class="hidden h-6 w-6 rounded-full bg-lime-500 lg:flex"></div>
											<span class="text-lime-500">Squeaky Clean</span>
										{:else if status === 'yellow'}
											<div class="hidden h-6 w-6 rounded-full bg-yellow-500 lg:flex"></div>
											<span class="text-yellow-500">Still Fresh</span>
										{:else if status === 'orange'}
											<div class="hidden h-6 w-6 rounded-full bg-orange-400 lg:flex"></div>
											<span class="text-orange-400">Kinda Funky</span>
										{:else}
											<div class="hidden h-5 w-5 rounded-full bg-red-700 lg:flex"></div>
											<span class="text-red-700">Wash Me!</span>
										{/if}
									</div>
								{:else}
									<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
								{/if}
							{/key}
						</div>
						<div class="grid justify-items-center p-4">
							<div>Last Washed</div>
							<div class="text-center text-2xl font-bold">
								{#key towelRecords}
									{#if towelRecords && towelRecords.length > 0}
										{@const formatted = dayjs(towelRecords[0].time).fromNow()}
										{formatted}
									{:else}
										<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
									{/if}
								{/key}
							</div>
						</div>
					</div>
				</div>

				<div class="border-base-content/5 bg-base-200/50 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Gaps for Past Washes</h2>
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
				{#key towels.data}
					<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
				{/key}
			</div>
		</div>
	</div>
</PageWrapper>

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
