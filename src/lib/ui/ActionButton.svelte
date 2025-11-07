<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { RecordModel, RecordOptions, RecordService } from 'pocketbase';
	import type { MouseEventHandler } from 'svelte/elements';
	import { addToast } from './ArkToaster.svelte';
	import type { RefetchQueryFilters } from '@tanstack/svelte-query';

	let {
		text,
		query,
		refetch
	}: {
		text: string;
		query: () => Promise<RecordModel>;
		refetch: () => Promise<void>;
	} = $props();

	let status: ButtonState = $state('default');

	async function addHandler() {
		status = 'loading';

		const result = await query();

		if (result.id) {
			addToast('success', 'Added successfully!');
			status = 'success';

			setTimeout(() => {
				status = 'default';
			}, 3000);
		}

		await refetch();
	}
</script>

<button
	class={[
		'btn btn-lg flex w-full items-center gap-2 rounded-full',
		status === 'default' && 'btn-primary',
		status === 'loading' && 'btn-primary',
		status === 'success' && 'btn-success'
	]}
	onclick={addHandler}
>
	{#if status === 'success'}
		<MaterialSymbolsCheck class="size-6" />Added!
	{/if}
	{#if status === 'loading'}
		<span class="loading loading-spinner loading-md"></span>
	{/if}
	{#if status !== 'success' && status !== 'loading'}
		{#if text}
			{text}
		{:else}
			Submit
		{/if}
	{/if}
</button>
