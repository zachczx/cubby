<script lang="ts">
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import timezone from 'dayjs/plugin/timezone';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { allEntriesQueryOptions } from '$lib/queries';
	import TwoColumnCard from '$lib/ui/TwoColumnCard.svelte';
	import StatusHeroImage from '$lib/ui/StatusHeroImage.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let { options }: { options: TrackerPageOptions } = $props();

	const allEntriesDb = createQuery(allEntriesQueryOptions);

	type TabPages = 'overview' | 'history';

	let currentTab = $state<TabPages>('overview');

	let tracker = $derived(options.tracker);

	let historicalRecords = $derived.by(() => {
		const t = options.tracker;
		if (!t) return null;

		const subscriptionStart = t.startDate ?? '';
		const historicalRecords: string[] = [subscriptionStart];
		const today = dayjs();
		let currentDateTime = dayjs(subscriptionStart);
		while (currentDateTime.add(t.interval, t.intervalUnit).isBefore(today)) {
			currentDateTime = currentDateTime.add(t.interval, t.intervalUnit);

			historicalRecords.unshift(currentDateTime.toISOString());
		}

		return historicalRecords;
	});

	let nextCharge = $derived.by(() => {
		if (!historicalRecords || historicalRecords.length === 0 || !tracker) return null;

		const latest = historicalRecords[0];

		return dayjs(latest).add(tracker.interval, tracker.intervalUnit);
	});

	let entriesByYears = $derived.by(() => {
		if (!historicalRecords) return null;
		const years = new Map<number, string[]>();
		for (const t of historicalRecords) {
			const y = dayjs(t).year();

			if (!years.has(y)) {
				years.set(y, []);
			}
			years.get(y)?.push(t);
		}

		return years;
	});

	let lifetimeSpend = $derived.by(() => {
		if (!historicalRecords || historicalRecords.length === 0 || !tracker) return null;

		const cost = tracker.cost ?? 0;

		return historicalRecords.length * cost;
	});

	let monthlyCost = $derived.by(() => {
		if (!historicalRecords || historicalRecords.length === 0 || !lifetimeSpend) return null;

		const firstEntryTime = historicalRecords[historicalRecords.length - 1];
		const totalDuration = Math.floor(dayjs().diff(dayjs(firstEntryTime), 'month', true));

		return totalDuration > 0 ? Math.floor(lifetimeSpend / totalDuration) : lifetimeSpend;
	});
</script>

<PageWrapper title={options.labels.pageTitle} {pb}>
	<main class="grid w-full max-w-xl content-start justify-items-center gap-4 justify-self-center">
		<div class="grid w-full content-start justify-items-center gap-4">
			<StatusHeroImage notification={{ level: 'ok', show: true }} kind="subscription" />
		</div>

		<div class="grid w-full content-start gap-8 pt-4 pb-8">
			<ul class="tracker-tabs grid w-full grid-cols-2 content-center justify-items-center">
				<li class="w-full" aria-current={currentTab === 'overview' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'overview')}
					>
						Overview</button
					>
				</li>

				<li class="w-full" aria-current={currentTab === 'history' ? 'page' : undefined}>
					<button
						class="w-full cursor-pointer text-center"
						onclick={() => (currentTab = 'history')}
					>
						History</button
					>
				</li>
			</ul>

			{#if historicalRecords && historicalRecords.length === 0}
				<div class="mx-4 mt-4 text-center">
					<p class="font-bold">Ready to track?</p>
					<p>Start by adding a log.</p>
				</div>
			{:else}
				<div class={['grid w-full gap-8 px-4', currentTab === 'overview' ? undefined : 'hidden']}>
					<div
						class="border-base-content/5 bg-primary/10 grid justify-items-center rounded-xl border p-4 shadow"
					>
						<h2 class="text-md text-center">Next Charge</h2>
						<div class="grid min-h-20 content-center justify-items-center">
							{#if nextCharge}
								{@const semantic = nextCharge.fromNow()}
								{@const formatted = nextCharge.format('D MMM YYYY')}
								<p class="text-primary text-2xl font-bold">
									{semantic}
								</p>
								<p>{formatted}</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-2xl font-bold">Nil</div>
							{/if}
						</div>
					</div>

					<TwoColumnCard leftTitle="Frequency" rightTitle="Last">
						{#snippet left()}
							{#if tracker}
								{@const plural = tracker.interval > 1 ? true : false}
								<p>
									{tracker.interval}&nbsp;{plural
										? tracker.intervalUnit + 's'
										: tracker.intervalUnit}
								</p>
								<p class="text-base-content/70 text-base font-normal">once</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-xl font-bold">Not set yet</div>
							{/if}
						{/snippet}

						{#snippet right()}
							{#if historicalRecords && historicalRecords.length > 0}
								{#if historicalRecords.length > 0}
									{@const formatted = dayjs(historicalRecords[0]).fromNow(true)}
									<p>
										{#if formatted === 'a few seconds'}
											seconds
										{:else}
											{formatted}
										{/if}
									</p>
									<p class="text-base-content/70 text-base font-normal">ago</p>
								{:else}
									<div class="flex min-h-20 items-center gap-4 text-xl font-bold">
										Never recorded
									</div>
								{/if}
							{/if}
							{#if allEntriesDb.isPending}
								<div class="custom-loader"></div>
							{/if}
						{/snippet}
					</TwoColumnCard>

					<TwoColumnCard
						leftTitle="Monthly Avg"
						leftIcon="material-symbols:calculate"
						rightTitle="Lifetime Spend"
						rightIcon="material-symbols:donut-small"
					>
						{#snippet left()}
							{#if historicalRecords && historicalRecords.length > 0}
								<p>
									${monthlyCost}/mth
								</p>
							{:else}
								<div class="flex min-h-20 items-center gap-4 text-xl font-bold">Not set yet</div>
							{/if}
						{/snippet}

						{#snippet right()}
							{#if historicalRecords && historicalRecords.length > 0}
								{#if historicalRecords.length > 0}
									<p>
										${lifetimeSpend}
									</p>
								{:else}
									<div class="flex min-h-20 items-center gap-4 text-xl font-bold">
										Never recorded
									</div>
								{/if}
							{/if}
							{#if allEntriesDb.isPending}
								<div class="custom-loader"></div>
							{/if}
						{/snippet}
					</TwoColumnCard>
				</div>

				<div
					class={[
						'grid w-full grid-cols-[minmax(0,1fr)] gap-4 overflow-hidden px-4',
						currentTab === 'history' ? undefined : 'hidden'
					]}
				>
					{#if historicalRecords && historicalRecords.length > 0}
						{#each entriesByYears as [key, year]}
							<div class="grid gap-2">
								<h3 class="text-base-content/50 text-sm font-bold">{key}</h3>

								{#each year as entry}
									{@const formatted = dayjs(entry).format('D MMM')}
									<div
										class="border-base-300 grid min-h-18 grid-cols-3 items-baseline gap-4 rounded-2xl border bg-white/70 px-2 py-2"
									>
										<div class="flex p-2 text-lg font-bold">
											{formatted}
										</div>
										<div class="flex justify-center">
											{#if entry === tracker?.startDate}
												<span class="text-primary">Sign-Up</span>
											{:else}
												<span class="text-neutral">Renewal</span>
											{/if}
										</div>
										<div class="text-base-content/60 flex justify-end p-2 font-semibold">
											{#if tracker}
												${tracker.cost}
											{/if}
										</div>
									</div>
								{/each}
							</div>
						{/each}
					{/if}
				</div>
			{/if}
		</div>
	</main>
</PageWrapper>
