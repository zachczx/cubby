<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { page } from '$app/state';
	import MarketPriceForm from '../MarketPriceForm.svelte';

	let { data } = $props();

	const duplicatePrice =
		(page.state as { duplicatePrice?: MarketPriceDB }).duplicatePrice ?? null;
</script>

<PageWrapper title="Add Price" focusedScreen={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:bg-base-50 lg:p-8 lg:shadow-md">
		<MarketPriceForm
			prefillCategory={data.prefillCategory}
			prefillItemName={data.prefillItemName}
			{duplicatePrice}
			onSuccess={() => history.back()}
		/>
	</div>
</PageWrapper>
