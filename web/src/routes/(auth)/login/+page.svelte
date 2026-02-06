<script>
	import { goto } from '$app/navigation';
	import Icon from '@iconify/svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Logo from '$lib/assets/logo.webp?w=600&enhanced';
	import { resolve } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';

	// if (pb.authStore.isValid) {
	// 	goto(resolve('/app'));
	// }

	let email = $state('');

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

<PageWrapper title="Login" {pb}>
	<form class="grid h-full w-full max-w-sm content-center justify-self-center">
		<div class="lg:bg-base-200 w-full rounded-2xl lg:p-8 lg:shadow-md">
			<enhanced:img src={Logo} alt="logo" />

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
		</div>
	</form>
</PageWrapper>
