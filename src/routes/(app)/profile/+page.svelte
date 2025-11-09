<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
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

	let mute = $derived.by(() => (user.isSuccess ? !user.data?.mute : undefined));

	let spinner = $state(false);

	async function submitHandler() {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			await pb.collection('users').update(user.data.id, {
				towelIntervalDays: towelIntervalDays,
				sprayIntervalDays: sprayIntervalDays,
				gummyIntervalDays: gummyIntervalDays,
				doggoChewableIntervalMonths: doggoChewableIntervalMonths,
				doggoBathIntervalDays: doggoBathIntervalDays,
				mute: mute === undefined ? false : !mute
			});
			addToast('success', 'Updated!');
			spinner = false;

			await tanstackClient.refetchQueries(createUserRefetchOptions());
		} catch (err) {
			console.log(err);
		}
	}

	async function addVacationDateHandler() {
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
</script>

<PageWrapper title="Profile" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Profile</h1>
			<ul class="grid w-full grid-cols-2 content-center justify-items-center pt-4">
				<li class="w-full">
					<button
						class={[
							currentTab === 'settings' ? 'border-b-primary font-semibold' : 'border-b-base-300',
							'w-full cursor-pointer border-b-2 text-center'
						]}
						onclick={() => {
							currentTab = 'settings';
						}}
					>
						Settings</button
					>
				</li>
				<li class="w-full">
					<button
						class={[
							currentTab === 'vacations' ? 'border-b-primary font-semibold' : 'border-b-base-300',
							'w-full cursor-pointer border-b-2 text-center'
						]}
						onclick={() => {
							currentTab = 'vacations';
						}}
					>
						Vacations</button
					>
				</li>
			</ul>

			<div class="mt-4 grid w-full content-start">
				{#if currentTab === 'settings'}
					<div class="px-1">
						<div class="border-b-base-300 grid border-b py-4 text-lg">
							<legend class="fieldset-legend">Days Per Towel Wash</legend>

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
						</div>

						<div class="border-b-base-300 grid border-b py-4 text-lg">
							<legend class="fieldset-legend">Days Per Nasal Spray</legend>

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
						</div>

						<div class="border-b-base-300 grid border-b py-4 text-lg">
							<legend class="fieldset-legend">Days Per Gummy</legend>

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
						</div>

						<div class="border-b-base-300 grid border-b py-4 text-lg">
							<legend class="fieldset-legend">Months Per Pet Chewable</legend>

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
						</div>

						<div class="border-b-base-300 grid border-b py-4 text-lg">
							<legend class="fieldset-legend">Days Per Pet Bath</legend>

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
						</div>

						<div class="border-b-base-300 flex items-center border-b py-4 text-lg">
							<legend class="fieldset-legend grow">Sound</legend>

							<input
								type="checkbox"
								name="mute"
								class="toggle toggle-lg toggle-primary"
								bind:checked={mute}
							/>
						</div>
					</div>
				{/if}

				{#if currentTab === 'vacations'}
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

					<div class="border-b-base-300 border-b py-4 text-lg">
						<legend class="fieldset-legend">Start Date</legend>

						<input type="date" name="vacationStart" class="input" bind:value={vacationStart} />
					</div>

					<div class="border-b-base-300 border-b py-4 text-lg">
						<legend class="fieldset-legend">End Date</legend>

						<input type="date" name="vacationEnd" class="input" bind:value={vacationEnd} />
					</div>
				{/if}
			</div>
		</div>
		{#if currentTab === 'settings'}
			<button class="btn btn-lg btn-primary mt-8" onclick={() => submitHandler()}>
				{#if !spinner}
					Save Settings
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
		{:else if currentTab === 'vacations'}
			<button class="btn btn-lg btn-primary mt-8" onclick={() => addVacationDateHandler()}>
				{#if !spinner}
					Add Vacation
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
		{/if}
	</div>
</PageWrapper>

<dialog bind:this={vacationsModal} class="modal">
	<div class="modal-box">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<h3 class="mb-4 text-lg font-bold uppercase">Recent Vacations</h3>
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
</dialog>
