<script lang="ts">
	import ActionButton from './ActionButton.svelte';
	import dayjs from 'dayjs';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		allTrackersQueryOptions,
		createLogsQuery,
		notificationQueryOptions,
		notificationRefetchOptions,
		trackerIdToName
	} from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';

	let { options }: { options: ActionCardOptions } = $props();

	let size = $derived(options.size ?? 'default');

	const tanstackClient = useQueryClient();
	const latestLogs = createQuery(notificationQueryOptions);
	const trackers = createQuery(allTrackersQueryOptions);
	let tracker = $derived.by(() =>
		trackers.data?.find((tracker) => tracker.name === options.collectionName)
	);
	const notification = $derived.by(() => {
		if (latestLogs.isSuccess && latestLogs.data) {
			const notif = latestLogs.data.find(
				(item) => trackerIdToName(item.tracker, trackers.data) === options.collectionName
			);
			return getTrackerStatus(notif);
		}

		return undefined;
	});

	let interval = $derived(tracker?.interval);
	let intervalUnit = $derived(tracker?.intervalUnit);
	const query = () =>
		createLogsQuery({
			trackerId: tracker?.id ?? '',
			interval: interval,
			intervalUnit: intervalUnit
		});

	const refetch = async () => {
		await tanstackClient.refetchQueries(notificationRefetchOptions());
	};
</script>

{#if size === 'compact'}
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
{/if}

{#if size === 'default'}
	<section
		class={[
			'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
			notification?.show ? 'bg-error/15 outline-error/30 outline' : 'bg-white/70'
		]}
	>
		<a href={options.route} class="flex items-center">
			<div class="flex grow items-center gap-4">
				<options.icon class="size-12 opacity-75" />
				<div>
					<p class="text-xl font-bold">{options.title}</p>
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
			</div>
			<div class="flex h-full items-center">
				<div class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75">
					<MaterialSymbolsChevronRight class="size-6" />
				</div>
			</div>
		</a>
		<ActionButton {query} {refetch} text={options.button.text} />
	</section>
{/if}
