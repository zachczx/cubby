<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { createTowelQueryOptions, createTowelRefetchOptions } from '$lib/queries';
	import { getNotificationStatus } from '$lib/notification';
	import { goto } from '$app/navigation';

	import PhTowelFill from '$lib/assets/svg/PhTowelFill.svelte';
	import MaterialSymbolsTimer from '$lib/assets/svg/MaterialSymbolsTimer.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelButtonStatus: ButtonState = $state('default');

	const towels = createQuery(createTowelQueryOptions);
	const tanstackClient = useQueryClient();

	let towelLast: string = $derived.by(() => {
		if (towels.isSuccess && towels.data.length > 0) {
			return dayjs(towels.data[0].time).fromNow();
		}
		return '';
	});

	const towelQuery = async () =>
		await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});
	const towelRefetch = async () => await tanstackClient.refetchQueries(createTowelRefetchOptions());

	let towelNotification = $derived.by(() => getNotificationStatus(towels));
</script>

<PageWrapper title="Sundry" back={false} {pb}>
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
							query: towelQuery,
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
						class="btn btn-lg btn-secondary flex w-full items-center gap-2 rounded-full"
						onclick={() => goto('/household/count?start=true')}
					>
						Start Timer
					</button>
				</div>
			</div>
		</div>
	</main>
</PageWrapper>
