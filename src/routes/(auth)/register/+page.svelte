<script lang="ts">
	import { goto } from '$app/navigation';
	import MaterialSymbolsVisibilityOffOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOffOutline.svelte';
	import MaterialSymbolsVisibilityOutline from '$lib/assets/svg/MaterialSymbolsVisibilityOutline.svelte';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { trackerDefaults } from './tracker-defaults';
	import Logo from '$lib/assets/logo.webp?w=600&enhanced';
	import { resolve } from '$app/paths';

	if (pb.authStore.isValid) {
		goto(resolve('/app'));
	}

	let newUser = $state<Record<string, string>>({
		email: '',
		password: '',
		name: ''
	});
	let togglePasswordStatus = $state(false);
	let newUserRecord = $state<Record<string, string>>();

	/**
	 * Submit handler needs to do a few things:
	 * 1) Register
	 * 2) Try logging in
	 * 3) Create family
	 * 4) Update user with that family id
	 * 5) Create default trackers
	 */
	async function submitHandler() {
		spinner = true;
		const cleanEmail = newUser.email.toLowerCase().trim();
		const cleanPassword = newUser.password.trim();
		const cleanName = newUser.name.trim();
		try {
			const userData = {
				name: cleanName,
				email: cleanEmail,
				emailVisibility: true,
				password: cleanPassword,
				passwordConfirm: cleanPassword,
				sound: true,
				generalTasksUpcomingDays: 14
			};

			newUserRecord = await pb.collection('users').create(userData);
		} catch (err) {
			console.log(err);
		}

		const authData = await pb.collection('users').authWithPassword(newUser.email, newUser.password);
		if (!authData.token || !newUserRecord) return;

		let newFamily;
		try {
			newFamily = await pb.collection('families').create({
				name: 'Default',
				'members+': newUserRecord?.id,
				owner: newUserRecord?.id,
				enabled: true
			});

			console.log('here', newFamily);

			const batch = pb.createBatch();

			for (const t of trackerDefaults) {
				const record = {
					user: newUserRecord?.id,
					family: newFamily.id,
					...t
				};

				batch.collection('trackers').create(record);
			}

			await batch.send();

			spinner = false;
			goto(resolve('/app'));
		} catch (err) {
			console.log(err);
			goto(resolve('/app/error'));
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
			<!-- <h1
				class="text-primary mb-4 text-center text-7xl font-bold tracking-tighter lg:mb-12 lg:text-9xl"
			>
				Cubby
			</h1> -->
			<enhanced:img src={Logo} alt="logo" />

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Email</legend>
				<input
					type="text"
					name="email"
					bind:value={newUser.email}
					class="input input-lg w-full"
					required
				/>
			</fieldset>

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Password</legend>
				<label class="input input-lg w-full gap-4">
					<input
						type={togglePasswordStatus ? 'text' : 'password'}
						bind:value={newUser.password}
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

			<fieldset class="fieldset mt-2">
				<legend class="fieldset-legend -mb-2 text-lg opacity-50">Name</legend>
				<input
					type="text"
					name="name"
					bind:value={newUser.name}
					class="input input-lg w-full"
					required
				/>
			</fieldset>

			<button
				class="btn btn-lg btn-primary mt-8 w-full rounded-full"
				onclick={() => submitHandler()}
			>
				{#if !spinner}
					Register
				{:else}
					<span class="loading loading-md loading-spinner"></span>
				{/if}
			</button>
			<div class="mt-8 text-center text-lg">
				Have an account? <a href="/login" class="text-primary font-bold">Login here.</a>
			</div>
		</div>
	</form>
</PageWrapper>
