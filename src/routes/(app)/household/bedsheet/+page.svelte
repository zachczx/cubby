<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'bedsheet',
		labels: {
			pageTitle: 'Bedsheet',
			ctaButtonText: 'Changed'
		},
		calculateGaps: calculateRecordsGapsWithVacation
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
