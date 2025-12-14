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

	dayjs.extend(relativeTime);

	const trackersDb = createQuery(allTrackersQueryOptions);
	const queryKey = allTrackersQueryOptions().queryKey;

	const tanstackClient = useQueryClient();

	function updateLocalCache(trackerId: string, updates: Partial<TrackerDB>) {
		tanstackClient.setQueryData<TrackerDB[]>(queryKey, (oldData: TrackerDB[] | undefined) => {
			if (!oldData) return [];
			return oldData.map((t) => (t.id === trackerId ? { ...t, ...updates } : t));
		});
	}

	async function checkHandler(tracker: TrackerDB) {
		const newValue = !tracker.show;
		updateLocalCache(tracker.id, { show: newValue });

		try {
			await pb.collection('trackers').update(tracker.id, { show: newValue });

			addToast('success', 'Saved!');
		} catch {
			addToast('error', 'Error toggling filter!');
		}
	}

	async function pinHandler(tracker: TrackerDB) {
		const newValue = !tracker.pinned;
		updateLocalCache(tracker.id, { pinned: newValue });

		try {
			await pb.collection('trackers').update(tracker.id, { pinned: newValue });
			console.log('finished');
			if (newValue) {
				addToast('success', 'Pinned!');
			} else {
				addToast('success', 'Unpinned!');
			}
		} catch {
			addToast('error', 'Error toggling filter!');
		}
	}

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

<PageWrapper title="Active Trackers" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Active Trackers</h1>

		<div class="grid gap-8">
			{#each Object.values(trackers) as trackerList}
				<section
					class="border-base-300 grid min-h-18 content-start rounded-2xl border bg-white/70 p-4"
				>
					{#if trackerList && trackerList.length > 0}
						<h2 class="flex items-center gap-4 text-xl font-bold">
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
		class="border-b-base-300/50 flex items-center gap-2 border-b pt-4 pb-6 last-of-type:border-b-0 last-of-type:pb-0"
	>
		<button
			aria-label="pin"
			class="group size-9 cursor-pointer"
			onclick={(evt) => {
				evt.preventDefault();
				pinHandler(tracker);
			}}
		>
			{#if tracker.pinned}
				<MaterialSymbolsKeep class="size-[1.5em]" />
			{:else}
				<MaterialSymbolsKeepOutline class="size-[1.5em] opacity-30 group-hover:opacity-75" />
			{/if}
		</button>
		<label class="flex w-full cursor-pointer items-center"
			><div class="grow font-medium">
				{tracker.display}
			</div>
			<input
				type="checkbox"
				class="peer hidden"
				checked={tracker.show}
				onchange={() => checkHandler(tracker)}
			/>
			<div
				class="hover:bg-neutral/20 active:bg-neutral/20 hidden rounded-lg p-2 peer-checked:block"
			>
				<MaterialSymbolsVisibilityOutline class="size-6" />
			</div>
			<div class="hover:bg-neutral/20 active:bg-neutral/20 rounded-lg p-2 peer-checked:hidden">
				<MaterialSymbolsVisibilityOffOutline class="size-6" />
			</div>
		</label>
	</div>
{/snippet}

{#snippet item(tracker: TrackerDB)}
	<li class="grid grid-cols-[auto_1fr] items-center gap-0">
		<button
			type="button"
			onclick={(evt) => {
				evt.preventDefault();
				pinHandler(tracker);
			}}
			class="flex h-full items-center justify-center"
		>
			{#if tracker.pinned}
				<MaterialSymbolsKeep class="opacity-75" />
			{:else}
				<MaterialSymbolsKeepOutline class="opacity-75" />
			{/if}
		</button>
	</li>
{/snippet}

<input type="checkbox" class="hidden" /><svg class="hidden"></svg><svg class="hidden"></svg>

<style>
	input[type='checkbox']:checked + svg {
		margin-inline-end: 3rem;
	}

	input[type='checkbox']:checked + svg + svg {
		display: none;
	}
</style>
