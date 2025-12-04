<script lang="ts">
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsHistory from '$lib/assets/svg/MaterialSymbolsHistory.svelte';
	import { feedQueryOptions } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';

	dayjs.extend(relativeTime);

	const feed = createQuery(feedQueryOptions);

	let feedItems = $derived(feed.isSuccess ? feed.data : undefined);
</script>

<div class="dropdown dropdown-end">
	<div tabindex="0" role="button" class="btn btn-ghost drawer-button px-2 py-0">
		<MaterialSymbolsHistory
			class={['size-6', feedItems?.items.length === 0 ? 'opacity-50' : undefined]}
		/>
	</div>
	<ul
		class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 w-72 p-2 shadow-lg"
	>
		{#if !feedItems || feedItems?.items?.length === 0}
			<li>
				<div class="flex items-center justify-center gap-2">
					<MaterialSymbolsCheckCircle class="size-[1.3em]" /><span
						>No records yet, start tracking now!</span
					>
				</div>
			</li>
		{:else if feedItems?.items && feedItems?.items.length > 0}
			<div class="text-base-content/70 px-3 py-2 font-medium uppercase">Recent Actions</div>
			{#each feedItems?.items as item}
				{@const fromNow = dayjs(item.time).fromNow()}
				<li class="border-b-base-content/10 border-b py-1 last:border-b-0">
					<div class="flex items-center">
						<div class="flex grow items-center gap-2">
							<span class="font-semibold">{item.expand?.tracker?.display}</span
							><MaterialSymbolsArrowRightAlt />
							{fromNow}
						</div>
						<!-- <div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div> -->
					</div>
				</li>
			{/each}
		{/if}
	</ul>
</div>
