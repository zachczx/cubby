<script lang="ts">
	import { page } from '$app/state';
	import type Client from 'pocketbase';
	import type { Snippet } from 'svelte';
	import MaterialSymbolsSettings from '$lib/assets/svg/MaterialSymbolsSettings.svelte';
	import { topLevelRoutes } from './nav';
	import FeedDropdown from '$lib/ui/FeedDropdown.svelte';
	import FilterDropdown from '$lib/ui/FilterDropdown.svelte';

	let {
		pb,
		children,
		title,
		back = true,
		largeScreenCenter
	}: {
		pb: Client;
		children: Snippet;
		title: string | undefined;
		back?: boolean;
		largeScreenCenter?: boolean;
	} = $props();

	let currentSection = $derived.by(() => {
		const p = page.url.pathname;

		if (p === '/') return 'home';

		const animation = topLevelRoutes.animation.find(
			(route) => route.href !== '/' && p.startsWith(route.href)
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
</script>

<svelte:head>
	<title>{title ? title : defaultTitle}</title>
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

			<a class="hidden text-xl font-bold lg:flex" href="/">Cubby</a>
			<span class="text-xl font-bold lg:hidden">{title ? title : defaultTitle}</span>
		</div>
		<div id="desktop-menu" class="navbar-center hidden lg:flex">
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
		<div class="navbar-end">
			<div id="mobile-hamburger" class="dropdown flex items-center lg:hidden">
				<FeedDropdown />
				<FilterDropdown />
				<a href="/profile" class="btn btn-ghost -me-4 px-2 py-0"
					><MaterialSymbolsSettings class="size-6" /></a
				>
			</div>
			<div id="desktop-logout" class="hidden items-center text-sm lg:flex">
				<FeedDropdown />
				<FilterDropdown />
				{#if pb.authStore.isValid}
					<a href="/profile" class="btn btn-ghost px-2 py-0"
						><MaterialSymbolsSettings class="size-6" /></a
					><a href="/logout" class="btn btn-outline btn-sm ms-2">Logout</a>
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
		{#each topLevelRoutes.animation as route}
			<a href={route.href} aria-current={currentSection === route.id ? 'page' : undefined}>
				<route.icon class="size-[1.5em]" />
				<span class="text-sm tracking-wider">{route.label}</span>
			</a>
		{/each}
	</nav>
</div>

<style>
	.dock {
		a[aria-current='page'] {
			font-weight: bold;
			color: var(--color-primary);

			&::after {
				view-transition-name: activepage;
				position: absolute;
				top: bottom;
				left: 50%;
				width: 2rem;
				transform: translateX(-50%);
				height: 0.25rem;
				background-color: var(--color-primary);
			}
		}
	}
</style>
