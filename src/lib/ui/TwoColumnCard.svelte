<script lang="ts">
	import type { Snippet } from 'svelte';

	interface ContentPart {
		title: string;
		content: Snippet;
	}

	let {
		left,
		leftTitle,
		right,
		rightTitle
	}: { left: Snippet; leftTitle: string; right: Snippet; rightTitle: string } = $props();
</script>

<div class="w-full">
	<div
		class="text-md border-base-content/5 grid grid-cols-2 content-center gap-4 rounded-lg border shadow"
	>
		<div class="border-r-base-content/15 grid justify-items-center border-r p-4">
			<div>{leftTitle}</div>
			<div class="flex min-h-20 items-center justify-center gap-4 text-center text-2xl font-bold">
				{@render left()}
			</div>
		</div>
		<div class="grid justify-items-center p-4">
			<div>{rightTitle}</div>
			<div class="text-center text-2xl font-bold">
				{@render right()}
			</div>
		</div>
	</div>
</div>
