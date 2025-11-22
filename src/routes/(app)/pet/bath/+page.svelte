<script lang="ts">
	import { calculateRecordsGapsWithVacation } from '$lib/overlap';
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'doggoBath',
		labels: {
			pageTitle: 'Bath',
			ctaButtonText: 'Bathed',
			noun: 'bath'
		},
		calculateGaps: calculateRecordsGapsWithVacation
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
