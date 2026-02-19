<script lang="ts">
	import { allTrackersQueryOptions } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';

	dayjs.extend(relativeTime);

	const trackersDb = createQuery(allTrackersQueryOptions);

	let trackers = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return [];

		let trackers: Record<string, TrackerDB[]> = {};

		for (const t of trackersDb.data) {
			if (!trackers[t.familyId]) {
				trackers[t.familyId] = [];
			}

			trackers[t.familyId].push(t);
		}

		return trackers;
	});

	$inspect(trackers);
</script>

<PageWrapper title="Edit Tracker">
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center">
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Edit Trackers</h1>

		<div class="grid gap-8">
			{#each Object.values(trackers) as trackerList}
				<section
					class="border-base-300 grid min-h-18 content-start rounded-2xl border bg-white/70 p-6"
				>
					{#if trackerList && trackerList.length > 0}
						<h2 class="mb-2 flex items-center gap-4 text-xl font-bold">
							{trackerList[0].familyName}
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
<Icon icon="material-symbols:filter-list" class="size-6" />

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
			<span><Icon icon="material-symbols:chevron-right" class="size-5 opacity-75" /></span>
		</a>
	</div>
{/snippet}
