<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { createQuery } from '@tanstack/svelte-query';
	import {
		allLogsQueryOptions,
		allTrackersQueryOptions,
		notificationQueryOptions
	} from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import FluentEmojiFlatAirplane from '$lib/assets/expressive-icons/FluentEmojiFlatAirplane.svelte';
	import { getColoredTrackers, getFamilyColor, getTrackerIcon } from '$lib/mapper';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const latestLogs = createQuery(notificationQueryOptions);
	const trackersDb = createQuery(allTrackersQueryOptions);
	const allLogsDb = createQuery(allLogsQueryOptions);

	let buttonStatuses = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		const statuses = <Record<string, ButtonState>>{};

		for (const t of trackersDb.data) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	let trackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return { pinned: [], general: [] };

		const coloredTrackers = getColoredTrackers(trackersDb.data);

		const pinned = coloredTrackers.filter((tracker) => tracker.pinned && tracker.show);
		const general = coloredTrackers.filter((tracker) => !tracker.pinned && tracker.show);

		return { pinned: pinned, general: general };
	});

	let logs = $derived.by(() => {
		if (!allLogsDb.isSuccess || !allLogsDb.data) return { pinned: [], general: [] };

		const pinned = [];
		const general = [];

		for (const t of trackers.pinned) {
			const logData = allLogsDb.data.filter((log) => t.id === log.tracker);
			const trackerData = trackers.pinned.find((tracker) => tracker.id === t.id);
			if (!trackerData) continue;

			const mergedData = {
				trackerName: t.name,
				trackerData: trackerData,
				logData: logData,
				notification: getTrackerStatus(logData)
			};
			pinned.push(mergedData);
		}

		for (const t of trackers.general) {
			const logData = allLogsDb.data.filter((log) => t.id === log.tracker);
			const trackerData = trackers.general.find((tracker) => tracker.id === t.id);
			if (!trackerData) continue;

			const mergedData = {
				trackerName: t.name,
				trackerData: trackerData,
				logData: logData,
				notification: getTrackerStatus(logData)
			};
			general.push(mergedData);
		}

		return { pinned: pinned, general: general };
	});
</script>

<PageWrapper title="Cubby" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Pinned</h2>

				{#if allLogsDb.isPending}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{:else if logs.pinned && logs.pinned.length > 0}
					{#each logs.pinned as log (log.trackerData?.id)}
						<ActionCard
							options={{
								tracker: log.trackerData,
								size: 'compact',
								title: log.trackerData?.display,
								route: `/${log.trackerData?.category}/${log.trackerData?.id}`,
								icon: getTrackerIcon(log.trackerName),
								button: {
									status: buttonStatuses?.[log.trackerName],
									text: log.trackerData?.actionLabel
								}
							}}
						></ActionCard>
					{/each}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No pinned tasks!</p>
					</div>
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Other Tasks</h2>

				<div class="border-base-300/50 rounded-2xl border bg-white/70">
					{#if allLogsDb.isSuccess && allLogsDb.data && logs.general && logs.general.length > 0}
						{#each logs.general as log, i (log.trackerData?.id)}
							<ActionCard
								options={{
									tracker: log.trackerData,
									size: 'list',
									title: log.trackerData?.display,
									route: `/${log.trackerData?.category}/${log.trackerData?.id}`,
									icon: getTrackerIcon(log.trackerName),
									lastChild: i === logs.general.length - 1 ? true : undefined,
									button: {
										status: buttonStatuses?.[log.trackerName],
										text: log.trackerData?.actionLabel
									}
								}}
							></ActionCard>
						{/each}
					{:else if allLogsDb.isSuccess && allLogsDb.data && logs.general && logs.general.length === 0}
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
