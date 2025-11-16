<script lang="ts">
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import type { RecordModel } from 'pocketbase';
	import type { Component } from 'svelte';
	import ActionButton from './ActionButton.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsAdd from '$lib/assets/svg/MaterialSymbolsAdd.svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import dayjs from 'dayjs';
	import { createQuery, useQueryClient, type CreateQueryResult } from '@tanstack/svelte-query';
	import {
		createLogsQuery,
		logsQueryOptions,
		logsRefetchOptions,
		trackerQueryOptions
	} from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';

	let { options }: { options: Options } = $props();

	interface Options {
		title: string;
		collectionName: CollectionName;
		icon: Component;
		route: string;
		button: {
			status: ButtonState;
			text: string;
		};
	}

	const tanstackClient = useQueryClient();
	const logs = createQuery(() => logsQueryOptions(options.collectionName));
	const notification = $derived.by(() => getNotificationStatus(logs));

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

	const refetch = async () =>
		await tanstackClient.refetchQueries(logsRefetchOptions(options.collectionName));
</script>

<section
	class={[
		'border-base-300 grid min-h-18 gap-4 rounded-3xl border px-4 py-2',
		notification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-base-100'
	]}
>
	<div class="flex items-center">
		<a href={options.route} class="flex grow items-center gap-4">
			<options.icon class="size-9 opacity-75" />
			<div>
				<p class="flex items-center gap-2 text-xl font-bold">
					{options.title}
					<MaterialSymbolsArrowRightAlt class="size-[1em]" />
				</p>
				{#if logs.isPending && !logs.data}
					<div class="custom-loader"></div>
				{/if}
				{#if logs.error}
					An error has occurred:
					{logs.error.message}
				{/if}
				{#if logs.isSuccess}
					{#if notification.show}
						{#if notification.level === 'overdue'}
							<span class="text-error font-medium tracking-tight">Overdue</span>
						{:else if notification.level === 'due'}
							<span class="text-error font-medium tracking-tight">Due</span>
						{/if}
					{:else}
						<span class="font-medium tracking-tight">Due {dayjs(notification.next).fromNow()}</span>
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
				color={notification.level === 'overdue' || notification.level === 'due'
					? 'primary'
					: 'neutral'}
				icon={MaterialSymbolsAdd}
			/>
		</div>
	</div>
</section>
