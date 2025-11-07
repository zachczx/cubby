<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { dirtyTowelDays } from '$lib/config';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createGummyQueryOptions,
		createGummyRefetchOptions,
		createSprayQueryOptions,
		createSprayRefetchOptions,
		createTowelQueryOptions,
		createTowelRefetchOptions,
		createUserQueryOptions
	} from '$lib/queries';
	import ReshotTimerOff from '$lib/assets/towsey-icons/ReshotTimerOff.svelte';
	import StreamlineKameleonColorMedicine from '$lib/assets/towsey-icons/StreamlineKameleonColorMedicine.svelte';
	import StreamlineKameleonColorTowel from '$lib/assets/towsey-icons/StreamlineKameleonColorTowel.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import { getNotificationStatus } from '$lib/notification';
	import { goto } from '$app/navigation';
	import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
	import IconParkSolidBottleOne from '$lib/assets/svg/IconParkSolidBottleOne.svelte';
	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import ActionButton from '$lib/ui/ActionButton.svelte';
	import type { Component } from 'svelte';
	import type { RecordModel } from 'pocketbase';

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

		return user.data?.defaultSprayInterval;
	});

	let gummyDaysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.defaultGummyInterval;
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
</script>

<svelte:head>
	<title>Towsey</title>
</svelte:head>

<PageWrapper title="Sundry" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				{@render actionCard({
					title: 'Wash Towel',
					query: towels,
					notification: towelNotification,
					route: '/towel',
					icon: PhTowelFill,
					last: towelLast,
					button: {
						query: towelQuery,
						refetch: towelRefetch,
						status: towelButtonStatus,
						text: 'Just Washed'
					}
				})}

				{@render actionCard({
					title: 'Spray Nose',
					query: sprays,
					notification: sprayNotification,
					route: '/spray',
					icon: IconParkSolidBottleOne,
					last: sprayLast,
					button: {
						query: sprayQuery,
						refetch: sprayRefetch,
						status: sprayButtonStatus,
						text: 'Just Sprayed'
					}
				})}

				{@render actionCard({
					title: 'Elderberry Gummy',
					query: gummies,
					notification: gummyNotification,
					route: '/gummy',
					icon: MaterialSymbolsHealthAndSafety,
					last: gummyLast,
					button: {
						query: gummyQuery,
						refetch: gummyRefetch,
						status: gummyButtonStatus,
						text: 'Just Ate Gummy'
					}
				})}

				<div class="bg-base-100 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4">
					<a href="/count" class="flex items-center">
						<div class="flex grow items-center gap-4">
							<MaterialSymbolsTimer class="size-12 opacity-75" />

							<div>
								<p class="text-xl font-bold">Timer</p>
							</div>
						</div>
						<div class="flex h-full items-center">
							<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
								><MaterialSymbolsChevronRight class="size-6" /></button
							>
						</div>
					</a>
					<button
						class="btn btn-lg btn-secondary flex w-full items-center gap-2 rounded-full"
						onclick={() => goto('/count?start=true')}
					>
						Start Timer
					</button>
				</div>
			</div>
		</div>
	</main>
</PageWrapper>

{#snippet actionCard(options: {
	title: string;
	query: Query;
	notification: NotificationStatus;
	icon: Component;
	route: string;
	last: string;
	button: {
		query: () => Promise<RecordModel>;
		refetch: () => Promise<void>;
		status: ButtonState;
		text: string;
	};
})}
	<section
		class={[
			'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
			options.notification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-base-100'
		]}
	>
		<a href={options.route} class="flex items-center">
			<div class="flex grow items-center gap-4">
				<options.icon class="size-12 opacity-75" />
				<div>
					<p class="text-xl font-bold">{options.title}</p>
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
		<ActionButton
			query={options.button.query}
			refetch={options.button.refetch}
			text={options.button.text}
		/>
	</section>
{/snippet}
