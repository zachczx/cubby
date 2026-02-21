<script lang="ts">
	import Logo from '$lib/assets/logo.webp?w=200&enhanced';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		userQueryOptions,
		userRefetchOptions,
		familyQueryOptions,
		familyRefetchOptions
	} from '$lib/queries';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';

	let { data } = $props();

	let isOnboarding = $derived(data.onboarding);

	const user = createQuery(userQueryOptions);
	const family = createQuery(familyQueryOptions);
	const tanstackClient = useQueryClient();

	let name = $state('');
	let familyName = $state('');
	let saving = $state(false);

	$effect(() => {
		if (user.isSuccess && user.data?.name && !name) {
			name = user.data.name;
		}
		if (family.isSuccess && family.data?.[0]?.name && !familyName) {
			familyName = family.data[0].name;
		}
	});

	async function onsubmit(evt: Event) {
		evt.preventDefault();
		saving = true;

		try {
			const response = await api.patch('users/me/account', {
				body: JSON.stringify({
					name,
					familyName
				})
			});

			if (response.status === 204) {
				addToast('success', 'Profile saved successfully!');
				await Promise.all([
					tanstackClient.refetchQueries(userRefetchOptions()),
					tanstackClient.refetchQueries(familyRefetchOptions())
				]);

				if (data.onboarding) goto('/app');
			}
		} catch (err) {
			console.error(err);
			addToast('error', 'Error saving profile!');
		} finally {
			saving = false;
		}
	}
</script>

<PageWrapper title="Account Info">
	<div
		class="grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-md lg:justify-self-center lg:bg-white/70 lg:p-8 lg:shadow-md"
	>
		<div class="w-full max-w-lg justify-self-center overflow-y-auto">
			{#if isOnboarding}
				<div class="mb-4 rounded-xl text-center">
					<enhanced:img src={Logo} class="justify-self-center" alt="Welcome" />
					<p class="text-base-content/80">Welcome! Let's finish setting up your profile.</p>
				</div>
			{:else}
				<h1 class="text-primary mb-8 text-center text-4xl font-bold max-lg:hidden">Account Info</h1>
			{/if}

			{#if user.isPending || family.isPending}
				<div class="flex justify-center p-8">
					<span class="loading loading-spinner text-primary"></span>
				</div>
			{:else}
				<form
					class="border-neutral/20 grid gap-6 rounded-xl border bg-white/70 p-8 shadow"
					{onsubmit}
				>
					<fieldset class="fieldset">
						<legend class="fieldset-legend text-base font-semibold">Name</legend>
						<input
							type="text"
							class="input input-bordered w-full"
							placeholder="Enter your name"
							bind:value={name}
							required
						/>
					</fieldset>

					<fieldset class="fieldset">
						<legend class="fieldset-legend text-base font-semibold">Group Name</legend>
						<input
							type="text"
							class="input input-bordered w-full"
							placeholder="e.g. The Smiths"
							bind:value={familyName}
							required
						/>
						<p class="fieldset-label mt-1 text-xs">
							This is where your shared trackers live. You can name it anything, e.g., Family, Home,
							Us.
						</p>
					</fieldset>

					<div class="mt-4">
						<button
							type="submit"
							class="btn btn-primary btn-lg w-full rounded-full"
							disabled={saving}
						>
							{#if saving}
								<span class="loading loading-spinner"></span> Saving...
							{:else}
								Save Profile
							{/if}
						</button>

						{#if isOnboarding}
							<div class="mt-4 text-center">
								<a href="/app" class="btn btn-ghost w-full rounded-full"
									><span class="opacity-75">Skip for now</span></a
								>
							</div>
						{/if}
					</div>
				</form>
			{/if}
		</div>
	</div>
</PageWrapper>
