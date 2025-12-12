<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { RecordModel, RecordOptions, RecordService } from 'pocketbase';
	import type { MouseEventHandler } from 'svelte/elements';
	import { addToast } from './ArkToaster.svelte';
	import { useQueryClient, type RefetchQueryFilters } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';
	import { allLogsRefetchOptions } from '$lib/queries';
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
		tanstackClient.setQueryData(
			[pb.authStore.record?.id, 'logs-all'],
			(oldLogs: LogsDB[] | undefined) => {
				if (!oldLogs) return [newLog];
				return [newLog, ...oldLogs];
			}
		);

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
</script>

{#if !compact}
	<button
		class={[
			'btn btn-lg flex w-full items-center gap-2',
			rounded === '3xl' && 'rounded-3xl',
			rounded === '2xl' && 'rounded-2xl',
			rounded === 'xl' && 'rounded-xl',
			rounded === 'lg' && 'rounded-lg',
			rounded === 'md' && 'rounded-md',
			rounded === 'sm' && 'rounded-sm',
			rounded === 'xs' && 'rounded-xs',
			rounded === 'full' && 'rounded-full',
			status === 'default' && color === 'primary' && 'btn-primary',
			status === 'loading' && color === 'primary' && 'btn-primary',
			status === 'default' && color === 'neutral' && 'btn-neutral',
			status === 'loading' && color === 'neutral' && 'btn-neutral',
			status === 'success' && 'btn-success',
			status === 'error' && 'btn-error',
			className
		]}
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
		class={[
			'btn btn-lg flex aspect-square w-full items-center gap-2 p-0',
			rounded === '3xl' && 'rounded-3xl',
			rounded === '2xl' && 'rounded-2xl',
			rounded === 'xl' && 'rounded-xl',
			rounded === 'lg' && 'rounded-lg',
			rounded === 'md' && 'rounded-md',
			rounded === 'sm' && 'rounded-sm',
			rounded === 'xs' && 'rounded-xs',
			rounded === 'full' && 'rounded-full',
			status === 'default' && color === 'primary' && 'btn-primary',
			status === 'loading' && color === 'primary' && 'btn-primary',
			status === 'default' && color === 'neutral' && 'btn-neutral',
			status === 'loading' && color === 'neutral' && 'btn-neutral',
			status === 'success' && 'btn-success',
			className
		]}
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
