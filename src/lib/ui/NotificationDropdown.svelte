<script lang="ts">
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import MaterialSymbolsNotificationImportant from '$lib/assets/svg/MaterialSymbolsNotificationImportant.svelte';
	import MaterialSymbolsNotifications from '$lib/assets/svg/MaterialSymbolsNotifications.svelte';
	import MaterialSymbolsWarning from '$lib/assets/svg/MaterialSymbolsWarning.svelte';
	import { notificationQueryOptions, trackerIdToName } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';

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
		{#if notifications?.count === 0}
			<MaterialSymbolsNotifications class="size-6" />
		{:else}
			<MaterialSymbolsNotificationImportant class="size-6" />
			<span class="absolute top-1 right-1 z-2 size-4 rounded-full bg-red-600 text-xs"
				>{notifications?.count}</span
			>
		{/if}
	</div>
	<ul
		class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 w-72 p-2 shadow-lg"
	>
		{#if notifications?.count === 0}
			<li>
				<div class="flex items-center justify-center gap-2">
					<MaterialSymbolsCheckCircle class="size-[1.3em]" /><span>No pending items</span>
				</div>
			</li>
		{:else}
			{#each notifications?.collections as collection}
				{#if collection.notif.show}
					<li>
						<a href={collection.notif.href} class="flex items-center">
							<div class="flex grow items-center gap-2">
								<MaterialSymbolsWarning class="size-[1.3em]" />{collection.notif.label}
							</div>
							<div>
								<MaterialSymbolsChevronRight class="size-5 opacity-50" />
							</div>
						</a>
					</li>
				{/if}
			{/each}
		{/if}
		<!-- {#if !sprayNotification.show && !towelNotification.show && !gummyNotification.show}
				
			{/if}

			{#if sprayNotification.show}
				
			{/if}

			{#if towelNotification.show}
				<li>
					<a href="/household/towel" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if gummyNotification.show}
				<li>
					<a href="/personal/gummy" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if bedsheetNotification.show}
				<li>
					<a href="/household/bedsheet" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if doggoBathNotification.show}
				<li>
					<a href="/pet/bath" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if doggoChewableNotification.show}
				<li>
					<a href="/pet/chewable" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if} -->
	</ul>
</div>
