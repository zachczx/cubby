<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { userQueryOptions, familyQueryOptions, allTrackersRefetchOptions } from '$lib/queries';
	import { goto } from '$app/navigation';
	import TrackerForm from '../TrackerForm.svelte';
	import { api } from '$lib/api';
	import { router } from '$lib/routes';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(userQueryOptions);
	const tanstackClient = useQueryClient();
	const families = createQuery(familyQueryOptions);
	let userOwnedFamily = $derived.by(() => {
		if (families.isSuccess && families.data) {
			return families.data.find((family) => family.isOwner)?.id;
		}
	});

	async function addTracker(inputTrackerDetails: TrackerInput) {
		if (!user.isSuccess) return;

		try {
			const formData = {
				...inputTrackerDetails,
				interval: Number(inputTrackerDetails.interval),
				family: userOwnedFamily
			};

			const response = await api
				.post('trackers', {
					body: JSON.stringify(formData)
				})
				.json<string>();

			if (response) {
				addToast('success', 'Added successfully!');
				await tanstackClient.refetchQueries(allTrackersRefetchOptions());
				goto(router.tracker(response));
			}
		} catch (err) {
			console.log(err);
			addToast('error', 'Failed to add tracker');
		}
	}
</script>

<PageWrapper title="New Tracker" largeScreenCenter={true} focusedScreen={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Add Tracker</h1>

		<TrackerForm onsubmit={addTracker} />
	</div>
</PageWrapper>
