<script lang="ts">
	import { onMount } from 'svelte';
	import { play } from '$lib/play';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { Capacitor } from '@capacitor/core';
	import { KeepAwake } from '@capacitor-community/keep-awake';
	import { userQueryOptions } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';

	const user = createQuery(userQueryOptions);
	const isAndroid = Capacitor.getPlatform() === 'android';
	const tickingSoundPath = '/soft-ticking.mp3';

	let character: Characters = $derived.by(() => {
		if (!user || !user.data || user.isPending || user.data.preferredCharacter === 'default')
			return 'robot';

		return user.data.preferredCharacter;
	});

	let successSoundPath = $derived(`/${character}/timesup.mp3`);

	let pauseTarget = $state(0);
	let targetMinutes = $state(5);
	let targetSeconds = $state(0);
	let targetTotalSeconds = $derived(targetMinutes * 60 + targetSeconds);

	let remainingSeconds = $state(0);
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
		remainingSeconds -= 1;
		if (remainingSeconds > 0 && playSound) {
			audioPlayer?.play();

			play(remainingSeconds, character);
		}

		const minSec = getMinSec(remainingSeconds);
		min = minSec.min;
		sec = minSec.sec;

		if (remainingSeconds === 0) {
			if (playSound) {
				new Audio(successSoundPath).play();
			}
			clearInterval(timerInterval);
			return;
		}
	}

	function start() {
		if (targetTotalSeconds === 0) return;

		remainingSeconds = targetTotalSeconds;
		started = true;

		tick();
		timerInterval = setInterval(tick, 1000);
	}

	function stop() {
		clearInterval(timerInterval);
		remainingSeconds = 0;
		pauseTarget = 0;

		started = false;
	}

	function pause() {
		pauseTarget = min * 60 + sec;
		clearInterval(timerInterval);
		started = false;
	}

	function resume() {
		remainingSeconds = pauseTarget;
		started = true;

		tick();
		timerInterval = setInterval(tick, 1000);
	}

	let playSound = $state(false);
	let playSoundInitialized = false;

	$effect(() => {
		if (user.isSuccess && !playSoundInitialized) {
			playSound = !!user.data.soundOn;
			playSoundInitialized = true;
		}
	});

	let audioPlayer: HTMLAudioElement | undefined = $state();

	let keepAwakeOption = $state(true);
	$effect(() => {
		if (keepAwakeOption) {
			KeepAwake.keepAwake();
		} else {
			KeepAwake.allowSleep();
		}
	});

	onMount(() => {
		if (Capacitor.getPlatform() !== 'android') {
			audioPlayer = new Audio(tickingSoundPath);
		}

		let timeout: NodeJS.Timeout;

		if (page.url.searchParams.get('start')) {
			timeout = setTimeout(() => start(), 700);
		}

		return () => {
			clearTimeout(timeout);
			KeepAwake.allowSleep();
			stop();
		};
	});

	function doubleDigits(num: number): string {
		if (num < 0) return '00';

		if (num < 10 && num >= 0) {
			return 0 + String(num);
		}

		return String(num);
	}

	async function changeCharacterHandler(char: string) {
		const response = await api.patch('users/me/character', {
			body: JSON.stringify({
				preferredCharacter: char
			})
		});

		if (response.status !== 204 && response.status !== 201) {
			addToast('error', 'Error saving character choice!');
		}
	}

	function userChangeMinutes(direction: 'increment' | 'decrement') {
		if (direction === 'increment') {
			targetMinutes += 1;
		}

		if (direction === 'decrement') {
			if (targetMinutes <= 0) return;

			targetMinutes -= 1;
		}
	}

	function userChangeSeconds(direction: 'increment' | 'decrement') {
		if (direction === 'increment') {
			if (targetSeconds >= 59) return;

			targetSeconds += 1;
		}

		if (direction === 'decrement') {
			if (targetSeconds <= 0) return;

			targetSeconds -= 1;
		}
	}
</script>

<svelte:head>
	<title>Count</title>
</svelte:head>

<PageWrapper title="Count">
	<div
		class={[
			'grid min-h-[calc(100vh-3.5rem-6rem)] w-full grid-rows-[auto_1fr_auto] justify-items-center lg:min-h-[calc(100vh-3.5rem-1rem)]'
		]}
	>
		<label
			class={[
				'btn w-36 justify-self-end',
				keepAwakeOption ? 'btn-soft btn-primary' : 'btn-neutral btn-soft opacity-75',
				!isAndroid && 'opacity-0'
			]}
		>
			<input type="checkbox" class="hidden" bind:checked={keepAwakeOption} />
			{#if keepAwakeOption}
				Screen: Awake
			{:else}
				Screen: Auto
			{/if}
		</label>

		<main class="grid h-full max-w-xl content-center justify-items-center gap-8 p-2">
			<div class="grid w-full max-w-lg grid-cols-2 justify-items-center gap-4">
				<div class="flex flex-col px-2 text-center lg:px-8">
					<button
						class="btn btn-ghost text-primary opacity-75"
						aria-label="+1"
						onclick={() => userChangeMinutes('increment')}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-up size-8"
							><path fill="currentColor" d="m12 10.8l-4.6 4.6L6 14l6-6l6 6l-1.4 1.4z" /></svg
						></button
					>
					<span class="text-8xl">
						{#if !started && !pauseTarget}
							<span
								style={`--value:${targetMinutes};`}
								aria-live="polite"
								aria-label={String(targetMinutes)}>{doubleDigits(targetMinutes)}</span
							>
						{:else}
							<span style={`--value:${min};`} aria-live="polite" aria-label={String(min)}
								>{doubleDigits(min)}</span
							>
						{/if}
					</span>

					<button
						class="btn btn-ghost text-primary opacity-75"
						aria-label="-1"
						onclick={() => userChangeMinutes('decrement')}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-down size-8"
							><path fill="currentColor" d="m12 15.4l-6-6L7.4 8l4.6 4.6L16.6 8L18 9.4z" /></svg
						></button
					>
				</div>
				<div class="flex flex-col px-2 text-center lg:px-8">
					<button
						class="btn btn-ghost text-primary opacity-75"
						aria-label="+1"
						onclick={() => userChangeSeconds('increment')}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-up size-8"
							><path fill="currentColor" d="m12 10.8l-4.6 4.6L6 14l6-6l6 6l-1.4 1.4z" /></svg
						></button
					>
					<span class="text-8xl">
						{#if !started && !pauseTarget}
							<span
								style={`--value:${targetSeconds};`}
								aria-live="polite"
								aria-label={String(targetSeconds)}>{doubleDigits(targetSeconds)}</span
							>
						{:else}
							<span style={`--value:${sec};`} aria-live="polite" aria-label={String(sec)}
								>{doubleDigits(sec)}</span
							>
						{/if}
					</span>
					<button
						class="btn btn-ghost text-primary opacity-75"
						aria-label="-1"
						onclick={() => userChangeSeconds('decrement')}
						><svg
							xmlns="http://www.w3.org/2000/svg"
							width="1em"
							height="1em"
							viewBox="0 0 24 24"
							class="material-symbols:keyboard-arrow-down size-8"
							><path fill="currentColor" d="m12 15.4l-6-6L7.4 8l4.6 4.6L16.6 8L18 9.4z" /></svg
						></button
					>
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
					disabled={targetTotalSeconds === 0}
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
				<input
					type="radio"
					bind:group={character}
					value="robot"
					name="character"
					onclick={() => changeCharacterHandler('robot')}
				/>Robot
			</label>
			<label>
				<input
					type="radio"
					bind:group={character}
					value="furnando"
					name="character"
					onclick={() => changeCharacterHandler('furnando')}
				/>Furnando
			</label>
		</SegmentedControl>
	</div>
</PageWrapper>
