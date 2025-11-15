<script lang="ts">
	import { createSprayQueryOptions, createSprayRefetchOptions } from '$lib/queries';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'spray',
		labels: {
			pageTitle: 'Spray',
			ctaButtonText: 'Sprayed'
		},
		queries: {
			query: createSprayQueryOptions,
			refetch: createSprayRefetchOptions
		}
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
