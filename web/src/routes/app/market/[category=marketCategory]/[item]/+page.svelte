<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import {
		filteredMarketPricesQueryOptions,
		filteredMarketInsightsQueryOptions,
		queryClient
	} from '$lib/queries';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { goto } from '$app/navigation';
	import { router } from '$lib/routes';
	import { titleCase } from '$lib/utils';
	import { marketStores } from '$lib/market';
	import type { marketStoresType } from '$lib/market';

	dayjs.extend(relativeTime);

	let { data } = $props();

	const pricesQuery = createQuery(() =>
		filteredMarketPricesQueryOptions({ category: data.category, item: data.item })
	);
	const insightsQuery = createQuery(() => filteredMarketInsightsQueryOptions(data.category));

	let deleteDialog = $state<HTMLDialogElement | null>(null);
	let pendingDeletePrice = $state<MarketPriceDB | null>(null);
	let isDeleting = $state(false);

	function requestDelete(price: any) {
		pendingDeletePrice = price;
		deleteDialog?.showModal();
	}

	async function confirmDelete() {
		if (!pendingDeletePrice) return;
		isDeleting = true;
		try {
			await api.delete(`market/prices/${pendingDeletePrice.id}`);
			addToast('success', 'Price deleted');
			queryClient.refetchQueries({ queryKey: ['cubby', 'market-prices'] });
			queryClient.refetchQueries({ queryKey: ['cubby', 'market-insights'] });
		} catch {
			addToast('error', 'Failed to delete');
		} finally {
			isDeleting = false;
			pendingDeletePrice = null;
			deleteDialog?.close();
		}
	}

	function unitPrice(price: number, quantity: number | null): number | null {
		if (!quantity || quantity <= 0) return null;
		return price / quantity;
	}

	let allPrices = $derived(pricesQuery.isSuccess ? (pricesQuery.data ?? []) : []);

	let countries = $derived.by(() => {
		const set = new Set<string>();
		for (const p of allPrices) {
			if (p.country) set.add(p.country);
		}
		return [...set].sort();
	});

	let selectedCountry = $state<string | null>(null);

	let prices = $derived(
		selectedCountry ? allPrices.filter((p) => p.country === selectedCountry) : allPrices
	);

	let insight = $derived.by(() => {
		if (!insightsQuery.isSuccess || !insightsQuery.data) return null;
		const items = insightsQuery.data.filter(
			(i) => i.itemName.toLowerCase() === data.item.toLowerCase()
		);
		if (selectedCountry) {
			return items.find((i) => i.country === selectedCountry) ?? null;
		}
		return items[0] ?? null;
	});

	let trackedUnit = $derived.by(() => {
		if (!prices.length) return null;
		return prices[0].unit || null;
	});
</script>

<PageWrapper title={titleCase(data.item)}>
	{#snippet subtitle()}
		{#if insight?.category || trackedUnit || (countries.length === 1 && countries[0])}
			<div class="text-base-content/50 flex items-center gap-2 text-sm">
				{#if countries.length === 1}
					<span>{countries[0]}</span>
					{#if insight?.category || trackedUnit}
						<span class="text-info/60">&bull;</span>
					{/if}
				{/if}
				{#if insight?.category}
					<span>{titleCase(insight.category)}</span>
				{/if}
				{#if insight?.category && trackedUnit}
					<span class="text-info/60">&bull;</span>
				{/if}
				{#if trackedUnit}
					<span>Tracked by {titleCase(trackedUnit)}</span>
				{/if}
			</div>
		{/if}
	{/snippet}
	<main class="h-full">
		<div class="grid w-full max-w-lg justify-self-center lg:text-base">
			<section class="grid gap-3 pt-4 pb-8">
				{#if countries.length > 1}
					<div class="flex flex-wrap items-center justify-start gap-2">
						{#each countries as country}
							<button
								class="btn btn-sm rounded-full {selectedCountry === country
									? 'btn-neutral'
									: 'btn-soft btn-neutral'}"
								onclick={() => (selectedCountry = selectedCountry === country ? null : country)}
							>
								{country}
							</button>
						{/each}
					</div>
				{/if}

				<a
					href={router.marketAdd({ category: data.category, item: data.item })}
					class="btn btn-primary btn-lg mt-4 w-full rounded-full"
				>
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price
				</a>
			</section>

			<div class="grid gap-8">
				{#if insightsQuery.isLoading}
					<div class="skeleton h-40 w-full rounded-2xl"></div>
				{:else if insight}
					<section>
						<div
							class="border-base-300/50 bg-base-50 grid grid-cols-2 gap-4 rounded-2xl border p-5"
						>
							<div class="flex flex-col gap-1">
								<span class="text-base-content/50 text-xs font-medium tracking-wider uppercase"
									>Lowest</span
								>
								<span class="text-success text-2xl font-bold"
									>${insight.lowestPrice.toFixed(2)}</span
								>
								<span class="text-base-content/40 text-xs">
									{#if insight.lowestUnit}per {trackedUnit || 'unit'}{/if}
									{#if insight.lowestUnit && insight.lowestStore}&nbsp;@&nbsp;{/if}
									{#if insight.lowestStore}{insight.lowestStore}{/if}
								</span>
							</div>
							<div class="flex flex-col gap-1">
								<span class="text-base-content/50 text-xs font-medium tracking-wider uppercase"
									>Latest</span
								>
								<span class="text-2xl font-bold">${insight.latestPrice.toFixed(2)}</span>
								<span class="text-base-content/40 text-xs">
									{#if insight.latestUnit}per {trackedUnit || 'unit'}{/if}
									{#if insight.latestUnit && insight.latestStore}&nbsp;@&nbsp;{/if}
									{#if insight.latestStore}{insight.latestStore}{/if}
								</span>
							</div>
						</div>
					</section>
				{/if}

				<section class="grid gap-4 py-2">
					<h2 class="text-base-content/70 text-lg font-bold">Price History</h2>

					{#if pricesQuery.isLoading}
						<div class="skeleton h-16 w-full rounded-2xl"></div>
					{:else if prices.length > 0}
						<div
							class="border-base-300/50 bg-base-50 divide-base-300/50 divide-y rounded-2xl border"
						>
							{#each prices as price}
								{@const up = unitPrice(price.price, price.quantity)}
								{@const storeLogo = marketStores[price.store as keyof marketStoresType]?.icon}
								<div
									class="hover:bg-base-200/50 flex items-center justify-between gap-3 p-4 transition-colors"
								>
									<!-- Left: Store + Time -->
									<div class="flex items-center gap-3">
										{#if price.store}
											<div class="avatar">
												<div class="w-8 rounded-full">
													<img src={storeLogo} alt="logo" />
												</div>
											</div>
										{/if}
										<div class="flex flex-col">
											{#if price.store}
												<span class="font-semibold">{price.store}</span>
											{/if}
											<div class="text-base-content/40 flex items-center gap-2 text-xs">
												<span>{dayjs(price.createdAt).fromNow()}</span>
												{#if price.isPromo}
													<span class="badge badge-warning badge-xs text-warning-content">Sale</span
													>
												{/if}
											</div>
										</div>
									</div>

									<div class="flex items-center gap-2">
										<div class="flex flex-col items-end">
											<span class="text-lg font-bold">${price.price.toFixed(2)}</span>
											{#if up}
												<span class="text-base-content/40 text-xs"
													>${up.toFixed(2)} / {price.unit || 'unit'}</span
												>
											{/if}
										</div>
										<button
											class="btn btn-ghost btn-sm btn-circle text-base-content/30 hover:text-primary"
											onclick={() =>
												goto(router.marketAdd(), {
													state: { duplicatePrice: price }
												})}
											aria-label="Duplicate price"
										>
											<Icon icon="material-symbols:content-copy-outline" class="size-4" />
										</button>
										<button
											class="btn btn-ghost btn-sm btn-circle text-base-content/30 hover:text-primary"
											onclick={() => goto(router.marketEdit(price.id), { state: { price } })}
											aria-label="Edit price"
										>
											<Icon icon="material-symbols:edit-outline" class="size-4" />
										</button>
										<button
											class="btn btn-ghost btn-sm btn-circle text-base-content/30 hover:text-error"
											onclick={() => requestDelete(price)}
											aria-label="Delete price"
										>
											<Icon icon="material-symbols:delete-outline" class="size-4" />
										</button>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<div class="border-base-300/50 bg-base-50 rounded-2xl border p-8 text-center">
							<p class="text-base-content/50">No prices tracked for {data.item} yet.</p>
						</div>
					{/if}
				</section>
			</div>
		</div>
	</main>
</PageWrapper>

<dialog bind:this={deleteDialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-8">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<div class="grid gap-4">
			<div
				class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center rounded-full"
			>
				<Icon icon="material-symbols:delete-outline" class="size-10" />
			</div>
			<h2 class="text-2xl font-bold">Delete this price?</h2>
			{#if pendingDeletePrice}
				<p class="text-base-content/60">
					This will permanently remove the
					<span class="text-base-content font-semibold">
						${pendingDeletePrice.price.toFixed(2)}
					</span>
					entry from {dayjs(pendingDeletePrice.createdAt).fromNow()}.
				</p>
			{/if}
		</div>
		<div class="grid gap-4">
			<button class="btn btn-error btn-lg" disabled={isDeleting} onclick={confirmDelete}>
				{#if isDeleting}<span class="loading loading-spinner loading-sm"></span>{/if}
				Delete
			</button>
			<button
				class="btn btn-outline btn-neutral btn-lg w-full"
				onclick={() => deleteDialog?.close()}>Cancel</button
			>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
