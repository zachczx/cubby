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
	import { allTrackersQueryOptions } from '$lib/queries';
	import { getColoredTrackers, getFamilyColor, getTrackerIcon } from '$lib/mapper.js';
	import SkeletonActionCard from '$lib/ui/SkeletonActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';

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

	let currentTrackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		return getColoredTrackers(trackersDb.data);
	});
</script>

<PageWrapper title={data.category.charAt(0).toUpperCase() + data.category.slice(1)} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				{#if trackersDb.isPending}
					<SkeletonActionCard />
					<SkeletonActionCard />
				{:else if currentTrackers && currentTrackers.length > 0}
					{#each currentTrackers as tracker}
						<ActionCard
							options={{
								tracker: tracker,
								title: tracker.display,
								route: `/${tracker.category}/${tracker.id}`,
								icon: getTrackerIcon(tracker.name),
								button: {
									text: tracker.actionLabel,
									status: buttonStatuses?.[tracker.name]
								}
							}}
						></ActionCard>
					{/each}

					{#if data.category === 'household'}
						<div class="border-base-300 grid min-h-24 gap-4 rounded-3xl border bg-white/70 p-4">
							<a href="/household/count" class="flex items-center">
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
								onclick={() => goto('/household/count?start=true')}
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
