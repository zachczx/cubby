<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		cleanEmail,
		userQueryOptions,
		userRefetchOptions,
		vacationQueryOptions,
		vacationRefetchOptions,
		familyQueryOptions,
		familyRefetchOptions,
		inviteQueryOptions
	} from '$lib/queries';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import type { ChangeEventHandler } from 'svelte/elements';
	import MaterialSymbolsDelete from '$lib/assets/svg/MaterialSymbolsDelete.svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import MaterialSymbolsPerson from '$lib/assets/svg/MaterialSymbolsPerson.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsFamilyGroup from '$lib/assets/svg/MaterialSymbolsFamilyGroup.svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsShare from '$lib/assets/svg/MaterialSymbolsShare.svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';
	import MdiAlertCircle from '$lib/assets/svg/MdiAlertCircle.svelte';
	import MaterialSymbolsMoreVert from '$lib/assets/svg/MaterialSymbolsMoreVert.svelte';
	import MaterialSymbolsPersonRemove from '$lib/assets/svg/MaterialSymbolsPersonRemove.svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(userQueryOptions);
	const families = createQuery(familyQueryOptions);
	const currentInvite = createQuery(inviteQueryOptions);
	const tanstackClient = useQueryClient();

	let userOwnedFamily = $derived.by(() => {
		if (families.isSuccess && families.data) {
			return families.data.find((family) => family.owner === pb.authStore.record?.id)?.id;
		}
	});

	let spinner = $state(false);
	let invited = $state(false);

	let invitee = $state('');

	let modals = $state<HTMLDialogElement[]>([]);

	async function addHandler() {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			// const start = dayjs.tz(vacationStart, 'Asia/Singapore');
			// const end = dayjs.tz(vacationEnd, 'Asia/Singapore');
			// const result = await pb.collection('vacation').create({
			// 	user: pb.authStore.record?.id,
			// 	startDateTime: start,
			// 	endDateTime: end
			// });
			// if (result.id) {
			// 	addToast('success', 'Added successfully!');
			// 	spinner = false;
			// 	await tanstackClient.refetchQueries(vacationRefetchOptions());
			// }
		} catch (err) {
			console.log(err);
		}
	}

	let inputTrackerDetails = $state({
		user: pb.authStore.record?.id,
		family: '',
		name: '',
		display: '',
		interval: 7,
		intervalUnit: 'day',
		category: '',
		actionLabel: '',
		pinned: false,
		show: true
	});

	$effect(() => {
		if (userOwnedFamily) {
			inputTrackerDetails.family = userOwnedFamily;
		}
	});

	$inspect(inputTrackerDetails);
</script>

<PageWrapper title="Add Tracker" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Add Tracker</h1>

		<div class="grid w-full content-start gap-8">
			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Name</legend>
				<input
					type="text"
					name="name"
					bind:value={inputTrackerDetails.name}
					class="input input-lg w-full"
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Display</legend>
				<input
					type="text"
					name="display"
					bind:value={inputTrackerDetails.display}
					class="input input-lg w-full"
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Interval</legend>
				<input
					type="text"
					name="interval"
					bind:value={inputTrackerDetails.interval}
					class="input input-lg w-full"
				/>
			</fieldset>
		</div>
	</div>
</PageWrapper>
