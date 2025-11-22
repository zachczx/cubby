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
	import StreamlineColorHotelLaundryFlat from '$lib/assets/expressive-icons/StreamlineColorHotelLaundryFlat.svelte';
	import FluentEmojiFlatLotionBottle from '$lib/assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
	import FluentEmojiFlatShield from '$lib/assets/expressive-icons/FluentEmojiFlatShield.svelte';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import FluentEmojiFlatAirplane from '$lib/assets/expressive-icons/FluentEmojiFlatAirplane.svelte';
	import FluentEmojiFlatBed from '$lib/assets/expressive-icons/FluentEmojiFlatBed.svelte';
	import FluentEmojiFlatShower from '$lib/assets/expressive-icons/FluentEmojiFlatShower.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $state<Record<string, ButtonState>>({
		towel: 'default',
		spray: 'default',
		gummy: 'default',
		bedsheet: 'default',
		doggoBath: 'default',
		doggoChewable: 'default'
	});

	const latestLogs = createQuery(notificationQueryOptions);

	const options: ActionCardCompactOptions[] = [
		{
			collectionName: 'towel',
			title: 'Wash Towel',
			route: '/household/towel',
			icon: StreamlineColorHotelLaundryFlat,
			button: {
				text: 'Washed',
				status: buttonStatuses.towel
			}
		},
		{
			collectionName: 'spray',
			title: 'Spray Nose',
			route: '/personal/spray',
			icon: FluentEmojiFlatLotionBottle,
			button: {
				text: 'Sprayed',
				status: buttonStatuses.spray
			}
		},
		{
			collectionName: 'gummy',
			title: 'Elderberry Gummy',
			route: '/personal/gummy',
			icon: FluentEmojiFlatShield,
			button: {
				text: 'Ate',
				status: buttonStatuses.gummy
			}
		},
		{
			collectionName: 'bedsheet',
			title: 'Bedsheet',
			route: '/household/bedsheet',
			icon: FluentEmojiFlatBed,
			button: {
				text: 'Changed',
				status: buttonStatuses.bedsheet
			}
		},
		{
			collectionName: 'doggoBath',
			title: 'Bath',
			route: '/pet/bath',
			icon: FluentEmojiFlatShower,
			button: {
				text: 'Bathed',
				status: buttonStatuses.doggoBath
			}
		},
		{
			collectionName: 'doggoChewable',
			title: 'Chewable',
			route: '/pet/chewable',
			icon: FluentEmojiFlatShield,
			button: {
				text: 'Fed',
				status: buttonStatuses.doggoChewable
			}
		}
	];

	let tasks = $derived.by(() => {
		if (!latestLogs.isSuccess || !latestLogs.data) return { important: [], general: [] };

		const important = ['towel', 'spray', 'gummy'];

		const importantTasks = options
			.filter((task) => important.includes(task.collectionName))
			.map((task) => {
				const data = latestLogs.data.find(
					(log) => log.tracker === trackerNameToId(task.collectionName)
				);

				return { ...task, notification: getNotificationStatus(data) };
			});

		const generalTasks = options
			.filter((task) => !important.includes(task.collectionName))
			.map((task) => {
				const data = latestLogs.data.find(
					(log) => log.tracker === trackerNameToId(task.collectionName)
				);

				return { ...task, notification: getNotificationStatus(data) };
			});

		return { important: importantTasks, general: generalTasks };
	});
</script>

<PageWrapper title="Sundry" back={false} {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">Important</h2>

				{#each tasks.important as task (task.collectionName)}
					<ActionCardCompact
						options={{
							collectionName: task.collectionName,
							title: task.title,
							route: task.route,
							icon: task.icon,
							button: {
								status: buttonStatuses[task.collectionName],
								text: task.button.text
							}
						}}
					></ActionCardCompact>
				{/each}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">General</h2>

				{#if tasks.general && tasks.general.length > 0}
					{@const show = tasks.general.filter((item) => item.notification.show).length}

					{#if show > 0}
						{#each tasks.general as task (task.collectionName)}
							{#if task.notification.show}
								<ActionCardCompact
									options={{
										collectionName: task.collectionName,
										title: task.title,
										route: task.route,
										icon: task.icon,
										button: {
											status: buttonStatuses[task.collectionName],
											text: task.button.text
										}
									}}
								></ActionCardCompact>
							{/if}
						{/each}
					{:else}
						<div class="justify-self-center">
							<enhanced:img src={EmptyCorgi} alt="nothing" />
							<p class="text-center">No other upcoming tasks!</p>
						</div>{/if}
				{/if}
			</section>

			<section class="grid gap-0 py-0">
				<h2 class="text-base-content/70 text-lg font-bold">Quick Links</h2>
				<div class="flex items-center gap-2">
					<a
						href="/household/count"
						class="active:bg-neutral/10 focus-within:bg-neutral/10 focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-24 content-center justify-items-center gap-1 rounded-2xl p-2 text-sm font-semibold"
					>
						<FluentEmojiFlatStopwatch class="size-8" />
						Stopwatch</a
					>
					<a
						href="/profile/vacation"
						class="active:bg-neutral/10 focus-within:bg-neutral/10 focus-within:text-base-content active:text-base-content text-neutral grid aspect-square w-24 content-center justify-items-center gap-1 rounded-2xl p-2 text-sm font-semibold"
					>
						<FluentEmojiFlatAirplane class="size-8" />
						Vacation</a
					>
				</div>
			</section>
		</div>
	</main>
</PageWrapper>
