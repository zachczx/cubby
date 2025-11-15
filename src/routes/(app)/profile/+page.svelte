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
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';

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

<PageWrapper title="Profile" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Profile</h1>
			<div class="grid w-full content-start">
				<a
					href="/profile/interval"
					class="border-b-base-300 flex items-center border-b pt-4 pb-6 text-lg"
				>
					<div class="grow font-semibold">Adjust Default Intervals</div>
					<div class="flex h-full items-center">
						<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
							><MaterialSymbolsChevronRight class="size-7" /></button
						>
					</div>
				</a>

				<a
					href="/profile/vacation"
					class="border-b-base-300 flex items-center border-b pt-4 pb-6 text-lg"
				>
					<div class="grow font-semibold">Manage Vacation Dates</div>
					<div class="flex h-full items-center">
						<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
							><MaterialSymbolsChevronRight class="size-7" /></button
						>
					</div>
				</a>

				<form class="border-b-base-300 flex items-center border-b pt-4 pb-6 text-lg" {onchange}>
					<legend class="fieldset-legend grow">Sound</legend>
					<input
						type="checkbox"
						name="sound"
						class="toggle toggle-lg toggle-primary"
						bind:checked={sound}
					/>
				</form>
			</div>
		</div>
		<div class="mb-4 text-lg">
			<a href="/logout" class="btn btn-neutral btn-lg w-full rounded-full">Logout</a>
		</div>
	</div>
</PageWrapper>
