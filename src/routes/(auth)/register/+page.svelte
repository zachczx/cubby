<script>
	import { goto } from '$app/navigation';
	import MaterialSymbolsVisibilityOffOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOffOutline.svelte';
	import MaterialSymbolsVisibilityOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOutline.svelte';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';

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
			const data = {
				email: cleanEmail,
				emailVisibility: true,
				password: cleanPassword,
				passwordConfirm: cleanPassword,
				sprayInterval: 2,
				gummyInterval: 2
			};

			const record = await pb.collection('users').create(data);
		} catch (err) {
			console.log(err);
		}

		try {
			const authData = await pb.collection('users').authWithPassword(email, password);
			if (authData.token) {
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
	<title>Register</title>
</svelte:head>

<PageWrapper title="Register" {pb}>
	<form class="grid h-full w-full max-w-sm content-center justify-self-center">
		<div class="lg:bg-base-200 w-full rounded-2xl lg:p-8 lg:shadow-md">
			<h1
				class="text-primary mb-4 text-center text-5xl font-bold tracking-tighter lg:mb-12 lg:text-9xl"
			>
				Sundry
			</h1>
			<fieldset class="fieldset">
				<input
					type="text"
					name="email"
					bind:value={email}
					class="input input-lg w-full"
					placeholder="Email"
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<label class="input input-lg w-full gap-4">
					<input
						type={togglePasswordStatus ? 'text' : 'password'}
						bind:value={password}
						placeholder="Password"
						required
					/>
					<button type="button" class="cursor-pointer" onclick={togglePassword}>
						{#if togglePasswordStatus}
							<MaterialSymbolsVisibilityOffOutline class="size-[1.3em] opacity-75" />
						{:else}
							<MaterialSymbolsVisibilityOutline class="size-[1.3em] opacity-75" />
						{/if}
					</button>
				</label>
			</fieldset>
			<button class="btn btn-lg btn-primary mt-4 w-full" onclick={() => submitHandler()}>
				{#if !spinner}
					Register
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
			<div class="mt-8 text-center">
				Have an account? <a href="/login" class="text-primary font-bold">Login here.</a>
			</div>
		</div>
	</form>
</PageWrapper>
