<script lang="ts">
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { onMount } from 'svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { defaultSprayInterval, dirtyTowelDays } from '$lib/config';
	import MaterialSymbolsAdd from '$lib/assets/svg/MaterialSymbolsAdd.svelte';
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	const towels = createQuery<TowelDB[]>(() => ({
		queryKey: ['towels'],
		queryFn: async () => await pb.collection('towel').getFullList({ sort: '-time' })
	}));

	const sprays = createQuery<SprayDB[]>(() => ({
		queryKey: ['sprays'],
		queryFn: async () => await pb.collection('spray').getFullList({ sort: '-time' })
	}));

	const userPref = createQuery<UserDB>(() => ({
		queryKey: ['user', pb.authStore?.record?.id],
		queryFn: async () => await pb.collection('user').getOne(String(pb.authStore?.record?.id))
	}));

	const tanstackClient = useQueryClient();

	let towelLast: string = $derived.by(() => {
		if (towels.isSuccess) {
			return dayjs(towels.data[0].time).fromNow();
		}
		return '';
	});

	let daysToNext = $derived.by(() => {
		if (userPref.isPending) {
			return undefined;
		}

		return userPref.data?.defaultSprayInterval;
	});

	let sprayLast: string = $derived.by(() => {
		if (sprays.isSuccess) return dayjs(sprays.data[0].time).fromNow();

		return '';
	});

	async function addTowelHandler() {
		const result = await pb.collection('towel').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore')
		});

		if (result.id) {
			addToast('success', 'Added successfully!');
		}

		await tanstackClient.refetchQueries({
			queryKey: ['towels'],
			type: 'active',
			exact: true
		});
	}

	async function addSprayHandler() {
		const result = await pb.collection('spray').create({
			user: pb.authStore.record?.id,
			time: dayjs.tz(new Date(), 'Asia/Singapore'),
			daysToNext: daysToNext
		});

		console.log(result);

		if (result.id) {
			addToast('success', 'Added successfully!');
		}

		await tanstackClient.refetchQueries({
			queryKey: ['sprays'],
			type: 'active',
			exact: true
		});
	}

	let sprayNotification = $derived.by(() => {
		if (!sprays?.isSuccess) return false;

		const lastSpray = sprays.data?.[0] ?? null;

		if (!lastSpray) {
			return false;
		}

		const now = dayjs();
		const leadTimeHours = 12;
		const intervalHours = lastSpray.daysToNext * 24;

		const hoursSinceLastSpray = now.diff(dayjs(lastSpray.time), 'hour', true);
		if (hoursSinceLastSpray >= intervalHours - leadTimeHours) {
			return true;
		}

		return false;
	});

	let towelNotification = $derived.by(() => {
		if (!towels?.isSuccess) return false;

		const lastWash = towels.data?.[0] ?? null;

		if (!lastWash) {
			return false;
		}

		const now = dayjs();
		const leadTimeHours = 12;
		const intervalHours = dirtyTowelDays * 24;

		const hoursSinceLastWash = now.diff(dayjs(lastWash.time), 'hour', true);

		if (hoursSinceLastWash >= intervalHours - leadTimeHours) {
			return true;
		}

		return false;
	});
</script>

<svelte:head>
	<title>Towsey</title>
</svelte:head>
<PageWrapper title="Towsey" {pb} {towelNotification} {sprayNotification}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-4 py-4">
				<section
					class={[
						'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
						towelNotification ? 'bg-error/30 outline-error/70 outline' : 'bg-primary/10'
					]}
				>
					<a href="/towelie" class="flex items-center">
						<div class="flex grow items-center gap-8">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="32"
								height="32"
								class="streamline-kameleon-color:towel size-24"
								viewBox="0 0 48 48"
								><!-- Icon from Kameleon color icons by Streamline - https://creativecommons.org/licenses/by/4.0/ --><g
									fill="none"
									><path
										fill="#fabc3d"
										d="M24 47.998c13.255 0 24-10.745 24-24C48 10.746 37.255 0 24 0S0 10.745 0 23.999s10.745 23.999 24 23.999"
									/><path
										fill="#3e3e3f"
										d="M37.087 22.93L24.87 15.66l-.01-1.866l-.576-.208c-.074-.027-.723-.28-.723-.921c0-.781 1.552-1.492 2.274-.046a.873.873 0 0 0 1.562-.78c-.686-1.373-2.025-2.048-3.41-1.72c-1.278.3-2.17 1.348-2.17 2.546c0 1.148.673 1.888 1.308 2.287v.707l-12.08 7.171l-.136.1a2.63 2.63 0 0 0-.703 2.886a2.63 2.63 0 0 0 2.446 1.685h22.691a2.63 2.63 0 0 0 2.446-1.685a2.63 2.63 0 0 0-.702-2.887m-.928 2.263a.88.88 0 0 1-.816.562h-22.69a.88.88 0 0 1-.816-.562a.88.88 0 0 1 .183-.913l11.978-7.11l11.977 7.11c.23.24.303.6.184.913"
									/><path
										fill="#f0f1f1"
										d="M33.602 24.883a.873.873 0 0 0-.873-.873H15.273a.873.873 0 0 0-.873.872v11.345h19.202z"
									/><path fill="#e2e4e5" d="M33.6 36.227H14.4v1.745h19.2z" /><path
										fill="#fff"
										d="M16.582 36.227h-.873v1.746h.872zm2.618 0h-.873v1.746h.873zm2.618 0h-.873v1.746h.873zm2.618 0h-.872v1.746h.872zm2.618 0h-.872v1.746h.872zm2.618 0H28.8v1.746h.872zm2.618 0h-.872v1.746h.873zM14.4 24.883v11.344h1.727L28.345 24.01H15.273a.873.873 0 0 0-.873.873"
									/></g
								></svg
							>
							<!-- <h3 class="text-sm lg:text-base">Towel Washed</h3> -->
							<div class="text-lg">
								{#if towelNotification}
									<span class="btn btn-error btn-xs mb-2 rounded-full">Overdue</span>
								{/if}
								<p class="font-semibold uppercase">Wash</p>
								{#if towels.isPending}
									<div class="custom-loader"></div>
								{/if}
								{#if towels.error}
									Error: {towels.error.message}
								{/if}
								{#if towels.isSuccess}
									{towelLast}
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
						class="btn btn-primary btn-lg flex w-full items-center gap-2 rounded-2xl"
						onclick={addTowelHandler}>Just Washed My Towel!</button
					>
				</section>

				<section
					class={[
						'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
						sprayNotification ? 'bg-error/30 outline-error/70 outline' : 'bg-primary/10'
					]}
				>
					<a href="/nosey" class="flex items-center">
						<div class="flex grow items-center gap-8">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="32"
								height="32"
								class="streamline-kameleon-color:medicine size-24"
								viewBox="0 0 48 48"
								><!-- Icon from Kameleon color icons by Streamline - https://creativecommons.org/licenses/by/4.0/ --><g
									fill="none"
									><path
										fill="#fabc3d"
										d="M24 47.998c13.255 0 24-10.745 24-24C48 10.746 37.255 0 24 0S0 10.745 0 23.999s10.745 23.999 24 23.999"
									/><path
										fill="#059bbf"
										d="M32.728 18.763a3.49 3.49 0 0 0-3.491-3.491H18.764a3.49 3.49 0 0 0-3.491 3.49v15.71a3.49 3.49 0 0 0 3.49 3.49h10.474a3.49 3.49 0 0 0 3.49-3.49z"
									/><path fill="#f0f1f1" d="M32.728 20.508H15.273v12.217h17.455z" /><path
										fill="#21b2d1"
										d="M24 15.272v5.236h-8.727v-1.745a3.49 3.49 0 0 1 3.49-3.491z"
									/><path fill="#fff" d="M24 20.508h-8.727v12.217H24z" /><path
										fill="#f0f1f1"
										d="M32.728 14.4a.873.873 0 0 1-.873.872h-15.71a.873.873 0 0 1-.872-.873v-3.49c0-.483.39-.873.873-.873h15.709c.482 0 .873.39.873.872z"
									/><path
										fill="#fff"
										d="M24 10.036h-7.854a.873.873 0 0 0-.873.872V14.4c0 .482.39.873.873.873H24z"
									/><path
										fill="#fa5655"
										d="M28.8 24.617H26v-2.8h-4v2.8h-2.8v4H22v2.8h4v-2.8h2.8z"
									/><path
										fill="#f0f1f1"
										d="M19.2 10.036h-1.309v5.236h1.31zm3.49 0h-1.308v5.236h1.309z"
									/><path
										fill="#e2e4e5"
										d="M26.618 10.036H25.31v5.236h1.31zm3.491 0h-1.31v5.236h1.31z"
									/><path
										fill="#21b2d1"
										d="M24 32.726v5.236h-5.236a3.49 3.49 0 0 1-3.491-3.491v-1.745z"
									/></g
								></svg
							>
							<!-- <h3 class="text-sm lg:text-base">Nose Sprayed</h3> -->
							<div class="text-lg">
								{#if sprayNotification}
									<span class="btn btn-error btn-xs mb-2 rounded-full">Overdue</span>
								{/if}
								<p class="font-semibold uppercase">Spray</p>
								{#if sprays.isPending}
									<div class="custom-loader"></div>
								{/if}
								{#if sprays.error}
									An error has occurred:
									{sprays.error.message}
								{/if}
								{#if sprays.isSuccess}
									{sprayLast}
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
						class="btn btn-primary btn-lg flex w-full items-center gap-2 rounded-2xl"
						onclick={addSprayHandler}>Just Sprayed Nose!</button
					>
				</section>

				<a
					href="/count"
					class="bg-primary/10 border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4"
				>
					<div class="flex items-center">
						<div class="flex grow items-center gap-8">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								xml:space="preserve"
								width="655.359"
								height="655.359"
								class="size-24"
								style="shape-rendering:geometricPrecision;text-rendering:geometricPrecision;image-rendering:optimizeQuality;fill-rule:evenodd;clip-rule:evenodd"
								viewBox="0 0 6.827 6.827"
							>
								<!-- customised from https://www.reshot.com/free-svg-icons/item/timer-off-KHCXLP4E3N/ -->
								<defs
									><clipPath id="id0"
										><path
											d="M3.413 0a3.413 3.413 0 1 1 0 6.827 3.413 3.413 0 0 1 0-6.827z"
										/></clipPath
									><style>
										.fil1 {
											fill: #fabc3d;
										}
										.fil3,
										.fil5 {
											fill: #424242;
											fill-rule: nonzero;
										}
										.fil3 {
											fill: #8d6e63;
										}
									</style></defs
								><g id="Layer_x0020_1"
									><path
										d="M3.413 0a3.413 3.413 0 1 1 0 6.827 3.413 3.413 0 0 1 0-6.827z"
										style="fill:#fabc3d"
									/><g style="clip-path:url(#id0)"
										><g id="_351278000"
											><path
												id="_3555835041"
												class="fil1"
												d="m3.936 1.89 3.787 3.786-.006.005-.005.005-.006.005-.006.004-.007.004-.006.004-.007.003-.007.004-.007.002-.007.003-.007.002-.008.002-.008.001-.007.001h-.016L3.826 1.937h.008l.008-.001.008-.001.008-.002.007-.001.008-.002.007-.003.007-.003.007-.003.007-.003.006-.004.006-.004.006-.004.006-.005.006-.005z"
											/><path
												id="_355583528"
												class="fil1"
												d="m3.18 1.936 3.786 3.786h-.18L3 1.936z"
											/><path
												id="_355584752"
												class="fil1"
												d="m4.621 4.727 3.786 3.786.42.42-3.785-3.786z"
											/><path
												id="_355584968"
												class="fil1"
												d="m5.042 5.147 3.786 3.787-.11.11-3.786-3.787z"
											/><path
												id="_355584608"
												class="fil1"
												d="M4.514 4.84 8.3 8.625l-.027.025-.028.025-.028.024-.029.023-.03.023-.03.022-.03.021-.03.02-.032.02-.032.02-.032.018-.033.018-.034.017-.034.016-.034.016-.034.014-.036.014-.035.013-.036.012L7.66 9l-.037.01-.037.01-.037.009-.038.008-.038.007-.038.006-.038.005-.04.005-.038.003-.04.002-.04.002H7.2L3.414 5.28h.04l.04-.002.04-.002.038-.004.04-.004.038-.005.038-.006.038-.007.038-.008.037-.01.037-.009.037-.01.036-.012.036-.012.035-.013.035-.014.035-.015.034-.015.034-.016.033-.017.033-.018.033-.018.032-.02.031-.02.031-.02.03-.022.03-.022.03-.022.029-.024.028-.024.028-.024z"
											/><path
												id="_355584944"
												class="fil1"
												d="M3.413 5.28 7.2 9.066h-.042l-.04-.002-.041-.002-.04-.004-.04-.005-.04-.005-.04-.007-.039-.007-.038-.009-.039-.009-.038-.01-.038-.011-.037-.013-.037-.012-.036-.014-.036-.015-.036-.016-.035-.016-.034-.017-.034-.018-.034-.019-.033-.02-.033-.02-.032-.021-.031-.022-.031-.023-.03-.023-.03-.024-.03-.025L2.37 4.891l.03.025.029.024.03.023.031.023.032.022.032.02.033.021.033.02.033.019.034.018.035.017.035.016.036.016.036.014.036.014.037.013.037.012.038.012.038.01.038.01.039.008.04.007.039.007.04.005.04.005.04.003.04.003.041.001z"
											/><path
												id="_355584224"
												class="fil1"
												d="m1.895 2 3.786 3.787.477.477-3.786-3.786z"
											/><path
												id="_355584536"
												class="fil1"
												d="m3 1.936 3.786 3.786h-.008l-.007-.001-.008-.001-.008-.001-.007-.002-.008-.002-.007-.003-.007-.002-.007-.004-.007-.003-.006-.004-.006-.004-.007-.004-.005-.005-.006-.005-.005-.005L2.89 1.89l.006.005.005.005.006.005.006.004.007.004.006.004.007.003.007.003.007.003.007.003.007.002.008.001.008.002h.007l.008.001z"
											/><path
												id="_355584776"
												class="fil1"
												d="m2.89 1.89 3.787 3.786-.006-.005-3.786-3.786z"
											/><path
												id="_355584848"
												class="fil1"
												d="m5.008 3.686 3.786 3.786v.037l-.002.037-.002.037-.003.036-.003.037-.005.036-.005.036-.006.035-.007.035-.008.035-.009.035-.009.034-.01.035-.01.033-.012.034-.012.033-.013.033-.013.032-.015.032-.014.032-.016.03-.016.032-.017.03-.018.03-.018.03-.019.03-.019.028-.02.028-.02.028-.022.028-.021.026-.023.027-3.786-3.786.022-.027.022-.027.021-.027.02-.028.02-.028.02-.03.019-.028.018-.03.018-.03.016-.03.017-.031.015-.032.015-.031.014-.032.014-.033.013-.032.012-.034.011-.033.011-.034.01-.034.01-.034.008-.035.007-.035.007-.035.006-.036.006-.036.004-.036.004-.036.003-.036.002-.037.001-.037z"
											/><path
												id="_355584296"
												class="fil1"
												d="m3.942 1.598 3.786 3.786.005.006.004.006.005.006.004.006.003.006.004.007.003.007.003.007.002.007.002.008.002.007.001.008.001.008.001.008v.007L3.982 1.702v-.016L3.98 1.68 3.98 1.67l-.002-.008-.002-.007-.002-.007-.003-.007-.003-.007-.004-.007-.004-.006-.004-.007-.004-.006-.005-.006z"
											/><path
												id="_355585496"
												class="fil1"
												d="m3.982 1.702 3.786 3.786v.078L3.982 1.78z"
											/><path
												id="_355585016"
												class="fil1"
												d="m3.982 1.78 3.786 3.786v.016l-.002.008v.008l-.003.007-.002.007-.002.008-.003.007-.003.007-.004.006-.003.007-.004.006-.005.006-.004.006-.005.006-.005.005L3.936 1.89l.006-.005.004-.006.005-.006.004-.006.004-.006.004-.007.004-.007.003-.007.003-.007.002-.007.002-.007.002-.008.001-.007.001-.008v-.008z"
											/><path
												id="_355585256"
												class="fil1"
												d="m3.826 1.936 3.787 3.786h-.18L3.647 1.936z"
											/><path
												id="_355585376"
												class="fil1"
												d="m3.647 1.936 3.786 3.786v.172L3.647 2.108z"
											/><path
												id="_355585424"
												class="fil1"
												d="M4.594 2.613 8.38 6.4l.025.028.025.03.024.03.023.03.023.03.022.032.02.032.021.033.02.033.019.034.018.034.017.035.016.035.016.035.014.036.014.037.013.036.012.038.012.037.01.038.01.039.008.038.007.04.007.039.005.04.005.04.003.04.003.04.002.042v.04L5.008 3.687v-.042l-.002-.04-.003-.041L5 3.523l-.005-.04-.006-.04-.006-.04-.008-.039-.008-.039-.01-.038-.01-.038-.01-.038-.013-.037-.013-.037-.014-.036-.015-.036-.015-.036-.016-.035-.018-.034-.018-.034-.018-.034-.02-.033-.02-.033-.022-.032-.021-.031-.023-.031-.023-.03-.024-.03-.025-.03z"
											/><path
												id="_351278336"
												class="fil1"
												d="M3.647 2.108a1.59 1.59 0 0 1 1.361 1.577c0 .398-.146.762-.387 1.042l.42.42-.11.11-.417-.418a1.59 1.59 0 0 1-1.1.441 1.59 1.59 0 0 1-1.128-.467 1.59 1.59 0 0 1-.467-1.127c0-.427.168-.815.44-1.1l-.474-.475.11-.11.477.477c.224-.193.502-.325.808-.37v-.172H3a.155.155 0 0 1-.155-.156v-.078A.155.155 0 0 1 3 1.547h.826a.155.155 0 0 1 .156.155v.078a.155.155 0 0 1-.156.156h-.18v.172z"
											/></g
										></g
									><path
										d="M3.413 0a3.413 3.413 0 1 1 0 6.827 3.413 3.413 0 0 1 0-6.827z"
										style="fill:none"
									/><path
										class="fil3"
										d="M3.336 1.858v.31H3.18v-.31zM3.647 1.858v.31H3.49v-.31z"
									/><path
										d="M3 1.547h.826a.155.155 0 0 1 .156.155v.078a.155.155 0 0 1-.156.156H3a.155.155 0 0 1-.155-.156v-.078A.155.155 0 0 1 3 1.547z"
										style="fill:#4e342e;fill-rule:nonzero"
									/><path class="fil5" d="M3 1.78zM3.826 1.78zM3 1.702zM3.826 1.702z" /><path
										d="M3.413 2.091c.44 0 .84.179 1.128.467.288.289.467.687.467 1.128 0 .44-.179.838-.467 1.127a1.59 1.59 0 0 1-1.128.467 1.59 1.59 0 0 1-1.127-.467 1.59 1.59 0 0 1-.467-1.127c0-.44.178-.84.467-1.128a1.59 1.59 0 0 1 1.127-.467z"
										style="fill:#6d4c41;fill-rule:nonzero"
									/><path
										d="M3.93 2.832a.078.078 0 0 0-.138-.07l-.3.597V2.48a.078.078 0 1 0-.156 0v1.206a.078.078 0 0 0 .15.03l.445-.884z"
										style="fill:#ffb74d;fill-rule:nonzero"
									/><path
										class="this-is-ticking-hand"
										style="fill:none;fill-rule:nonzero"
										d="m1.895 2 3.147 3.147-.11.11-3.147-3.146z"
									/></g
								></svg
							>

							<div class="text-lg font-semibold uppercase">Count</div>
						</div>
						<div class="flex h-full items-center">
							<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
								><MaterialSymbolsChevronRight class="size-6" /></button
							>
						</div>
					</div></a
				>
			</div>
		</div>
	</main>
</PageWrapper>
