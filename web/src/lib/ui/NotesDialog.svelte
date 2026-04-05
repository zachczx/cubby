<script lang="ts">
	import Dialog from '$lib/ui/Dialog.svelte';

	let {
		open = $bindable(false),
		isLoading = false,
		onsave
	}: {
		open?: boolean;
		isLoading?: boolean;
		onsave: (notes: string) => void;
	} = $props();

	let notes = $state('');

	export function prefill(value: string) {
		notes = value;
	}
</script>

<Dialog bind:open title="Notes">
	<div class="grid gap-4">
		<textarea
			class="textarea textarea-bordered w-full resize-none"
			rows="4"
			placeholder="Add notes about this workout..."
			bind:value={notes}
		></textarea>
		<button
			class="btn btn-primary btn-lg w-full rounded-full"
			disabled={isLoading}
			onclick={() => onsave(notes)}
		>
			{#if isLoading}<span class="loading loading-spinner loading-sm"></span>{/if}
			Save
		</button>
	</div>
</Dialog>
