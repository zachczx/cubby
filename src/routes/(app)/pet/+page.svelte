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
		createDoggoBathRefetchOptions,
		createDoggoChewableQueryOptions,
		createDoggoChewableRefetchOptions,
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
	import MaterialSymbolsShower from '$lib/assets/svg/MaterialSymbolsShower.svelte';
	import MaterialSymbolsPill from '$lib/assets/svg/MaterialSymbolsPill.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let doggoBathButtonStatus: ButtonState = $state('default');
	let doggoChewableButtonStatus: ButtonState = $state('default');

	const doggoBaths = createQuery(createDoggoBathQueryOptions);
	const doggoChewables = createQuery(createDoggoChewableQueryOptions);
	const user = createQuery(createUserQueryOptions);
	const tanstackClient = useQueryClient();

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
	const doggoBathRefetch = async () =>
		await tanstackClient.refetchQueries(createDoggoBathRefetchOptions());

	const doggoChewableCreateQuery = async () =>
		await pb.collection('doggoChewable').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			monthsToNext: doggoChewableMonthsToNext
		});
	const doggoChewableRefetch = async () =>
		await tanstackClient.refetchQueries(createDoggoChewableRefetchOptions());

	let doggoBathNotification = $derived.by(() => getNotificationStatus(doggoBaths));
	let doggoChewableNotification = $derived.by(() => getNotificationStatus(doggoChewables));
</script>

<PageWrapper title="Sundry" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<ActionCard
					options={{
						title: 'Bath',
						query: doggoBaths,
						notification: doggoBathNotification,
						route: '/pet/bath',
						icon: MaterialSymbolsShower,
						last: doggoBathLast,
						button: {
							query: doggoBathCreateQuery,
							refetch: doggoBathRefetch,
							status: doggoBathButtonStatus,
							text: 'Bathed'
						}
					}}
				></ActionCard>

				<ActionCard
					options={{
						title: 'Chewable',
						query: doggoChewables,
						notification: doggoChewableNotification,
						route: '/pet/chewable',
						icon: MaterialSymbolsPill,
						last: doggoChewableLast,
						button: {
							query: doggoChewableCreateQuery,
							refetch: doggoChewableRefetch,
							status: doggoChewableButtonStatus,
							text: 'Fed'
						}
					}}
				></ActionCard>
			</div>
		</div>
	</main>
</PageWrapper>
