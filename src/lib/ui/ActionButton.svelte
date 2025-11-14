<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { RecordModel, RecordOptions, RecordService } from 'pocketbase';
	import type { MouseEventHandler } from 'svelte/elements';
	import { addToast } from './ArkToaster.svelte';
	import type { RefetchQueryFilters } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';

	let {
		text,
		query,
		refetch,
		compact = false,
		color = 'primary',
		icon: CustomIcon
	}: {
		text: string;
		query: () => Promise<RecordModel>;
		refetch: () => Promise<void>;
		compact?: boolean;
		color?: 'primary' | 'neutral';
		icon?: Component;
	} = $props();

	let status: ButtonState = $state('default');

	async function addHandler() {
		status = 'loading';

		try {
			const result = await query();

			if (result.id) {
				addToast('success', 'Added successfully!');
				status = 'success';

				setTimeout(() => {
					status = 'default';
				}, 3000);
			}

			await refetch();
		} catch (err) {
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
			'btn btn-lg flex w-full items-center gap-2 rounded-full',
			status === 'default' && color === 'primary' && 'btn-primary',
			status === 'loading' && color === 'primary' && 'btn-primary',
			status === 'default' && color === 'neutral' && 'btn-neutral',
			status === 'loading' && color === 'neutral' && 'btn-neutral',
			status === 'success' && 'btn-success',
			status === 'error' && 'btn-error'
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
			'btn btn-lg flex aspect-square w-full items-center gap-2 rounded-full p-0',
			status === 'default' && color === 'primary' && 'btn-primary',
			status === 'loading' && color === 'primary' && 'btn-primary',
			status === 'default' && color === 'neutral' && 'btn-neutral',
			status === 'loading' && color === 'neutral' && 'btn-neutral',
			status === 'success' && 'btn-success'
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
