<script lang="ts">
	import { QueryClient, useQueryClient } from '@tanstack/svelte-query';
	import { pb } from './pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from './ui/ArkToaster.svelte';
	import { createTowelRefetchOptions } from './queries';
	import MaterialSymbolsCheck from './assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsArrowRightAlt from './assets/svg/MaterialSymbolsArrowRightAlt.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let {
		collectionName,
		tanstackClient
	}: { collectionName: 'towel' | 'spray' | 'gummy'; tanstackClient: QueryClient } = $props();

	let buttonStatus: 'default' | 'loading' | 'success' = $state('default');

	let dialog = $state() as HTMLDialogElement;
	let date = $state('');
	let time = $state('');
	let timestamp = $derived.by(() => {
		if (!date || !time) return '';

		const ts = dayjs(date + 'T' + time);

		const tzTime = dayjs.tz(ts, 'Asia/Singapore');

		if (tzTime.diff(dayjs(), 'hours') > 0) {
			const now = dayjs();
			return now;
		}

		return ts;
	});

	async function addHandler() {
		buttonStatus = 'loading';

		const result = await pb.collection(collectionName).create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(timestamp, 'Asia/Singapore')
		});

		if (result.id) {
			dialog.close();
			addToast('success', 'Added successfully!');
			buttonStatus = 'success';

			setTimeout(() => {
				buttonStatus = 'default';
			}, 3000);
		}

		await tanstackClient.refetchQueries(createTowelRefetchOptions());
	}
</script>

<button
	class="text-base-content/70 flex cursor-pointer items-center gap-1 py-2 font-medium"
	aria-label="custom time"
	onclick={() => dialog.showModal()}
>
	Custom Time<MaterialSymbolsArrowRightAlt class="size-5 opacity-[0.7]" /></button
>

<dialog bind:this={dialog} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<h3 class="mb-8 text-xl font-bold">Add Custom Date</h3>

		<fieldset class="fieldset mb-8">
			<div class="grid grid-cols-2 gap-2">
				<input type="date" class="input input-lg w-full" bind:value={date} required />
				<input type="time" class="input input-lg w-full" bind:value={time} required />
			</div>
		</fieldset>

		<button
			class={[
				'btn btn-lg flex w-full items-center gap-2 rounded-full',
				buttonStatus === 'default' && 'btn-primary',
				buttonStatus === 'loading' && 'btn-primary',
				buttonStatus === 'success' && 'btn-success'
			]}
			onclick={addHandler}
		>
			{#if buttonStatus === 'success'}
				<MaterialSymbolsCheck class="size-6" />Added!
			{:else if buttonStatus === 'loading'}
				<span class="loading loading-spinner loading-md"></span>
			{:else}
				{#if collectionName === 'towel'}
					Add Wash
				{/if}
				{#if collectionName === 'spray'}
					Add Spray
				{/if}
				{#if collectionName === 'gummy'}
					Add Gummy
				{/if}
			{/if}
		</button>

		<form method="dialog">
			<!-- if there is a button in form, it will close the modal -->
			<div class="flex justify-center">
				<button class="btn btn-lg btn-ghost rounded-full">Close</button>
			</div>
		</form>
	</div>
</dialog>
