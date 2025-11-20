<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		createLogsQuery,
		logsQueryOptions,
		logsRefetchOptions,
		notificationQueryOptions,
		trackerNameToId,
		trackerQueryOptions
	} from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import MaterialSymbolsShower from '$lib/assets/svg/MaterialSymbolsShower.svelte';
	import MaterialSymbolsPill from '$lib/assets/svg/MaterialSymbolsPill.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let doggoBathButtonStatus: ButtonState = $state('default');
	let doggoChewableButtonStatus: ButtonState = $state('default');

	const doggoBaths = createQuery(() => logsQueryOptions('doggoBath'));
	const doggoBathTracker = createQuery(() => trackerQueryOptions('doggoBath'));
	let doggoBathInterval = $derived.by(() =>
		doggoBathTracker.isSuccess ? doggoBathTracker.data?.interval : undefined
	);
	let doggoBathIntervalUnit = $derived.by(() =>
		doggoBathTracker.isSuccess ? doggoBathTracker.data?.intervalUnit : undefined
	);
	const createDoggoBathRecord = () =>
		createLogsQuery({
			collectionName: 'doggoBath',
			interval: doggoBathInterval,
			intervalUnit: doggoBathIntervalUnit
		});
	const doggoBathRefetch = async () =>
		await tanstackClient.refetchQueries(logsRefetchOptions('doggoBath'));

	const doggoChewables = createQuery(() => logsQueryOptions('doggoChewable'));
	const doggoChewableTracker = createQuery(() => trackerQueryOptions('doggoChewable'));
	let doggoChewableInterval = $derived.by(() =>
		doggoChewableTracker.isSuccess ? doggoChewableTracker.data?.interval : undefined
	);
	let doggoChewableIntervalUnit = $derived.by(() =>
		doggoChewableTracker.isSuccess ? doggoChewableTracker.data?.intervalUnit : undefined
	);
	const createDoggoChewableRecord = () =>
		createLogsQuery({
			collectionName: 'doggoChewable',
			interval: doggoChewableInterval,
			intervalUnit: doggoChewableIntervalUnit
		});
	const doggoChewableRefetch = async () =>
		await tanstackClient.refetchQueries(logsRefetchOptions('doggoChewable'));

	const tanstackClient = useQueryClient();

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

	const latestLogs = createQuery(notificationQueryOptions);
	let doggoBathNotification = $derived.by(() => {
		if (!doggoBaths.isSuccess || !doggoBaths || doggoBaths.data.length === 0) return;
		const latest = doggoBaths.data[0];
		return getNotificationStatus(latest);
	});

	let doggoChewableNotification = $derived.by(() => {
		if (!doggoChewables.isSuccess || !doggoChewables || doggoChewables.data.length === 0) return;
		const latest = doggoChewables.data[0];
		return getNotificationStatus(latest);
	});
</script>

<PageWrapper title="Pet" {pb}>
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
							query: createDoggoBathRecord,
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
							query: createDoggoChewableRecord,
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
