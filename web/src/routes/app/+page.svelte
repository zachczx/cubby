<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { createQuery } from '@tanstack/svelte-query';
	import { allEntriesQueryOptions, allTrackersQueryOptions, userQueryOptions } from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import FluentEmojiFlatAirplane from '$lib/assets/expressive-icons/FluentEmojiFlatAirplane.svelte';
	import { getColoredTrackers, getTrackerIcon, generateSubscriptionEntries } from '$lib/mapper';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import { calculateStreak } from '$lib/streaks';
	import { onMount, type Component } from 'svelte';
	import { router } from '$lib/routes';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const trackersDb = createQuery(allTrackersQueryOptions);
	const allEntriesDb = createQuery(allEntriesQueryOptions);
	const userOptions = createQuery(userQueryOptions);

	let generalTasksUpcomingDays = $derived.by(() => {
		if (!userOptions.isSuccess || !userOptions.data) return 14;

		return userOptions.data.taskLookaheadDays;
	});

	let buttonStatuses = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		const statuses = <Record<string, ButtonState>>{};

		for (const t of trackersDb.data) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	let trackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return { pinned: [], general: [] };

		const coloredTrackers = getColoredTrackers(trackersDb.data);

		const pinned = coloredTrackers.filter(
			(tracker) => tracker.pinned && tracker.show && tracker.kind === 'task'
		);
		const general = coloredTrackers.filter(
			(tracker) => !tracker.pinned && tracker.show && tracker.kind === 'task'
		);

		return { pinned: pinned, general: general };
	});

	let subscriptions = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return;

		const coloredTrackers = getColoredTrackers(trackersDb.data);

		return coloredTrackers
			.filter((tracker) => tracker.show && tracker.kind === 'subscription')
			.map((sub) => {
				return {
					...sub,
					entryData: generateSubscriptionEntries(sub, userOptions.data!.id)
				};
			});
	});

	let entries = $derived.by(() => {
		if (!allEntriesDb.isSuccess) return { pinned: [], general: [] };

		return {
			pinned: classifyTrackers(trackers.pinned, allEntriesDb.data, 'pinned'),
			general: classifyTrackers(trackers.general, allEntriesDb.data, 'general')
		};
	});

	function classifyTrackers(trackers: TrackerDB[], entries: EntryDB[], kind: 'general' | 'pinned') {
		const data = [];

		for (const t of trackers) {
			const entryData = entries?.filter((entry) => t.id === entry.trackerId) ?? [];
			const trackerData = trackers.find((tracker) => tracker.id === t.id);
			if (!trackerData) continue;

			const mergedData = {
				trackerName: t.name,
				trackerData: trackerData,
				entries: entryData,
				notification: getTrackerStatus(entryData),
				streak: calculateStreak(entryData, trackerData)
			};

			if (
				kind === 'general' &&
				dayjs(mergedData.notification.next).diff(dayjs(), 'day', true) > generalTasksUpcomingDays
			) {
				continue;
			}
			data.push(mergedData);
		}

		return data;
	}

	const viewButtons = [
		{ days: 7, description: '1 week' },
		{ days: 14, description: '2 weeks' },
		{ days: 31, description: '1 month' },
		{ days: 183, description: '6 months' },
		{ days: 9999, description: 'All' }
	];

	async function generalTasksViewBtnHandler(numberDays: number) {
		generalTasksUpcomingDays = numberDays;
		console.log(numberDays, typeof numberDays);
		try {
			await api.patch('users/me/task-lookahead', {
				body: JSON.stringify({
					taskDays: numberDays
				})
			});
		} catch (err) {
			console.error(err);
			addToast('error', 'Error!');
		}
	}

	onMount(async () => {
		const response = await api.get('check');
		if (response.status !== 204) {
			goto(resolve('/login'));
		}
	});
</script>

<PageWrapper title="Cubby" back={false}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Pinned</h2>

				{#if allEntriesDb.isPending}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{:else if allEntriesDb.isSuccess}
					{#each entries.pinned as entry (entry.trackerData?.id)}
						<ActionCard
							options={{
								tracker: entry.trackerData,
								size: 'compact',
								title: entry.trackerData?.display,
								entries: entry.entries,
								route: router.tracker(entry.trackerData?.id),
								icon: getTrackerIcon(entry.trackerData?.icon),
								button: {
									status: buttonStatuses?.[entry.trackerName],
									text: entry.trackerData?.actionLabel
								},
								streak: entry.streak
							}}
						></ActionCard>
					{/each}
				{:else}
					Error!
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Other Tasks</h2>
				<div class="flex items-center gap-2">
					{#each viewButtons as btn}
						<button
							class={[
								'btn-soft btn btn-sm rounded-full',
								generalTasksUpcomingDays === btn.days && 'btn-primary'
							]}
							onclick={() => generalTasksViewBtnHandler(btn.days)}>{btn.description}</button
						>
					{/each}
				</div>

				{#if allEntriesDb.isPending}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{:else if allEntriesDb.isSuccess}
					<div class="border-base-300/50 rounded-2xl border bg-white/70">
						{#each entries.general as entry, i (entry.trackerData?.id)}
							<ActionCard
								options={{
									tracker: entry.trackerData,
									size: 'list',
									title: entry.trackerData?.display,
									entries: entry.entries,
									route: router.tracker(entry.trackerData?.id),
									icon: getTrackerIcon(entry.trackerData?.icon),
									lastChild: i === entries.general.length - 1 ? true : undefined,
									button: {
										status: buttonStatuses?.[entry.trackerName],
										text: entry.trackerData?.actionLabel
									},
									streak: entry.streak
								}}
							></ActionCard>
						{/each}
					</div>
					<!-- {:else if allEntriesDb.isSuccess && allEntriesDb.data && entries.general && entries.general.length === 0}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No tasks!</p>
					</div> -->
				{:else}
					Error!
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Subscriptions</h2>

				{#if allEntriesDb.isPending}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
				{#if allEntriesDb.isSuccess && subscriptions && subscriptions.length > 0}
					<div class="border-base-300/50 rounded-2xl border bg-white/70">
						{#each subscriptions as sub, i (sub.id)}
							<ActionCard
								options={{
									tracker: sub,
									size: 'list',
									title: sub.display,
									entries: sub.entryData,
									route: `/app/${sub.id}`,
									icon: getTrackerIcon(sub.icon),
									lastChild: i === entries.general.length - 1 ? true : undefined,
									button: {
										status: buttonStatuses?.[sub.name],
										text: sub.actionLabel
									}
								}}
							></ActionCard>
						{/each}
					</div>
				{:else if allEntriesDb.isSuccess && subscriptions && subscriptions.length === 0}
					<div class="justify-self-center">
						<enhanced:img src={EmptyCorgi} alt="nothing" />
						<p class="text-center">No subscriptions!</p>
					</div>
				{/if}
			</section>

			<section class="grid gap-0 py-0">
				<h2 class="text-base-content/70 text-lg font-bold">Quick Links</h2>
				<div class="flex items-center gap-2">
					{@render quickLink('Stopwatch', '/app/count', FluentEmojiFlatStopwatch, 'size-8')}
					{@render quickLink(
						'Vacation',
						'/app/profile/vacation',
						FluentEmojiFlatAirplane,
						'size-8'
					)}
				</div>
			</section>
		</div>
	</main>
</PageWrapper>

{#snippet quickLink(name: string, href: string, Icon: Component, size: string)}
	<a
		{href}
		class="active:bg-neutral/10 focus-within:bg-neutral/10 focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-24 content-center justify-items-center gap-1 rounded-2xl p-2 text-sm font-semibold"
	>
		<Icon class={size}></Icon>
		{name}</a
	>
{/snippet}
