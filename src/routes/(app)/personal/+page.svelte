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
	import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
	import IconParkSolidBottleOne from '$lib/assets/svg/IconParkSolidBottleOne.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let sprayButtonStatus: ButtonState = $state('default');
	let gummyButtonStatus: ButtonState = $state('default');

	const sprays = createQuery(() => logsQueryOptions('spray'));
	const sprayTracker = createQuery(() => trackerQueryOptions('spray'));
	let sprayInterval = $derived.by(() =>
		sprayTracker.isSuccess ? sprayTracker.data?.interval : undefined
	);
	let sprayIntervalUnit = $derived.by(() =>
		sprayTracker.isSuccess ? sprayTracker.data?.intervalUnit : undefined
	);
	let sprayLast: string = $derived.by(() => {
		if (sprays.isSuccess && sprays.data.length > 0) return dayjs(sprays.data[0].time).fromNow();

		return '';
	});
	const createSprayRecord = () =>
		createLogsQuery({
			collectionName: 'spray',
			interval: sprayInterval,
			intervalUnit: sprayIntervalUnit
		});
	const sprayRefetch = async () => await tanstackClient.refetchQueries(logsRefetchOptions('spray'));

	const gummies = createQuery(() => logsQueryOptions('gummy'));
	const gummyTracker = createQuery(() => trackerQueryOptions('gummy'));
	let gummyInterval = $derived.by(() =>
		gummyTracker.isSuccess ? gummyTracker.data?.interval : undefined
	);
	let gummyIntervalUnit = $derived.by(() =>
		gummyTracker.isSuccess ? gummyTracker.data?.intervalUnit : undefined
	);
	let gummyLast: string = $derived.by(() => {
		if (gummies.isSuccess && gummies.data.length > 0) return dayjs(gummies.data[0].time).fromNow();

		return '';
	});
	const createGummyRecord = () =>
		createLogsQuery({
			collectionName: 'gummy',
			interval: gummyInterval,
			intervalUnit: gummyIntervalUnit
		});
	const gummyRefetch = async () => await tanstackClient.refetchQueries(logsRefetchOptions('gummy'));

	const tanstackClient = useQueryClient();

	const latestLogs = createQuery(notificationQueryOptions);
	let sprayNotification = $derived.by(() => {
		if (!latestLogs.isSuccess) return;

		const spray = latestLogs.data.find((item) => item.tracker === trackerNameToId('spray'));
		return getNotificationStatus(spray);
	});

	let gummyNotification = $derived.by(() => {
		if (!latestLogs.isSuccess) return;

		const gummy = latestLogs.data.find((item) => item.tracker === trackerNameToId('gummy'));
		return getNotificationStatus(gummy);
	});
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
							query: createSprayRecord,
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
							query: createGummyRecord,
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
