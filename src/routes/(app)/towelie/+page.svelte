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

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelDB: TowelDB[] | undefined = $state([]);
	let towelRecords: TowelRecord[] | undefined = $derived.by(() => {
		if (!towelDB) return [];

		return towelDB.map((record, i, allRecords) => {
			const nextRecord = allRecords[i + 1];
			const gap = nextRecord ? dayjs(record.time).diff(nextRecord.time, 'day', true) : 0;
			return { ...record, gap };
		});
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

	let towelDirty: number | undefined = $derived.by(() => {
		let towelDirty;
		if (towelRecords && towelRecords.length > 0) {
			const lastWashDate = dayjs(towelRecords[0].time);
			const today = dayjs();
			/**
			 * https://stackoverflow.com/questions/36560806/the-left-hand-side-of-an-arithmetic-operation-must-be-of-type-any-number-or
			 */
			towelDirty = today.diff(lastWashDate, 'hours', true);
			return towelDirty;
		}
	});
	let spinner = $state(false);

	let lineChartEl = $state() as HTMLCanvasElement;
	let lineChart: Chart | undefined = $state();

	onMount(async () => {
		// if (!pb.authStore.isValid) {
		// 	goto('/login');
		// }

		if (pb.authStore.isValid) {
			towelDB = await pb.collection('towel').getFullList({
				sort: '-time'
			});

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
		return getValue(towelDirty);
	});

	let modal = $state() as HTMLDialogElement;

	function getValue(val: number): string {
		const day = 24;

		if (val > 0 && val <= 2 * day) return 'green';

		if (val > 2 * day && val <= 4 * day) return 'yellow';

		if (val > 4 * day && val <= 6 * day) return 'orange';

		if (val > 6 * day && val <= 999 * day) return 'red';

		return '';
	}

	async function addHandler() {
		const result = await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			spinner = false;
		}

		towelDB = await pb.collection('towel').getFullList({
			sort: '-time'
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

		if (gapDays > 0 && gapDays <= 2) return 'lime-400';

		if (gapDays > 2 && gapDays <= 4) return 'yellow-500';

		if (gapDays > 4 && gapDays <= 6) return 'orange-400';

		if (gapDays > 6 && gapDays <= 999) return 'red-700';

		return '';
	}
</script>

<svelte:head>
	<title>Towelie</title>
</svelte:head>

<PageWrapper title="Towelie" {pb}>
	<div
		class="grid max-w-[1200px] content-center justify-items-center gap-8 justify-self-center p-2 lg:grid-cols-2"
	>
		<div class="grid content-start justify-items-center gap-4">
			{#key towelRecords}
				{#if status === 'green'}
					<enhanced:img src={Green} alt="Clean" class="rounded-3xl" />
				{:else if status === 'yellow'}
					<enhanced:img src={Yellow} alt="Still Good" class="rounded-3xl" />
				{:else if status === 'orange'}
					<enhanced:img src={Orange} alt="Ripening" class="rounded-3xl" />
				{:else}
					<enhanced:img src={Red} alt="Time to Wash!" class="rounded-3xl" />
				{/if}
			{/key}
			<form class="w-full">
				<button class="btn btn-xl btn-primary flex w-full items-center gap-2" onclick={addHandler}
					>Just Washed My Towel!</button
				>
			</form>
		</div>
		<div class="grid w-full grid-rows-[auto_1fr_auto] content-start justify-items-center gap-8">
			<div class="w-full">
				<div
					class="text-md border-base-content/5 bg-base-200/50 grid grid-cols-2 content-center gap-4 rounded-lg border shadow"
				>
					<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
						<div>Status</div>
						{#key towelRecords}
							<button
								class="flex min-h-20 cursor-pointer items-center gap-4 py-8 text-2xl font-bold"
								onclick={() => modal.showModal()}
								>{#if status === 'green'}
									<div class="hidden h-6 w-6 rounded-full bg-lime-400 lg:flex"></div>
									<span class="text-lime-400">Squeaky Clean</span>
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
							</button>
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
									Nil
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
	</div>
</PageWrapper>

<!-- <dialog bind:this={modal} class="modal">
	<div class="modal-box grid gap-8">
		<h3 class="text-4xl font-bold">More Details</h3>
		<div class="grid min-h-16 w-full grid-cols-4">
			<div
				class="flex items-center justify-center rounded-l-lg bg-lime-400 {status === 'green'
					? 'border-base-content border-4 font-bold'
					: undefined}"
			>
				Clean
			</div>
			<div
				class="flex items-center justify-center bg-yellow-500 {status === 'yellow'
					? 'border-base-content border-4 font-bold'
					: undefined}"
			>
				Still Fresh
			</div>
			<div
				class="flex items-center justify-center bg-orange-400 {status === 'orange'
					? 'border-base-content border-4 font-bold'
					: undefined}"
			>
				Kinda Funky
			</div>
			<div
				class="flex items-center justify-center rounded-r-lg bg-red-700 {status === 'yikes'
					? 'border-base-content border-4 font-bold'
					: undefined}"
			>
				Wash Me!
			</div>
		</div>

		<div class="w-full text-sm">
			<h2 class="mb-2 text-lg font-semibold">Past 5 Washes</h2>
			{#if results && results.length > 0}
				{#each results as towel, i}
					{#if i <= 4}
						{@const formatted = dayjs(results[i].time).format('D MMM, h:mma')}
						{@const semantic = stripNegative(
							i === results.length - 1
								? 0
								: dayjs(results[i].time).diff(results[i + 1].time, 'day', true)
						)}
						<li class="border-b-base-300 ms-1 border-b py-2">
							{formatted}&nbsp;&nbsp;<span class="border-neutral/30 rounded-full border px-2"
								>+{semantic}</span
							>
						</li>
					{/if}
				{/each}
			{:else}
				No washes!
			{/if}
		</div>
		<div class="modal-action">
			<form method="dialog">
				
				<button class="btn">Close</button>
			</form>
		</div>
	</div>
</dialog> -->
