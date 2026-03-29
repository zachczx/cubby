<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { marketInsightsQueryOptions } from '$lib/queries';
	import { router } from '$lib/routes';
	import AddPriceLog from '../AddPriceLog.svelte';

	let { data } = $props();

	const insightsQuery = createQuery(marketInsightsQueryOptions);

	let isModalOpen = $state(false);
	let editPrice = $state<MarketPriceDB | null>(null);

	function openAddModal() {
		editPrice = null;
		isModalOpen = true;
	}

	function handleCloseModal() {
		isModalOpen = false;
		editPrice = null;
		insightsQuery.refetch();
	}

	function deltaLabel(d: number): string {
		if (d <= 0) return 'At lowest';
		return `+${d.toFixed(1)}% vs lowest`;
	}

	function matchesCategory(cat: string | null): boolean {
		if (!cat) return false;
		return cat.toLowerCase() === data.category.toLowerCase();
	}

	let filteredInsights = $derived.by(() => {
		if (!insightsQuery.isSuccess || !insightsQuery.data) return [];
		return insightsQuery.data.filter((i) => matchesCategory(i.category));
	});
</script>

<PageWrapper title={data.category.charAt(0).toUpperCase() + data.category.slice(1)}>
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<button class="btn btn-primary btn-lg w-full rounded-full" onclick={openAddModal}>
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price
				</button>
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Items</h2>

				{#if insightsQuery.isLoading}
					<div class="grid gap-3 sm:grid-cols-2">
						{#each Array(4) as _}
							<div class="skeleton h-32 w-full rounded-2xl"></div>
						{/each}
					</div>
				{:else if insightsQuery.isError}
					<div class="border-base-300/50 bg-base-50 text-error rounded-2xl border p-4">
						Failed to load items.
					</div>
				{:else if filteredInsights.length > 0}
					<div class="grid gap-3 sm:grid-cols-2">
						{#each filteredInsights as insight}
							<a
								href={router.marketItem(data.category, insight.itemName)}
								class="border-base-300/50 bg-base-50 hover:border-primary/30 flex flex-col gap-2 rounded-2xl border p-4 transition-colors"
							>
								<div class="flex items-start justify-between">
									<div>
										<h3 class="text-lg font-semibold">{insight.itemName}</h3>
										{#if insight.country}
											<span class="text-base-content/50 text-xs tracking-wider uppercase"
												>{insight.country}</span
											>
										{/if}
									</div>
									{#if insight.deltaPercent <= 0}
										<div class="badge badge-success text-success-content badge-sm">Best Price</div>
									{:else if insight.deltaPercent <= 10}
										<div class="badge badge-warning text-warning-content badge-sm">Near Lowest</div>
									{/if}
								</div>
								<div class="mt-2 grid grid-cols-2 gap-4">
									<div class="flex flex-col">
										<span class="text-base-content/60 text-xs">Current</span>
										<span class="font-medium">${insight.latestPrice.toFixed(2)}</span>
										{#if insight.latestUnit}
											<span class="text-base-content/50 text-xs"
												>${insight.latestUnit.toFixed(2)} / unit</span
											>
										{/if}
										{#if insight.latestStore}
											<span class="text-base-content/50 text-xs">{insight.latestStore}</span>
										{/if}
									</div>
									<div class="flex flex-col">
										<span class="text-base-content/60 text-xs">Lowest</span>
										<span class="text-success font-medium">${insight.lowestPrice.toFixed(2)}</span>
										{#if insight.lowestUnit}
											<span class="text-success/70 text-xs"
												>${insight.lowestUnit.toFixed(2)} / unit</span
											>
										{/if}
										{#if insight.lowestStore}
											<span class="text-base-content/50 text-xs">{insight.lowestStore}</span>
										{/if}
									</div>
								</div>
								<span
									class="text-xs {insight.deltaPercent <= 0
										? 'text-success'
										: insight.deltaPercent <= 10
											? 'text-warning'
											: 'text-base-content/50'}"
								>
									{deltaLabel(insight.deltaPercent)}
								</span>
							</a>
						{/each}
					</div>
				{:else}
					<div class="border-base-300/50 bg-base-50 rounded-2xl border p-8 text-center">
						<p class="text-base-content/50">No items in {data.category} yet.</p>
					</div>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>

{#if isModalOpen}
	<AddPriceLog paramCategory={data.category} onClose={handleCloseModal} {editPrice} />
{/if}
