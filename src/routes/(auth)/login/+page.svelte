<script>
	import { goto } from '$app/navigation';
	import MaterialSymbolsVisibilityOffOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOffOutline.svelte';
	import MaterialSymbolsVisibilityOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOutline.svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Logo from '$lib/assets/logo.webp?w=600&enhanced';

	if (pb.authStore.isValid) {
		goto('/app');
	}

	let email = $state('');
	let password = $state('');
	let togglePasswordStatus = $state(false);

	async function submitHandler() {
		spinner = true;
		const cleanEmail = email.toLowerCase().trim();
		const cleanPassword = password.trim();

		try {
			const authData = await pb.collection('users').authWithPassword(cleanEmail, cleanPassword);
			if (authData.token) {
				addToast('success', 'Logged in successfully!');
				spinner = false;
				goto('/app');
			}
		} catch (err) {
			console.log(err);
		}
	}

	let spinner = $state(false);

	function togglePassword() {
		togglePasswordStatus = !togglePasswordStatus;
	}
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<PageWrapper title="Login" {pb}>
	<form class="grid h-full w-full max-w-sm content-center justify-self-center">
		<div class="lg:bg-base-200 w-full rounded-2xl lg:p-8 lg:shadow-md">
			<!-- <h1
				class="text-primary mb-4 text-center text-7xl font-bold tracking-tighter lg:mb-12 lg:text-9xl"
			>
				Cubby
			</h1> -->
			<enhanced:img src={Logo} alt="logo" />

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Email</legend>
				<input type="text" name="email" bind:value={email} class="input input-lg w-full" />
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Password</legend>
				<label class="input validator input-lg w-full gap-4">
					<input type={togglePasswordStatus ? 'text' : 'password'} bind:value={password} required />
					<button type="button" class="cursor-pointer" onclick={togglePassword}>
						{#if togglePasswordStatus}
							<MaterialSymbolsVisibilityOffOutline class="size-[1.3em] opacity-75" />
						{:else}
							<MaterialSymbolsVisibilityOutline class="size-[1.3em] opacity-75" />
						{/if}
					</button>
				</label>
			</fieldset>
			<button
				class="btn btn-lg btn-primary full mt-8 w-full rounded-full"
				onclick={() => submitHandler()}
			>
				{#if !spinner}
					Login
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
			<div class="mt-8 text-center text-lg">
				Don't have an account? <a href="/register" class="text-primary font-bold">Create one.</a>
			</div>
		</div>
	</form>
</PageWrapper>
