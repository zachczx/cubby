<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { marketPricesQueryOptions } from '$lib/queries';
	import { router } from '$lib/routes';
	import { marketCategories } from '$lib/market';
	import AddPriceLog from './AddPriceLog.svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';

	dayjs.extend(relativeTime);

	const pricesQuery = createQuery(marketPricesQueryOptions);

	let isModalOpen = $state(false);
	let editPrice = $state<any | null>(null);
	let searchQuery = $state('');

	function openAddModal() {
		editPrice = null;
		isModalOpen = true;
	}

	function handleCloseModal() {
		isModalOpen = false;
		editPrice = null;
		pricesQuery.refetch();
	}

	function getCategoryIcon(catValue: string | null): string {
		if (!catValue) return 'material-symbols:category';
		const found = marketCategories.find((c) => c.value === catValue.toLowerCase());
		return found?.icon ?? 'material-symbols:category';
	}

	let categoryCounts = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data) return [];
		const counts: Record<string, number> = {};
		for (const p of pricesQuery.data) {
			if (!p.category) continue;
			const key = p.category.toLowerCase();
			counts[key] = (counts[key] || 0) + 1;
		}
		return marketCategories
			.map((cat) => ({
				...cat,
				count: counts[cat.value] ?? 0
			}))
			.filter((c) => c.count > 0)
			.sort((a, b) => b.count - a.count);
	});

	let topCategories = $derived(categoryCounts.slice(0, 5));

	let recentPrices = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data) return [];
		return [...pricesQuery.data]
			.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
			.slice(0, 10);
	});

	let searchResults = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data || !searchQuery.trim()) return [];
		const lower = searchQuery.toLowerCase();
		return [...pricesQuery.data]
			.filter(
				(p) =>
					p.itemName.toLowerCase().includes(lower) ||
					(p.store && p.store.toLowerCase().includes(lower))
			)
			.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
	});

	let displayedPrices = $derived(searchQuery.trim() ? searchResults : recentPrices);
</script>

<PageWrapper title="Market">
	<main class="h-full">
		<div class="grid w-full max-w-lg gap-6 justify-self-center lg:text-base">
			<section class="grid gap-3 py-2">
				<button class="btn btn-primary btn-lg w-full rounded-full" onclick={openAddModal}>
					<Icon icon="material-symbols:add" class="size-6" />
					Add Price
				</button>

				<label class="input input-bordered input-lg flex w-full items-center gap-2 rounded-full">
					<Icon icon="material-symbols:search" class="text-base-content/40 size-5" />
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search items or stores..."
						class="grow"
					/>
					{#if searchQuery}
						<button
							class="btn btn-ghost btn-xs btn-circle"
							onclick={() => (searchQuery = '')}
							aria-label="Clear search"
						>
							<Icon icon="material-symbols:close" class="size-4" />
						</button>
					{/if}
				</label>
			</section>

			{#if !searchQuery.trim()}
				<section class="grid gap-3">
					<div class="flex items-center justify-between">
						<h2 class="text-base-content/70 text-sm font-bold tracking-wide uppercase">
							Quick Categories
						</h2>
						<a
							href="/app/market/categories"
							class="btn btn-ghost btn-xs text-base-content/70 gap-1"
						>
							View All
							<Icon icon="material-symbols:chevron-right" class="size-4" />
						</a>
					</div>

					{#if pricesQuery.isLoading}
						<div class="flex gap-3 overflow-hidden">
							{#each Array(5) as _}
								<div class="skeleton h-20 w-16 shrink-0 rounded-2xl"></div>
							{/each}
						</div>
					{:else if topCategories.length > 0}
						<div
							class="flex gap-3 overflow-x-auto pb-1 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
						>
							{#each topCategories as cat}
								<a
									href={router.marketCategory(cat.value)}
									class="flex shrink-0 flex-col items-center gap-1.5"
								>
									<div
										class="text-primary/80 bg-primary/10 flex size-14 items-center justify-center rounded-2xl"
									>
										<Icon icon={cat.icon} class="size-7" />
									</div>
									<span class="text-xs font-medium">{cat.label}</span>
								</a>
							{/each}
						</div>
					{/if}
				</section>
			{/if}

			<section class="grid gap-3">
				<h2 class="text-base-content/70 text-sm font-bold tracking-wide uppercase">
					{searchQuery.trim() ? 'Results' : 'Recent'}
				</h2>

				{#if pricesQuery.isLoading}
					<div class="grid gap-2">
						{#each Array(5) as _}
							<div class="skeleton h-16 w-full rounded-2xl"></div>
						{/each}
					</div>
				{:else if pricesQuery.isError}
					<div class="border-base-300/50 bg-base-50 text-error rounded-2xl border p-4">
						Failed to load market data.
					</div>
				{:else if displayedPrices.length > 0}
					<div class="grid gap-2">
						{#each displayedPrices as price}
							{@const unitLabel =
								price.quantity && price.quantity > 1
									? `${price.quantity}${price.unit ?? ''}`
									: (price.unit ?? '')}
							<a
								href={price.category ? router.marketItem(price.category, price.itemName) : '#'}
								class="border-base-300/50 bg-base-50 hover:border-primary/30 flex items-center gap-3 rounded-2xl border p-3.5 transition-colors"
							>
								<div
									class="text-primary/80 bg-primary/10 flex size-10 shrink-0 items-center justify-center rounded-xl"
								>
									<Icon icon={getCategoryIcon(price.category)} class="size-5" />
								</div>
								<div class="min-w-0 flex-1">
									<p class="truncate font-semibold">{price.itemName}</p>
									{#if price.store}
										<p class="text-base-content/65 truncate text-xs">{price.store}</p>
									{/if}
								</div>
								<div class="shrink-0 text-right">
									<p class="font-semibold">
										${price.price.toFixed(2)}{#if unitLabel}
											<span class="text-base-content/70 text-sm font-normal">/ {unitLabel}</span
											>{/if}
									</p>
									<p class="text-base-content/65 text-xs">{dayjs(price.createdAt).fromNow()}</p>
								</div>
							</a>
						{/each}
					</div>
				{:else if pricesQuery.isSuccess}
					<div class="border-base-300/50 bg-base-50 rounded-2xl border p-8 text-center">
						<p class="text-base-content/50">
							{searchQuery.trim() ? `No results for "${searchQuery}"` : 'No price logs yet.'}
						</p>
					</div>
				{/if}
			</section>
		</div>
	</main>
</PageWrapper>

{#if isModalOpen}
	<AddPriceLog onClose={handleCloseModal} {editPrice} />
{/if}
