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
	import {
		trackerQueryOptions,
		vacationQueryOptions,
		createLogsQuery,
		allLogsQueryOptions
	} from '$lib/queries';
	import { getCalendarEntries } from '$lib/calendar';
	import CustomDateModal from '$lib/ui/CustomDateModal.svelte';
	import StatusDescriptions from '$lib/ui/StatusDescriptions.svelte';
	import TwoColumnCard from '$lib/ui/TwoColumnCard.svelte';
	import StatusHeroImage from '$lib/ui/StatusHeroImage.svelte';
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import SingleDayModal from '$lib/ui/SingleDayModal.svelte';
	import { getTrackerStatus } from '$lib/notification';

	dayjs.extend(calendar);
	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let { options }: { options: TrackerPageOptions } = $props();

	const allLogsDb = createQuery(allLogsQueryOptions);

	let currentTrackerLogs = $derived.by(() => {
		if (!allLogsDb.isSuccess || !allLogsDb.data || !options.tracker) return [];

		return allLogsDb.data.filter((log) => log.tracker === options.tracker?.id);
	});

	const tracker = createQuery(() => trackerQueryOptions(options.tracker?.id));
	let interval = $derived.by(() => (tracker.isSuccess ? tracker.data?.interval : undefined));
	let intervalUnit = $derived.by(() =>
		tracker.isSuccess ? tracker.data?.intervalUnit : undefined
	);
	const query = () =>
		createLogsQuery({
			trackerId: tracker.data?.id ?? '',
			interval: interval,
			intervalUnit: intervalUnit
		});

	let notification = $derived.by(() => getTrackerStatus(currentTrackerLogs));

	type TabPages = 'overview' | 'stats' | 'history';

	let currentTab = $state<TabPages>('history');

	let records: LogsRecord[] | undefined = $state([]);

	$effect(() => {
		if (currentTrackerLogs && currentTrackerLogs.length > 0) {
			records = currentTrackerLogs.map((record, i, allRecords) => {
				const nextRecord = allRecords[i + 1];
				const gap = nextRecord ? dayjs(record.time).diff(nextRecord.time, 'day', true) : 0;
				return { ...record, gap };
			});
		}
	});

	let logsByYears = $derived.by(() => {
		const years = new Map<number, LogsDB[]>();
		for (const t of currentTrackerLogs) {
			const y = dayjs(t.time).year();

			if (!years.has(y)) {
				years.set(y, []);
			}
			years.get(y)?.push(t);
		}

		return years;
	});

	let signUpFirstRecord: LogsDB | null = $derived.by(() => {
		if (!currentTrackerLogs || currentTrackerLogs.length === 0) return null;

		return currentTrackerLogs[currentTrackerLogs.length - 1];
	});
</script>

<PageWrapper title={options.labels.pageTitle} {pb}>
	<main class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			{#if currentTrackerLogs}
				<StatusHeroImage {notification} kind="subscription" />
			{:else}
				<div class="avatar relative mt-2 mb-4">
					<div class="skeleton aspect-square w-40 rounded-full shadow-md"></div>
				</div>
			{/if}

			<!-- <ActionButton {query} text={options.labels.ctaButtonText} /> -->

			<!-- <div class="flex justify-start">
				<CustomDateModal tracker={options.tracker} {interval} {intervalUnit} />
			</div> -->
		</div>

		<div class="grid w-full content-start gap-8 pt-4 pb-8">
			<ul class="tracker-tabs grid w-full grid-cols-2 content-center justify-items-center">
				<li class="w-full" aria-current={currentTab === 'overview' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'overview')}
					>
						Overview</button
					>
				</li>

				<!-- <li class="w-full" aria-current={currentTab === 'stats' ? 'page' : undefined}>
					<button class="w-full cursor-pointer text-center" onclick={() => (currentTab = 'stats')}>
						Stats</button
					>
				</li> -->

				<li class="w-full" aria-current={currentTab === 'history' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'history')}
					>
						History</button
					>
				</li>
			</ul>

			{#if currentTrackerLogs && currentTrackerLogs.length === 0}
				<div class="mx-4 mt-4 text-center">
					<p class="font-bold">Ready to track?</p>
					<p>Start by adding a log.</p>
				</div>
			{:else}
				<div class={['grid w-full gap-8 px-4', currentTab === 'overview' ? undefined : 'hidden']}>
					<div
						class="border-base-content/5 bg-primary/10 grid justify-items-center rounded-xl border p-4 shadow"
					>
						<h2 class="text-md text-center">Next Charge</h2>
						<div class="grid min-h-20 content-center justify-items-center">
							{#if notification && notification.level}
								<!-- {@const semantic = dayjs(notification.next).calendar(dayjs(), dayjsCalendarOptions)} -->
								{@const semantic = dayjs(notification.next).fromNow()}
								{@const formatted = dayjs(notification.next).format('D MMM YYYY')}
								<p class="text-primary text-2xl font-bold">
									{semantic}
								</p>
								<p>{formatted}</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						</div>
					</div>

					<TwoColumnCard leftTitle="Frequency" rightTitle="Last">
						{#snippet left()}
							{#if tracker.isSuccess && tracker.data}
								{@const plural = tracker.data.interval > 1 ? true : false}
								<p>
									{tracker.data.interval}&nbsp;{plural
										? tracker.data.intervalUnit + 's'
										: tracker.data.intervalUnit}
								</p>
								<p class="text-base-content/70 text-base font-normal">once</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-xl font-bold">Not set yet</div>
							{/if}
						{/snippet}

						{#snippet right()}
							{#if currentTrackerLogs && currentTrackerLogs.length > 0}
								{#if currentTrackerLogs.length > 0}
									{@const formatted = dayjs(currentTrackerLogs[0].time).fromNow(true)}
									<p>
										{#if formatted === 'a few seconds'}
											seconds
										{:else}
											{formatted}
										{/if}
									</p>
									<p class="text-base-content/70 text-base font-normal">ago</p>
								{:else}
									<div class="flex min-h-20 items-center gap-4 text-xl font-bold">
										Never recorded
									</div>
								{/if}
							{/if}
							{#if allLogsDb.isPending}
								<div class="custom-loader"></div>
							{/if}
						{/snippet}
					</TwoColumnCard>
				</div>

				<div class={['grid w-full gap-8 px-4', currentTab === 'stats' ? undefined : 'hidden']}>
					<div class="border-base-content/5 w-full rounded-lg border p-4 shadow">
						<h2 class="text-md text-center">Trend</h2>
						<div></div>
					</div>

					<div
						class="border-base-content/5 grid w-full grid-cols-2 content-center gap-4 rounded-lg border shadow"
					>
						<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
							<h2 class="text-md">Longest Gap</h2>
							<div class="text-center text-2xl font-bold"></div>
						</div>

						<div class="grid justify-items-center p-4">
							<h2 class="text-md">Average Gap</h2>
							<div class="text-center text-2xl font-bold"></div>
						</div>
					</div>
				</div>

				<div
					class={[
						'grid w-full grid-cols-[minmax(0,1fr)] gap-4 overflow-hidden px-4',
						currentTab === 'history' ? undefined : 'hidden'
					]}
				>
					{#each logsByYears as [key, year]}
						<div class="grid gap-2">
							<h3 class="text-base-content/50 text-sm font-bold">{key}</h3>

							{#each year as log}
								{@const formatted = dayjs(log.time).format('D MMM')}
								<div
									class="border-base-300 grid min-h-18 grid-cols-3 items-baseline gap-4 rounded-2xl border bg-white/70 px-2 py-2"
								>
									<div class="flex p-2 text-lg font-bold">
										{formatted}
									</div>
									<div class="flex justify-center">
										{#if signUpFirstRecord?.id === log.id}
											<span class="text-primary">Sign-Up</span>
										{:else}
											<span class="text-info">Renewal</span>{/if}
									</div>
									<div class="text-base-content/60 flex justify-end p-2 font-semibold">
										${log.expand?.tracker?.cost}
									</div>
								</div>
							{/each}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</main>
</PageWrapper>
