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
	import { router } from '$lib/routes';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Icon from '@iconify/svelte';

	const categories = [
		{ id: 'all', label: 'All', icon: 'material-symbols:checklist' },
		{ id: 'personal', label: 'Personal', icon: 'material-symbols:person' },
		{ id: 'household', label: 'Household', icon: 'material-symbols:laundry' },
		{ id: 'pet', label: 'Pet', icon: 'material-symbols:pets' }
	] as const;

	let selectedCategory = $state<string>('all');

	const viewButtons = [
		{ days: 7, description: '1 week' },
		{ days: 14, description: '2 weeks' },
		{ days: 31, description: '1 month' },
		{ days: 183, description: '6 months' },
		{ days: 9999, description: 'All' }
	];

	let generalTasksUpcomingDays = $derived.by(() => {
		if (!userOptions.isSuccess || !userOptions.data) return 14;
		return userOptions.data.taskLookaheadDays;
	});

	async function generalTasksViewBtnHandler(numberDays: number) {
		generalTasksUpcomingDays = numberDays;
		try {
			await api.patch('users/me/task-lookahead', {
				body: JSON.stringify({ taskDays: numberDays })
			});
		} catch (err) {
			console.error(err);
			addToast('error', 'Error!');
		}
	}

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const trackersDb = createQuery(allTrackersQueryOptions);
	const allEntriesDb = createQuery(allEntriesQueryOptions);
	const userOptions = createQuery(userQueryOptions);

	let buttonStatuses = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		const statuses = <Record<string, ButtonState>>{};
		for (const t of trackersDb.data) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	let currentTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return { pinned: [], general: [] };

		const filtered =
			selectedCategory === 'all'
				? trackersDb.data
				: trackersDb.data.filter((tracker) => tracker.category === selectedCategory);

		const colored = getColoredTrackers(filtered);

		return {
			pinned: colored.filter((t) => t.pinned && t.show && t.kind === 'task'),
			general: colored.filter((t) => !t.pinned && t.show && t.kind === 'task')
		};
	});

	let subscriptions = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return [];

		const filtered =
			selectedCategory === 'all'
				? trackersDb.data
				: trackersDb.data.filter((tracker) => tracker.category === selectedCategory);

		const colored = getColoredTrackers(filtered);

		return colored
			.filter((t) => t.show && t.kind === 'subscription')
			.map((sub) => ({
				...sub,
				entryData: generateSubscriptionEntries(sub, userOptions.data!.id)
			}));
	});

	let entries = $derived.by(() => {
		if (!allEntriesDb.isSuccess || !allEntriesDb.data)
			return { pinned: [], general: [] };

		return {
			pinned: classifyTrackers(currentTrackers.pinned, allEntriesDb.data),
			general: classifyTrackers(
				currentTrackers.general,
				allEntriesDb.data,
				generalTasksUpcomingDays
			)
		};
	});

	function classifyTrackers(
		trackers: TrackerDB[],
		entries: EntryDB[],
		filterDays?: number
	) {
		const data = [];

		for (const t of trackers) {
			const entryData = entries.filter((entry) => t.id === entry.trackerId);
			const notification = getTrackerStatus(entryData);

			if (
				filterDays !== undefined &&
				notification.next &&
				dayjs(notification.next).diff(dayjs(), 'day', true) > filterDays
			) {
				continue;
			}

			data.push({
				trackerName: t.name,
				trackerData: t,
				entries: entryData,
				notification,
				streak: calculateStreak(entryData, t)
			});
		}

		return data;
	}
</script>

<PageWrapper title="Tasks">
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<nav class="border-b-base-300/50 flex border-b">
				{#each categories as cat (cat.id)}
					<button
						class="text-base-content/50 relative flex grow items-center justify-center gap-1.5 px-3 py-2.5 text-sm font-medium transition-colors {selectedCategory === cat.id ? 'text-primary !font-semibold' : 'hover:text-base-content/70'}"
						onclick={() => (selectedCategory = cat.id)}
					>
						<Icon icon={cat.icon} class="size-4" />
						{cat.label}
						{#if selectedCategory === cat.id}
							<span class="bg-primary absolute bottom-0 left-0 h-0.5 w-full rounded-full"></span>
						{/if}
					</button>
				{/each}
			</nav>

			<div class="flex items-center gap-2 px-1">
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

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Pinned</h2>

				{#if allEntriesDb.isSuccess && trackersDb.isSuccess && userOptions.isSuccess}
					{#if entries.pinned.length > 0}
						{#each entries.pinned as entry (entry.trackerData.id)}
							<ActionCard
								options={{
									tracker: entry.trackerData,
									size: 'compact',
									title: entry.trackerData.display,
									entries: entry.entries,
									route: router.tracker(entry.trackerData.id),
									icon: getTrackerIcon(entry.trackerData.icon),
									button: {
										status: buttonStatuses?.[entry.trackerName],
										text: entry.trackerData.actionLabel
									},
									streak: entry.streak
								}}
							></ActionCard>
						{/each}
					{:else}
						<p class="text-base-content/50 py-4 text-center">No pinned tasks</p>
					{/if}
				{:else if allEntriesDb.isError || trackersDb.isError || userOptions.isError}
					Error!
				{:else}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Other Tasks</h2>

				{#if allEntriesDb.isSuccess && trackersDb.isSuccess && userOptions.isSuccess}
					{#if entries.general.length > 0}
						<div class="border-base-300/50 rounded-2xl border bg-base-50">
							{#each entries.general as entry, i (entry.trackerData.id)}
								<ActionCard
									options={{
										tracker: entry.trackerData,
										size: 'list',
										title: entry.trackerData.display,
										entries: entry.entries,
										route: router.tracker(entry.trackerData.id),
										icon: getTrackerIcon(entry.trackerData.icon),
										lastChild: i === entries.general.length - 1 ? true : undefined,
										button: {
											status: buttonStatuses?.[entry.trackerName],
											text: entry.trackerData.actionLabel
										},
										streak: entry.streak
									}}
								></ActionCard>
							{/each}
						</div>
					{:else}
						<div class="justify-self-center">
							<enhanced:img src={EmptyCorgi} alt="nothing" />
							<p class="text-center">Nothing being tracked!</p>
						</div>
					{/if}
				{:else if allEntriesDb.isError || trackersDb.isError || userOptions.isError}
					Error!
				{:else}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
			</section>

			{#if trackersDb.isSuccess && userOptions.isSuccess && subscriptions.length > 0}
				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Subscriptions</h2>

					<div class="border-base-300/50 rounded-2xl border bg-base-50">
						{#each subscriptions as sub, i (sub.id)}
							<ActionCard
								options={{
									tracker: sub,
									size: 'list',
									title: sub.display,
									entries: sub.entryData,
									route: router.tracker(sub.id),
									icon: getTrackerIcon(sub.icon),
									lastChild: i === subscriptions.length - 1 ? true : undefined,
									button: {
										status: buttonStatuses?.[sub.name],
										text: sub.actionLabel
									}
								}}
							></ActionCard>
						{/each}
					</div>
				</section>
			{/if}
		</div>
	</main>
</PageWrapper>
