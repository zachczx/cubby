<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient, type CreateQueryResult } from '@tanstack/svelte-query';
	import {
		createUserQueryOptions,
		logsQueryOptions,
		notificationQueryOptions,
		trackerNameToId
	} from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';
	import { goto } from '$app/navigation';
	import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
	import IconParkSolidBottleOne from '$lib/assets/svg/IconParkSolidBottleOne.svelte';
	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import ActionCardCompact from '$lib/ui/ActionCardCompact.svelte';
	import { type RecordModel } from 'pocketbase';
	import { type Component } from 'svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsPill from '$lib/assets/svg/MaterialSymbolsPill.svelte';
	import MaterialSymbolsShower from '$lib/assets/svg/MaterialSymbolsShower.svelte';
	import MaterialSymbolsFlight from '$lib/assets/svg/MaterialSymbolsFlight.svelte';
	import EmptyState from '$lib/assets/svg/EmptyState.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelButtonStatus: ButtonState = $state('default');
	let sprayButtonStatus: ButtonState = $state('default');
	let gummyButtonStatus: ButtonState = $state('default');

	const user = createQuery(createUserQueryOptions);
	const tanstackClient = useQueryClient();

	const latestLogs = createQuery(notificationQueryOptions);

	// Upcoming
	const doggoBath = $derived.by(() => {
		if (latestLogs.isSuccess) {
			return latestLogs.data.find((item) => item.tracker === trackerNameToId('doggoBath'));
		}
		return undefined;
	});
	const doggoChewable = $derived.by(() => {
		if (latestLogs.isSuccess) {
			return latestLogs.data.find((item) => item.tracker === trackerNameToId('doggoChewable'));
		}
		return undefined;
	});

	let doggoBathDaysToNext = $derived(doggoBath?.interval);

	let doggoChewableMonthsToNext = $derived(doggoChewable?.interval);

	let doggoBathLast: string = $derived.by(() => {
		if (doggoBath && doggoBath.time) return dayjs(doggoBath.time).fromNow();

		return '';
	});

	let doggoChewableLast: string = $derived.by(() => {
		if (doggoChewable && doggoChewable.time) return dayjs(doggoChewable.time).fromNow();

		return '';
	});

	const doggoBathCreateQuery = async () =>
		await pb.collection('doggoBath').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: doggoBathDaysToNext
		});

	const doggoChewableCreateQuery = async () =>
		await pb.collection('doggoChewable').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			monthsToNext: doggoChewableMonthsToNext
		});

	let doggoBathNotification = $derived.by(() => getNotificationStatus(doggoBath));
	let doggoChewableNotification = $derived.by(() => getNotificationStatus(doggoChewable));
</script>

<PageWrapper title="Sundry" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Key Actions</h2>
				<ActionCardCompact
					options={{
						collectionName: 'towel',
						title: 'Wash Towel',
						route: '/household/towel',
						icon: PhTowelFill,
						button: {
							status: towelButtonStatus,
							text: 'Washed'
						}
					}}
				></ActionCardCompact>

				<ActionCardCompact
					options={{
						collectionName: 'spray',
						title: 'Spray Nose',
						route: '/personal/spray',
						icon: IconParkSolidBottleOne,
						button: {
							status: sprayButtonStatus,
							text: 'Sprayed'
						}
					}}
				></ActionCardCompact>

				<ActionCardCompact
					options={{
						collectionName: 'gummy',
						title: 'Elderberry Gummy',
						route: '/personal/gummy',
						icon: MaterialSymbolsHealthAndSafety,
						button: {
							status: gummyButtonStatus,
							text: 'Ate'
						}
					}}
				></ActionCardCompact>
			</section>

			<section class="grid gap-0 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">General Tasks</h2>

				{#if doggoBathNotification.show}
					{@render upcomingCard({
						title: 'Bath',
						data: doggoBath,
						route: '/pet/bath',
						icon: MaterialSymbolsShower,
						last: doggoBathLast,
						notification: doggoBathNotification
					})}
				{/if}

				{#if doggoChewableNotification.show}
					{@render upcomingCard({
						title: 'Chewable',
						data: doggoChewable,
						route: '/pet/chewable',
						icon: MaterialSymbolsPill,
						last: doggoChewableLast,
						notification: doggoChewableNotification
					})}
				{/if}

				{#if !doggoBathNotification.show && !doggoChewableNotification.show}
					<enhanced:img src={EmptyCorgi} alt="nothing" class="justify-self-center" />
					<p class="text-center">No other upcoming tasks!</p>
				{/if}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Quick Links</h2>
				<div class="flex items-center gap-6">
					<a
						href="/household/count"
						class="group focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-20 content-center justify-items-center gap-1 rounded-3xl p-2 text-sm font-semibold"
					>
						<div
							class="text-base-100 bg-neutral/80 group-active:bg-neutral/90 group-focus-within:bg-neutral/90 rounded-xl p-2"
						>
							<MaterialSymbolsTimer class="size-7" />
						</div>
						Stopwatch</a
					>
					<a
						href="/profile/vacation"
						class="group focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-20 content-center justify-items-center gap-1 rounded-3xl p-2 text-sm font-semibold"
					>
						<div
							class="text-base-100 bg-neutral/80 group-active:bg-neutral/90 group-focus-within:bg-neutral/90 rounded-xl p-2"
						>
							<MaterialSymbolsFlight class="size-7" />
						</div>
						Vacation</a
					>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>

{#snippet upcomingCard(options: {
	title: string;
	data: LogsDB | undefined;
	notification: NotificationStatus;
	icon: Component;
	route: string;
	last: string;
})}
	<section
		class={[
			'border-base-300 grid gap-4 rounded-3xl border px-4 py-2',
			options.notification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-base-100'
		]}
	>
		<a href={options.route} class="flex items-center">
			<div class="flex grow items-center gap-4">
				<options.icon class="size-9 opacity-75" />
				<div>
					<p class="text-xl font-bold">
						{options.title}
					</p>

					{#if options.data}
						{#if options.notification.show}
							<span class="text-error font-medium tracking-tight">
								{#if options.notification.level === 'overdue'}
									Overdue
								{:else if options.notification.level === 'due'}
									Due
								{/if}
							</span>&nbsp;&nbsp;•&nbsp;&nbsp;
						{/if}<span>{options.last}</span>
					{/if}
				</div>
			</div>
			<div class="flex h-full items-center">
				<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
					><MaterialSymbolsChevronRight class="size-6" /></button
				>
			</div>
		</a>
	</section>
{/snippet}
