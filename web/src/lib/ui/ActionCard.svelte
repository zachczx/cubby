<script lang="ts">
	import ActionButton from './ActionButton.svelte';
	import dayjs from 'dayjs';
	import { createQuery } from '@tanstack/svelte-query';
	import { allTrackersQueryOptions, createEntryQuery } from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';

	import Icon from '@iconify/svelte';

	import { onMount } from 'svelte';
	import { error } from '@sveltejs/kit';

	let { options }: { options: ActionCardOptions } = $props();

	onMount(() => {
		if (!options.title) {
			error(500);
		}
	});

	let size = $derived(options.size ?? 'default');
	let entries = $derived(options.entries);

	const trackers = createQuery(allTrackersQueryOptions);
	let tracker = $derived.by(() =>
		trackers.data?.find((tracker) => tracker.name === options.tracker?.name)
	);

	let notification = $derived.by(() => (entries ? getTrackerStatus(entries) : undefined));

	let interval = $derived(tracker?.interval);
	let intervalUnit = $derived(tracker?.intervalUnit);
	const query = () =>
		createEntryQuery({
			trackerId: tracker?.id ?? '',
			interval: Number(interval),
			intervalUnit: intervalUnit
		});
</script>

{#if size === 'list'}
	<a
		href={options.route}
		class={['grid gap-4 px-2 py-2', !options.lastChild && 'border-b-base-300/50 border-b']}
	>
		<div class="flex items-center">
			<div class="flex grow items-center gap-4">
				<div class="flex size-9 items-center justify-center">
					<options.icon class="size-7 opacity-75" />
				</div>
				<div class="grid gap-0">
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{@render ownerBadge(options)}
					</div>
					{@render notificationLogic()}
				</div>
			</div>

			<div class="flex h-full items-center">
				<Icon icon="material-symbols:chevron-right" class="size-5 opacity-75" />
			</div>
		</div>
	</a>
{/if}

{#if size === 'compact'}
	<section class="border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 px-2 py-2">
		<div class="flex items-center">
			<a href={options.route} class="flex grow items-center gap-4">
				<options.icon class="size-9 opacity-75" />
				<div>
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{@render ownerBadge(options)}
					</div>
					{@render notificationLogic()}
				</div>
			</a>

			<div class="flex h-full items-center">
				<ActionButton
					{query}
					text={options.button.text}
					compact={true}
					color={'primary'}
					icon="material-symbols:check"
					rounded="2xl"
				/>
			</div>
		</div>
	</section>
{/if}

{#if size === 'default'}
	<section class="border-base-300 grid min-h-24 gap-4 rounded-3xl border bg-white/70 p-4">
		<a href={options.route} class="flex items-center">
			<div class="flex grow items-center gap-4">
				<options.icon class="size-12 opacity-75" />
				<div>
					<div
						class="text-base-content/90 flex items-center gap-2 align-baseline text-lg font-bold"
					>
						{@render ownerBadge(options)}
					</div>
					{@render notificationLogic()}
				</div>
			</div>
			<div class="flex h-full items-center">
				<div class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75">
					<Icon icon="material-symbols:chevron-right" class="size-6" />
				</div>
			</div>
		</a>
		<ActionButton {query} text={options.button.text} />
	</section>
{/if}

{#snippet ownerBadge(options: ActionCardOptions)}
	{options.title}
	<!-- {#if options.tracker?.expand?.family?.expand?.owner}
		{@const owner = options.tracker?.expand?.family?.expand?.owner}
		{#if owner.id !== pb.authStore.record?.id}
			<span
				class={[
					'flex aspect-square size-5 items-center justify-center rounded-full p-0 text-sm font-bold',
					bgColor,
					textColor
				]}
			>
				{owner.name.charAt(0)}
			</span>
		{/if}
	{/if} -->
	{#if options.streak && options.streak > 0}
		<div
			class="text-success bg-success/10 flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-bold"
			title="Current streak: {options.streak}"
		>
			{options.streak}x
		</div>
	{/if}
{/snippet}

{#snippet notificationLogic()}
	{#if entries}
		<div class="flex items-center gap-3">
			{#if notification?.show}
				{#if notification.level === 'overdue'}
					<span class="text-error font-bold tracking-tight">Overdue</span>
				{:else if notification.level === 'due'}
					<span class="text-warning font-medium tracking-tight"
						>Due {dayjs(notification?.next).fromNow()}</span
					>
				{/if}
			{:else}
				<span class="text-neutral/70 font-medium tracking-tight"
					>Due {dayjs(notification?.next).fromNow()}</span
				>
			{/if}
		</div>
	{/if}
{/snippet}
