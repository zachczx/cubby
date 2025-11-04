<script lang="ts">
	import { page } from '$app/state';
	import { SafeArea } from '@capacitor-community/safe-area';

	import type Client from 'pocketbase';
	import type { Snippet } from 'svelte';
	import MaterialSymbolsClose from './assets/svg/MaterialSymbolsClose.svelte';
	import MaterialSymbolsLogin from './assets/svg/MaterialSymbolsLogin.svelte';
	import MaterialSymbolsLogout from './assets/svg/MaterialSymbolsLogout.svelte';
	import MaterialSymbolsSettings from './assets/svg/MaterialSymbolsSettings.svelte';
	import MaterialSymbolsNotifications from './assets/svg/MaterialSymbolsNotifications.svelte';
	import MaterialSymbolsWarning from './assets/svg/MaterialSymbolsWarning.svelte';
	import MaterialSymbolsChevronRight from './assets/svg/MaterialSymbolsChevronRight.svelte';
	import MaterialSymbolsCheckCircle from './assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsNotificationImportant from './assets/svg/MaterialSymbolsNotificationImportant.svelte';
	import { createSprayQueryOptions, createTowelQueryOptions } from './queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import { dirtyTowelDays } from './config';
	import { getNotificationStatus } from './notification';

	let {
		pb,
		children,
		title,
		back,
		noPadding = false,
		// towelNotification,
		// sprayNotification,
		largeScreenCenter
	}: {
		pb: Client;
		children: Snippet;
		title: string;
		back?: boolean;
		noPadding?: boolean;
		// towelNotification?: boolean;
		// sprayNotification?: boolean;
		largeScreenCenter?: boolean;
	} = $props();

	let drawerToggle = $state() as HTMLInputElement;

	SafeArea.enable({
		config: {
			customColorsForSystemBars: true,
			statusBarColor: '#00000000', // transparent
			statusBarContent: 'light',
			navigationBarColor: '#00000000', // transparent
			navigationBarContent: 'light',
			offset: 0
		}
	});

	const towels = createQuery(createTowelQueryOptions);
	const sprays = createQuery(createSprayQueryOptions);

	let sprayNotification = $derived.by(() => getNotificationStatus(sprays));

	let towelNotification = $derived.by(() => getNotificationStatus(towels));
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>

<div class="grid h-dvh w-full content-start justify-items-center">
	<div
		id="safe-area-topnav"
		class={[
			'navbar border-b-base-300/50 bg-neutral text-neutral-content fixed top-0 z-1 min-h-14 items-center border-b pe-4',
			back ? 'lg:ps-4' : 'ps-4'
		]}
	>
		<div class="navbar-start">
			{#if back}
				<button
					aria-label="go back"
					class="cursor-pointer p-2 max-lg:me-4 lg:hidden"
					onclick={() => {
						if (window) {
							history.back();
						}
					}}
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="32"
						height="32"
						class="lucide:arrow-left size-6"
						viewBox="0 0 24 24"
						><path
							fill="none"
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="m12 19l-7-7l7-7m7 7H5"
						/></svg
					>
				</button>
			{/if}

			<a class="hidden text-xl font-bold lg:flex" href="/">Nosey</a>
			<span class="text-xl font-bold lg:hidden">{title}</span>
		</div>
		<div id="desktop-menu" class="navbar-center hidden lg:flex">
			<ul class="menu menu-horizontal gap-8 px-1 text-lg">
				<li>
					<a
						href="/towelie"
						class={[
							'border-b-2',
							page.route.id?.includes('towelie') &&
								'border-b-neutral-content rounded-none border-b-2 font-bold',
							!page.route.id?.includes('towelie') && 'border-b-transparent hover:bg-white/20'
						]}>Towelie</a
					>
				</li>
				<li>
					<a
						href="/nosey"
						class={[
							'border-b-2',
							page.route.id?.includes('nosey') && 'border-b-neutral-content rounded-none font-bold',
							!page.route.id?.includes('nosey') && 'border-b-transparent hover:bg-white/20'
						]}>Nosey</a
					>
				</li>
				<li>
					<a
						href="/count"
						class={[
							'border-b-2',
							page.route.id?.includes('count') &&
								'border-b-neutral-content rounded-none border-b-2 font-bold',
							!page.route.id?.includes('count') && 'border-b-transparent hover:bg-white/20'
						]}>Count</a
					>
				</li>
			</ul>
		</div>
		<div class="navbar-end">
			<div id="mobile-hamburger" class="dropdown flex items-center gap-4 lg:hidden">
				{@render notification(sprayNotification, towelNotification)}

				<div class="drawer drawer-end">
					<input id="side-drawer" type="checkbox" class="drawer-toggle" bind:this={drawerToggle} />
					<div class="drawer-content">
						<label for="side-drawer" class="btn btn-ghost drawer-button px-2 py-0">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="32"
								height="32"
								class="material-symbols:menu h-6 w-6"
								viewBox="0 0 24 24"
								><path fill="currentColor" d="M3 18v-2h18v2zm0-5v-2h18v2zm0-5V6h18v2z" /></svg
							>
						</label>
					</div>
					<div id="mobile-drawer" class="drawer-side">
						<label for="side-drawer" aria-label="close sidebar" class="drawer-overlay"></label>

						<ul
							class="menu bg-base-200 text-base-content grid min-h-full w-80 grid-rows-[auto_1fr_auto] p-4"
						>
							<div class="flex items-center justify-end">
								<button
									class="btn btn-ghost"
									onclick={() => {
										drawerToggle.checked = false;
									}}
									onmousedown={() => {
										drawerToggle.checked = false;
									}}><MaterialSymbolsClose class="size-5" /></button
								>
							</div>
							<ul>
								<li>
									<a
										href="/"
										class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											width="1em"
											height="1em"
											class="material-symbols:home-outline size-[1.3em] opacity-75"
											viewBox="0 0 24 24"
											><path
												fill="currentColor"
												d="M6 19h3v-6h6v6h3v-9l-6-4.5L6 10zm-2 2V9l8-6l8 6v12h-7v-6h-2v6zm8-8.75"
											/></svg
										>Home</a
									>
								</li>
								<li>
									<a
										href="/towelie"
										class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
										><svg
											xmlns="http://www.w3.org/2000/svg"
											width="1em"
											height="1em"
											class="lucide-lab:towel-rack size-[1.3em] opacity-75"
											viewBox="0 0 24 24"
											><path
												fill="none"
												stroke="currentColor"
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M8 6H2m4-4h12a2 2 0 0 1 2 2v18H8V4a2 2 0 0 0-4 0v15h4M22 6h-2M8 18h12"
											/></svg
										>Towelie</a
									>
								</li>
								<li>
									<a
										href="/nosey"
										class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											width="1em"
											height="1em"
											class="material-symbols:health-and-safety-outline size-[1.3em] opacity-75"
											viewBox="0 0 24 24"
											><path
												fill="currentColor"
												d="M10.5 15.5h3V13H16v-3h-2.5V7.5h-3V10H8v3h2.5zM12 22q-3.475-.875-5.738-3.988T4 11.1V5l8-3l8 3v6.1q0 3.8-2.262 6.913T12 22m0-2.1q2.6-.825 4.3-3.3t1.7-5.5V6.375l-6-2.25l-6 2.25V11.1q0 3.025 1.7 5.5t4.3 3.3m0-7.9"
											/></svg
										>Nosey</a
									>
								</li>
								<li>
									<a
										href="/count"
										class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											width="32"
											height="32"
											class="material-symbols:timer-outline size-[1.3em] opacity-75"
											viewBox="0 0 24 24"
											><path
												fill="currentColor"
												d="M9 3V1h6v2zm2 11h2V8h-2zm1 8q-1.85 0-3.488-.712T5.65 19.35t-1.937-2.863T3 13t.713-3.488T5.65 6.65t2.863-1.937T12 4q1.55 0 2.975.5t2.675 1.45l1.4-1.4l1.4 1.4l-1.4 1.4Q20 8.6 20.5 10.025T21 13q0 1.85-.713 3.488T18.35 19.35t-2.863 1.938T12 22m0-2q2.9 0 4.95-2.05T19 13t-2.05-4.95T12 6T7.05 8.05T5 13t2.05 4.95T12 20m0-7"
											/></svg
										>Count</a
									>
								</li>
							</ul>
							<ul>
								{#if pb.authStore.isValid}
									<li>
										<a
											href="/profile"
											class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
											><MaterialSymbolsSettings class="size-[1.3em] opacity-75" />Profile</a
										>
									</li>
									<li>
										<a
											href="/logout"
											class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
											><MaterialSymbolsLogout class="size-[1.3em] opacity-75" />Logout</a
										>
									</li>
								{/if}
								{#if !pb.authStore.isValid}
									<li>
										<a
											href="/login"
											class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
											><MaterialSymbolsLogin class="size-[1.3em] opacity-75" />Login</a
										>
									</li>
									<li>
										<a
											href="/register"
											class="text-base-content flex w-full gap-4 py-4 text-lg font-semibold"
											><svg
												xmlns="http://www.w3.org/2000/svg"
												width="32"
												height="32"
												class="material-symbols:signature size-[1.3em] opacity-75"
												viewBox="0 0 24 24"
												><path
													fill="currentColor"
													d="M14.075 11.725q1.825-1.35 2.85-2.962T17.95 5.55q0-.8-.262-1.175T16.975 4Q15.8 4 14.9 5.988t-.9 4.487q0 .35.013.663t.062.587M3 21v-2h2v2zm4 0v-2h2v2zm4 0v-2h2v2zm4 0v-2h2v2zm4 0v-2h2v2zM3.4 17L2 15.6L3.6 14L2 12.4L3.4 11L5 12.6L6.6 11L8 12.4L6.4 14L8 15.6L6.6 17L5 15.4zm12.05-1q-.75 0-1.375-.288T13 14.776q-.625.35-1.287.625t-1.363.55l-.7-1.875q.7-.25 1.338-.537t1.237-.613q-.125-.55-.187-1.2t-.063-1.4q0-3.6 1.425-5.962T16.975 2q1.3 0 2.125.963t.825 2.687q0 2.15-1.362 4.25t-3.788 3.775q.175.175.363.263t.387.087q.65 0 1.513-.825t1.562-2.175l1.825.85q-.175.425-.275 1.025t.025 1.05q.25-.125.588-.425t.687-.75L23.025 14q-.65.9-1.5 1.45T19.95 16q-.525 0-.937-.312t-.688-.963q-.7.625-1.425.95T15.45 16"
												/></svg
											>
											Register</a
										>
									</li>
								{/if}
							</ul>
						</ul>
					</div>
				</div>
			</div>
			<div id="desktop-logout" class="hidden items-center text-sm lg:flex">
				{@render notification(sprayNotification, towelNotification)}
				{#if pb.authStore.isValid}
					<a href="/profile" class="btn btn-ghost"
						><MaterialSymbolsSettings class="size-[1.5em] opacity-75" /></a
					><a href="/logout" class="btn btn-outline ms-1">Logout</a>
				{:else}
					<a href="/register" class="underline">Register</a>&nbsp;&nbsp;|&nbsp;&nbsp;<a
						href="/login"
						class="underline">Login</a
					>
				{/if}
			</div>
		</div>
	</div>

	<div
		class={[
			'bg-pattern mt-14 min-h-[calc(100dvh-3.5rem)] w-full lg:grid',
			!noPadding && 'px-4 pt-4 pb-8 lg:px-12 lg:pt-12',
			largeScreenCenter && 'lg:content-center'
		]}
	>
		{@render children?.()}
	</div>
</div>

{#snippet notification(
	sprayNotification: NotificationStatus,
	towelNotification: NotificationStatus
)}
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost drawer-button px-2 py-0">
			{#if !sprayNotification.show && !towelNotification.show}
				<MaterialSymbolsNotifications class="size-6" />
			{:else}
				{@const count = sprayNotification.show && towelNotification.show ? 2 : 1}
				<MaterialSymbolsNotificationImportant class="size-6" />
				<span class="absolute top-1 right-1 z-2 size-4 rounded-full bg-red-600 text-xs"
					>{count}</span
				>
			{/if}
		</div>
		<ul
			class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 w-72 p-2 shadow-lg"
		>
			{#if !sprayNotification.show && !towelNotification.show}
				<li>
					<div class="flex items-center justify-center gap-2">
						<MaterialSymbolsCheckCircle class="size-[1.3em]" /><span>No pending items</span>
					</div>
				</li>
			{/if}

			{#if sprayNotification.show}
				<li>
					<a href="/nosey" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />Spray your nose!
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if towelNotification.show}
				<li>
					<a href="/towelie" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />Wash your towel!
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}
		</ul>
	</div>
{/snippet}

<style>
	#safe-area-topnav,
	#mobile-drawer {
		padding-top: calc(var(--safe-area-inset-top) + 0.5rem) !important;
	}
	/* 
	.bg-pattern {
		background-color: red;
	} */

	.bg-pattern {
		background:
			linear-gradient(rgba(255, 255, 255, 0.7), rgba(255, 255, 255, 0.7)), url('/gplay.webp'); /* #008c8c; */
	}

	/* .content-height {
		min-height: calc(100vh - 57px);
	} */
</style>
