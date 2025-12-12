<script lang="ts">
	import { QueryClient, useQueryClient } from '@tanstack/svelte-query';
	import { pb } from '../pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from './ArkToaster.svelte';
	import MaterialSymbolsCheck from '../assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsArrowRightAlt from '../assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import type { RecordModel } from 'pocketbase';
	import { getAllLogsQueryKey } from '../queries';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let {
		tracker,
		interval,
		intervalUnit
	}: {
		tracker: TrackerDB | undefined;
		interval: number | undefined;
		intervalUnit: string | undefined;
	} = $props();

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

	const tanstackClient = useQueryClient();
	export const insertNewLogToCache = (newLog: RecordModel) =>
		tanstackClient.setQueryData(getAllLogsQueryKey(), (oldLogs: LogsDB[] | undefined) => {
			if (!oldLogs) return [newLog];
			return [newLog, ...oldLogs];
		});

	async function addHandler() {
		buttonStatus = 'loading';

		const result = await pb.collection('logs').create({
			tracker: tracker?.id,
			user: pb.authStore.record?.id,
			time: dayjs.tz(timestamp, 'Asia/Singapore'),
			interval: interval,
			intervalUnit: intervalUnit
		});

		if (result.id) {
			dialog.close();
			addToast('success', 'Added successfully!');
			buttonStatus = 'success';

			await insertNewLogToCache(result);

			setTimeout(() => {
				buttonStatus = 'default';
			}, 3000);
		}
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
				{tracker?.actionLabel}
			{/if}
		</button>

		<form method="dialog">
			<!-- if there is a button in form, it will close the modal -->
			<div class="flex justify-center">
				<button class="btn btn-lg btn-ghost mt-4 rounded-full">Close</button>
			</div>
		</form>
	</div>
</dialog>
