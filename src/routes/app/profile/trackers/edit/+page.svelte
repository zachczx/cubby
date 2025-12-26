<script lang="ts">
	import MaterialSymbolsFilterList from '$lib/assets/svg/MaterialSymbolsFilterList.svelte';
	import { pb } from '$lib/pb';
	import { allTrackersQueryOptions } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';

	dayjs.extend(relativeTime);

	const trackersDb = createQuery(allTrackersQueryOptions);

	let trackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return [];

		let trackers: Record<string, TrackerDB[]> = {};

		for (const t of trackersDb.data) {
			if (!trackers[t.family]) {
				trackers[t.family] = [];
			}

			trackers[t.family].push(t);
		}

		return trackers;
	});
</script>

<PageWrapper title="Edit Tracker" {pb} largeScreenCenter={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8">
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Edit Trackers</h1>

		<div class="grid gap-8">
			{#each Object.values(trackers) as trackerList}
				<section
					class="border-base-300 grid min-h-18 content-start rounded-2xl border bg-white/70 p-6"
				>
					{#if trackerList && trackerList.length > 0}
						<h2 class="mb-2 flex items-center gap-4 text-xl font-bold">
							{trackerList[0].expand?.family?.name}
						</h2>
					{/if}

					{#each trackerList as tracker (tracker.id)}
						{@render menuItem(tracker)}
					{/each}
				</section>
			{/each}
		</div>
	</div>
</PageWrapper>
<MaterialSymbolsFilterList class="size-6" />

{#snippet menuItem(tracker: TrackerDB)}
	<div
		class="border-b-base-300/50 flex items-center gap-2 border-b py-6 last-of-type:border-b-0 last-of-type:pb-0"
	>
		<a
			href="/app/profile/trackers/edit/{tracker.id}"
			aria-label="edit"
			class="flex w-full cursor-pointer items-center"
		>
			<div class="grow font-medium">
				{tracker.display}
			</div>
			<span><MaterialSymbolsChevronRight class="size-5 opacity-75" /></span>
		</a>
	</div>
{/snippet}
