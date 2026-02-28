<script>
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Logo from '$lib/assets/logo.webp?w=600&enhanced';
	import { resolve } from '$app/paths';

	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { createQuery } from '@tanstack/svelte-query';
	import OtpInput from '$lib/ui/OtpInput.svelte';

	let email = $state('');
	let isLoading = $state(true);
	let isLoggedIn = $state(false);
	let showOtpInput = $state(false);
	let spinner = $state(false);
	let otp = $state('');
	let methodId = $state('');

	$effect(() => {
		if (otp.length === 6) {
			submitOtpHandler();
		}
	});

	onMount(async () => {
		try {
			const response = await api.get('check');
			if (response.status === 204) {
				isLoggedIn = true;
				goto('/app');
				return;
			}
		} catch {
			console.error('403: unauthenticated');
		}

		isLoading = false;
	});

	async function getOtpHandler() {
		spinner = true;
		const cleanEmail = email.toLowerCase().trim();
		const formData = new URLSearchParams();
		formData.append('email', cleanEmail);

		try {
			const response = await api
				.post('otp/send', {
					body: formData
				})
				.json();

			if (response.methodId) {
				methodId = response.methodId;
				addToast('success', 'Check your email!');
				spinner = false;
				showOtpInput = true;
			}
		} catch (err) {
			spinner = false;
			console.log(err);
		}
	}

	async function submitOtpHandler() {
		spinner = true;

		try {
			const response = await api
				.post('otp/verify', {
					body: JSON.stringify({ methodId: methodId, otp: otp })
				})
				.json();

			if (response.status === 'ok') {
				spinner = false;
				if (response.onboarding) {
					goto('/app/profile/account?onboarding=true');
				} else {
					goto('/app');
				}
			} else {
				spinner = false;
				addToast('error', 'Failed to login!');
			}
		} catch (err) {
			spinner = false;
			console.log(err);
		}
	}
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
			{:else if !showOtpInput}
				<fieldset class="fieldset mt-6">
					<legend class="fieldset-legend -mb-2 text-lg opacity-50">Email</legend>
					<input
						id="email"
						type="text"
						name="email"
						bind:value={email}
						class="input input-lg w-full"
					/>
				</fieldset>

				<button
					class="btn btn-lg btn-primary full mt-4 w-full rounded-full"
					onclick={() => getOtpHandler()}
				>
					{#if !spinner}
						Get One-Time Password
					{:else}
						<span class="loading loading-md loading-spinner"></span>
					{/if}
				</button>
			{:else if showOtpInput}
				<div class="grid max-w-96 content-start justify-items-center justify-self-center">
					<OtpInput bind:otp />
					<button
						class="btn btn-lg btn-primary full mt-4 w-full rounded-full"
						onclick={() => submitOtpHandler()}
					>
						{#if !spinner}
							Login
						{:else}
							<span class="loading loading-md loading-spinner"></span>
						{/if}
					</button>
				</div>
			{/if}
		</div>
	</form>
</PageWrapper>
