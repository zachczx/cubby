<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { filteredMarketPricesQueryOptions, filteredMarketInsightsQueryOptions, queryClient } from '$lib/queries';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { goto } from '$app/navigation';
	import { router } from '$lib/routes';
	import Sparkline from '$lib/ui/Sparkline.svelte';
	import { titleCase } from '$lib/utils';

	dayjs.extend(relativeTime);

	let { data } = $props();

	const pricesQuery = createQuery(() => filteredMarketPricesQueryOptions({ category: data.category, item: data.item }));
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

	let insight = $derived.by(() => {
		if (!insightsQuery.isSuccess || !insightsQuery.data) return null;
		return insightsQuery.data.find((i) => i.itemName.toLowerCase() === data.item.toLowerCase()) ?? null;
	});

	let prices = $derived(pricesQuery.isSuccess ? pricesQuery.data ?? [] : []);

	let sparklineData = $derived.by(() => {
		if (prices.length < 2) return [];
		return [...prices]
			.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime())
			.map((p) => ({ date: p.createdAt, price: p.price }));
	});
</script>

<PageWrapper title={titleCase(data.item)}>
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<a href={router.marketAdd({ category: data.category, item: data.item })} class="btn btn-primary btn-lg w-full rounded-full">
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price
				</a>
			</section>

			{#if insightsQuery.isLoading}
				<div class="skeleton h-40 w-full rounded-2xl"></div>
			{:else if insight}
				<section>
					<div class="border-base-300/50 bg-base-50 flex flex-col gap-3 rounded-2xl border p-4">
						<div class="flex items-center gap-2">
							{#if insight.country}
								<span class="text-base-content/50 text-xs tracking-wider uppercase"
									>{insight.country}</span
								>
							{/if}
							{#if insight.category}
								<span class="text-base-content/50 text-xs">{insight.category}</span>
							{/if}
						</div>
						<div class="grid grid-cols-2 gap-4">
							<div class="flex flex-col">
								<span class="text-base-content/60 text-xs">Current</span>
								<span class="text-xl font-semibold">${insight.latestPrice.toFixed(2)}</span>
								{#if insight.latestUnit}
									<span class="text-base-content/50 text-xs"
										>${insight.latestUnit.toFixed(2)} / unit</span
									>
								{/if}
								{#if insight.latestStore}
									<span class="text-base-content/50 text-xs">@ {insight.latestStore}</span>
								{/if}
							</div>
							<div class="flex flex-col">
								<span class="text-base-content/60 text-xs">Lowest</span>
								<span class="text-success text-xl font-semibold"
									>${insight.lowestPrice.toFixed(2)}</span
								>
								{#if insight.lowestUnit}
									<span class="text-success/70 text-xs"
										>${insight.lowestUnit.toFixed(2)} / unit</span
									>
								{/if}
								{#if insight.lowestStore}
									<span class="text-base-content/50 text-xs">@ {insight.lowestStore}</span>
								{/if}
							</div>
						</div>
					</div>
				</section>
			{/if}

			{#if sparklineData.length > 0}
				<section>
					<div class="border-base-300/50 bg-base-50 rounded-2xl border p-4">
						<Sparkline points={sparklineData} />
					</div>
				</section>
			{/if}

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Price History</h2>

				{#if pricesQuery.isLoading}
					<div class="skeleton h-16 w-full rounded-2xl"></div>
				{:else if prices.length > 0}
					<div class="border-base-300/50 bg-base-50 divide-base-300/50 divide-y rounded-2xl border">
						{#each prices as price}
							{@const up = unitPrice(price.price, price.quantity)}
							<div
								class="hover:bg-base-200/50 flex items-center justify-between p-4 transition-colors"
							>
								<div class="flex flex-col gap-0.5">
									<div class="flex items-center gap-2">
										{#if price.isPromo}
											<span class="badge badge-warning badge-xs text-warning-content">Sale</span>
										{/if}
									</div>
									<div class="text-base-content/60 flex flex-wrap items-center gap-2 text-xs">
										{#if price.quantity && price.unit}
											<span>{price.quantity} {price.unit}</span>
										{:else if price.unit}
											<span>/ {price.unit}</span>
										{/if}
										{#if price.store}
											<span class="flex items-center gap-1"
												><Icon icon="material-symbols:storefront" /> {price.store}</span
											>
										{/if}
										{#if price.country}
											<span class="flex items-center gap-1"
												><Icon icon="material-symbols:public" /> {price.country}</span
											>
										{/if}
										<span>· {dayjs(price.createdAt).fromNow()}</span>
									</div>
									{#if price.remarks}
										<span class="text-base-content/40 text-xs italic">{price.remarks}</span>
									{/if}
								</div>
								<div class="flex items-center gap-3">
									<div class="flex flex-col items-end">
										<span class="text-lg font-semibold">${price.price.toFixed(2)}</span>
										{#if up}
											<span class="text-base-content/50 text-xs"
												>${up.toFixed(2)}/{price.unit || 'unit'}</span
											>
										{/if}
									</div>

									<button
										class="btn btn-ghost btn-sm btn-circle text-base-content/40 hover:text-primary"
										onclick={() => goto(router.marketEdit(price.id), { state: { price } })}
										aria-label="Edit price"
									>
										<Icon icon="material-symbols:edit-outline" class="size-4" />
									</button>
									<button
										class="btn btn-ghost btn-sm btn-circle text-base-content/40 hover:text-error"
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
