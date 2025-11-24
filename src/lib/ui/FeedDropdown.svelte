<script lang="ts">
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import MaterialSymbolsHistory from '$lib/assets/svg/MaterialSymbolsHistory.svelte';
	import MaterialSymbolsNotificationImportant from '$lib/assets/svg/MaterialSymbolsNotificationImportant.svelte';
	import MaterialSymbolsNotifications from '$lib/assets/svg/MaterialSymbolsNotifications.svelte';
	import MaterialSymbolsWarning from '$lib/assets/svg/MaterialSymbolsWarning.svelte';
	import { feedQueryOptions, notificationQueryOptions, trackerIdToName } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';

	const feed = createQuery(feedQueryOptions);

	let feedItems = $derived(feed.isSuccess ? feed.data : undefined);

	//Old
	const latestLogs = createQuery(notificationQueryOptions);

	interface CollectionNotification {
		name: CollectionName | null;
		notif: NotificationStatus;
	}

	let notifications = $derived.by(() => {
		if (!latestLogs.isSuccess) return undefined;

		let collections: CollectionNotification[] = [];

		// for (const ll of latestLogs.data) {
		// 	collections.push({
		// 		name: trackerIdToName(ll.tracker),
		// 		notif: getTrackerStatus(ll)
		// 	});
		// }

		// const notifCount = collections.reduce((total, item) => {
		// 	if (item.notif.show) {
		// 		return (total += 1);
		// 	}
		// 	return total;
		// }, 0);

		// return {
		// 	count: notifCount,
		// 	collections: collections
		// };
	});

	type Notifications = typeof notifications;
</script>

<!-- {#if notifications}
					{@render notification(notifications)}
				{/if} -->
<!-- {@render notification(notifications)} -->

<div class="dropdown dropdown-end">
	<div tabindex="0" role="button" class="btn btn-ghost drawer-button px-2 py-0">
		<MaterialSymbolsHistory
			class={['size-6', feedItems?.items.length === 0 ? 'opacity-50' : undefined]}
		/>
	</div>
	<ul
		class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 w-72 p-2 shadow-lg"
	>
		{#if feedItems?.items?.length === 0}
			<li>
				<div class="flex items-center justify-center gap-2">
					<MaterialSymbolsCheckCircle class="size-[1.3em]" /><span>No pending items</span>
				</div>
			</li>
		{:else}
			{#each feedItems?.items as item}
				<li class="border-b-base-content/10 border-b py-1 last:border-b-0">
					<div class="flex items-center">
						<div class="flex grow items-center gap-2">
							<span class="font-semibold">{item.expand?.tracker?.display}</span
							><MaterialSymbolsArrowRightAlt />
							{dayjs(item.time).fromNow()}
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
