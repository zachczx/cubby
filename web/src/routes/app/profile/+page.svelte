<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import Icon from '@iconify/svelte';
</script>

<PageWrapper title="Profile" hideTitle>
	<div
		class="lg:bg-base-50 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-md lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="mb-4 text-center text-3xl font-bold max-lg:hidden">Profile</h1>
			<div class="grid w-full content-start">
				{@render menuItem('account', 'Account Info')}
				{@render menuItem('trackers', 'Manage Trackers')}
				{@render menuItem('family', 'Manage Family')}
				{@render menuItem('vacation', 'Manage Vacation Dates')}
			</div>
		</div>
		<div class="mt-4 text-lg">
			<a href="/logout" class="btn btn-neutral btn-lg w-full rounded-full">Logout</a>
		</div>
	</div>
</PageWrapper>

{#snippet menuItem(href: string, title: string)}
	<a
		href={'/app/profile/' + href}
		class="border-b-base-300 flex items-center border-b py-6 text-lg"
	>
		<div class="grow font-semibold">{title}</div>
		<div class="flex h-full items-center">
			<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
				><Icon icon="material-symbols:chevron-right" class="size-7" /></button
			>
		</div>
	</a>
{/snippet}
