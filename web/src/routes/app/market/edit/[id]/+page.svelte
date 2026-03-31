<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { page } from '$app/state';
	import MarketPriceForm from '../../MarketPriceForm.svelte';

	let price = $derived((page.state as { price?: MarketPriceDB }).price ?? null);
</script>

<PageWrapper title="Edit Price" focusedScreen={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:bg-base-50 lg:p-8 lg:shadow-md">
		{#if price}
			<MarketPriceForm editPrice={price} onSuccess={() => history.back()} />
		{:else}
			<p class="text-base-content/50 p-8 text-center">Price not found.</p>
		{/if}
	</div>
</PageWrapper>
