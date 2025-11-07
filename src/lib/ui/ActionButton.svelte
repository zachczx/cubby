<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { MouseEventHandler } from 'svelte/elements';

	let {
		status,
		handler,
		text
	}: { status: ButtonState; handler: MouseEventHandler<HTMLButtonElement>; text: string } =
		$props();
</script>

<button
	class={[
		'btn btn-lg flex w-full items-center gap-2 rounded-full',
		status === 'default' && 'btn-primary',
		status === 'loading' && 'btn-primary',
		status === 'success' && 'btn-success'
	]}
	onclick={handler}
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
