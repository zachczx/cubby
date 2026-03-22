<script lang="ts">
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import Icon from '@iconify/svelte';
	import { topLevelRoutes } from './nav';
	import Refresher from '$lib/components/Refresher.svelte';
	import { invalidateAll } from '$app/navigation';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { getAllTrackersQueryKey } from '$lib/queries';
	import Logo from '$lib/assets/transparent-top.webp?w=200&enhanced';

	let {
		children,
		title,
		back = true,
		focusedScreen = false,
		showSettingsIcon = true,
		trackerMuteToggle
	}: {
		children: Snippet;
		title: string | undefined;
		back?: boolean;
		focusedScreen?: boolean;
		focusedScreenAction?: Snippet;
		showSettingsIcon?: boolean;
		trackerMuteToggle?: { trackerId: string | undefined; isMuted: boolean | undefined };
	} = $props();

	let currentSection = $derived.by(() => {
		const p = page.url.pathname;

		if (p === '/app' || p === '/app/') return 'home';

		const animation = topLevelRoutes.animation.find(
			(route) => route.href !== '/app' && p.startsWith(route.href)
		);

		if (animation) {
			return animation.id;
		}

		const noAnimation = topLevelRoutes.noAnimation.find((route) => route.href === p);

		if (noAnimation) {
			return noAnimation.id;
		}

		return '';
	});

	const defaultTitle = 'Cubby';
	const tanstackClient = useQueryClient();

	async function muteHandler(changeTo: 'mute' | 'unmute') {
		const trackerId = trackerMuteToggle?.trackerId;
		if (!trackerId) return;
		const response = await api.post(`trackers/${trackerId}/toggle-mute`, {
			body: JSON.stringify({
				isMuted: changeTo === 'mute'
			})
		});

		if (response.status === 204) {
			const message = changeTo === 'mute' ? 'Muted notifications' : 'Enabled notifications';
			addToast('success', message);

			tanstackClient.setQueryData(
				getAllTrackersQueryKey(),
				(oldTrackers: TrackerDB[] | undefined) =>
					oldTrackers?.map((t) => {
						if (t.id === trackerId) {
							console.log('here');
							return { ...t, isMuted: changeTo === 'mute' };
						}
						return t;
					})
			);
		}
	}
</script>

<svelte:head>
	<title>{title ? title : defaultTitle}</title>
</svelte:head>

<div class="grid min-h-dvh w-full justify-items-center">
	<header
		class={[
			'bg-base-100 text-base-content fixed top-0 z-1 w-full px-4 lg:hidden',
			back ? 'pt-2 pb-3' : 'py-2'
		]}
		style="view-transition-name: top-nav"
	>
		<div class="flex min-h-10 items-center justify-between">
			{#if back}
				<button
					aria-label="go back"
					class="-ml-2 cursor-pointer rounded-full p-2"
					onclick={() => {
						if (window) {
							history.back();
						}
					}}
				>
					<Icon icon="material-symbols:arrow-back" class="size-5" />
				</button>
			{:else}
				<div class="flex items-center gap-2">
					<div class="avatar">
						<div class="bg-primary/50 border-primary/50 w-10 rounded-full border-2">
							<enhanced:img src={Logo} alt="Cubby" class="scale-150" />
						</div>
					</div>

					<h1 class="text-xl font-bold">{title ? title : defaultTitle}</h1>
				</div>
			{/if}

			<div class="bg-base-content/5 flex items-center gap-1 rounded-full px-1">
				{#key trackerMuteToggle}
					{#if trackerMuteToggle}
						{#if trackerMuteToggle?.isMuted}
							<button class="btn btn-ghost btn-sm px-2" onclick={() => muteHandler('unmute')}>
								<Icon icon="material-symbols:notifications-off" class="size-5 opacity-80" />
							</button>
						{:else}
							<button class="btn btn-ghost btn-sm px-2" onclick={() => muteHandler('mute')}>
								<Icon icon="material-symbols:notifications" class="size-5" />
							</button>
						{/if}
					{/if}
				{/key}
				{#if showSettingsIcon}
					<a href="/app/profile" class="btn btn-ghost btn-sm px-2">
						<Icon icon="material-symbols:settings" class="size-5" />
					</a>
				{/if}
			</div>
		</div>

		{#if back}
			<h1 class="mt-1 text-[1.75rem] leading-tight font-bold">
				{title ? title : defaultTitle}
			</h1>
		{/if}
	</header>

	<div
		class={[
			'navbar border-b-base-300/50 bg-base-nav text-primary-content fixed top-0 z-1 hidden min-h-14 items-center border-b pe-4 lg:flex',
			back ? 'lg:ps-4' : 'ps-4'
		]}
		style="view-transition-name: top-nav-desktop"
	>
		<div class="navbar-start grow">
			<a class="text-xl font-bold" href="/app">Cubby</a>
		</div>
		<div id="desktop-menu" class="navbar-center">
			<ul class="menu menu-horizontal gap-8 px-1 text-lg">
				{#each topLevelRoutes.animation as route}
					{#if route.desktopNav}
						<li>
							<a
								href={route.href}
								class={[
									'px-4 py-2',
									currentSection === route.id && 'rounded-full bg-white/30 font-bold',
									currentSection !== route.id && 'rounded-full hover:bg-white/20'
								]}>{route.label}</a
							>
						</li>
					{/if}
				{/each}
			</ul>
		</div>
		<div class="navbar-end flex">
			{#key trackerMuteToggle}
				{#if trackerMuteToggle}
					{#if trackerMuteToggle?.isMuted}
						<button class="btn btn-ghost px-2 py-0" onclick={() => muteHandler('unmute')}>
							<Icon icon="material-symbols:notifications-off" class="size-6 opacity-80" />
						</button>
					{:else}
						<button class="btn btn-ghost px-2 py-0" onclick={() => muteHandler('mute')}>
							<Icon icon="material-symbols:notifications" class="size-6" />
						</button>
					{/if}
				{/if}
			{/key}
			{#if showSettingsIcon}
				<div id="desktop-logout" class="flex items-center text-sm">
					<a href="/app/profile" class="btn btn-ghost px-2 py-0"
						><Icon icon="material-symbols:settings" class="size-6" /></a
					><a href="/logout" class="btn btn-outline btn-sm ms-2">Logout</a>
				</div>
			{/if}
		</div>
	</div>

	<Refresher onRefresh={() => invalidateAll()}>
		<div
			class={[
				'bg-base-100 w-full p-4 max-lg:pb-24 lg:grid lg:px-12',
				back
					? 'max-lg:mt-24 max-lg:min-h-[calc(100vh-6rem-6rem)]'
					: 'max-lg:mt-14 max-lg:min-h-[calc(100vh-3.5rem-6rem)]',
				'lg:mt-14 lg:min-h-[calc(100vh-3.5rem-1rem)] lg:pt-16'
			]}
			style="view-transition-name: content;"
		>
			{@render children?.()}
		</div>
	</Refresher>

	{#if !focusedScreen}
		<nav
			class={[
				'fixed bottom-4 left-1/2 z-10 -translate-x-1/2 lg:hidden',
				'bg-base-content/85 flex items-center gap-1 rounded-full px-3 py-2 shadow-lg backdrop-blur-md',
				title === 'Login' ? 'hidden' : undefined
			]}
			style="view-transition-name: bottom-nav"
		>
			{#each topLevelRoutes.animation as route}
				<a
					href={route.href}
					aria-current={currentSection === route.id ? 'page' : undefined}
					class={[
						'flex flex-col items-center gap-0.5 rounded-full px-5 py-2 transition-colors',
						currentSection === route.id ? 'text-primary bg-white font-bold' : 'text-white/80'
					]}
				>
					{#if route.icon}
						<Icon icon={route.icon} class={[currentSection === route.id ? 'size-6' : 'size-5']} />
					{/if}
					<span class={['text-xs', currentSection === route.id ? 'tracking-wider' : undefined]}
						>{route.label}</span
					>
				</a>
			{/each}
		</nav>
	{/if}
</div>
