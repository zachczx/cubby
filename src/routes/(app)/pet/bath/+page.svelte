<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import { createDoggoBathQueryOptions, createDoggoBathRefetchOptions } from '$lib/queries';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'doggoBath',
		labels: {
			pageTitle: 'Bath',
			ctaButtonText: 'Bathed'
		},
		queries: {
			query: createDoggoBathQueryOptions,
			refetch: createDoggoBathRefetchOptions
		},
		calculateGaps: calculateRecordsGapsWithVacation
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
