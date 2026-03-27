<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { marketPricesQueryOptions } from '$lib/queries';
	import { router } from '$lib/routes';
	import AddPriceLog from './AddPriceLog.svelte';

	const categories = [
		{ name: 'Fruit', icon: 'material-symbols:nutrition' },
		{ name: 'Vegetable', icon: 'material-symbols:eco' },
		{ name: 'Dairy', icon: 'material-symbols:egg-alt' },
		{ name: 'Meat', icon: 'material-symbols:restaurant' },
		{ name: 'Seafood', icon: 'material-symbols:set-meal' },
		{ name: 'Bakery', icon: 'material-symbols:bakery-dining' },
		{ name: 'Pantry', icon: 'material-symbols:shelves' },
		{ name: 'Frozen', icon: 'material-symbols:ac-unit' },
		{ name: 'Beverage', icon: 'material-symbols:local-cafe' },
		{ name: 'Snack', icon: 'material-symbols:cookie' },
		{ name: 'Other', icon: 'material-symbols:category' }
	];

	const pricesQuery = createQuery(marketPricesQueryOptions);

	let isModalOpen = $state(false);
	let editPrice = $state<any | null>(null);

	function openAddModal() {
		editPrice = null;
		isModalOpen = true;
	}

	function handleCloseModal() {
		isModalOpen = false;
		editPrice = null;
		pricesQuery.refetch();
	}

	let itemCounts = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data) return {};
		const counts: Record<string, number> = {};
		const seen = new Set<string>();
		for (const p of pricesQuery.data) {
			if (!p.category) continue;
			const key = `${p.category.toLowerCase()}::${p.itemName.toLowerCase()}`;
			if (!seen.has(key)) {
				seen.add(key);
				counts[p.category] = (counts[p.category] || 0) + 1;
			}
		}
		return counts;
	});
</script>

<PageWrapper title="Market">
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<button class="btn btn-primary btn-lg w-full rounded-full" onclick={openAddModal}>
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price Log
				</button>
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Categories</h2>

				{#if pricesQuery.isLoading}
					<div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
						{#each Array(10) as _}
							<div class="skeleton h-29.5 w-full rounded-2xl"></div>
						{/each}
					</div>
				{:else if pricesQuery.isError}
					<div class="border-base-300/50 bg-base-50 text-error rounded-2xl border p-4">
						Failed to load market data.
					</div>
				{:else}
					<div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
						{#each categories as cat}
							{@const count = itemCounts[cat.name] ?? 0}
							<a
								href={router.marketCategory(cat.name)}
								class="border-base-300/50 bg-base-50 hover:border-primary/30 flex flex-col gap-2 rounded-2xl border p-4 transition-colors"
							>
								<div class="flex items-center gap-2">
									<div class="bg-primary/10 text-primary rounded-xl p-2">
										<Icon icon={cat.icon} class="size-5" />
									</div>
								</div>
								<div>
									<span class="font-semibold">{cat.name}</span>
									<p class="text-base-content/50 text-xs">
										{count === 0 ? 'No items' : `${count} item${count === 1 ? '' : 's'}`}
									</p>
								</div>
							</a>
						{/each}
					</div>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>

{#if isModalOpen}
	<AddPriceLog onClose={handleCloseModal} {editPrice} />
{/if}
