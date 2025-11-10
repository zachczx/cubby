<script lang="ts">
	import { createGummyQueryOptions, createSprayRefetchOptions } from '$lib/queries';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'gummy',
		labels: {
			pageTitle: 'Elderberry',
			statusLabels: { ok: 'Ate', due: 'Due', overdue: 'Overdue' },
			ctaButtonText: 'Ate'
		},
		queries: {
			query: createGummyQueryOptions,
			refetch: createSprayRefetchOptions
		}
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
