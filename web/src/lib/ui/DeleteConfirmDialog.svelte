<script lang="ts">
	import type { Snippet } from 'svelte';
	import Icon from '@iconify/svelte';
	import Dialog from '$lib/ui/Dialog.svelte';

	let {
		open = $bindable(false),
		title,
		confirmLabel = 'Delete',
		isLoading = false,
		onconfirm,
		children
	}: {
		open?: boolean;
		title: string;
		confirmLabel?: string;
		isLoading?: boolean;
		onconfirm: () => void;
		children: Snippet;
	} = $props();
</script>

<Dialog bind:open {title}>
	<div class="grid gap-8">
		<div
			class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center rounded-full"
		>
			<Icon icon="material-symbols:delete-outline" class="size-10" />
		</div>
		{@render children()}
		<div class="grid gap-4">
			<button class="btn btn-error btn-lg" disabled={isLoading} onclick={onconfirm}>
				{#if isLoading}<span class="loading loading-spinner loading-sm"></span>{/if}
				{confirmLabel}
			</button>
			<button
				class="btn btn-outline btn-neutral btn-lg w-full"
				onclick={() => (open = false)}>Cancel</button
			>
		</div>
	</div>
</Dialog>
