<script lang="ts">
	import { goto } from '$app/navigation';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import corgiEmail from '$lib/assets/corgi_email.webp?w=600;400;300&enhanced';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
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
			const response: Record<'methodId', string> = await api
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

	interface VerifyResponse {
		status: string;
		onboarding: string;
	}

	async function submitOtpHandler() {
		spinner = true;

		try {
			const response: VerifyResponse = await api
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

<PageWrapper title="Login" focusedScreen showSettingsIcon={false} back={false}>
	<form class={['grid h-full w-full max-w-sm content-center justify-self-center']}>
		<div class="lg:bg-base-200 w-full rounded-2xl lg:p-8 lg:shadow-md">
			<div class="w-full max-w-72 justify-self-center lg:max-w-lg">
				<enhanced:img
					src={corgiEmail}
					alt="email"
					sizes="(min-width:1920px) 1000px, (min-width:1080px) 600px, (min-width:768px) 400px, 300px"
				/>
			</div>

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
				<div class="mt-4 grid max-w-96 content-start justify-items-center justify-self-center">
					<div class="text-center">
						<h2 class="text-2xl font-bold">OTP Sent</h2>
						<span>Enter the code sent to</span>&nbsp;<span class="font-bold">{email}</span>
					</div>
					<OtpInput bind:otp />
					<button
						class="btn btn-lg btn-primary full mt-4 w-full rounded-full"
						onclick={() => submitOtpHandler()}
					>
						{#if !spinner}
							Verify
						{:else}
							<span class="loading loading-md loading-spinner"></span>
						{/if}
					</button>
				</div>
			{/if}
		</div>
	</form>
</PageWrapper>
