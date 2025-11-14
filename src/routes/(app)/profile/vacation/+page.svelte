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
	import MaterialSymbolsDelete from '$lib/assets/svg/MaterialSymbolsDelete.svelte';
	import { invalidateAll } from '$app/navigation';

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

	async function addHandler() {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			const start = dayjs.tz(vacationStart, 'Asia/Singapore');
			const end = dayjs.tz(vacationEnd, 'Asia/Singapore');

			const result = await pb.collection('vacation').create({
				user: pb.authStore.record?.id,
				startDateTime: start,
				endDateTime: end
			});
			if (result.id) {
				addToast('success', 'Added successfully!');
				spinner = false;

				await tanstackClient.refetchQueries(createVacationRefetchOptions());
			}
		} catch (err) {
			console.log(err);
		}
	}

	async function deleteHandler(deleteId: string) {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			const result = await pb.collection('vacation').delete(deleteId);
			if (result) {
				addToast('success', 'Deleted successfully!');
				spinner = false;

				await tanstackClient.refetchQueries(createVacationRefetchOptions());
				invalidateAll();
			}
		} catch (err) {
			console.log(err);
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
				<div class="bg-base-200 lg:bg-base-300 mt-4 rounded-xl p-4">
					<div class="mb-2 flex items-center">
						<h3 class="text-base-content/70 grow font-semibold uppercase">Recent Vacations</h3>
						<button
							class="btn btn-ghost -me-4 flex items-center gap-2 font-normal"
							onclick={() => vacationsModal.showModal()}
							aria-label="see all"
						>
							View all
							<!-- <MaterialSymbolsOpenInNew class="size-[1.3em] opacity-75" /> -->
						</button>
					</div>
					<ul class="list-disc">
						{#if vacations.isSuccess}
							{#each vacations.data as v}
								{@const dateTime = formatTime(v.startDateTime, v.endDateTime)}
								<li class="ms-6 py-0.5">
									{dateTime}
								</li>
							{/each}
						{/if}
					</ul>
				</div>

				<div class="border-b-base-300 border-b pt-4 pb-6 text-lg">
					<legend class="fieldset-legend pb-4">Start Date</legend>

					<input
						type="date"
						name="vacationStart"
						class="input input-lg"
						bind:value={vacationStart}
					/>
				</div>

				<div class="border-b-base-300 border-b pt-4 pb-6 text-lg">
					<legend class="fieldset-legend pb-4">End Date</legend>

					<input type="date" name="vacationEnd" class="input input-lg" bind:value={vacationEnd} />
				</div>
			</div>
		</div>
		<button class="btn btn-primary btn-lg w-full rounded-full" onclick={addHandler}
			>Add Vacation</button
		>
	</div>
</PageWrapper>

<dialog bind:this={vacationsModal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<h3 class="mb-4 text-lg font-bold uppercase">Recent Vacations</h3>
		<ul class="list-disc">
			{#if vacations.isSuccess}
				{#each vacations.data as v}
					{@const dateTime = formatTime(v.startDateTime, v.endDateTime)}
					<li class="border-b-base-300 ms-6 border-b py-4">
						<div class="flex items-center gap-4">
							<div class="grow">{dateTime}</div>
							<button
								class="btn btn-error btn-sm"
								onclick={() => deleteHandler(v.id)}
								aria-label="delete"><MaterialSymbolsDelete class="size-5" /></button
							>
						</div>
					</li>
				{/each}
			{/if}
		</ul>
	</div>
</dialog>
