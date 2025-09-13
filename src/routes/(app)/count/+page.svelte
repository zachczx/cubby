<script lang="ts">
	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import { play } from '$lib/play';
	import PageWrapper from '$lib/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { Capacitor } from '@capacitor/core';
	dayjs.extend(duration);

	const tickingSoundPath = '/soft-ticking.mp3';
	const successSoundPath = '/timesup.mp3';

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

			if (totalSec % 60 === 0) {
				play(totalSec / 60);
			}
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

		started = false;
	}

	$inspect(min, sec);
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

	let playSound = $state(false);
	let audioPlayer: HTMLAudioElement | undefined = $state();

	onMount(() => {
		if (Capacitor.getPlatform() !== 'android') {
			audioPlayer = new Audio(tickingSoundPath);
		}

		return () => {
			if (timerInterval) {
				clearInterval(timerInterval);
			}
		};
	});
</script>

<svelte:head>
	<title>Count</title>
</svelte:head>
<PageWrapper title="Count" {pb}>
	<main class="grid h-full content-center justify-items-center gap-8 p-2">
		<!-- <h1 class="w-full max-w-lg pb-2 text-center text-4xl font-medium text-secondary/60">Countdown</h1> -->
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
					<span class="countdown-custom text-8xl">
						{#if !started && !pauseTarget}
							<span
								style={`--value:${targetMin};`}
								aria-live="polite"
								aria-label={String(targetMin)}>{targetMin}</span
							>
						{:else}
							<span style={`--value:${min};`} aria-live="polite" aria-label={String(min)}
								>{min}</span
							>
						{/if}
					</span>
					min
				</div>
				<div class="flex flex-col px-2 text-center lg:px-8">
					<span class="countdown-custom text-8xl">
						{#if !started && !pauseTarget}
							<span
								style={`--value:${targetSec};`}
								aria-live="polite"
								aria-label={String(targetSec)}>{targetSec}</span
							>
						{:else}
							<span style={`--value:${sec};`} aria-live="polite" aria-label={String(sec)}
								>{sec}</span
							>
						{/if}
					</span>
					sec
				</div>
			</div>
		</div>

		<div class="grid w-full max-w-lg grid-cols-2 gap-2">
			<button
				class="btn btn-lg btn-primary w-full"
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
			<button class="btn btn-lg btn-neutral w-full" onclick={() => stop()}>Reset</button>
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
</PageWrapper>

<style>
	.countdown-custom {
		display: inline-flex;

		& > * {
			display: inline-block;
			overflow-y: hidden;
			height: 1em;
			&:before {
				position: relative;
				content: '00\A 01\A 02\A 03\A 04\A 05\A 06\A 07\A 08\A 09\A 10\A 11\A 12\A 13\A 14\A 15\A 16\A 17\A 18\A 19\A 20\A 21\A 22\A 23\A 24\A 25\A 26\A 27\A 28\A 29\A 30\A 31\A 32\A 33\A 34\A 35\A 36\A 37\A 38\A 39\A 40\A 41\A 42\A 43\A 44\A 45\A 46\A 47\A 48\A 49\A 50\A 51\A 52\A 53\A 54\A 55\A 56\A 57\A 58\A 59\A 60\A 61\A 62\A 63\A 64\A 65\A 66\A 67\A 68\A 69\A 70\A 71\A 72\A 73\A 74\A 75\A 76\A 77\A 78\A 79\A 80\A 81\A 82\A 83\A 84\A 85\A 86\A 87\A 88\A 89\A 90\A 91\A 92\A 93\A 94\A 95\A 96\A 97\A 98\A 99\A';
				white-space: pre;
				top: calc(var(--value) * -1em);
				text-align: center;
			}
		}
	}
</style>
