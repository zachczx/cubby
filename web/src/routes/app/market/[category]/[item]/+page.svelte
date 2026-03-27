<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { marketPricesQueryOptions, marketInsightsQueryOptions } from '$lib/queries';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import AddPriceLog from '../../AddPriceLog.svelte';

	dayjs.extend(relativeTime);

	let { data } = $props();

	const pricesQuery = createQuery(marketPricesQueryOptions);
	const insightsQuery = createQuery(marketInsightsQueryOptions);

	let isModalOpen = $state(false);
	let editPrice = $state<any | null>(null);
	let deletingId = $state<string | null>(null);

	function openAddModal() {
		editPrice = null;
		isModalOpen = true;
	}

	function openEditModal(price: any) {
		editPrice = price;
		isModalOpen = true;
	}

	function handleCloseModal() {
		isModalOpen = false;
		editPrice = null;
		pricesQuery.refetch();
		insightsQuery.refetch();
	}

	async function handleDelete(id: string) {
		deletingId = id;
		try {
			await api.delete(`market/prices/${id}`);
			addToast('success', 'Price deleted');
			pricesQuery.refetch();
			insightsQuery.refetch();
		} catch {
			addToast('error', 'Failed to delete');
		} finally {
			deletingId = null;
		}
	}

	function deltaLabel(d: number): string {
		if (d <= 0) return 'At lowest';
		return `+${d.toFixed(1)}% vs lowest`;
	}

	function unitPrice(price: number, quantity: number | null): number | null {
		if (!quantity || quantity <= 0) return null;
		return price / quantity;
	}

	function matchesItem(name: string | null): boolean {
		if (!name) return false;
		return name.toLowerCase() === data.item.toLowerCase();
	}

	let insight = $derived.by(() => {
		if (!insightsQuery.isSuccess || !insightsQuery.data) return null;
		return insightsQuery.data.find((i) => matchesItem(i.itemName)) ?? null;
	});

	let filteredPrices = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data) return [];
		return pricesQuery.data.filter((p) => matchesItem(p.itemName));
	});
</script>

<PageWrapper title={data.item}>
	<main class="h-full relative pb-24">
		<div class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<!-- Insight Summary -->
			{#if insightsQuery.isLoading}
				<div class="skeleton h-40 w-full rounded-2xl"></div>
			{:else if insight}
				<section>
					<div class="border-base-300/50 bg-base-50 rounded-2xl border p-4 flex flex-col gap-3">
						<div class="flex justify-between items-start">
							<div class="flex gap-2 items-center">
								{#if insight.country}
									<span class="text-xs text-base-content/50 uppercase tracking-wider">{insight.country}</span>
								{/if}
								{#if insight.category}
									<span class="text-xs text-base-content/50">{insight.category}</span>
								{/if}
							</div>
							{#if insight.deltaPercent <= 0}
								<div class="badge badge-success text-success-content badge-sm">Best Price</div>
							{:else if insight.deltaPercent <= 10}
								<div class="badge badge-warning text-warning-content badge-sm">Near Lowest</div>
							{:else}
								<div class="badge badge-ghost badge-sm">{deltaLabel(insight.deltaPercent)}</div>
							{/if}
						</div>
						<div class="grid grid-cols-2 gap-4">
							<div class="flex flex-col">
								<span class="text-xs text-base-content/60">Current</span>
								<span class="font-semibold text-xl">${insight.latestPrice.toFixed(2)}</span>
								{#if insight.latestUnit}
									<span class="text-xs text-base-content/50">${insight.latestUnit.toFixed(2)} / unit</span>
								{/if}
								{#if insight.latestStore}
									<span class="text-xs text-base-content/50">@ {insight.latestStore}</span>
								{/if}
							</div>
							<div class="flex flex-col">
								<span class="text-xs text-base-content/60">Lowest</span>
								<span class="font-semibold text-xl text-success">${insight.lowestPrice.toFixed(2)}</span>
								{#if insight.lowestUnit}
									<span class="text-xs text-success/70">${insight.lowestUnit.toFixed(2)} / unit</span>
								{/if}
								{#if insight.lowestStore}
									<span class="text-xs text-base-content/50">@ {insight.lowestStore}</span>
								{/if}
							</div>
						</div>
					</div>
				</section>
			{/if}

			<!-- Price History -->
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Price History</h2>

				{#if pricesQuery.isLoading}
					<div class="skeleton h-16 w-full rounded-2xl"></div>
				{:else if filteredPrices.length > 0}
					<div class="border-base-300/50 rounded-2xl border bg-base-50 divide-y divide-base-300/50">
						{#each filteredPrices as price}
							{@const up = unitPrice(price.price, price.quantity)}
							<div class="p-4 flex justify-between items-center hover:bg-base-200/50 transition-colors">
								<div class="flex flex-col gap-0.5">
									<div class="flex items-center gap-2">
										{#if price.isPromo}
											<span class="badge badge-warning badge-xs text-warning-content">Sale</span>
										{/if}
									</div>
									<div class="flex gap-2 text-xs text-base-content/60 items-center flex-wrap">
										{#if price.quantity && price.unit}
											<span>{price.quantity} {price.unit}</span>
										{:else if price.unit}
											<span>/ {price.unit}</span>
										{/if}
										{#if price.store}
											<span class="flex items-center gap-1"><Icon icon="material-symbols:storefront" /> {price.store}</span>
										{/if}
										{#if price.country}
											<span class="flex items-center gap-1"><Icon icon="material-symbols:public" /> {price.country}</span>
										{/if}
										<span>· {dayjs(price.createdAt).fromNow()}</span>
									</div>
								</div>
								<div class="flex items-center gap-3">
									<div class="flex flex-col items-end">
										<span class="font-semibold text-lg">${price.price.toFixed(2)}</span>
										{#if up}
											<span class="text-xs text-base-content/50">${up.toFixed(2)}/{price.unit || 'unit'}</span>
										{/if}
									</div>
									<button
										class="btn btn-ghost btn-sm btn-circle text-base-content/40 hover:text-error"
										onclick={() => handleDelete(price.id)}
										disabled={deletingId === price.id}
										aria-label="Delete price"
									>
										{#if deletingId === price.id}
											<span class="loading loading-spinner loading-xs"></span>
										{:else}
											<Icon icon="material-symbols:delete-outline" class="size-4" />
										{/if}
									</button>
									<button
										class="btn btn-ghost btn-sm btn-circle text-base-content/40 hover:text-primary"
										onclick={() => openEditModal(price)}
										aria-label="Edit price"
									>
										<Icon icon="material-symbols:edit-outline" class="size-4" />
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

		<button
			class="btn btn-primary btn-circle fixed bottom-24 right-6 shadow-xl z-40 sm:bottom-12 sm:right-12"
			onclick={openAddModal}
			aria-label="Add Price Log"
		>
			<Icon icon="material-symbols:add" class="size-6 text-primary-content" />
		</button>
	</main>
</PageWrapper>

{#if isModalOpen}
	<AddPriceLog onClose={handleCloseModal} {editPrice} />
{/if}
