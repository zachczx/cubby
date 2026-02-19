<script lang="ts">
	import Icon from '@iconify/svelte';
	import { allTrackersQueryOptions } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { api } from '$lib/api';

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

	async function visibilityHandler(tracker: TrackerDB) {
		const newValue = !tracker.show;
		updateLocalCache(tracker.id, { show: newValue });

		try {
			await api.patch(`trackers/${tracker.id}/show`, {
				body: JSON.stringify({
					show: newValue
				})
			});
		} catch (err) {
			console.error(err);
			addToast('error', 'Error toggling pin!');
		}
	}

	async function pinHandler(tracker: TrackerDB) {
		const newValue = !tracker.pinned;
		updateLocalCache(tracker.id, { pinned: newValue });

		try {
			await api.patch(`trackers/${tracker.id}/pinned`, {
				body: JSON.stringify({
					pinned: newValue
				})
			});
		} catch (err) {
			console.error(err);
			addToast('error', 'Error toggling pin!');
		}
	}

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
</script>

<PageWrapper title="Active Trackers">
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center">
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Active Trackers</h1>

		{#if trackersDb.isPending}
			<div class="grid min-h-32 content-center justify-center">
				<span class="loading loading-md loading-spinner"></span>
			</div>
		{:else}
			<div class="grid gap-8">
				{#each Object.values(trackers) as trackerList}
					<section
						class="border-base-300 grid min-h-18 content-start rounded-2xl border bg-white/70 p-6"
					>
						{#if trackerList && trackerList.length > 0}
							<h2 class="flex items-center gap-4 text-xl font-bold">
								{trackerList[0].familyName}
							</h2>
						{/if}

						{#each trackerList as tracker (tracker.id)}
							{@render menuItem(tracker)}
						{/each}
					</section>
				{/each}
			</div>
		{/if}
	</div>
</PageWrapper>
<Icon icon="material-symbols:filter-list" class="size-6" />

{#snippet menuItem(tracker: TrackerDB)}
	<div
		class="border-b-base-300/50 flex items-center gap-2 border-b py-6 last-of-type:border-b-0 last-of-type:pb-0"
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
				<Icon icon="material-symbols:keep" class="size-[1.5em]" />
			{:else}
				<Icon
					icon="material-symbols:keep-outline"
					class="size-[1.5em] opacity-30 group-hover:opacity-75"
				/>
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
				onchange={() => visibilityHandler(tracker)}
			/>
			<div
				class="hover:bg-neutral/20 active:bg-neutral/20 hidden rounded-lg p-2 peer-checked:block"
			>
				<Icon icon="material-symbols:visibility-outline" class="size-6" />
			</div>
			<div class="hover:bg-neutral/20 active:bg-neutral/20 rounded-lg p-2 peer-checked:hidden">
				<Icon icon="material-symbols:visibility-off-outline" class="size-6" />
			</div>
		</label>
	</div>
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
