<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { filteredMarketInsightsQueryOptions } from '$lib/queries';
	import { router } from '$lib/routes';
	import { titleCase } from '$lib/utils';

	let { data } = $props();

	const insightsQuery = createQuery(() => filteredMarketInsightsQueryOptions(data.category));
</script>

<PageWrapper title={titleCase(data.category)}>
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<a href={router.marketAdd({ category: data.category })} class="btn btn-primary btn-lg w-full rounded-full">
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price
				</a>
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
				{:else if insightsQuery.data && insightsQuery.data.length > 0}
					<div class="grid gap-3 sm:grid-cols-2">
						{#each insightsQuery.data as insight}
							<a
								href={router.marketItem(data.category, insight.itemName)}
								class="border-base-300/50 bg-base-50 hover:border-primary/30 flex flex-col gap-2 overflow-hidden rounded-2xl border p-4 transition-colors"
							>
								<div class="flex items-center gap-2">
									<div class="min-w-0">
										<h3 class="truncate text-lg font-semibold">{insight.itemName}</h3>
										{#if insight.country}
											<span class="text-base-content/50 text-xs tracking-wider uppercase"
												>{insight.country}</span
											>
										{/if}
									</div>
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

