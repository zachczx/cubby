<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import { createTowelQueryOptions, createTowelRefetchOptions } from '$lib/queries';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'towel',
		labels: {
			pageTitle: 'Towel',
			ctaButtonText: 'Washed'
		},
		queries: {
			query: createTowelQueryOptions,
			refetch: createTowelRefetchOptions
		},
		calculateGaps: calculateRecordsGapsWithVacation
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
