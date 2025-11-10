<script lang="ts">
	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import { play } from '$lib/play';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { Capacitor } from '@capacitor/core';
	import { createUserQueryOptions } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import { beforeNavigate } from '$app/navigation';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	dayjs.extend(duration);

	const user = createQuery(createUserQueryOptions);
	const tanstackClient = useQueryClient();

	const tickingSoundPath = '/soft-ticking.mp3';

	let character: Characters = $state('robot');

	let successSoundPath = $derived(`/${character}/timesup.mp3`);

	let target = $state(5);
	let pauseTarget = $state(0);
	let targetMin = $derived(target);
	let targetSec = $derived(0);

	let totalSec = $state(0);
	let min: number = $state(0);
	let sec: number = $state(0);

	let started = $state(false);

	let timerInterval: ReturnType<typeof setTimeout> | undefined = $state();

	interface MinSec {
		min: number;
		sec: number;
	}

	function getMinSec(total: number): MinSec {
		const m = Math.floor(total / 60);
		const s = total % 60;
		return { min: m, sec: s };
	}

	function tick() {
		totalSec -= 1;
		if (totalSec > 0 && playSound) {
			audioPlayer?.play();

			play(totalSec, character);
		}

		const minSec = getMinSec(totalSec);
		min = minSec.min;
		sec = minSec.sec;

		if (totalSec === 0) {
			if (playSound) {
				new Audio(successSoundPath).play();
			}
			clearInterval(timerInterval);
			return;
		}
	}

	function start() {
		totalSec = target * 60;
		started = true;

		tick();
		timerInterval = setInterval(tick, 1000);
	}

	function stop() {
		clearInterval(timerInterval);
		totalSec = 0;
		pauseTarget = 0;
		target = 5;

		started = false;
	}

	function pause() {
		pauseTarget = min * 60 + sec;
		clearInterval(timerInterval);
		started = false;
	}

	function resume() {
		totalSec = pauseTarget;
		started = true;

		tick();
		timerInterval = setInterval(tick, 1000);
	}

	let playSound = $derived.by(() => {
		if (user.isSuccess && !user.data.mute) {
			return true;
		}
		return false;
	});
	let audioPlayer: HTMLAudioElement | undefined = $state();

	onMount(async () => {
		if (Capacitor.getPlatform() !== 'android') {
			audioPlayer = new Audio(tickingSoundPath);
		}

		if (page.url.searchParams.get('start')) {
			setTimeout(() => start(), 700);
		}

		beforeNavigate(() => stop());
	});

	function doubleDigits(num: number): string {
		if (num < 10) {
			return 0 + String(num);
		}
		return String(num);
	}
</script>

<svelte:head>
	<title>Count</title>
</svelte:head>

<PageWrapper title="Count" {pb}>
	<div class="grid h-full max-w-xl grid-rows-[1fr_auto] justify-self-center">
		<main class="grid h-full content-center justify-items-center gap-8 p-2">
			<div class="grid w-full max-w-lg grid-cols-[auto_1fr] justify-items-center gap-y-8">
				<div class="grid content-center opacity-75">
					<button
						class="btn btn-ghost"
						aria-label="+1"
						onclick={() => {
							target += 1;
						}}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-up size-6 lg:size-8"
							><path fill="currentColor" d="m12 10.8l-4.6 4.6L6 14l6-6l6 6l-1.4 1.4z" /></svg
						></button
					>
					<button
						class="btn btn-ghost"
						aria-label="-1"
						onclick={() => {
							target -= 1;
						}}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-down size-6 lg:size-8"
							><path fill="currentColor" d="m12 15.4l-6-6L7.4 8l4.6 4.6L16.6 8L18 9.4z" /></svg
						></button
					>
				</div>
				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col px-2 text-center lg:px-8">
						<span class="text-8xl">
							{#if !started && !pauseTarget}
								<span
									style={`--value:${targetMin};`}
									aria-live="polite"
									aria-label={String(targetMin)}>{doubleDigits(targetMin)}</span
								>
							{:else}
								<span style={`--value:${min};`} aria-live="polite" aria-label={String(min)}
									>{doubleDigits(min)}</span
								>
							{/if}
						</span>
						min
					</div>
					<div class="flex flex-col px-2 text-center lg:px-8">
						<span class="text-8xl">
							{#if !started && !pauseTarget}
								<span
									style={`--value:${targetSec};`}
									aria-live="polite"
									aria-label={String(targetSec)}>{doubleDigits(targetSec)}</span
								>
							{:else}
								<span style={`--value:${sec};`} aria-live="polite" aria-label={String(sec)}
									>{doubleDigits(sec)}</span
								>
							{/if}
						</span>
						sec
					</div>
				</div>
			</div>

			<div class="grid w-full max-w-lg grid-cols-2 gap-2">
				<button
					class="btn btn-lg btn-primary w-full rounded-full"
					onclick={() => {
						if (!started && !pauseTarget) {
							start();
							return;
						}

						if (started) {
							pause();
						} else {
							resume();
						}
					}}
				>
					{#if !started && !pauseTarget}
						Start
					{:else if started}
						Pause
					{:else}
						Resume
					{/if}
				</button>
				<button
					class="btn btn-lg btn-neutral btn-outline w-full rounded-full"
					onclick={() => stop()}>Reset</button
				>
			</div>

			<div>
				<button
					aria-label="toggle sound"
					class="btn btn-ghost btn-xl"
					onclick={() => {
						playSound = !playSound;
					}}
				>
					{#if playSound}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:volume-up"
							><path
								fill="currentColor"
								d="M14 20.725v-2.05q2.25-.65 3.625-2.5t1.375-4.2t-1.375-4.2T14 5.275v-2.05q3.1.7 5.05 3.138T21 11.975t-1.95 5.613T14 20.725M3 15V9h4l5-5v16l-5-5zm11 1V7.95q1.175.55 1.838 1.65T16.5 12q0 1.275-.663 2.363T14 16"
							/></svg
						>
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:volume-off"
							><path
								fill="currentColor"
								d="m19.8 22.6l-3.025-3.025q-.625.4-1.325.688t-1.45.462v-2.05q.35-.125.688-.25t.637-.3L12 14.8V20l-5-5H3V9h3.2L1.4 4.2l1.4-1.4l18.4 18.4zm-.2-5.8l-1.45-1.45q.425-.775.638-1.625t.212-1.75q0-2.35-1.375-4.2T14 5.275v-2.05q3.1.7 5.05 3.138T21 11.975q0 1.325-.363 2.55T19.6 16.8m-3.35-3.35L14 11.2V7.95q1.175.55 1.838 1.65T16.5 12q0 .375-.062.738t-.188.712M12 9.2L9.4 6.6L12 4z"
							/></svg
						>
					{/if}
				</button>
			</div>
		</main>

		<SegmentedControl items={2}>
			<label>
				<input type="radio" bind:group={character} value="robot" name="character" />Robot
			</label>
			<label>
				<input type="radio" bind:group={character} value="furnando" name="character" />Furnando
			</label>
		</SegmentedControl>
	</div>
</PageWrapper>
