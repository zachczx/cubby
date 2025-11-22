<script lang="ts">
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import type { RecordModel } from 'pocketbase';
	import type { Component } from 'svelte';
	import ActionButton from './ActionButton.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsAdd from '$lib/assets/svg/MaterialSymbolsAdd.svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import dayjs from 'dayjs';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createLogsQuery,
		logsRefetchOptions,
		notificationQueryOptions,
		notificationRefetchOptions,
		trackerIdToName,
		trackerQueryOptions
	} from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';
	import MaterialSymbolsRefresh from '$lib/assets/svg/MaterialSymbolsRefresh.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';

	let { options }: { options: ActionCardCompactOptions } = $props();

	const tanstackClient = useQueryClient();
	const latestLogs = createQuery(notificationQueryOptions);
	const notification = $derived.by(() => {
		if (latestLogs.isSuccess && latestLogs.data) {
			const notif = latestLogs.data.find(
				(item) => trackerIdToName(item.tracker) === options.collectionName
			);
			return getNotificationStatus(notif);
		}

		return undefined;
	});

	const tracker = createQuery(() => trackerQueryOptions(options.collectionName));
	let interval = $derived.by(() => (tracker.isSuccess ? tracker.data?.interval : undefined));
	let intervalUnit = $derived.by(() =>
		tracker.isSuccess ? tracker.data?.intervalUnit : undefined
	);
	const query = () =>
		createLogsQuery({
			collectionName: options.collectionName,
			interval: interval,
			intervalUnit: intervalUnit
		});

	const refetch = async () => {
		await tanstackClient.refetchQueries(notificationRefetchOptions());
	};
</script>

<section
	class={[
		'border-base-300 grid min-h-18 gap-4 rounded-2xl border px-2 py-2',
		notification?.show ? 'bg-error/15 outline-error/30 outline' : 'bg-white/70'
	]}
>
	<div class="flex items-center">
		<a href={options.route} class="flex grow items-center gap-4">
			<options.icon class="size-9 opacity-75" />
			<div>
				<p class="text-xl font-bold">
					{options.title}
				</p>
				{#if latestLogs.isPending && !latestLogs.data}
					<div class="custom-loader"></div>
				{/if}
				{#if latestLogs.error}
					An error has occurred:
					{latestLogs.error.message}
				{/if}
				{#if latestLogs.isSuccess && latestLogs.data?.length > 0}
					{#if notification?.show}
						{#if notification.level === 'overdue'}
							<span class="text-error font-semibold tracking-tight">Overdue</span>
						{:else if notification.level === 'due'}
							<span class="text-error font-semibold tracking-tight">Due</span>
						{/if}
					{:else}
						<span class="text-neutral/70 font-medium tracking-tight"
							>Due {dayjs(notification?.next).fromNow()}</span
						>
					{/if}
				{/if}
			</div>
		</a>
		<div class="flex h-full items-center">
			<ActionButton
				{query}
				{refetch}
				text={options.button.text}
				compact={true}
				color={'primary'}
				icon={MaterialSymbolsCheck}
				rounded="2xl"
			/>
		</div>
	</div>
</section>
