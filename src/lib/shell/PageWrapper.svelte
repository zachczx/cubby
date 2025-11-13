<script lang="ts">
	import { page } from '$app/state';
	import { SafeArea } from '@capacitor-community/safe-area';

	import type Client from 'pocketbase';
	import type { Snippet } from 'svelte';
	import MaterialSymbolsClose from '$lib/assets/svg/MaterialSymbolsClose.svelte';
	import MaterialSymbolsLogin from '$lib/assets/svg/MaterialSymbolsLogin.svelte';
	import MaterialSymbolsLogout from '$lib/assets/svg/MaterialSymbolsLogout.svelte';
	import MaterialSymbolsSettings from '$lib/assets/svg/MaterialSymbolsSettings.svelte';
	import MaterialSymbolsNotifications from '$lib/assets/svg/MaterialSymbolsNotifications.svelte';
	import MaterialSymbolsWarning from '$lib/assets/svg/MaterialSymbolsWarning.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsNotificationImportant from '$lib/assets/svg/MaterialSymbolsNotificationImportant.svelte';
	import {
		createGummyQueryOptions,
		createSprayQueryOptions,
		createTowelQueryOptions
	} from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import { getNotificationCount, getNotificationStatus } from '$lib/notification';
	import MaterialSymbolsHome from '$lib/assets/svg/MaterialSymbolsHome.svelte';
	import MaterialSymbolsHealthCross from '$lib/assets/svg/MaterialSymbolsHealthCross.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
	import LucideLabTowelRack from '$lib/assets/svg/LucideLabTowelRack.svelte';
	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import IconParkSolidBottleOne from '$lib/assets/svg/IconParkSolidBottleOne.svelte';
	import MaterialSymbolsCleaningServices from '$lib/assets/svg/MaterialSymbolsCleaningServices.svelte';
	import MaterialSymbolsPets from '$lib/assets/svg/MaterialSymbolsPets.svelte';
	import { topLevelRoutes } from './nav';

	let {
		pb,
		children,
		title,
		back = true,
		largeScreenCenter
	}: {
		pb: Client;
		children: Snippet;
		title: string;
		back?: boolean;
		largeScreenCenter?: boolean;
	} = $props();

	let currentSection = $derived.by(() => {
		const p = page.url.pathname;

		if (p.includes('personal')) {
			return 'personal';
		} else if (p.includes('pet')) {
			return 'pet';
		} else if (p.includes('household')) {
			return 'household';
		} else if (p.endsWith('')) {
			return 'home';
		} else {
			return '';
		}
	});

	const towels = createQuery(createTowelQueryOptions);
	const sprays = createQuery(createSprayQueryOptions);
	const gummies = createQuery(createGummyQueryOptions);
	let sprayNotification = $derived.by(() => getNotificationStatus(sprays));
	let towelNotification = $derived.by(() => getNotificationStatus(towels));
	let gummyNotification = $derived.by(() => getNotificationStatus(gummies));
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>

<div class="grid h-dvh w-full justify-items-center">
	<div
		class={[
			'navbar border-b-base-300/50 bg-base-nav text-primary-content fixed top-0 z-1 min-h-14 items-center border-b pe-4',
			back ? 'lg:ps-4' : 'ps-4'
		]}
		style="view-transition-name: top-nav"
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

			<a class="hidden text-xl font-bold lg:flex" href="/">Sundry</a>
			<span class="text-xl font-bold lg:hidden">{title}</span>
		</div>
		<div id="desktop-menu" class="navbar-center hidden lg:flex">
			<ul class="menu menu-horizontal gap-8 px-1 text-lg">
				{#each topLevelRoutes as route}
					<li>
						<a
							href="/"
							class={[
								'px-4 py-2',
								currentSection === route.id && 'rounded-full bg-white/30 font-bold',
								currentSection !== route.id && 'rounded-full hover:bg-white/20'
							]}>{route.label}</a
						>
					</li>
				{/each}
			</ul>
		</div>
		<div class="navbar-end">
			<div id="mobile-hamburger" class="dropdown flex items-center lg:hidden">
				{@render notification(sprayNotification, towelNotification, gummyNotification)}
				<a href="/profile" class="btn btn-ghost -me-4"
					><MaterialSymbolsSettings class="size-[1.5em]" /></a
				>
			</div>
			<div id="desktop-logout" class="hidden items-center text-sm lg:flex">
				{@render notification(sprayNotification, towelNotification, gummyNotification)}
				{#if pb.authStore.isValid}
					<a href="/profile" class="btn btn-ghost"><MaterialSymbolsSettings class="size-6" /></a><a
						href="/logout"
						class="btn btn-outline btn-sm ms-2">Logout</a
					>
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
			'bg-base-100 max-lg:min-h-[calc(100vh - 3.5rem - 6rem)] lg:min-h-[calc(100vh - 3.5rem - 1rem)] mt-14 w-full p-4 max-lg:pb-24 lg:grid lg:px-12 lg:pt-12',
			largeScreenCenter && 'lg:content-center'
		]}
		style="view-transition-name: content;"
	>
		{@render children?.()}
	</div>

	<nav
		class={[
			'dock border-t-base-content/15 fixed h-20 border-t-2 pb-2 lg:hidden',
			title === 'Login' ? 'hidden' : undefined
		]}
		style="view-transition-name: bottom-nav"
	>
		<a href="/" class={[currentSection === 'home' && 'text-primary font-semibold']}>
			<MaterialSymbolsHome class="size-[1.5em]" />
			<span class="text-sm tracking-wider">Home</span>
		</a>

		<a href="/personal" class={[currentSection === 'personal' && 'text-primary font-semibold']}>
			<MaterialSymbolsHealthAndSafety class="size-[1.5em]" />
			<span class="text-sm tracking-wider">Personal</span>
		</a>

		<a href="/household" class={[currentSection === 'household' && 'text-primary font-semibold']}>
			<MaterialSymbolsCleaningServices class="size-[1.5em]" />
			<span class="text-sm tracking-wider">Household</span>
		</a>

		<a href="/pet" class={[currentSection === 'pet' && 'text-primary font-semibold']}>
			<MaterialSymbolsPets class="size-[1.5em]" />
			<span class="text-sm tracking-wider">Pet</span>
		</a>
	</nav>
</div>

{#snippet notification(
	sprayNotification: NotificationStatus,
	towelNotification: NotificationStatus,
	gummyNotification: NotificationStatus
)}
	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost drawer-button px-2 py-0">
			{#if !sprayNotification.show && !towelNotification.show}
				<MaterialSymbolsNotifications class="size-6" />
			{:else}
				{@const count = getNotificationCount(
					sprayNotification,
					towelNotification,
					gummyNotification
				)}
				<MaterialSymbolsNotificationImportant class="size-6" />
				<span class="absolute top-1 right-1 z-2 size-4 rounded-full bg-red-600 text-xs"
					>{count}</span
				>
			{/if}
		</div>
		<ul
			class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 w-72 p-2 shadow-lg"
		>
			{#if !sprayNotification.show && !towelNotification.show && !gummyNotification.show}
				<li>
					<div class="flex items-center justify-center gap-2">
						<MaterialSymbolsCheckCircle class="size-[1.3em]" /><span>No pending items</span>
					</div>
				</li>
			{/if}

			{#if sprayNotification.show}
				<li>
					<a href="/spray" class="flex items-center">
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
					<a href="/towel" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />Wash your towel!
						</div>
						<div>
							<MaterialSymbolsChevronRight class="size-5 opacity-50" />
						</div>
					</a>
				</li>
			{/if}

			{#if gummyNotification.show}
				<li>
					<a href="/gummy" class="flex items-center">
						<div class="flex grow items-center gap-2">
							<MaterialSymbolsWarning class="size-[1.3em]" />Eat your gummy!
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
	.dock {
		& > * {
			transition: none !important;

			&:focus-within {
				transform: none;
			}
		}
	}
</style>
