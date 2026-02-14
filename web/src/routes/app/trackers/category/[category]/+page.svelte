<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { allEntriesQueryOptions, allTrackersQueryOptions, userQueryOptions } from '$lib/queries';
	import { getColoredTrackers, getTrackerIcon, generateSubscriptionEntries } from '$lib/mapper.js';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import { getTrackerStatus } from '$lib/notification.js';
	import { calculateStreak } from '$lib/streaks';
	import { router } from '$lib/routes.js';

	let { data } = $props();

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		const statuses = <Record<string, ButtonState>>{};

		for (const t of trackersDb.data) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	const trackersDb = createQuery(allTrackersQueryOptions);
	const allEntriesDB = createQuery(allEntriesQueryOptions);
	const userOptions = createQuery(userQueryOptions);

	let allTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return [];
		return getColoredTrackers(trackersDb.data);
	});

	let categoryTrackers = $derived.by(() => {
		return allTrackers.filter((tracker) => tracker.category === data.category);
	});

	let entries = $derived.by(() => {
		if (!allEntriesDB.isSuccess || !allEntriesDB.data) return { tasks: [], subscriptions: [] };

		return {
			tasks: classifyTrackers(categoryTrackers, allEntriesDB.data, 'task'),
			subscriptions: classifyTrackers(categoryTrackers, allEntriesDB.data, 'subscription')
		};
	});

	function classifyTrackers(
		trackers: TrackerDB[],
		entries: EntryDB[],
		kind: 'task' | 'subscription'
	) {
		const data = [];

		for (const t of trackers) {
			if (t.kind !== kind) continue;

			// Filter only entries for this specific tracker using trackerId
			const entryData = entries.filter((entry) => t.id === entry.trackerId);

			const mergedData = {
				trackerName: t.name,
				trackerData: t,
				entries:
					kind === 'subscription'
						? generateSubscriptionEntries(t, userOptions.data!.id)
						: entryData,
				notification: getTrackerStatus(entryData),
				streak: calculateStreak(entryData, t)
			};

			data.push(mergedData);
		}

		return data;
	}

	let latestEntries: EntryDB[] = $derived.by(() => {
		if (
			!allEntriesDB.isSuccess ||
			!allEntriesDB.data ||
			categoryTrackers.length === 0 ||
			!allTrackers.length
		)
			return [];

		const trackerMap = new Map(allTrackers.map((t) => [t.id, t]));
		const categoryTrackerIds = new Set(categoryTrackers.map((t) => t.id));

		// Use flatMap to transform logs AND filter out any without matching trackers in one step.
		return allEntriesDB.data
			.filter((log) => categoryTrackerIds.has(log.trackerId))
			.slice(0, 5)
			.flatMap((log) => {
				const tracker = trackerMap.get(log.trackerId);
				if (!tracker) return [];

				return [{ ...log, expand: { tracker: { ...tracker } } }];
			});
	});
</script>

<PageWrapper title={data.category.charAt(0).toUpperCase() + data.category.slice(1)}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-4">
				{#if trackersDb.isPending}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{:else if entries.tasks && entries.tasks.length > 0}
					{#each entries.tasks as task (task.trackerData.id)}
						<ActionCard
							options={{
								size: 'compact',
								tracker: task.trackerData,
								title: task.trackerData.display,
								route: router.tracker(task.trackerData.id),
								entries: task.entries,
								icon: getTrackerIcon(task.trackerData.icon),
								button: {
									text: task.trackerData.actionLabel,
									status: buttonStatuses?.[task.trackerData.name]
								},
								streak: task.streak
							}}
						></ActionCard>
					{/each}
				{:else}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">Nothing being tracked!</p>
					</div>
				{/if}
			</section>

			{#if entries.subscriptions && entries.subscriptions.length > 0}
				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Subscriptions</h2>

					{#each entries.subscriptions as subscription (subscription.trackerData.id)}
						<ActionCard
							options={{
								size: 'compact',
								tracker: subscription.trackerData,
								title: subscription.trackerData.display,
								route: router.tracker(subscription.trackerData.id),
								entries: subscription.entries,
								icon: getTrackerIcon(subscription.trackerData.icon),
								button: {
									text: subscription.trackerData.actionLabel,
									status: buttonStatuses?.[subscription.trackerData.name]
								}
							}}
						></ActionCard>
					{/each}
				</section>
			{/if}

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Recent Activity</h2>

				{#if allEntriesDB.isPending}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
				{#if latestEntries && latestEntries.length > 0}
					<div class="border-base-300/50 rounded-2xl border bg-white/70">
						{#each latestEntries as entry (entry.id)}
							{@const fromNow = dayjs(entry.performedAt).fromNow()}
							<div class={['border-b-base-300/50 grid gap-4 border-b px-2 py-1']}>
								<div class="flex items-center p-2">
									<div class="flex grow items-center gap-4">
										<div class="flex items-center gap-2 align-baseline">
											{entry.expand?.tracker?.display}
										</div>
									</div>

									<div class="text-base-content/70 flex h-full items-center">{fromNow}</div>
								</div>
							</div>
						{/each}
					</div>
				{:else if !latestEntries || latestEntries.length === 0}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No tasks!</p>
					</div>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>
