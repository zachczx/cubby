<script lang="ts">
	import TrackerPage from '$lib/shell/TrackerPage.svelte';

	const options = {
		collectionName: 'gummy',
		labels: {
			pageTitle: 'Elderberry',
			ctaButtonText: 'Ate',
			noun: 'gummy'
		}
	} satisfies TrackerPageOptions;
</script>

<TrackerPage {options}></TrackerPage>
