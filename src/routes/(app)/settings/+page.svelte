<script lang="ts">
	import { goto } from '$app/navigation';
	import MaterialSymbolsVisibilityOffOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOffOutline.svelte';
	import MaterialSymbolsVisibilityOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOutline.svelte';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { onMount } from 'svelte';

	let user: UserDB | undefined = $state();

	if (!pb.authStore.isValid) {
		goto('/login');
	}

	onMount(async () => {
		if (pb.authStore.isValid && pb.authStore.record) {
			user = await pb.collection('users').getOne(pb.authStore.record.id);
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
</script>

<PageWrapper title="Settings" {pb}>
	<div
		class="grid h-full w-full max-w-lg justify-self-center max-lg:grid-rows-[auto_1fr] max-lg:content-start lg:content-center"
	>
		<div class="lg:bg-base-200 w-full rounded-2xl shadow-md lg:p-8">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold lg:mb-12">Settings</h1>

			<div class="grid w-full grid-rows-[auto_auto_1fr_auto] content-start">
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

				<div></div>

				<button class="btn btn-lg btn-primary mt-8" onclick={() => submitHandler()}>
					{#if !spinner}
						Save Settings
					{:else}
						<span class="loading loading-md loading-spinner"></span>
					{/if}
				</button>
			</div>
		</div>
	</div>
</PageWrapper>
