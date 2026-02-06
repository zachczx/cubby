<script lang="ts">
	import { allTrackersQueryOptions } from '$lib/queries';
	import SubscriptionPage from '$lib/shell/SubscriptionPage.svelte';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	let { data } = $props();

	const trackersDb = createQuery(allTrackersQueryOptions);

	let currentTracker = $derived.by(() => {
		if (!trackersDb.isSuccess || !trackersDb.data) return;

		return trackersDb.data.find((tracker) => tracker.id === data.trackerId);
	});

	let options = $derived.by(() => {
		return {
			tracker: currentTracker,
			labels: {
				pageTitle: currentTracker?.display ?? 'Loading...',
				ctaButtonText: currentTracker?.actionLabel ?? 'Loading...',
				noun: currentTracker?.name ?? 'Loading...'
			}
		};
	});
</script>

{#if trackersDb.isPending}
	Loading..............
{:else if currentTracker?.kind === 'subscription'}
	<SubscriptionPage {options}></SubscriptionPage>
{:else}
	<TrackerPage {options}></TrackerPage>
{/if}
