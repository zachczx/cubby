<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		userQueryOptions,
		userRefetchOptions,
		vacationQueryOptions,
		vacationRefetchOptions
	} from '$lib/queries';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import type { ChangeEventHandler } from 'svelte/elements';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { MenuItem } from '@ark-ui/svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(userQueryOptions);

	const tanstackClient = useQueryClient();

	let sound = $derived.by(() => (user.isSuccess ? user.data?.sound : undefined));

	async function onchange(evt: Event) {
		const target = evt.target;

		if (target instanceof HTMLInputElement && pb.authStore.record?.id) {
			try {
				const data = {
					sound: target.checked ? true : false
				};

				const response = await pb.collection('users').update(pb.authStore.record.id, data);

				if (!response.status) {
					addToast('success', 'Saved!');
					await tanstackClient.refetchQueries(userRefetchOptions());
				}
			} catch (err) {
				addToast('error', 'Error saving!');
			}
		} else {
			addToast('error', 'Error saving!');
		}
	}
</script>

<PageWrapper title="Profile" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<div class="overflow-y-auto">
			<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Profile</h1>
			<div class="grid w-full content-start">
				{@render menuItem('active', 'Active Trackers')}
				{@render menuItem('schedule', 'Tracker Schedules')}
				{@render menuItem('family', 'Manage Family')}
				{@render menuItem('vacation', 'Manage Vacation Dates')}

				<form class="border-b-base-300 flex items-center border-b pt-4 pb-6 text-lg" {onchange}>
					<legend class="fieldset-legend grow">Sound</legend>
					<input
						type="checkbox"
						name="sound"
						class="toggle toggle-lg toggle-primary"
						bind:checked={sound}
					/>
				</form>
			</div>
		</div>
		<div class="text-lg">
			<a href="/logout" class="btn btn-neutral btn-lg w-full rounded-full">Logout</a>
		</div>
	</div>
</PageWrapper>

{#snippet menuItem(href: string, title: string)}
	<a
		href={'/profile/' + href}
		class="border-b-base-300 flex items-center border-b pt-4 pb-6 text-lg"
	>
		<div class="grow font-semibold">{title}</div>
		<div class="flex h-full items-center">
			<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
				><MaterialSymbolsChevronRight class="size-7" /></button
			>
		</div>
	</a>
{/snippet}
