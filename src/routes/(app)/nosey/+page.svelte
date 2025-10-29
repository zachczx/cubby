<script lang="ts">
	import Red from '$lib/images/red.webp?w=600&enhanced';
	import Orange from '$lib/images/orange.webp?w=600&enhanced';
	import Yellow from '$lib/images/yellow.webp?w=600&enhanced';
	import Green from '$lib/images/green.webp?w=600&enhanced';
	import { pb } from '$lib/pb';
	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import timezone from 'dayjs/plugin/timezone';
	import { Calendar, DayGrid, Interaction } from '@event-calendar/core';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsAdd from '$lib/assets/svg/MaterialSymbolsAdd.svelte';
	import Chart from 'chart.js/auto';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let singleDay: SprayDB[] | undefined = $state([]);
	let singleDayModal = $state() as HTMLDialogElement;

	const sprays = createQuery<SprayDB[]>(() => ({
		queryKey: ['sprays', pb.authStore?.record?.id],
		queryFn: async () => await pb.collection('spray').getFullList({ sort: '-time' })
	}));

	const userPref = createQuery<UserDB>(() => ({
		queryKey: ['user', pb.authStore?.record?.id],
		queryFn: async () => await pb.collection('users').getOne(String(pb.authStore?.record?.id))
	}));

	const tanstackClient = useQueryClient();

	async function handleClick() {
		spinner = true;
		const result = await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: daysToNext
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			spinner = false;
		}

		await tanstackClient.refetchQueries({
			queryKey: ['sprays'],
			type: 'active',
			exact: true
		});
	}

	let times = $derived.by(() => {
		let times: Calendar.EventInput[] = [];
		if (sprays.isSuccess) {
			for (const r of sprays.data) {
				//Adding the timezone here creates a problem for mobile devices, not sure why => .tz('Asia/Singapore');
				const n = dayjs.utc(r.time);
				times.push({ start: n.toDate(), end: n.toDate(), title: `— Sprayed` });
			}
		}
		return times;
	});

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
				if (sprays.isSuccess) {
					singleDay = sprays.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.date).get('date');
					});
					singleDayModal.showModal();
				}
			},
			eventClick: async (info) => {
				if (sprays.isSuccess) {
					singleDay = sprays.data.filter((day) => {
						return dayjs(day.time).get('date') == dayjs(info.event.start).get('date');
					});
					singleDayModal.showModal();
				}
			}
		};
	});

	let spinner = $state(false);

	function getStatusColorFromValue(val: number): string {
		if (!daysToNext) return '';

		const day = 24;

		if (val > 0 && val <= daysToNext * day) return 'green';

		if (val > 2 * day && val <= 4 * day) return 'yellow';

		if (val > 4 * day && val <= 6 * day) return 'orange';

		if (val > daysToNext * day && val <= 999 * day) return 'red';

		return '';
	}

	let lastSpray: number | undefined = $derived.by(() => {
		let lastSpray;
		if (sprays.isSuccess) {
			const lastWashDate = dayjs(sprays.data[0].time);
			const today = dayjs();
			lastSpray = today.diff(lastWashDate, 'hours', true);
			return lastSpray;
		}
	});

	let status = $derived.by(() => {
		if (!lastSpray) {
			return '';
		}
		return getStatusColorFromValue(lastSpray);
	});

	let currentTab = $state('overview');

	let daysToNext = $derived.by(() => {
		if (userPref.isPending) {
			return undefined;
		}

		return userPref.data?.defaultSprayInterval;
	});

	// For Statuses

	let sprayRecords: SprayRecord[] | undefined = $derived.by(() => {
		if (sprays.isSuccess) {
			return sprays.data.map((record, i, allRecords) => {
				const nextRecord = allRecords[i + 1];
				const gap = nextRecord ? dayjs(record.time).diff(nextRecord.time, 'day', true) : 0;
				return { ...record, gap };
			});
		}

		return [];
	});

	let longestGap: SprayRecord | undefined = $derived.by(() => {
		if (sprayRecords.length <= 1) return;
		const avoidMutatingOriginalArray = [...sprayRecords];
		const sorted = avoidMutatingOriginalArray.sort((a, b) => b.gap - a.gap);
		return sorted[0];
	});
	let averageGap: number = $derived.by(() => {
		const total = sprayRecords.reduce((accumulator, entry) => {
			return (accumulator += entry.gap);
		}, 0);
		const average = total / sprayRecords.length;
		return average;
	});

	let gaps = $derived.by(() => {
		const gaps: number[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (!sprayRecords[i]) {
				gaps.unshift(0);
			} else {
				gaps.unshift(sprayRecords[i].gap);
			}
		}
		return gaps;
	});

	let gapsDates = $derived.by(() => {
		const gapsDates: string[] = [];
		const numberOfRecords = 10;
		for (let i = 0; i < numberOfRecords; i++) {
			if (!sprayRecords[i]) {
				gapsDates.unshift('-');
			} else {
				gapsDates.unshift(dayjs(sprayRecords[i].time).format('D/M'));
			}
		}
		return gapsDates;
	});

	let lineChartEl = $state() as HTMLCanvasElement;
	let lineChart: Chart | undefined = $state();

	onMount(async () => {
		if (pb.authStore.isValid) {
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
		}
	});
</script>

<svelte:head>
	<title>Nosey</title>
</svelte:head>

<PageWrapper title="Spray Logs" back={true} {pb} noPadding={true}>
	<main
		class="grid w-full max-w-[1200px] content-start justify-items-center gap-4 justify-self-center lg:grid-cols-2 lg:pt-8"
	>
		<div class="grid content-start justify-items-center gap-4 px-4 pt-8">
			{#if sprays.isSuccess}
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
			{/if}

			<button
				class="btn btn-xl btn-primary flex w-full min-w-54 grow items-center gap-2 rounded-2xl"
				onclick={handleClick}
			>
				{#if !spinner}
					Just Sprayed!
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
		</div>

		<div class="grid w-full content-start gap-8 py-4">
			<ul class="grid w-full grid-cols-3 content-center justify-items-center">
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
				<li class="w-full">
					<button
						class={[
							currentTab === 'settings' ? 'border-b-primary font-semibold' : 'border-b-base-300',
							'w-full cursor-pointer border-b-2 text-center'
						]}
						onclick={() => {
							currentTab = 'settings';
						}}
					>
						Settings</button
					>
				</li>
			</ul>

			<div class="{currentTab === 'overview' ? 'grid' : 'hidden'} w-full gap-8 px-4">
				<div
					class="text-md border-base-content/5 bg-base-200/50 grid grid-cols-2 content-center gap-4 rounded-lg border shadow"
				>
					<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
						<div>Status</div>
						{#if sprays.isSuccess}
							{#if status}
								<div
									class="flex min-h-20 items-center justify-center gap-4 text-center text-2xl font-bold"
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
									{:else}
										<div class="hidden h-5 w-5 rounded-full bg-red-700 lg:flex"></div>
										<span class="text-red-700">Overdue</span>
									{/if}
								</div>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						{/if}
					</div>
					<div class="grid justify-items-center p-4">
						<div>Last Sprayed</div>
						<div class="text-center text-2xl font-bold">
							{#if sprays.isSuccess}
								{#if sprays.data.length > 0}
									{@const formatted = dayjs(sprays.data[0].time).fromNow()}
									{formatted}
								{:else}
									<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
								{/if}
							{/if}
							{#if sprays.isPending}
								<div class="custom-loader"></div>
							{/if}
						</div>
					</div>
				</div>

				<div class="border-base-content/5 bg-base-200/50 w-full rounded-lg border p-4 shadow">
					<h2 class="text-md text-center">Dose Intervals</h2>
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
				{#key sprays.data}
					<Calendar plugins={[DayGrid, Interaction]} options={calendarOptions} />
				{/key}
			</div>

			<div class="{currentTab === 'settings' ? 'grid' : 'hidden'} w-full px-4">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Days Between Doses</legend>
					<select bind:value={daysToNext} class="select">
						<option value={1}>1</option>
						<option value={2}>2</option>
						<option value={3}>3</option>
						<option value={4}>4</option>
						<option value={5}>5</option>
					</select>
				</fieldset>
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
