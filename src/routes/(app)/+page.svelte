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
		allTrackersQueryOptions,
		createUserQueryOptions,
		logsQueryOptions,
		notificationQueryOptions,
		trackerNameToId
	} from '$lib/queries';
	import { getTrackerStatus } from '$lib/notification';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import EmptyCorgi from '$lib/assets/empty.webp?w=200&enhanced';
	import FluentEmojiFlatStopwatch from '$lib/assets/expressive-icons/FluentEmojiFlatStopwatch.svelte';
	import FluentEmojiFlatAirplane from '$lib/assets/expressive-icons/FluentEmojiFlatAirplane.svelte';
	import {
		bedsheet,
		doggoBath,
		doggoChewable,
		gummy,
		spray,
		towel
	} from '$lib/ui/ActionCardDefaults';

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
	const trackers = createQuery(allTrackersQueryOptions);

	const options: ActionCardOptions[] = [
		{ ...towel, size: 'compact', button: { ...towel.button, status: buttonStatuses.towel } },
		{ ...spray, size: 'compact', button: { ...spray.button, status: buttonStatuses.spray } },
		{ ...gummy, size: 'compact', button: { ...gummy.button, status: buttonStatuses.gummy } },
		{
			...bedsheet,
			size: 'compact',
			button: { ...bedsheet.button, status: buttonStatuses.bedsheet }
		},
		{
			...doggoBath,
			size: 'compact',
			button: { ...doggoBath.button, status: buttonStatuses.doggoBath }
		},
		{
			...doggoChewable,
			size: 'compact',
			button: { ...doggoChewable.button, status: buttonStatuses.doggoChewable }
		}
	];

	let tasks = $derived.by(() => {
		if (!latestLogs.isSuccess || !latestLogs.data) return { important: [], general: [] };

		const important = ['towel', 'spray', 'gummy'];

		const importantTasks = options
			.filter((task) => important.includes(task.collectionName))
			.map((task) => {
				const data = latestLogs.data.find(
					(log) => log.tracker === trackerNameToId(task.collectionName, trackers.data)
				);

				return { ...task, notification: getTrackerStatus(data) };
			});

		const generalTasks = options
			.filter((task) => !important.includes(task.collectionName))
			.map((task) => {
				const data = latestLogs.data.find(
					(log) => log.tracker === trackerNameToId(task.collectionName, trackers.data)
				);

				return { ...task, notification: getTrackerStatus(data) };
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
					<ActionCard
						options={{
							collectionName: task.collectionName,
							size: task.size ? task.size : undefined,
							title: task.title,
							route: task.route,
							icon: task.icon,
							button: {
								status: buttonStatuses[task.collectionName],
								text: task.button.text
							}
						}}
					></ActionCard>
				{/each}
			</section>

			<section class="grid gap-4 py-2">
				<h2 class="text-base-content/70 text-lg font-bold">General</h2>

				{#if tasks.general && tasks.general.length > 0}
					{@const show = tasks.general.filter((item) => item.notification.show).length}

					{#if show > 0}
						{#each tasks.general as task (task.collectionName)}
							{#if task.notification.show}
								<ActionCard
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
								></ActionCard>
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
