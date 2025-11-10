<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
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
	import { getNotificationStatus } from '$lib/notification';
	import { goto } from '$app/navigation';
	import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
	import IconParkSolidBottleOne from '$lib/assets/svg/IconParkSolidBottleOne.svelte';
	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';

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

		return user.data?.sprayIntervalDays;
	});

	let gummyDaysToNext = $derived.by(() => {
		if (user.isPending) {
			return undefined;
		}

		return user.data?.gummyIntervalDays;
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

<PageWrapper title="Personal" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<ActionCard
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
				></ActionCard>

				<ActionCard
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
				></ActionCard>
			</div>
		</div>
	</main>
</PageWrapper>
