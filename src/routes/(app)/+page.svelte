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

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelButtonStatus: ButtonStatus = $state('default');
	let sprayButtonStatus: ButtonStatus = $state('default');
	let gummyButtonStatus: ButtonStatus = $state('default');

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

	async function addTowelHandler() {
		towelButtonStatus = 'loading';

		const result = await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			towelButtonStatus = 'success';

			setTimeout(() => {
				towelButtonStatus = 'default';
			}, 3000);
		}

		await tanstackClient.refetchQueries(createTowelRefetchOptions());
	}

	async function addSprayHandler() {
		sprayButtonStatus = 'loading';

		const result = await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: sprayDaysToNext
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			sprayButtonStatus = 'success';

			setTimeout(() => {
				sprayButtonStatus = 'default';
			}, 3000);
		}

		await tanstackClient.refetchQueries(createSprayRefetchOptions());
	}

	async function addGummyHandler() {
		gummyButtonStatus = 'loading';

		const result = await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: gummyDaysToNext
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			gummyButtonStatus = 'success';

			setTimeout(() => {
				gummyButtonStatus = 'default';
			}, 3000);
		}

		await tanstackClient.refetchQueries(createGummyRefetchOptions());
	}

	let sprayNotification = $derived.by(() => getNotificationStatus(sprays));
	let towelNotification = $derived.by(() => getNotificationStatus(towels));
	let gummyNotification = $derived.by(() => getNotificationStatus(gummies));
</script>

<svelte:head>
	<title>Towsey</title>
</svelte:head>

<PageWrapper title="Towsey" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<section
					class={[
						'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
						towelNotification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-primary/10'
					]}
				>
					<a href="/towelie" class="flex items-center">
						<div class="flex grow items-center gap-4">
							<PhTowelFill class="size-12 opacity-75" />
							<!-- <h3 class="text-sm lg:text-base">Towel Washed</h3> -->
							<div>
								<p class="text-xl font-bold">Wash Towel</p>
								{#if towels.isPending && !towels.data}
									<div class="custom-loader"></div>
								{/if}
								{#if towels.error}
									Error: {towels.error.message}
								{/if}

								{#if towels.isSuccess}
									{#if towelNotification.show}
										<span class="text-error font-medium tracking-tight">
											{#if towelNotification.level === 'overdue'}
												Overdue
											{:else if towelNotification.level === 'due'}
												Due
											{/if}
										</span>&nbsp;&nbsp;•&nbsp;&nbsp;
									{/if}<span>{towelLast}</span>
								{/if}
							</div>
						</div>
						<div class="flex h-full items-center">
							<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
								><MaterialSymbolsChevronRight class="size-6" /></button
							>
						</div>
					</a>
					<button
						class={[
							'btn btn-lg flex w-full items-center gap-2 rounded-full',
							towelButtonStatus === 'default' && 'btn-primary',
							towelButtonStatus === 'loading' && 'btn-primary',
							towelButtonStatus === 'success' && 'btn-success'
						]}
						onclick={addTowelHandler}
					>
						{#if towelButtonStatus === 'success'}
							<MaterialSymbolsCheck class="size-6" />Added!
						{:else if towelButtonStatus === 'loading'}
							<span class="loading loading-spinner loading-md"></span>
						{:else}
							Just Washed
						{/if}
					</button>
				</section>

				<section
					class={[
						'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
						sprayNotification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-primary/10'
					]}
				>
					<a href="/nosey" class="flex items-center">
						<div class="flex grow items-center gap-4">
							<IconParkSolidBottleOne class="size-12 opacity-75" />
							<div>
								<p class="text-xl font-bold">Spray Nose</p>
								{#if sprays.isPending && !sprays.data}
									<div class="custom-loader"></div>
								{/if}
								{#if sprays.error}
									An error has occurred:
									{sprays.error.message}
								{/if}
								{#if sprays.isSuccess}
									{#if sprayNotification.show}
										<span class="text-error font-medium tracking-tight">
											{#if sprayNotification.level === 'overdue'}
												Overdue
											{:else if sprayNotification.level === 'due'}
												Due
											{/if}
										</span>&nbsp;&nbsp;•&nbsp;&nbsp;
									{/if}<span>{sprayLast}</span>
								{/if}
							</div>
						</div>
						<div class="flex h-full items-center">
							<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
								><MaterialSymbolsChevronRight class="size-6" /></button
							>
						</div>
					</a>
					<button
						class={[
							'btn btn-lg flex w-full items-center gap-2 rounded-full',
							sprayButtonStatus === 'default' && 'btn-primary',
							sprayButtonStatus === 'loading' && 'btn-primary',
							sprayButtonStatus === 'success' && 'btn-success'
						]}
						onclick={addSprayHandler}
					>
						{#if sprayButtonStatus === 'success'}
							<MaterialSymbolsCheck class="size-6" />Added!
						{:else if sprayButtonStatus === 'loading'}
							<span class="loading loading-spinner loading-md"></span>
						{:else}
							Just Sprayed
						{/if}
					</button>
				</section>

				<section
					class={[
						'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
						sprayNotification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-primary/10'
					]}
				>
					<a href="/gummy" class="flex items-center">
						<div class="flex grow items-center gap-4">
							<MaterialSymbolsHealthAndSafety class="size-12 opacity-75" />
							<div>
								<p class="text-xl font-bold">Elderberry Gummy</p>
								{#if gummies.isPending && !gummies.data}
									<div class="custom-loader"></div>
								{/if}
								{#if gummies.error}
									An error has occurred:
									{gummies.error.message}
								{/if}
								{#if gummies.isSuccess}
									{#if gummyNotification.show}
										<span class="text-error font-medium tracking-tight">
											{#if gummyNotification.level === 'overdue'}
												Overdue
											{:else if gummyNotification.level === 'due'}
												Due
											{/if}
										</span>&nbsp;&nbsp;•&nbsp;&nbsp;
									{/if}<span>{gummyLast}</span>
								{/if}
							</div>
						</div>
						<div class="flex h-full items-center">
							<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
								><MaterialSymbolsChevronRight class="size-6" /></button
							>
						</div>
					</a>
					<button
						class={[
							'btn btn-lg flex w-full items-center gap-2 rounded-full',
							gummyButtonStatus === 'default' && 'btn-primary',
							gummyButtonStatus === 'loading' && 'btn-primary',
							gummyButtonStatus === 'success' && 'btn-success'
						]}
						onclick={addSprayHandler}
					>
						{#if gummyButtonStatus === 'success'}
							<MaterialSymbolsCheck class="size-6" />Added!
						{:else if gummyButtonStatus === 'loading'}
							<span class="loading loading-spinner loading-md"></span>
						{:else}
							Just Ate Gummy
						{/if}
					</button>
				</section>

				<div class="bg-primary/10 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4">
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
