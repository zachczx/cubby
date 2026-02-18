<script>
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Logo from '$lib/assets/logo.webp?w=600&enhanced';
	import { resolve } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { createQuery } from '@tanstack/svelte-query';

	let email = $state('');
	let isLoading = $state(true);
	let isLoggedIn = $state(false);

	onMount(async () => {
		const response = await api.get('check');

		if (response.status === 204) {
			isLoggedIn = true;
			goto('/app');
			return;
		}

		isLoading = false;
	});

	async function submitHandler() {
		spinner = true;
		const cleanEmail = email.toLowerCase().trim();
		const formData = new URLSearchParams();
		formData.append('email', cleanEmail);

		try {
			const response = await fetch(PUBLIC_API_URL + '/magic-link', {
				method: 'post',
				body: formData
			});

			if (response.ok) {
				addToast('success', 'Check your email!');
				spinner = false;
			}
		} catch (err) {
			console.log(err);
		}
	}

	let spinner = $state(false);
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<PageWrapper title="Login">
	<form class={['grid h-full w-full max-w-sm content-center justify-self-center']}>
		<div class="lg:bg-base-200 w-full rounded-2xl lg:p-8 lg:shadow-md">
			<enhanced:img src={Logo} alt="logo" />

			{#if isLoading}
				<div class="grid min-h-24 content-center justify-center">
					<span class="loading loading-md loading-spinner"></span>
				</div>
			{:else if !isLoading && isLoggedIn}
				You're logged in, redirecting to app...
			{:else}
				<fieldset class="fieldset mt-6">
					<legend class="fieldset-legend -mb-2 text-lg opacity-50">Email</legend>
					<input type="text" name="email" bind:value={email} class="input input-lg w-full" />
				</fieldset>

				<button
					class="btn btn-lg btn-primary full mt-4 w-full rounded-full"
					onclick={() => submitHandler()}
				>
					{#if !spinner}
						Get Login Link
					{:else}
						<span class="loading loading-md loading-spinner"></span>
					{/if}
				</button>
			{/if}
		</div>
	</form>
</PageWrapper>
