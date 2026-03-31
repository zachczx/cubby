<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { marketPriceQueryOptions } from '$lib/queries';
	import MarketPriceForm from '../../MarketPriceForm.svelte';

	let { data } = $props();

	const statePrice = (page.state as { price?: MarketPriceDB }).price ?? null;
	const priceQuery = createQuery(() => marketPriceQueryOptions(data.id));

	let price: MarketPriceDB | null = $derived(statePrice ?? (priceQuery.isSuccess ? priceQuery.data : null));
</script>

<PageWrapper title="Edit Price" focusedScreen={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:bg-base-50 lg:p-8 lg:shadow-md">
		{#if price}
			<MarketPriceForm editPrice={price} onSuccess={() => history.back()} />
		{:else if priceQuery.isLoading}
			<div class="flex justify-center p-8">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{:else}
			<p class="text-base-content/50 p-8 text-center">Price not found.</p>
		{/if}
	</div>
</PageWrapper>
