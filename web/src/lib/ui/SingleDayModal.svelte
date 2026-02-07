<script lang="ts">
	import Icon from '@iconify/svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import { addToast } from './ArkToaster.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=150&enhanced';
	import type { RecordModel } from 'pocketbase';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { getAllEntriesQueryKey } from '$lib/queries';

	let {
		modal = $bindable(),
		singleDay
	}: {
		modal: HTMLDialogElement | undefined;
		singleDay: EntryDB[] | undefined;
	} = $props();

	let entries = $derived.by(() => {
		if (!singleDay) return [];

		return singleDay.sort((a, b) => dayjs(a.performedAt).diff(dayjs(b.performedAt)));
	});

	async function deleteHandler(id: string) {
		if (!id) return;

		const success = await pb.collection('entries').delete(id);

		if (success) {
			deleteEntryFromCache(id);
			addToast('success', 'Deleted!');
			modal?.close();
		}
	}

	const tanstackClient = useQueryClient();
	const deleteEntryFromCache = (deletedEntryId: string) =>
		tanstackClient.setQueryData(getAllEntriesQueryKey(), (oldEntries: EntryDB[] | undefined) => {
			if (!oldEntries) return [oldEntries];
			return oldEntries.filter((entry) => entry.id !== deletedEntryId);
		});

	let editMode: Record<string, boolean> = $state({});

	async function saveEdit(evt: SubmitEvent, entry: EntryDB) {
		evt.preventDefault();

		try {
			const formData = new FormData(evt.currentTarget as HTMLFormElement);
			const newTime = String(formData.get('newTime'));

			const newDateTimeString = `${dayjs(entry.performedAt).format('YYYY-MM-DD')}T${newTime}`;
			const newDateTime = dayjs.tz(
				dayjs(newDateTimeString).format('YYYY-MM-DD HH:mm:ss'),
				'Asia/Singapore'
			);
			const response = await pb.collection('entries').update(entry.id, {
				time: newDateTime
			});

			if (response.id) {
				modifyEntryInCache(response);
				addToast('success', 'Updated!');
				editMode[entry.id] = false;
			}
		} catch (err) {
			addToast('error', 'Error updating!');
		}
	}

	const modifyEntryInCache = (newEntry: RecordModel) =>
		tanstackClient.setQueryData(getAllEntriesQueryKey(), (oldEntries: EntryDB[] | undefined) => {
			if (!oldEntries) return [newEntry];
			return oldEntries.map((entry) => {
				if (entry.id === newEntry.id) {
					return newEntry;
				}

				return entry;
			});
		});
</script>

<dialog bind:this={modal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		{#if entries && entries.length > 0}
			{@const theDay = dayjs(entries[0].performedAt).format('ddd, D MMM')}
			<div class="flex items-center">
				<h3 class="grow text-lg font-bold">{theDay}</h3>
				<div class="flex items-center justify-end">
					<form method="dialog">
						<button class="btn btn-ghost -me-2"
							><Icon icon="material-symbols:close" class="size-5" /></button
						>
					</form>
				</div>
			</div>
			<div class="grid gap-2 py-4">
				{#each entries as entry}
					{@const formatted = dayjs(entry.performedAt).format('hh:mma')}
					{#if editMode[entry.id]}
						<div
							class="border-base-300 grid min-h-18 content-center gap-2 rounded-2xl border bg-white/70 px-2 py-2"
						>
							<form class="grid w-full content-center" onsubmit={(evt) => saveEdit(evt, entry)}>
								<input
									type="time"
									name="newTime"
									value={dayjs(entry.performedAt).format('HH:mm')}
									class="input input-lg w-full"
								/>
								<button class="btn btn-primary btn-lg mt-2 w-full rounded-full"> Save </button>
							</form>
							<button
								class="btn btn-ghost btn-lg w-full"
								onclick={() => (editMode[entry.id] = false)}>Cancel</button
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
								<button
									class="btn btn-ghost"
									onclick={() => (editMode[entry.id] = !editMode[entry.id])}
									><Icon icon="material-symbols:edit" class="size-[1.3em]" /></button
								>
								<button class="btn btn-error btn-soft" onclick={() => deleteHandler(entry.id)}
									><Icon icon="material-symbols:delete" class="size-[1.3em]" /></button
								>
							</div>
						</div>
					{/if}
				{/each}
			</div>
		{:else}
			<div class="flex items-center justify-end">
				<form method="dialog">
					<button class="btn btn-ghost -me-2"
						><Icon icon="material-symbols:close" class="size-5" /></button
					>
				</form>
			</div>
			<div class="justify-self-center">
				<enhanced:img src={EmptyCorgi} alt="nothing" />
				<p class="text-center text-lg">Nothing here!</p>
			</div>
		{/if}
	</div>
</dialog>
