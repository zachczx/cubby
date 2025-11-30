<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient, type CreateQueryResult } from '@tanstack/svelte-query';
	import {
		allLogsQueryOptions,
		allTrackersQueryOptions,
		createUserQueryOptions,
		logsQueryOptions,
		notificationQueryOptions,
		trackerNameToId
	} from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import FluentEmojiFlatAirplane from '$lib/assets/expressive-icons/FluentEmojiFlatAirplane.svelte';
	import { bedsheet, bath, chewable, gummy, spray, towel } from '$lib/ui/ActionCardDefaults';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $state<Record<string, ButtonState>>({
		towel: 'default',
		spray: 'default',
		gummy: 'default',
		bedsheet: 'default',
		bath: 'default',
		chewable: 'default'
	});

	const latestLogs = createQuery(notificationQueryOptions);
	const trackersDb = createQuery(allTrackersQueryOptions);
	const allLogsDb = createQuery(allLogsQueryOptions);

	const options: ActionCardOptions[] = [
		{ ...towel, size: 'compact', button: { ...towel.button, status: buttonStatuses.towel } },
		{ ...spray, size: 'compact', button: { ...spray.button, status: buttonStatuses.spray } },
		{ ...gummy, size: 'compact', button: { ...gummy.button, status: buttonStatuses.gummy } },
		{
			...bedsheet,
			size: 'compact',
			button: { ...bedsheet.button, status: buttonStatuses.bedsheet }
		},
		{
			...bath,
			size: 'compact',
			button: { ...bath.button, status: buttonStatuses.bath }
		},
		{
			...chewable,
			size: 'compact',
			button: { ...chewable.button, status: buttonStatuses.chewable }
		}
	];

	// let tasks = $derived.by(() => {
	// 	if (!latestLogs.isSuccess || !latestLogs.data) return { important: [], general: [] };

	// 	const importantTasks = trackersDb
	// 		.filter((tracker) => tracker.pinned)
	// 		.map((task) => {
	// 			const data = latestLogs.data.find(
	// 				(log) => log.tracker === trackerNameToId(task.collectionName, trackers.data)
	// 			);

	// 			return { ...task, notification: getTrackerStatus(data) };
	// 		});

	// 	const generalTasks = options
	// 		.filter((task) => !important.includes(task.collectionName))
	// 		.map((task) => {
	// 			const data = latestLogs.data.find(
	// 				(log) => log.tracker === trackerNameToId(task.collectionName, trackers.data)
	// 			);

	// 			return { ...task, notification: getTrackerStatus(data) };
	// 		});

	// 	return { important: importantTasks, general: generalTasks };
	// });

	//////////////////////////////////////
	//////////////////////////////////////
	//////////////////////////////////////
	//////////////////////////////////////
	//////////////////////////////////////

	let trackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return { pinned: [], general: [] };

		const pinned = trackersDb.data.filter((tracker) => tracker.pinned);
		const general = trackersDb.data.filter((tracker) => !tracker.pinned);

		return { pinned: pinned, general: general };
	});

	let logs = $derived.by(() => {
		if (!allLogsDb.isSuccess || !allLogsDb.data) return { pinned: [], general: [] };

		const pinned = [];
		const general = [];

		for (const t of trackers.pinned) {
			const logData = allLogsDb.data.filter((log) => t.id === log.tracker);
			const trackerData = trackers.pinned.find((tracker) => tracker.id === t.id);
			const config = options.find((opt) => opt.collectionName === trackerData?.name);
			const mergedData = {
				trackerName: t.name,
				trackerData: trackerData,
				logData: logData,
				config: config,
				notification: getTrackerStatus(logData)
			};
			pinned.push(mergedData);
		}

		for (const t of trackers.general) {
			const logData = allLogsDb.data.filter((log) => t.id === log.tracker);
			const trackerData = trackers.general.find((tracker) => tracker.id === t.id);
			const config = options.find((opt) => opt.collectionName === trackerData?.name);
			const mergedData = {
				trackerName: t.name,
				trackerData: trackerData,
				logData: logData,
				config: config,
				notification: getTrackerStatus(logData)
			};
			general.push(mergedData);
		}

		return { pinned: pinned, general: general };
	});
</script>

<PageWrapper title="Sundry" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Pinned</h2>

				{#if logs.pinned && logs.pinned.length > 0}
					{#each logs.pinned as log (log.trackerName)}
						{#if log.trackerData?.show}
							<ActionCard
								options={{
									collectionName: log.trackerName,
									size: 'compact',
									title: log.trackerData?.display ?? '',
									route: `/${log.trackerData.category}/${log.trackerData.id}`,
									icon: log.config?.icon ?? '',
									button: {
										status: buttonStatuses[log.trackerName],
										text: log.trackerData?.actionLabel ?? ''
									}
								}}
							></ActionCard>
						{/if}
					{/each}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No pinned tasks!</p>
					</div>
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">General</h2>

				{#if logs.general && logs.general.length > 0}
					{#each logs.general as log (log.trackerName)}
						{#if log.trackerData?.show}
							<ActionCard
								options={{
									collectionName: log.trackerName,
									size: 'compact',
									title: log.trackerData?.display ?? '',
									route: `/${log.trackerData.category}/${log.trackerData.id}`,
									icon: log.config?.icon ?? '',
									button: {
										status: buttonStatuses[log.trackerName],
										text: log.trackerData?.actionLabel ?? ''
									}
								}}
							></ActionCard>
						{/if}
					{/each}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No tasks!</p>
					</div>
				{/if}
			</section>

			<section class="grid gap-0 py-0">
				<h2 class="text-base-content/70 text-lg font-bold">Quick Links</h2>
				<div class="flex items-center gap-2">
					<a
						href="/household/count"
						class="active:bg-neutral/10 focus-within:bg-neutral/10 focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-24 content-center justify-items-center gap-1 rounded-2xl p-2 text-sm font-semibold"
					>
						<FluentEmojiFlatStopwatch class="size-8" />
						Stopwatch</a
					>
					<a
						href="/profile/vacation"
						class="active:bg-neutral/10 focus-within:bg-neutral/10 focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-24 content-center justify-items-center gap-1 rounded-2xl p-2 text-sm font-semibold"
					>
						<FluentEmojiFlatAirplane class="size-8" />
						Vacation</a
					>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>
