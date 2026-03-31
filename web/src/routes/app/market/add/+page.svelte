<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import MarketPriceForm from '../MarketPriceForm.svelte';

	let { data } = $props();
</script>

<PageWrapper title="Add Price" focusedScreen={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:bg-base-50 lg:p-8 lg:shadow-md">
		<MarketPriceForm
			prefillCategory={data.prefillCategory}
			prefillItemName={data.prefillItemName}
			onSuccess={() => history.back()}
		/>
	</div>
</PageWrapper>
