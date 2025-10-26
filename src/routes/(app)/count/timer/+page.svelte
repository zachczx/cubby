<script lang="ts">
	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	dayjs.extend(duration);

	let startDateTime: dayjs.Dayjs | undefined = $state();
	let df: duration.Duration | undefined = $state();

	let sec = $state(0);
	let min = $state(0);

	let timerInterval: number | undefined = $state();

	function tick() {
		const now = dayjs();
		df = dayjs.duration(now.diff(startDateTime));
		sec = df.seconds();
		min = df.minutes();

		if (sec > 0 && playSound) {
			audioPlayer?.play();
		}
	}

	function start() {
		startDateTime = dayjs();

		tick();

		timerInterval = setInterval(tick, 1000);
	}

	function stop() {
		clearInterval(timerInterval);
		startDateTime = undefined;
		sec = 0;
		min = 0;
	}

	const soundPath = '/soft-ticking.mp3';

	let playSound = $state(true);
	let audioPlayer: HTMLAudioElement | undefined = $state();
	onMount(() => {
		audioPlayer = new Audio(soundPath);

		return () => {
			if (timerInterval) {
				clearInterval(timerInterval);
			}
		};
	});
</script>

<main class="grid min-h-dvh content-center justify-items-center gap-8">
	<!-- <h1 class="w-full max-w-lg pb-2 text-center text-4xl font-medium text-secondary/60">Count</h1> -->
	<div class="grid max-w-lg auto-cols-max grid-flow-col gap-12 text-center">
		<div class="flex flex-col">
			<span class="countdown text-8xl">
				<span style={`--value:${min};`} aria-live="polite" aria-label={String(min)}>{min}</span>
			</span>
			min
		</div>
		<div class="flex flex-col">
			<span class="countdown text-8xl transition-all ease-out">
				<span style={`--value:${sec};`} aria-live="polite" aria-label={String(sec)}>{sec}</span>
			</span>
			sec
		</div>
	</div>

	<div class="grid w-full max-w-lg grid-cols-2 gap-4">
		<button class="btn btn-lg btn-primary w-full" onclick={() => start()}>Start</button>
		<button class="btn btn-lg btn-error w-full" onclick={() => stop()}>Reset</button>
	</div>

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
</main>

<style>
	.countdown > * {
		&::before {
			transition: all 0.3s cubic-bezier(1, 0, 0, 1) !important;
		}
	}
</style>
