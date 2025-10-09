<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { onMount } from 'svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let towelDB: TowelDB[] | undefined = $state([]);
	let sprayDB: SprayDB[] | undefined = $state([]);

	onMount(async () => {
		if (pb.authStore.isValid) {
			towelDB = await pb.collection('towel').getFullList({
				sort: '-time'
			});

			sprayDB = await pb.collection('spray').getFullList({
				sort: '-created'
			});
		}
	});

	let towelLast: string | undefined = $derived.by(() => {
		if (towelDB && towelDB.length > 0) {
			return dayjs(towelDB[0].time).fromNow();
		}
		return undefined;
	});
	let spinner = $state(false);

	let sprayLast: string | undefined = $derived.by(() => {
		if (sprayDB && sprayDB.length > 0) {
			return dayjs(sprayDB[0].time).fromNow();
		}
		return undefined;
	});

	async function addTowelHandler() {
		const result = await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			spinner = false;
		}

		towelDB = await pb.collection('towel').getFullList({
			sort: '-time'
		});
	}

	async function addSprayHandler() {
		spinner = true;
		const result = await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
			spinner = false;
		}

		sprayDB = await pb.collection('spray').getFullList({
			sort: '-time'
		});
	}
</script>

<svelte:head>
	<title>Towsey</title>
</svelte:head>
<PageWrapper title="Towsey" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<header class="grid justify-self-start text-start">
				{#if pb.authStore && pb.authStore.record}
					<span class="opacity-[0.9]">Welcome back,</span><span class="font-semibold"
						>{pb.authStore.record.email}!</span
					>
				{/if}
			</header>

			<section class="grid gap-4">
				<div class="bg-primary/10 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4">
					<h3 class="text-sm lg:text-base">Towel Washed</h3>
					<div class="text-2xl font-bold">
						{#key towelLast}
							{#if towelLast}
								{towelLast}
							{:else}
								<div class="custom-loader"></div>
							{/if}
						{/key}
					</div>
					<button
						class="btn btn-primary lg:btn-lg flex w-full items-center gap-2 rounded-2xl"
						onclick={addTowelHandler}>Just Washed My Towel!</button
					>
				</div>

				<div class="bg-secondary/10 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4">
					<h3 class="text-sm lg:text-base">Nose Sprayed</h3>
					<div class="text-2xl font-bold">
						{#key sprayLast}
							{#if sprayLast}
								{sprayLast}
							{:else}
								<div class="custom-loader"></div>
							{/if}
						{/key}
					</div>
					<button
						class="btn btn-primary lg:btn-lg flex w-full items-center gap-2 rounded-2xl"
						onclick={addSprayHandler}>Just Sprayed Nose!</button
					>
				</div>
			</section>

			<section class="border-base-300 grid min-h-24 gap-4 rounded-3xl border-0 py-4">
				<h3 class="text-sm font-bold lg:text-base">Jump To</h3>
				<ul class="justify-items-around grid grid-cols-3 gap-2">
					<li>
						<a
							href="/towelie"
							class="bg-base-200/30 focus:bg-primary/10 active:bg-primary/10 border-primary/20 hover:bg-primary/10 text-primary grid aspect-square w-full content-center justify-items-center gap-1 rounded-3xl border p-2 font-semibold lg:text-lg"
							><svg
								xmlns="http://www.w3.org/2000/svg"
								width="1em"
								height="1em"
								class="lucide-lab:towel-rack size-[1.3em] opacity-50 lg:size-[1.7em]"
								viewBox="0 0 24 24"
								><path
									fill="none"
									stroke="currentColor"
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M8 6H2m4-4h12a2 2 0 0 1 2 2v18H8V4a2 2 0 0 0-4 0v15h4M22 6h-2M8 18h12"
								/></svg
							>Towelie</a
						>
					</li>
					<li>
						<a
							href="/nosey"
							class="bg-base-200/30 focus:bg-primary/10 active:bg-primary/10 border-primary/20 hover:bg-primary/10 text-primary grid aspect-square w-full content-center justify-items-center gap-1 rounded-3xl border p-2 font-semibold lg:text-lg"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="1em"
								height="1em"
								class="material-symbols:health-and-safety-outline size-[1.3em] opacity-50 lg:size-[1.7em]"
								viewBox="0 0 24 24"
								><path
									fill="currentColor"
									d="M10.5 15.5h3V13H16v-3h-2.5V7.5h-3V10H8v3h2.5zM12 22q-3.475-.875-5.738-3.988T4 11.1V5l8-3l8 3v6.1q0 3.8-2.262 6.913T12 22m0-2.1q2.6-.825 4.3-3.3t1.7-5.5V6.375l-6-2.25l-6 2.25V11.1q0 3.025 1.7 5.5t4.3 3.3m0-7.9"
								/></svg
							>Nosey</a
						>
					</li>
					<li>
						<a
							href="/count"
							class="bg-base-200/30 focus:bg-primary/10 active:bg-primary/10 border-primary/20 hover:bg-primary/10 text-primary grid aspect-square w-full content-center justify-items-center gap-1 rounded-3xl border p-2 font-semibold lg:text-lg"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="32"
								height="32"
								class="material-symbols:timer-outline size-[1.3em] opacity-50 lg:size-[1.7em]"
								viewBox="0 0 24 24"
								><path
									fill="currentColor"
									d="M9 3V1h6v2zm2 11h2V8h-2zm1 8q-1.85 0-3.488-.712T5.65 19.35t-1.937-2.863T3 13t.713-3.488T5.65 6.65t2.863-1.937T12 4q1.55 0 2.975.5t2.675 1.45l1.4-1.4l1.4 1.4l-1.4 1.4Q20 8.6 20.5 10.025T21 13q0 1.85-.713 3.488T18.35 19.35t-2.863 1.938T12 22m0-2q2.9 0 4.95-2.05T19 13t-2.05-4.95T12 6T7.05 8.05T5 13t2.05 4.95T12 20m0-7"
								/></svg
							>Count</a
						>
					</li>
				</ul>
			</section>
		</div>
	</main>
</PageWrapper>
