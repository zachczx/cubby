<script lang="ts">
	import { Dialog } from 'bits-ui';
	import type { Snippet } from 'svelte';

	let {
		open = $bindable(false),
		title,
		children,
		onOpenChange
	}: {
		open?: boolean;
		title?: string;
		children: Snippet;
		onOpenChange?: (open: boolean) => void;
	} = $props();
</script>

<Dialog.Root bind:open {onOpenChange}>
	<Dialog.Portal>
		<Dialog.Overlay
			class="fixed inset-0 z-50 bg-black/40 backdrop-blur-[2px] data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
		/>
		<Dialog.Content
			class="bg-base-100 fixed z-50 grid w-full max-w-lg gap-4 rounded-t-2xl p-6 shadow-lg duration-200 max-sm:inset-x-0 max-sm:bottom-0 sm:top-1/2 sm:left-1/2 sm:-translate-x-1/2 sm:-translate-y-1/2 sm:rounded-2xl data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 max-sm:data-[state=closed]:slide-out-to-bottom max-sm:data-[state=open]:slide-in-from-bottom sm:data-[state=closed]:zoom-out-95 sm:data-[state=open]:zoom-in-95"
		>
			{#if title}
				<div class="flex items-center justify-between">
					<Dialog.Title class="text-lg font-bold">{title}</Dialog.Title>
					<Dialog.Close class="btn btn-ghost btn-sm btn-square">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="20"
							height="20"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M18 6 6 18" />
							<path d="m6 6 12 12" />
						</svg>
					</Dialog.Close>
				</div>
			{/if}
			{@render children()}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
