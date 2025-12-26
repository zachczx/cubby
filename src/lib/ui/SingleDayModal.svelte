<script lang="ts">
	import MaterialSymbolsClose from '$lib/assets/svg/MaterialSymbolsClose.svelte';
	import MaterialSymbolsDelete from '$lib/assets/svg/MaterialSymbolsDelete.svelte';
	import MaterialSymbolsEdit from '$lib/assets/svg/MaterialSymbolsEdit.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import { addToast } from './ArkToaster.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=150&enhanced';
	import type { RecordModel } from 'pocketbase';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { getAllLogsQueryKey } from '$lib/queries';

	let {
		modal = $bindable(),
		singleDay
	}: {
		modal: HTMLDialogElement | undefined;
		singleDay: LogsDB[] | undefined;
	} = $props();

	let logs = $derived.by(() => {
		if (!singleDay) return [];

		return singleDay.sort((a, b) => dayjs(a.time).diff(dayjs(b.time)));
	});

	async function deleteHandler(id: string) {
		if (!id) return;

		const success = await pb.collection('logs').delete(id);

		if (success) {
			deleteLogFromCache(id);
			addToast('success', 'Deleted!');
			modal?.close();
		}
	}

	const tanstackClient = useQueryClient();
	const deleteLogFromCache = (deletedLogId: string) =>
		tanstackClient.setQueryData(getAllLogsQueryKey(), (oldLogs: LogsDB[] | undefined) => {
			if (!oldLogs) return [oldLogs];
			return oldLogs.filter((log) => log.id !== deletedLogId);
		});

	let editMode: Record<string, boolean> = $state({});

	async function saveEdit(evt: SubmitEvent, log: LogsDB) {
		evt.preventDefault();

		try {
			const formData = new FormData(evt.currentTarget as HTMLFormElement);
			const newTime = String(formData.get('newTime'));

			const newDateTimeString = `${dayjs(log.time).format('YYYY-MM-DD')}T${newTime}`;
			const newDateTime = dayjs.tz(
				dayjs(newDateTimeString).format('YYYY-MM-DD HH:mm:ss'),
				'Asia/Singapore'
			);
			const response = await pb.collection('logs').update(log.id, {
				time: newDateTime
			});

			if (response.id) {
				modifyLogInCache(response);
				addToast('success', 'Updated!');
				editMode[log.id] = false;
			}
		} catch (err) {
			addToast('error', 'Error updating!');
		}
	}

	const modifyLogInCache = (newLog: RecordModel) =>
		tanstackClient.setQueryData(getAllLogsQueryKey(), (oldLogs: LogsDB[] | undefined) => {
			if (!oldLogs) return [newLog];
			return oldLogs.map((log) => {
				if (log.id === newLog.id) {
					return newLog;
				}

				return log;
			});
		});
</script>

<dialog bind:this={modal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		{#if logs && logs.length > 0}
			{@const theDay = dayjs(logs[0].time).format('ddd, D MMM')}
			<div class="flex items-center">
				<h3 class="grow text-lg font-bold">{theDay}</h3>
				<div class="flex items-center justify-end">
					<form method="dialog">
						<button class="btn btn-ghost -me-2"><MaterialSymbolsClose class="size-5" /></button>
					</form>
				</div>
			</div>
			<div class="grid gap-2 py-4">
				{#each logs as log}
					{@const formatted = dayjs(log.time).format('hh:mma')}
					{#if editMode[log.id]}
						<div
							class="border-base-300 grid min-h-18 content-center gap-2 rounded-2xl border bg-white/70 px-2 py-2"
						>
							<form class="grid w-full content-center" onsubmit={(evt) => saveEdit(evt, log)}>
								<input
									type="time"
									name="newTime"
									value={dayjs(log.time).format('HH:mm')}
									class="input input-lg w-full"
								/>
								<button class="btn btn-primary btn-lg mt-2 w-full rounded-full">
									<!-- <MaterialSymbolsSave class="size-[1.3em]" /> -->
									Save
								</button>
							</form>
							<button class="btn btn-ghost btn-lg w-full" onclick={() => (editMode[log.id] = false)}
								>Cancel</button
							>
						</div>
					{:else}
						<div
							class="border-base-300 grid min-h-18 grid-cols-[1fr_auto] content-center gap-4 rounded-2xl border bg-white/70 px-2 py-2"
						>
							<div class="flex items-center p-2 text-lg font-bold">
								{formatted}
							</div>
							<div class="flex items-center">
								<button class="btn btn-ghost" onclick={() => (editMode[log.id] = !editMode[log.id])}
									><MaterialSymbolsEdit class="size-[1.3em]" /></button
								>
								<button class="btn btn-error btn-soft" onclick={() => deleteHandler(log.id)}
									><MaterialSymbolsDelete class="size-[1.3em]" /></button
								>
							</div>
						</div>
					{/if}
				{/each}
			</div>
		{:else}
			<div class="flex items-center justify-end">
				<form method="dialog">
					<button class="btn btn-ghost -me-2"><MaterialSymbolsClose class="size-5" /></button>
				</form>
			</div>
			<div class="justify-self-center">
				<enhanced:img src={EmptyCorgi} alt="nothing" />
				<p class="text-center text-lg">Nothing here!</p>
			</div>
		{/if}
	</div>
</dialog>
