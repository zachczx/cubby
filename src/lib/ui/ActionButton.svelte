<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { RecordModel, RecordOptions, RecordService } from 'pocketbase';
	import type { MouseEventHandler } from 'svelte/elements';
	import { addToast } from './ArkToaster.svelte';
	import { useQueryClient, type RefetchQueryFilters } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';
	import { getAllLogsQueryKey } from '$lib/queries';
	import { pb } from '$lib/pb';

	let {
		text,
		query,
		compact = false,
		color = 'primary',
		rounded = 'full',
		class: className,
		icon: CustomIcon
	}: {
		text: string | undefined;
		query: () => Promise<RecordModel>;
		compact?: boolean;
		color?: 'primary' | 'neutral';
		rounded?: 'full' | '3xl' | '2xl' | 'xl' | 'lg' | 'md' | 'sm' | 'xs';
		class?: string;
		icon?: Component;
	} = $props();

	let status: ButtonState = $state('default');

	const tanstackClient = useQueryClient();
	export const insertNewLogToCache = (newLog: RecordModel) =>
		tanstackClient.setQueryData(getAllLogsQueryKey(), (oldLogs: LogsDB[] | undefined) => {
			if (!oldLogs) return [newLog];
			return [newLog, ...oldLogs];
		});

	async function addHandler() {
		status = 'loading';

		try {
			const result = await query();
			if (result) {
				addToast('success', 'Added successfully!');
				status = 'success';
				await insertNewLogToCache(result);

				setTimeout(() => {
					status = 'default';
				}, 3000);
			}
		} catch (err) {
			console.log(err);
			status = 'error';
			addToast('error', 'Error creating!');
			setTimeout(() => {
				status = 'default';
			}, 3000);
		}
	}

	function getButtonClasses() {
		return [
			'btn btn-lg flex items-center gap-2',
			`rounded-${rounded}`,
			(status === 'default' || status === 'loading') && `btn-${color}`,
			status === 'success' && 'btn-success',
			status === 'error' && 'btn-error',
			className
		];
	}
</script>

{#if !compact}
	<button
		class={['btn btn-lg flex w-full items-center gap-2', getButtonClasses()]}
		onclick={addHandler}
	>
		{#if status === 'success'}
			<MaterialSymbolsCheck class="size-6" />Added!
		{/if}
		{#if status === 'loading'}
			<span class="loading loading-spinner loading-md"></span>
		{/if}
		{#if status === 'error'}
			Error!
		{/if}
		{#if status !== 'success' && status !== 'loading' && status !== 'error'}
			{#if text}
				{text}
			{:else}
				Submit
			{/if}
		{/if}
	</button>
{:else}
	<button
		class={['btn btn-lg flex aspect-square w-full items-center gap-2 p-0', getButtonClasses()]}
		onclick={addHandler}
	>
		{#if status === 'loading'}
			<span class="loading loading-spinner loading-md"></span>
		{:else if CustomIcon && status !== 'success'}
			<CustomIcon class="size-8" />
		{:else}
			<MaterialSymbolsCheck class="size-8" />
		{/if}
	</button>
{/if}
