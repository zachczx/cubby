<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createUserQueryOptions,
		createUserRefetchOptions,
		createVacationQueryOptions,
		createVacationRefetchOptions
	} from '$lib/queries';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import type { ChangeEventHandler } from 'svelte/elements';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(createUserQueryOptions);

	const vacations = createQuery(createVacationQueryOptions);

	const tanstackClient = useQueryClient();

	let towelIntervalDays = $derived.by(() =>
		user.isSuccess ? user.data?.towelIntervalDays : undefined
	);
	let sprayIntervalDays = $derived.by(() =>
		user.isSuccess ? user.data?.sprayIntervalDays : undefined
	);
	let gummyIntervalDays = $derived.by(() =>
		user.isSuccess ? user.data?.gummyIntervalDays : undefined
	);
	let doggoChewableIntervalMonths = $derived.by(() =>
		user.isSuccess ? user.data?.doggoChewableIntervalMonths : undefined
	);
	let doggoBathIntervalDays = $derived.by(() =>
		user.isSuccess ? user.data?.doggoBathIntervalDays : undefined
	);

	let sound = $derived.by(() => (user.isSuccess ? user.data?.sound : undefined));

	let spinner = $state(false);

	let currentTab = $state('settings');
	let param = $derived(page.url.searchParams.get('p'));

	$effect(() => {
		if (param === 'vacations') {
			currentTab = 'vacations';
		}
	});

	let vacationStart = $state('');
	let vacationEnd = $state('');
	let vacationsModal = $state() as HTMLDialogElement;

	function formatTime(startDateTime: string, endDateTime: string) {
		if (!startDateTime || !endDateTime) return;

		const s = dayjs(startDateTime);
		const e = dayjs(endDateTime);

		if (s.get('month') === e.get('month')) {
			if (s.get('date') === e.get('date')) {
				return s.format('D MMM YYYY');
			}
			return s.format('D') + ' – ' + e.format('D MMM YYYY');
		}
		return s.format('D MMM YYYY') + ' – ' + e.format('D MMM YYYY');
	}

	async function onchange(evt: Event) {
		const target = evt.target;

		if (target instanceof HTMLInputElement && pb.authStore.record?.id) {
			const name = target.name;

			try {
				const data: { [key: string]: string | boolean } = {};

				if (name === 'sound') {
					data[name] = target.checked ? true : false;
				} else {
					data[name] = target.value;
				}

				const response = await pb.collection('users').update(pb.authStore.record.id, data);

				if (!response.status) {
					addToast('success', 'Saved!');
					await tanstackClient.refetchQueries(createUserRefetchOptions());
				}
			} catch (err) {
				addToast('error', 'Error saving!');
			}
		} else {
			addToast('error', 'Error saving!');
		}
	}
</script>

<PageWrapper title="Intervals" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Profile</h1>
			<div class="grid w-full content-start">
				<div class="px-1">
					<form class="border-b-base-300 grid border-b pt-2 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Nasal Spray</legend>

						<SegmentedControl items={5}>
							<label>
								<input
									type="radio"
									bind:group={sprayIntervalDays}
									value={1}
									name="sprayIntervalDays"
								/>1
							</label>
							<label>
								<input
									type="radio"
									bind:group={sprayIntervalDays}
									value={2}
									name="sprayIntervalDays"
								/>2
							</label>
							<label>
								<input
									type="radio"
									bind:group={sprayIntervalDays}
									value={3}
									name="sprayIntervalDays"
								/>3
							</label>
							<label>
								<input
									type="radio"
									bind:group={sprayIntervalDays}
									value={4}
									name="sprayIntervalDays"
								/>4
							</label>
							<label>
								<input
									type="radio"
									bind:group={sprayIntervalDays}
									value={5}
									name="sprayIntervalDays"
								/>5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Gummy</legend>

						<SegmentedControl items={5}>
							<label>
								<input
									type="radio"
									bind:group={gummyIntervalDays}
									value={1}
									name="gummyIntervalDays"
								/>1
							</label>
							<label>
								<input
									type="radio"
									bind:group={gummyIntervalDays}
									value={2}
									name="gummyIntervalDays"
								/>2
							</label>
							<label>
								<input
									type="radio"
									bind:group={gummyIntervalDays}
									value={3}
									name="gummyIntervalDays"
								/>3
							</label>
							<label>
								<input
									type="radio"
									bind:group={gummyIntervalDays}
									value={4}
									name="gummyIntervalDays"
								/>4
							</label>
							<label>
								<input
									type="radio"
									bind:group={gummyIntervalDays}
									value={5}
									name="gummyIntervalDays"
								/>5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Towel Wash</legend>

						<SegmentedControl items={5}>
							<label>
								<input
									type="radio"
									bind:group={towelIntervalDays}
									value={1}
									name="towelIntervalDays"
								/>1
							</label>
							<label>
								<input
									type="radio"
									bind:group={towelIntervalDays}
									value={2}
									name="towelIntervalDays"
								/>2
							</label>
							<label>
								<input
									type="radio"
									bind:group={towelIntervalDays}
									value={3}
									name="towelIntervalDays"
								/>3
							</label>
							<label>
								<input
									type="radio"
									bind:group={towelIntervalDays}
									value={4}
									name="towelIntervalDays"
								/>4
							</label>
							<label>
								<input
									type="radio"
									bind:group={towelIntervalDays}
									value={5}
									name="towelIntervalDays"
								/>5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Months Per Pet Chewable</legend>

						<SegmentedControl items={3}>
							<label>
								<input
									type="radio"
									bind:group={doggoChewableIntervalMonths}
									value={1}
									name="doggoChewableIntervalMonths"
								/>Monthly
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoChewableIntervalMonths}
									value={3}
									name="doggoChewableIntervalMonths"
								/>Quarterly
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoChewableIntervalMonths}
									value={6}
									name="doggoChewableIntervalMonths"
								/>Half-Yearly
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Pet Bath</legend>

						<SegmentedControl items={5}>
							<label>
								<input
									type="radio"
									bind:group={doggoBathIntervalDays}
									value={3}
									name="doggoBathIntervalDays"
								/>3
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoBathIntervalDays}
									value={5}
									name="doggoBathIntervalDays"
								/>5
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoBathIntervalDays}
									value={7}
									name="doggoBathIntervalDays"
								/>7
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoBathIntervalDays}
									value={14}
									name="doggoBathIntervalDays"
								/>14
							</label>
							<label>
								<input
									type="radio"
									bind:group={doggoBathIntervalDays}
									value={30}
									name="doggoBathIntervalDays"
								/>30
							</label>
						</SegmentedControl>
					</form>
				</div>
			</div>
		</div>
	</div>
</PageWrapper>
