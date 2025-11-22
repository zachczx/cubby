<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'doggoChewable',
		labels: {
			pageTitle: 'Chewable',
			ctaButtonText: 'Fed',
			noun: 'chewable'
		},
		calculateGaps: calculateRecordsGapsWithVacation
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
