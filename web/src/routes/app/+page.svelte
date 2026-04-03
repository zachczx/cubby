<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		allEntriesQueryOptions,
		allTrackersQueryOptions,
		allWorkoutsQueryOptions,
		getAllWorkoutsQueryKey,
		gymSummaryQueryOptions,
		userQueryOptions
	} from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import { calculateStreak } from '$lib/streaks';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import GymLaunchpad from '$lib/ui/GymLaunchpad.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import { getColoredTrackers, getTrackerIcon, generateSubscriptionEntries } from '$lib/mapper';
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { router } from '$lib/routes';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { getDashboardTasks, getActiveWorkout } from '$lib/dashboard';
	import { goto } from '$app/navigation';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const queryClient = useQueryClient();
	const trackersDb = createQuery(allTrackersQueryOptions);
	const allEntriesDb = createQuery(allEntriesQueryOptions);
	const userOptions = createQuery(userQueryOptions);
	const workoutsDb = createQuery(allWorkoutsQueryOptions);
	const gymSummary = createQuery(gymSummaryQueryOptions);

	let buttonStatuses = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;
		const statuses = <Record<string, ButtonState>>{};
		for (const t of trackersDb.data) {
			statuses[t.name] = 'default';
		}
		return statuses;
	});

	let pinnedTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return [];
		return getColoredTrackers(trackersDb.data).filter(
			(t) => t.pinned && t.show && t.kind === 'task'
		);
	});

	let generalTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return [];
		return getColoredTrackers(trackersDb.data).filter(
			(t) => !t.pinned && t.show && t.kind === 'task'
		);
	});

	let subscriptionTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data || !userOptions.isSuccess || !userOptions.data)
			return [];
		return getColoredTrackers(trackersDb.data)
			.filter((t) => t.show && t.kind === 'subscription')
			.map((sub) => ({
				...sub,
				syntheticEntries: generateSubscriptionEntries(sub, userOptions.data!.id)
			}));
	});

	let pinnedEntries = $derived.by(() => {
		if (!allEntriesDb.isSuccess || !allEntriesDb.data) return [];
		return classifyTrackers(pinnedTrackers, allEntriesDb.data);
	});

	let urgentTasks = $derived.by(() => {
		if (!allEntriesDb.isSuccess || !allEntriesDb.data) return [];

		const tasks = getDashboardTasks(generalTrackers, allEntriesDb.data);

		// Treat subscriptions as upcoming tasks too
		const subs = getDashboardTasks(
			subscriptionTrackers,
			subscriptionTrackers.flatMap((s) => s.syntheticEntries)
		);

		return [...tasks, ...subs];
	});

	let activeWorkout = $derived(
		workoutsDb.isSuccess ? getActiveWorkout(workoutsDb.data) : undefined
	);

	const gymMonthlyTarget = 8;
	let gymThisMonth = $derived(
		gymSummary.isSuccess ? (gymSummary.data?.totalWorkoutsThisMonth ?? 0) : 0
	);

	function classifyTrackers(trackers: TrackerColored[], entries: EntryDB[]) {
		const data = [];
		for (const t of trackers) {
			const entryData = entries.filter((entry) => t.id === entry.trackerId);
			data.push({
				trackerName: t.name,
				trackerData: t,
				entries: entryData,
				notification: getTrackerStatus(entryData),
				streak: calculateStreak(entryData, t)
			});
		}
		return data;
	}

	async function startWorkout() {
		const response = await api.post('gym/workouts');
		if (response.status === 201) {
			const workout = await response.json<WorkoutDB>();
			addToast('success', 'Workout started!');
			queryClient.invalidateQueries({ queryKey: getAllWorkoutsQueryKey() });
			goto(router.gym(workout.id));
		} else {
			addToast('error', 'Failed to start workout');
		}
	}

	onMount(async () => {
		const response = await api.get('check');
		if (response.status !== 204) {
			return;
		}
	});
</script>

<PageWrapper title="Dashboard" back={false}>
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			<section class="grid gap-3">
				{#if allEntriesDb.isSuccess && trackersDb.isSuccess && userOptions.isSuccess}
					<a
						href="/app/trackers"
						class={[
							'border-base-300 flex items-center gap-3 rounded-2xl border px-4 py-3 transition-shadow hover:shadow',
							urgentTasks.length > 0
								? 'text-error-content bg-error'
								: 'text-base-content bg-base-50'
						]}
					>
						<span class={['text-3xl font-bold']}>{urgentTasks.length}</span>
						<span class="grow text-sm font-medium opacity-85"
							>Urgent Task{urgentTasks.length !== 1 ? 's' : ''}</span
						>
						<Icon icon="material-symbols:chevron-right" class="size-5 opacity-60" />
					</a>
				{:else}
					<div
						class="border-base-300/50 bg-base-50 flex items-center gap-3 rounded-2xl border px-4 py-3"
					>
						<div class="skeleton h-9 w-8 rounded-full"></div>
						<div class="skeleton h-4 w-24 grow-0 rounded-full"></div>
						<div class="grow"></div>
						<Icon icon="material-symbols:chevron-right" class="text-base-content/20 size-5" />
					</div>
				{/if}
			</section>

			{#if workoutsDb.isSuccess && gymSummary.isSuccess}
				<GymLaunchpad
					{activeWorkout}
					{gymThisMonth}
					{gymMonthlyTarget}
					onStartWorkout={startWorkout}
				/>
			{:else}
				<div class="border-primary/20 bg-primary/5 grid h-28.5 gap-3 rounded-2xl border p-4">
					<div class="flex items-center justify-between">
						<div class="skeleton h-4 w-24 rounded-full"></div>
						<div class="flex items-center gap-1">
							{#each Array(gymMonthlyTarget) as _}
								<div class="skeleton size-3.5 rounded-sm"></div>
							{/each}
						</div>
					</div>
					<div class="skeleton h-12 w-full rounded-full"></div>
				</div>
			{/if}

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Pinned</h2>

				{#if allEntriesDb.isSuccess && trackersDb.isSuccess && userOptions.isSuccess}
					{#each pinnedEntries as entry (entry.trackerData?.id)}
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
				{:else if allEntriesDb.isError}
					Error!
				{:else}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Upcoming</h2>

				{#if allEntriesDb.isSuccess && trackersDb.isSuccess && userOptions.isSuccess}
					{#if urgentTasks.length > 0}
						<div class="border-base-300/50 bg-base-50 rounded-2xl border">
							{#each urgentTasks as entry, i (entry.trackerData?.id)}
								<ActionCard
									options={{
										tracker: entry.trackerData,
										size: 'list',
										title: entry.trackerData?.display,
										entries: entry.entries,
										route: router.tracker(entry.trackerData?.id),
										icon: getTrackerIcon(entry.trackerData?.icon),
										lastChild: i === urgentTasks.length - 1 ? true : undefined,
										button: {
											status: buttonStatuses?.[entry.trackerName],
											text: entry.trackerData?.actionLabel
										},
										streak: entry.streak
									}}
								></ActionCard>
							{/each}
						</div>
					{:else}
						<div class="justify-self-center">
							<enhanced:img src={EmptyCorgi} alt="nothing" />
							<p class="text-center">Nothing urgent!</p>
						</div>
					{/if}
				{:else if allEntriesDb.isError}
					Error!
				{:else}
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
					<SkeletonActionCard size="compact" />
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>
