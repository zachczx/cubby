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
		logsQueryOptions,
		logsRefetchOptions,
		createLogsQuery,
		trackerQueryOptions,
		notificationQueryOptions,
		trackerNameToId
	} from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';
	import { goto } from '$app/navigation';

	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelButtonStatus: ButtonState = $state('default');

	const towels = createQuery(() => logsQueryOptions('towel'));
	const towelTracker = createQuery(() => trackerQueryOptions('towel'));
	let towelInterval = $derived.by(() =>
		towelTracker.isSuccess ? towelTracker.data?.interval : undefined
	);
	let towelIntervalUnit = $derived.by(() =>
		towelTracker.isSuccess ? towelTracker.data?.intervalUnit : undefined
	);
	const createTowelRecord = () =>
		createLogsQuery({
			collectionName: 'towel',
			interval: towelInterval,
			intervalUnit: towelIntervalUnit
		});
	const towelRefetch = async () => await tanstackClient.refetchQueries(logsRefetchOptions('towel'));
	const tanstackClient = useQueryClient();

	let towelLast: string = $derived.by(() => {
		if (towels.isSuccess && towels.data.length > 0) {
			return dayjs(towels.data[0].time).fromNow();
		}
		return '';
	});

	const latestLogs = createQuery(notificationQueryOptions);
	let towelNotification = $derived.by(() => {
		if (!latestLogs.isSuccess) return;

		const towel = latestLogs.data.find((item) => item.tracker === trackerNameToId('towel'));
		return getNotificationStatus(towel);
	});
</script>

<PageWrapper title="Household" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<ActionCard
					options={{
						title: 'Wash Towel',
						query: towels,
						notification: towelNotification,
						route: '/household/towel',
						icon: PhTowelFill,
						last: towelLast,
						button: {
							query: createTowelRecord,
							refetch: towelRefetch,
							status: towelButtonStatus,
							text: 'Washed'
						}
					}}
				></ActionCard>

				<div class="bg-base-100 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4">
					<a href="/household/count" class="flex items-center">
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
						class="btn btn-lg btn-primary flex w-full items-center gap-2 rounded-full"
						onclick={() => goto('/household/count?start=true')}
					>
						Start Timer
					</button>
				</div>
			</div>
		</div>
	</main>
</PageWrapper>
