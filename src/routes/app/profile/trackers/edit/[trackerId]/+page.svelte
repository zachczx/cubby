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
		familyQueryOptions,
		allTrackersRefetchOptions,
		allTrackersQueryOptions
	} from '$lib/queries';
	import { goto } from '$app/navigation';
	import TrackerForm from '../../TrackerForm.svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	let { data } = $props();

	const user = createQuery(userQueryOptions);
	const tanstackClient = useQueryClient();
	const families = createQuery(familyQueryOptions);
	let userOwnedFamily = $derived.by(() => {
		if (families.isSuccess && families.data) {
			return families.data.find((family) => family.owner === pb.authStore.record?.id)?.id;
		}
	});
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
</script>

<PageWrapper title="New Tracker" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Add Tracker</h1>

		{#if trackersDb.isPending}
			Loading
		{:else if trackersDb.isSuccess && trackersDb.data && currentTracker}
			<TrackerForm onsubmit={editTracker} {currentTracker} />
		{:else}
			Error!
		{/if}
	</div>
</PageWrapper>
