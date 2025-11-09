<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createDoggoBathQueryOptions,
		createDoggoChewableQueryOptions,
		createGummyQueryOptions,
		createGummyRefetchOptions,
		createSprayQueryOptions,
		createSprayRefetchOptions,
		createTowelQueryOptions,
		createTowelRefetchOptions,
		createUserQueryOptions
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

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelButtonStatus: ButtonState = $state('default');
	let sprayButtonStatus: ButtonState = $state('default');
	let gummyButtonStatus: ButtonState = $state('default');

	const towels = createQuery(createTowelQueryOptions);
	const sprays = createQuery(createSprayQueryOptions);
	const gummies = createQuery(createGummyQueryOptions);
	const user = createQuery(createUserQueryOptions);
	const tanstackClient = useQueryClient();

	let towelLast: string = $derived.by(() => {
		if (towels.isSuccess && towels.data.length > 0) {
			return dayjs(towels.data[0].time).fromNow();
		}
		return '';
	});

	let sprayDaysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.sprayInterval;
	});

	let gummyDaysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.gummyInterval;
	});

	let sprayLast: string = $derived.by(() => {
		if (sprays.isSuccess && sprays.data.length > 0) return dayjs(sprays.data[0].time).fromNow();

		return '';
	});

	let gummyLast: string = $derived.by(() => {
		if (gummies.isSuccess && gummies.data.length > 0) return dayjs(gummies.data[0].time).fromNow();

		return '';
	});

	const towelQuery = async () =>
		await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});
	const towelRefetch = async () => await tanstackClient.refetchQueries(createTowelRefetchOptions());

	const sprayQuery = async () =>
		await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: sprayDaysToNext
		});
	const sprayRefetch = async () => await tanstackClient.refetchQueries(createSprayRefetchOptions());

	const gummyQuery = async () =>
		await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: gummyDaysToNext
		});
	const gummyRefetch = async () => await tanstackClient.refetchQueries(createGummyRefetchOptions());

	let sprayNotification = $derived.by(() => getNotificationStatus(sprays));
	let towelNotification = $derived.by(() => getNotificationStatus(towels));
	let gummyNotification = $derived.by(() => getNotificationStatus(gummies));

	// Upcoming
	const doggoBaths = createQuery(createDoggoBathQueryOptions);
	const doggoChewables = createQuery(createDoggoChewableQueryOptions);

	let doggoBathDaysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.doggoBathInterval;
	});

	let doggoChewableMonthsToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.doggoChewableIntervalMonths;
	});

	let doggoBathLast: string = $derived.by(() => {
		if (doggoBaths.isSuccess && doggoBaths.data.length > 0)
			return dayjs(doggoBaths.data[0].time).fromNow();

		return '';
	});

	let doggoChewableLast: string = $derived.by(() => {
		if (doggoChewables.isSuccess && doggoChewables.data.length > 0)
			return dayjs(doggoChewables.data[0].time).fromNow();

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

	let doggoBathNotification = $derived.by(() => getNotificationStatus(doggoBaths));
	let doggoChewableNotification = $derived.by(() => getNotificationStatus(doggoChewables));
</script>

<PageWrapper title="Sundry" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Key Actions</h2>
				<ActionCardCompact
					options={{
						title: 'Wash Towel',
						query: towels,
						notification: towelNotification,
						route: '/household/towel',
						icon: PhTowelFill,
						last: towelLast,
						button: {
							query: towelQuery,
							refetch: towelRefetch,
							status: towelButtonStatus,
							text: 'Washed'
						}
					}}
				></ActionCardCompact>

				<ActionCardCompact
					options={{
						title: 'Spray Nose',
						query: sprays,
						notification: sprayNotification,
						route: '/personal/spray',
						icon: IconParkSolidBottleOne,
						last: sprayLast,
						button: {
							query: sprayQuery,
							refetch: sprayRefetch,
							status: sprayButtonStatus,
							text: 'Sprayed'
						}
					}}
				></ActionCardCompact>

				<ActionCardCompact
					options={{
						title: 'Elderberry Gummy',
						query: gummies,
						notification: gummyNotification,
						route: '/personal/gummy',
						icon: MaterialSymbolsHealthAndSafety,
						last: gummyLast,
						button: {
							query: gummyQuery,
							refetch: gummyRefetch,
							status: gummyButtonStatus,
							text: 'Ate'
						}
					}}
				></ActionCardCompact>
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Upcoming</h2>

				{#if doggoBathNotification.show}
					{@render upcomingCard({
						title: 'Bath',
						query: doggoBaths,
						route: '/pet/bath',
						icon: MaterialSymbolsShower,
						last: doggoBathLast,
						notification: doggoBathNotification
					})}
				{/if}

				{#if doggoChewableNotification.show}
					{@render upcomingCard({
						title: 'Chewable',
						query: doggoChewables,
						route: '/pet/chewable',
						icon: MaterialSymbolsPill,
						last: doggoChewableLast,
						notification: doggoChewableNotification
					})}
				{/if}
			</section>

			<section class="grid gap-2 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Quick Links</h2>
				<div class="flex items-center gap-4">
					<a
						href="/count"
						class="border-base-300 grid aspect-square w-24 content-center justify-items-center gap-0.5 rounded-3xl border p-4 text-xs font-semibold"
					>
						<MaterialSymbolsTimer class="size-7" />
						Stopwatch</a
					>
					<a
						href="/profile?p=vacations"
						class="border-base-300 grid aspect-square w-24 content-center justify-items-center gap-0.5 rounded-3xl border p-4 text-xs font-semibold"
					>
						<MaterialSymbolsFlight class="size-7" />Vacation</a
					>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>

{#snippet upcomingCard(options: {
	title: string;
	query: Query;
	notification: NotificationStatus;
	icon: Component;
	route: string;
	last: string;
	// button: {
	// 	query: () => Promise<RecordModel>;
	// 	refetch: () => Promise<void>;
	// 	status: ButtonState;
	// 	text: string;
	// };
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
					<p class="flex items-center gap-2 text-xl font-bold">
						{options.title}
						<MaterialSymbolsArrowRightAlt class="size-[1em]" />
					</p>
					{#if options.query.isPending && !options.query.data}
						<div class="custom-loader"></div>
					{/if}
					{#if options.query.error}
						An error has occurred:
						{options.query.error.message}
					{/if}
					{#if options.query.isSuccess}
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
