<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import dayjs from 'dayjs';

	let {
		modal,
		singleDay
	}: {
		modal: HTMLDialogElement | undefined;
		singleDay: CollectionsTrackingTime[] | undefined;
	} = $props();
</script>

<dialog bind:this={modal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		{#if singleDay && singleDay.length > 0}
			{@const theDay = dayjs(singleDay[0].time).format('D MMM YYYY')}
			<h3 class="text-lg font-bold">{theDay}</h3>
			<div class="py-4">
				<ul class="list-disc">
					{#each singleDay as day}
						{@const formatted = dayjs(day.time).format('h:mma')}
						<li class="flex items-center gap-2">
							<MaterialSymbolsCheck class="text-success size-[1.3em]" />
							{formatted}
						</li>
					{/each}
				</ul>
			</div>
		{:else}
			<h3 class="text-lg font-bold">Nothing here!</h3>
		{/if}
		<div class="modal-action">
			<form method="dialog">
				<button class="btn">Close</button>
			</form>
		</div>
	</div>
</dialog>
