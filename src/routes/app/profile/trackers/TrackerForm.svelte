<script lang="ts">
	import MaterialSymbolsAttachMoney from '$lib/assets/svg/MaterialSymbolsAttachMoney.svelte';
	import { getAllTrackerIcons } from '$lib/mapper';
	import { pb } from '$lib/pb';
	import { userQueryOptions } from '$lib/queries';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import NumberInput from '$lib/ui/NumberInput.svelte';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { slide } from 'svelte/transition';

	let {
		onsubmit,
		currentTracker
	}: {
		onsubmit: (trackerDetails: TrackerInput) => Promise<TrackerDB> | Promise<void>;
		currentTracker?: TrackerDB;
	} = $props();

	let intervalString = $state('1');
	const user = createQuery(userQueryOptions);

	function toCamelCase(text: string): string {
		const words = text.split(' ').map((word, idx) => {
			if (idx === 0) return word.toLowerCase();

			return word.charAt(0).toUpperCase() + word.slice(1);
		});

		return words.join('');
	}
	let loadedTrackerId = $state('');

	let inputTrackerDetails = $state<TrackerInput>({
		user: pb.authStore.record?.id ?? '',
		family: '',
		name: '',
		display: '',
		interval: '1',
		intervalUnit: 'day',
		category: 'personal',
		kind: 'task',
		icon: '', //todo
		startDate: '',
		actionLabel: '',
		pinned: false,
		show: true
	});

	let camelCaseName = $derived(toCamelCase(inputTrackerDetails.display));

	let clean = $derived.by(() => ({ ...inputTrackerDetails, name: camelCaseName }));

	async function handleSubmission() {
		try {
			await onsubmit(clean);
		} catch (err) {
			console.log(err);
		}
	}

	$effect(() => {
		if (currentTracker && loadedTrackerId === '') {
			const { id, ...currentTrackerLessId } = currentTracker;
			inputTrackerDetails = {
				...currentTrackerLessId,
				interval: String(currentTrackerLessId.interval)
			};
			intervalString = String(currentTracker.interval);
			loadedTrackerId = currentTracker.id;
		}
	});

	let icons = getAllTrackerIcons();
</script>

<form class="grid w-full content-start gap-4" onsubmit={() => handleSubmission()}>
	<fieldset class="fieldset">
		<legend class="fieldset-legend -mb-2 text-lg">Display Name</legend>
		<input
			type="text"
			name="display"
			placeholder="Give your new tracker a name"
			bind:value={inputTrackerDetails.display}
			autocomplete="off"
			class="input input-lg w-full"
		/>
	</fieldset>

	<fieldset class="fieldset mt-2">
		<legend class="fieldset-legend -mb-2 text-lg">Category</legend>

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
		<legend class="fieldset-legend -mb-2 text-lg">Type of Tracker</legend>

		<SegmentedControl items={2}>
			<label>
				<input type="radio" bind:group={inputTrackerDetails.kind} value="task" name="kind" />Task
			</label>
			<label>
				<input
					type="radio"
					bind:group={inputTrackerDetails.kind}
					value="subscription"
					name="kind"
				/>Subscription
			</label>
		</SegmentedControl>
	</fieldset>

	{#if inputTrackerDetails.kind === 'subscription'}
		<fieldset transition:slide={{ duration: 200 }} class="fieldset mt-2">
			<legend class="fieldset-legend -mb-2 text-lg">Cost (optional)</legend>
			<label class="input input-lg flex w-full items-center">
				<MaterialSymbolsAttachMoney class="size-[1.3em] opacity-50" />
				<input
					type="number"
					min="0"
					step="0.01"
					name="cost"
					bind:value={inputTrackerDetails.cost}
					class="grow"
					autocomplete="off"
				/></label
			>
		</fieldset>
	{/if}

	<fieldset class="fieldset mt-2">
		<legend class="fieldset-legend -mb-2 text-lg">Interval</legend>
		<NumberInput bind:value={intervalString} />
	</fieldset>

	<fieldset class="fieldset mt-2">
		<legend class="fieldset-legend -mb-2 text-lg">Frequency</legend>

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
		<legend class="fieldset-legend -mb-2 text-lg">Icon</legend>

		<div class="grid grid-cols-4 justify-items-center">
			{#each Object.entries(icons) as [key, IconComponent]}
				<label class="cursor-pointer">
					<input
						type="radio"
						bind:group={inputTrackerDetails.icon}
						value={key}
						name="intervalUnit"
						class="peer hidden"
					/>
					<div
						class="peer-checked:bg-bg-checked/60 aspect-square rounded-xl p-4 text-4xl not-peer-checked:saturate-[0.2] peer-checked:outline peer-checked:outline-green-200 hover:bg-white/80"
					>
						<IconComponent />
					</div>
				</label>
			{/each}
		</div>
	</fieldset>

	<fieldset class="fieldset mt-2">
		<legend class="fieldset-legend -mb-2 text-lg">Button Action Label (optional)</legend>
		<input
			type="text"
			name="actionLabel"
			placeholder="E.g., Washed, Fed, Paid."
			bind:value={inputTrackerDetails.actionLabel}
			class="input input-lg w-full"
			autocomplete="off"
		/>
	</fieldset>

	<label class="flex items-center py-2">
		<div class="fieldset-legend grow text-lg font-bold">Pinned</div>
		<input
			type="checkbox"
			class="toggle toggle-lg"
			bind:checked={inputTrackerDetails.pinned}
			name="pinned"
		/>
	</label>

	<button class="btn btn-primary btn-lg mt-4 w-full rounded-full">Save</button>
</form>
