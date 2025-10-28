<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';

	import { onMount } from 'svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	let user: UserDB | undefined = $state();
	let vacations: VacationDB[] | undefined = $state([]);

	onMount(async () => {
		if (pb.authStore.isValid && pb.authStore.record) {
			user = await pb.collection('users').getOne(pb.authStore.record.id);
			vacations = await pb.collection('vacation').getFullList({ sort: '-startDateTime' });
		}
	});

	let sprayInterval = $derived(user?.defaultSprayInterval);
	let mute = $derived(!user?.mute);

	let spinner = $state(false);

	async function submitHandler() {
		if (!user || !user.id) return;
		spinner = true;

		try {
			const result = await pb.collection('users').update(user?.id, {
				defaultSprayInterval: sprayInterval,
				mute: mute === undefined ? false : !mute
			});
			if (result.status === 200) {
				addToast('success', 'Updated!');
				spinner = false;
			}
		} catch (err) {
			console.log(err);
		}
	}

	async function addVacationDateHandler() {
		if (!user || !user.id) return;
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
				vacations = await pb.collection('vacation').getFullList({ sort: '-startDateTime' });
			}
		} catch (err) {
			console.log(err);
		}
	}

	let currentTab = $state('settings');
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

<PageWrapper title="Settings" back={true} {pb}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold">Profile</h1>

			<ul class="grid w-full grid-cols-2 content-center justify-items-center pt-8">
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
					<div class="">
						<div class="border-b-base-300 flex items-center border-b py-4 text-lg">
							<legend class="fieldset-legend mb-2 grow">Days Between Nasal Sprays</legend>

							<select bind:value={sprayInterval} class="select w-auto shrink text-lg">
								<option value={1}>1</option>
								<option value={2}>2</option>
								<option value={3}>3</option>
								<option value={4}>4</option>
								<option value={5}>5</option>
							</select>
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
							{#each vacations as v}
								{@const dateTime = formatTime(v.startDateTime, v.endDateTime)}
								<li class="ms-6 py-0.5">
									{dateTime}
								</li>
							{/each}
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
			{#each vacations as v}
				{@const dateTime = formatTime(v.startDateTime, v.endDateTime)}
				<li class="ms-6 py-0.5">
					{dateTime}
				</li>
			{/each}
		</ul>
	</div>
</dialog>
