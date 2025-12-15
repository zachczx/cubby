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
		inviteQueryOptions,
		allTrackersRefetchOptions
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
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import NumberInput from '$lib/ui/NumberInput.svelte';

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

	async function addTracker() {
		if (!user.isSuccess) return;
		spinner = true;
		const clean = { ...inputTrackerDetails, name: camelCaseName, interval: intervalString };

		try {
			const result: TrackerDB = await pb.collection('trackers').create(clean);
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

	let inputTrackerDetails = $state({
		user: pb.authStore.record?.id,
		family: '',
		name: '',
		display: '',
		interval: '',
		intervalUnit: 'day',
		category: 'personal',
		actionLabel: '',
		pinned: false,
		show: true
	});

	$effect(() => {
		if (userOwnedFamily) {
			inputTrackerDetails.family = userOwnedFamily;
		}
	});

	function toCamelCase(text: string): string {
		const words = text.split(' ').map((word, idx) => {
			if (idx === 0) return word.toLowerCase();

			return word.charAt(0).toUpperCase() + word.slice(1);
		});

		return words.join('');
	}

	let camelCaseName = $derived(toCamelCase(inputTrackerDetails.display));
	let intervalString = $state('1');
</script>

<PageWrapper title="Add Tracker" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Add Tracker</h1>

		<form class="grid w-full content-start gap-4" onsubmit={() => addTracker()}>
			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Display Name</legend>
				<input
					type="text"
					name="display"
					placeholder="Give your new tracker a name"
					bind:value={inputTrackerDetails.display}
					class="input input-lg w-full"
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Category</legend>

				<SegmentedControl items={3}>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.category}
							value="personal"
							name="category"
						/>Personal
					</label>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.category}
							value="household"
							name="category"
						/>Household
					</label>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.category}
							value="pet"
							name="category"
						/>Pet
					</label>
				</SegmentedControl>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Interval</legend>
				<NumberInput bind:value={intervalString} />
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Frequency</legend>

				<SegmentedControl items={3}>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.intervalUnit}
							value="day"
							name="intervalUnit"
						/>Day
					</label>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.intervalUnit}
							value="month"
							name="intervalUnit"
						/>Month
					</label>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.intervalUnit}
							value="year"
							name="intervalUnit"
						/>Year
					</label>
				</SegmentedControl>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50"
					>Button Action Label (optional)</legend
				>
				<input
					type="text"
					name="actionLabel"
					placeholder="E.g., Washed, Fed, Paid."
					bind:value={inputTrackerDetails.actionLabel}
					class="input input-lg w-full"
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Pinned</legend>
				<SegmentedControl items={2}>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.pinned}
							value={true}
							name="pinned"
						/>Yes
					</label>
					<label>
						<input
							type="radio"
							bind:group={inputTrackerDetails.pinned}
							value={false}
							name="pinned"
						/>No
					</label>
				</SegmentedControl>
			</fieldset>
			<button class="btn btn-primary btn-lg mt-8 w-full rounded-full">Save</button>
		</form>
	</div>
</PageWrapper>
