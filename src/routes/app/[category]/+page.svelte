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
</script>

<PageWrapper title={data.category.charAt(0).toUpperCase() + data.category.slice(1)} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-4 py-4">
				{#if trackersDb.isPending}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{:else if logs && logs.length > 0}
					{#each logs as log}
						<ActionCard
							options={{
								size: 'compact',
								tracker: log.trackerData,
								title: log.trackerData.display,
								route: `/app/${log.trackerData.category}/${log.trackerData.id}`,
								logs: log.logData,
								icon: getTrackerIcon(log.trackerData.name),
								button: {
									text: log.trackerData.actionLabel,
									status: buttonStatuses?.[log.trackerData.name]
								}
							}}
						></ActionCard>
					{/each}

					{#if data.category === 'household'}
						<div class="border-base-300 grid min-h-24 gap-4 rounded-3xl border bg-white/70 p-4">
							<a href="/app/household/count" class="flex items-center">
								<div class="flex grow items-center gap-4">
									<FluentEmojiFlatStopwatch class="size-12 opacity-75" />

									<div>
										<p class="text-xl font-bold">Timer</p>
									</div>
								</div>
								<div class="flex h-full items-center">
									<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
										><MaterialSymbolsChevronRight class="size-6" /></button
									>
								</div>
							</a>
							<button
								class="btn btn-lg btn-primary flex w-full items-center gap-2 rounded-full"
								onclick={() => goto('/app/household/count?start=true')}
							>
								Start Timer
							</button>
						</div>
					{/if}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">Nothing being tracked!</p>
					</div>
				{/if}
			</div>
		</div>
	</main>
</PageWrapper>
