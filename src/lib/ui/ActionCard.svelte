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
	import MaterialSymbolsWarning from '$lib/assets/svg/MaterialSymbolsWarning.svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';
	import { onMount } from 'svelte';
	import { error } from '@sveltejs/kit';
	import { pb } from '$lib/pb';

	let { options }: { options: ActionCardOptions } = $props();

	onMount(() => {
		if (!options.title) {
			error(500);
		}
	});

	let size = $derived(options.size ?? 'default');

	const tanstackClient = useQueryClient();
	const latestLogs = createQuery(notificationQueryOptions);
	const trackers = createQuery(allTrackersQueryOptions);
	let tracker = $derived.by(() =>
		trackers.data?.find((tracker) => tracker.name === options.tracker?.name)
	);
	const notification = $derived.by(() => {
		if (latestLogs.isSuccess && latestLogs.data) {
			const notif = latestLogs.data.find(
				(item) => trackerIdToName(item.tracker, trackers.data) === options.tracker?.name
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

{#if size === 'list'}
	<a
		href={options.route}
		class={['grid gap-4 px-2 py-2', !options.lastChild && 'border-b-base-300/50 border-b']}
	>
		<div class="flex items-center">
			<div class="flex grow items-center gap-4">
				<div class="flex size-9 items-center justify-center">
					<options.icon class="size-7 opacity-75" />
				</div>
				<div class="grid gap-0">
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{options.title}
						{#if options.tracker?.expand?.family?.expand?.owner}
							{@const owner = options.tracker?.expand?.family?.expand?.owner}
							{#if owner.id !== pb.authStore.record?.id}
								<span
									class={[
										'flex aspect-square size-5 items-center justify-center rounded-full bg-blue-200 p-0 text-sm font-bold text-blue-600',
										'btn-neutral opacity-75'
									]}
								>
									{owner.name.charAt(0)}
								</span>
							{/if}
						{/if}
					</div>
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
								<span class="text-error font-bold tracking-tight">Overdue</span>
							{:else if notification.level === 'due'}
								<span class="text-warning font-medium tracking-tight">Due</span>
							{/if}
						{:else}
							<span class="text-neutral/70 text-sm font-medium tracking-tight"
								>Due {dayjs(notification?.next).fromNow()}</span
							>
						{/if}
					{/if}
				</div>
			</div>

			<div class="flex h-full items-center">
				<MaterialSymbolsChevronRight class="size-5 opacity-75" />
			</div>
		</div>
	</a>
{/if}

{#if size === 'compact'}
	<section class="border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 px-2 py-2">
		<div class="flex items-center">
			<a href={options.route} class="flex grow items-center gap-4">
				<options.icon class="size-9 opacity-75" />
				<div>
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{options.title}
						{#if options.tracker?.expand?.family?.expand?.owner}
							{@const owner = options.tracker?.expand?.family?.expand?.owner}
							{#if owner.id !== pb.authStore.record?.id}
								{@const bgColor =
									'color' in options.tracker ? 'bg-' + options.tracker.color + '-200' : undefined}
								{@const textColor =
									'color' in options.tracker ? 'text-' + options.tracker.color + '-600' : undefined}
								<span
									class={[
										'flex aspect-square size-5 items-center justify-center rounded-full p-0 text-sm font-bold',
										bgColor,
										textColor
									]}
								>
									{owner.name.charAt(0)}
								</span>
							{/if}
						{/if}
					</div>
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
								<span class="text-error font-bold tracking-tight">Overdue</span>
							{:else if notification.level === 'due'}
								<span class="text-warning font-medium tracking-tight">Due</span>
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
	<section class="border-base-300 grid min-h-24 gap-4 rounded-3xl border bg-white/70 p-4">
		<a href={options.route} class="flex items-center">
			<div class="flex grow items-center gap-4">
				<options.icon class="size-12 opacity-75" />
				<div>
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{options.title}
						{#if options.tracker?.expand?.family?.expand?.owner}
							{@const owner = options.tracker?.expand?.family?.expand?.owner}
							{#if owner.id !== pb.authStore.record?.id}
								{@const bgColor =
									'color' in options.tracker ? 'bg-' + options.tracker.color + '-200' : undefined}
								{@const textColor =
									'color' in options.tracker ? 'text-' + options.tracker.color + '-600' : undefined}
								<span
									class={[
										'flex aspect-square size-5 items-center justify-center rounded-full p-0 text-sm font-bold',
										bgColor,
										textColor
									]}
								>
									{owner.name.charAt(0)}
								</span>
							{/if}
						{/if}
					</div>
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
								<span class="text-error font-bold tracking-tight">Overdue</span>
							{:else if notification.level === 'due'}
								<span class="text-warning font-medium tracking-tight">Due</span>
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
