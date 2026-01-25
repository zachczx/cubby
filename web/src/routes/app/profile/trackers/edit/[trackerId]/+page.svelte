<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		userQueryOptions,
		allTrackersRefetchOptions,
		allTrackersQueryOptions
	} from '$lib/queries';
	import { goto } from '$app/navigation';
	import TrackerForm from '../../TrackerForm.svelte';
	import Icon from '@iconify/svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	let { data } = $props();

	const user = createQuery(userQueryOptions);
	const tanstackClient = useQueryClient();

	const trackersDb = createQuery(allTrackersQueryOptions);
	let currentTracker = $derived.by(() => {
		if (trackersDb.isSuccess && trackersDb.data) {
			return trackersDb.data.find((t) => t.id === data.trackerId);
		}
	});

	let spinner = $state(false);

	async function editTracker(inputTrackerDetails: TrackerInput) {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			const result: TrackerDB = await pb
				.collection('trackers')
				.update(data.trackerId, { ...inputTrackerDetails });
			if (result.id) {
				addToast('success', 'Added successfully!');
				await tanstackClient.refetchQueries(allTrackersRefetchOptions());
				spinner = false;
				goto(`/app/${result.category}`);
			}
		} catch (err) {
			console.log(err);
		}
	}

	async function deleteTracker() {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			const result = await pb.collection('trackers').delete(data.trackerId);
			if (result) {
				addToast('success', 'Deleted successfully!');
				await tanstackClient.refetchQueries(allTrackersRefetchOptions());
				spinner = false;
				goto(`/app/profile/trackers/edit`);
			}
		} catch (err) {
			console.log(err);
		}
	}

	let deleteModal = $state<HTMLDialogElement>();
</script>

<PageWrapper title="Edit Tracker" {pb} largeScreenCenter={true}>
	<div class="grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8">
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Edit Tracker</h1>

		{#if trackersDb.isPending}
			Loading
		{:else if trackersDb.isSuccess && trackersDb.data && currentTracker}
			<TrackerForm onsubmit={editTracker} {currentTracker} />
		{:else}
			Error!
		{/if}

		<button
			class="btn btn-error btn-lg btn-soft mt-4 w-full rounded-full"
			onclick={() => deleteModal?.showModal()}>Delete Tracker</button
		>
	</div>
</PageWrapper>

<dialog bind:this={deleteModal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box grid gap-8">
		<div
			class="bg-primary/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center overflow-hidden rounded-full"
		>
			<Icon icon="material-symbols:delete" class="size-12" />
		</div>
		<h2 class="text-2xl font-bold">Confirm Deletion</h2>

		<ul class="ms-6 list-disc space-y-2">
			<li>You will delete all of your logs in this tracker.</li>
			<li>This action is permanent.</li>
		</ul>
		<div class="grid grid-cols-1 gap-4">
			<button
				class="btn btn-error btn-lg rounded-full"
				onclick={() => {
					spinner = true;
					deleteTracker();
				}}>Confirm Deletion</button
			>
			<form method="dialog" class="">
				<button class="btn btn-outline btn-primary btn-lg w-full rounded-full">Cancel</button>
			</form>
		</div>
	</div>
</dialog>
