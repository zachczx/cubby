<script lang="ts">
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsFilterList from '$lib/assets/svg/MaterialSymbolsFilterList.svelte';
	import MaterialSymbolsHistory from '$lib/assets/svg/MaterialSymbolsHistory.svelte';
	import { pb } from '$lib/pb';
	import {
		allLogsRefetchOptions,
		allTrackersQueryOptions,
		allTrackersRefetchOptions,
		feedQueryOptions
	} from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import MaterialSymbolsKeep from '$lib/assets/svg/MaterialSymbolsKeep.svelte';
	import MaterialSymbolsKeepOutline from '$lib/assets/svg/MaterialSymbolsKeepOutline.svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { preventDefault } from 'svelte/legacy';
	import MaterialSymbolsVisibilityOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOutline.svelte';
	import MaterialSymbolsVisibilityOffOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOffOutline.svelte';
	import MaterialSymbolsEdit from '$lib/assets/svg/MaterialSymbolsEdit.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';

	dayjs.extend(relativeTime);

	const trackersDb = createQuery(allTrackersQueryOptions);
	const queryKey = allTrackersQueryOptions().queryKey;

	const tanstackClient = useQueryClient();

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
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Active Trackers</h1>

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
