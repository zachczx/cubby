<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { goto } from '$app/navigation';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { allLogsQueryOptions, allTrackersQueryOptions } from '$lib/queries';
	import { getColoredTrackers, getFamilyColor, getTrackerIcon } from '$lib/mapper.js';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import { getTrackerStatus } from '$lib/notification.js';

	let { data } = $props();

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $derived.by(() => {
		if (!currentTrackers) return;

		const statuses = <Record<string, ButtonState>>{};

		for (const t of currentTrackers) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	const trackersDb = createQuery(allTrackersQueryOptions);
	const allLogsDb = createQuery(allLogsQueryOptions);

	let currentTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		const categoryTrackers = trackersDb.data.filter(
			(tracker) => tracker.category === data.category
		);

		return getColoredTrackers(categoryTrackers);
	});

	let logsByTracker = $derived.by(() => {
		const map = new Map();
		if (!allLogsDb.isSuccess || !allLogsDb.data) return map;

		for (const log of allLogsDb.data) {
			if (!map.has(log.tracker)) {
				map.set(log.tracker, []);
			}

			map.get(log.tracker).push(log);
		}

		return map;
	});

	let latestLogs: LogsDB[] = $derived.by(() => {
		if (!allLogsDb.isSuccess || !allLogsDb.data || !currentTrackers) return [];

		const trackerMap = new Map(currentTrackers.map((t) => [t.id, t]));

		// Use flatMap to transform logs AND filter out any without matching trackers in one step.
		// Returns [] to skip logs where tracker is undefined (flatmap flattens empty [], so nothing added), so I dont need a separate
		// .filter() for null values + verbose type guards. TypeScript automatically infers the correct type.
		return allLogsDb.data
			.filter((log) => trackerMap.has(log.tracker))
			.slice(0, 5)
			.flatMap((log) => {
				const tracker = currentTrackers?.find((tracker) => tracker.id === log.tracker);
				if (!tracker) return [];

				return [{ ...log, expand: { tracker: { ...tracker } } }];
			});
	});

	let logs = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !allLogsDb.isSuccess || !allLogsDb.data)
			return;

		return currentTrackers?.map((tracker) => {
			const associatedLogs = logsByTracker.get(tracker.id) || [];

			return {
				trackerName: tracker.name,
				trackerData: tracker,
				logData: associatedLogs,
				notification: getTrackerStatus(associatedLogs)
			};
		});
	});

	let tasks = $derived.by(() => {
		if (!logs) return;
		return logs.filter((log) => log.trackerData.kind === 'task');
	});

	let subscriptions = $derived.by(() => {
		if (!logs) return;
		return logs.filter((log) => log.trackerData.kind === 'subscription');
	});
</script>

<PageWrapper title={data.category.charAt(0).toUpperCase() + data.category.slice(1)} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-4">
				{#if trackersDb.isPending}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{:else if tasks && tasks.length > 0}
					{#each tasks as task (task.trackerData.id)}
						{#if task.trackerData.kind === 'task'}
							<ActionCard
								options={{
									size: 'compact',
									tracker: task.trackerData,
									title: task.trackerData.display,
									route: `/app/${task.trackerData.category}/${task.trackerData.id}`,
									logs: task.logData,
									icon: getTrackerIcon(task.trackerData.icon),
									button: {
										text: task.trackerData.actionLabel,
										status: buttonStatuses?.[task.trackerData.name]
									}
								}}
							></ActionCard>
						{/if}
					{/each}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">Nothing being tracked!</p>
					</div>
				{/if}
			</section>

			{#if subscriptions && subscriptions.length > 0}
				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Subscriptions</h2>

					{#each subscriptions as subscription (subscription.trackerData.id)}
						{#if subscription.trackerData.kind === 'subscription'}
							<ActionCard
								options={{
									size: 'compact',
									tracker: subscription.trackerData,
									title: subscription.trackerData.display,
									route: `/app/${subscription.trackerData.category}/${subscription.trackerData.id}`,
									logs: subscription.logData,
									icon: getTrackerIcon(subscription.trackerData.icon),
									button: {
										text: subscription.trackerData.actionLabel,
										status: buttonStatuses?.[subscription.trackerData.name]
									}
								}}
							></ActionCard>
						{/if}
					{/each}
				</section>
			{/if}

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Recent Activity</h2>

				<div class="border-base-300/50 rounded-2xl border bg-white/70">
					{#if latestLogs && latestLogs.length > 0}
						{#each latestLogs as log, i (log.id)}
							{@const fromNow = dayjs(log.time).fromNow()}
							<div class={['border-b-base-300/50 grid gap-4 border-b px-2 py-1']}>
								<div class="flex items-center p-2">
									<div class="flex grow items-center gap-4">
										<div class="flex items-center gap-2 align-baseline">
											{log.expand?.tracker?.display}
										</div>
									</div>

									<div class="text-base-content/70 flex h-full items-center">{fromNow}</div>
								</div>
							</div>
						{/each}
					{:else if !latestLogs || latestLogs.length === 0}
						<div class="justify-self-center">
							<enhanced:img src={EmptyCorgi} alt="nothing" />
							<p class="text-center">No tasks!</p>
						</div>
					{:else}
						<SkeletonActionCard size="compact" />
						<SkeletonActionCard size="compact" />
						<SkeletonActionCard size="compact" />
					{/if}
				</div>
			</section>
		</div>
	</main>
</PageWrapper>
