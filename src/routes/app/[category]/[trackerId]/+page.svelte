<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import { allTrackersQueryOptions } from '$lib/queries';
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
			},
			calculateGaps: calculateRecordsGapsWithVacation
		};
	});
</script>

<TrackerPage {options}></TrackerPage>
