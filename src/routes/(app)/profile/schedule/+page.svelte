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
		userRefetchOptions,
		vacationQueryOptions,
		vacationRefetchOptions
	} from '$lib/queries';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import type { ChangeEventHandler } from 'svelte/elements';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(userQueryOptions);

	const vacations = createQuery(vacationQueryOptions);

	const tanstackClient = useQueryClient();

	const trackers = createQuery(() => ({
		queryKey: ['profile', pb.authStore?.record?.id],
		queryFn: async (): Promise<TrackerDB[]> =>
			await pb.collection('trackers').getFullList({ filter: `user="${pb.authStore?.record?.id}"` })
	}));

	let towelInterval = $derived.by(() => {
		if (trackers.isSuccess) {
			const rec = trackers.data?.find((item) => item.name === 'towel');
			if (rec) {
				return rec.interval;
			}
		}
	});

	let sprayInterval = $derived.by(() => {
		if (trackers.isSuccess) {
			const rec = trackers.data?.find((item) => item.name === 'spray');
			if (rec) {
				return rec.interval;
			}
		}
	});

	let gummyInterval = $derived.by(() => {
		if (trackers.isSuccess) {
			const rec = trackers.data?.find((item) => item.name === 'gummy');
			if (rec) {
				return rec.interval;
			}
		}
	});

	let petChewableInterval = $derived.by(() => {
		if (trackers.isSuccess) {
			const rec = trackers.data?.find((item) => item.name === 'petChewable');
			if (rec) {
				return rec.interval;
			}
		}
	});

	let petBathInterval = $derived.by(() => {
		if (trackers.isSuccess) {
			const rec = trackers.data?.find((item) => item.name === 'petBath');
			if (rec) {
				return rec.interval;
			}
		}
	});

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
			const recordId = trackers.data?.find((item) => item.name === name)?.id ?? '';

			try {
				const data = {
					interval: target.value
				};

				const response = await pb.collection('trackers').update(recordId, data);

				if (!response.status) {
					addToast('success', 'Saved!');
					await tanstackClient.refetchQueries(userRefetchOptions());
				}
			} catch (err) {
				addToast('error', 'Error saving!');
			}
		} else {
			addToast('error', 'Error saving!');
		}
	}
</script>

<PageWrapper title="Tracker Schedules" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">
				Tracker Schedules
			</h1>
			<div class="grid w-full content-start">
				<div class="px-1">
					<form class="border-b-base-300 grid border-b pt-2 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Nasal Spray</legend>

						<SegmentedControl items={5}>
							<label>
								<input type="radio" bind:group={sprayInterval} value={1} name="spray" />1
							</label>
							<label>
								<input type="radio" bind:group={sprayInterval} value={2} name="spray" />2
							</label>
							<label>
								<input type="radio" bind:group={sprayInterval} value={3} name="spray" />3
							</label>
							<label>
								<input type="radio" bind:group={sprayInterval} value={4} name="spray" />4
							</label>
							<label>
								<input type="radio" bind:group={sprayInterval} value={5} name="spray" />5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Gummy</legend>

						<SegmentedControl items={5}>
							<label>
								<input type="radio" bind:group={gummyInterval} value={1} name="gummy" />1
							</label>
							<label>
								<input type="radio" bind:group={gummyInterval} value={2} name="gummy" />2
							</label>
							<label>
								<input type="radio" bind:group={gummyInterval} value={3} name="gummy" />3
							</label>
							<label>
								<input type="radio" bind:group={gummyInterval} value={4} name="gummy" />4
							</label>
							<label>
								<input type="radio" bind:group={gummyInterval} value={5} name="gummy" />5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Towel Wash</legend>

						<SegmentedControl items={5}>
							<label>
								<input type="radio" bind:group={towelInterval} value={1} name="towel" />1
							</label>
							<label>
								<input type="radio" bind:group={towelInterval} value={2} name="towel" />2
							</label>
							<label>
								<input type="radio" bind:group={towelInterval} value={3} name="towel" />3
							</label>
							<label>
								<input type="radio" bind:group={towelInterval} value={4} name="towel" />4
							</label>
							<label>
								<input type="radio" bind:group={towelInterval} value={5} name="towel" />5
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Months Per Pet petChewable</legend>

						<SegmentedControl items={3}>
							<label>
								<input
									type="radio"
									bind:group={petChewableInterval}
									value={1}
									name="petChewable"
								/>Monthly
							</label>
							<label>
								<input
									type="radio"
									bind:group={petChewableInterval}
									value={3}
									name="petChewable"
								/>Quarterly
							</label>
							<label>
								<input
									type="radio"
									bind:group={petChewableInterval}
									value={6}
									name="petChewable"
								/>Half-Yearly
							</label>
						</SegmentedControl>
					</form>

					<form class="border-b-base-300 grid border-b pt-4 pb-8 text-lg" {onchange}>
						<legend class="fieldset-legend pb-4">Days Per Pet petBath</legend>

						<SegmentedControl items={5}>
							<label>
								<input type="radio" bind:group={petBathInterval} value={3} name="petBath" />3
							</label>
							<label>
								<input type="radio" bind:group={petBathInterval} value={5} name="petBath" />5
							</label>
							<label>
								<input type="radio" bind:group={petBathInterval} value={7} name="petBath" />7
							</label>
							<label>
								<input type="radio" bind:group={petBathInterval} value={14} name="petBath" />14
							</label>
							<label>
								<input type="radio" bind:group={petBathInterval} value={30} name="petBath" />30
							</label>
						</SegmentedControl>
					</form>
				</div>
			</div>
		</div>
	</div>
</PageWrapper>
