<script lang="ts">
	import { onMount } from 'svelte';
	import { play } from '$lib/play';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { Capacitor } from '@capacitor/core';
	import { KeepAwake } from '@capacitor-community/keep-awake';
	import {
		userQueryOptions,
		timerProfilesQueryOptions,
		timerProfilesRefetchOptions
	} from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Icon from '@iconify/svelte';
	import QuickTimePicker from './QuickTimePicker.svelte';
	import TimerDisplay from './TimerDisplay.svelte';
	import SegmentList from './SegmentList.svelte';
	import SegmentProgress from './SegmentProgress.svelte';
	import ProgressRing from './ProgressRing.svelte';

	const user = createQuery(userQueryOptions);
	const profilesQuery = createQuery(timerProfilesQueryOptions);
	const queryClient = useQueryClient();
	const isAndroid = Capacitor.getPlatform() === 'android';
	const tickingSoundPath = '/soft-ticking.mp3';

	let character: Characters = $derived.by(() => {
		if (!user || !user.data || user.isPending || user.data.preferredCharacter === 'default')
			return 'robot';
		return user.data.preferredCharacter;
	});

	let successSoundPath = $derived(`/${character}/timesup.mp3`);

	let mode: 'quick' | 'profile' = $state('quick');

	let targetMinutes = $state(5);
	let targetSeconds = $state(0);
	let targetTotalSeconds = $derived(targetMinutes * 60 + targetSeconds);

	let profiles = $derived(profilesQuery.data ?? []);
	let selectedProfileId = $state('');
	let selectedProfile = $derived(profiles.find((p) => p.id === selectedProfileId));
	let segmentDurations: number[] = $state([]);
	let currentSegmentIndex = $state(0);

	$effect(() => {
		if (profiles.length > 0 && !selectedProfileId) {
			const defaultProfile = profiles.find((p) => p.isDefault);
			selectedProfileId = defaultProfile?.id ?? profiles[0].id;
		}
	});

	$effect(() => {
		if (selectedProfile) {
			segmentDurations = selectedProfile.segments.map((s) => s.defaultSeconds);
			currentSegmentIndex = 0;
		}
	});

	let profileModal: HTMLDialogElement | undefined = $state();
	let editingProfile: TimerProfileDB | null = $state(null);
	let profileName = $state('');
	let profileSegments: TimerSegmentDef[] = $state([]);
	let profileIsDefault = $state(false);
	let saving = $state(false);

	function openNewProfile() {
		editingProfile = null;
		profileName = '';
		profileSegments = [{ label: 'Focus', defaultSeconds: 1500 }];
		profileIsDefault = false;
		profileModal?.showModal();
	}

	function openEditProfile(profile: TimerProfileDB) {
		editingProfile = profile;
		profileName = profile.name;
		profileSegments = profile.segments.map((s) => ({ ...s }));
		profileIsDefault = profile.isDefault;
		profileModal?.showModal();
	}

	function addSegment() {
		profileSegments = [...profileSegments, { label: '', defaultSeconds: 300 }];
	}

	function removeSegment(index: number) {
		profileSegments = profileSegments.filter((_, i) => i !== index);
	}

	async function saveProfile() {
		if (!profileName.trim() || profileSegments.length === 0) return;
		saving = true;

		const body: TimerProfileInput = {
			name: profileName.trim(),
			segments: profileSegments.filter((s) => s.label.trim()),
			isDefault: profileIsDefault
		};

		if (editingProfile) {
			const res = await api.patch(`timer-profiles/${editingProfile.id}`, {
				body: JSON.stringify(body)
			});
			if (res.status === 204) {
				addToast('success', 'Profile updated!');
			} else {
				addToast('error', 'Failed to update profile');
			}
		} else {
			const res = await api.post('timer-profiles', {
				body: JSON.stringify(body)
			});
			if (res.status === 201) {
				const created = await res.json<TimerProfileDB>();
				selectedProfileId = created.id;
				addToast('success', 'Profile created!');
			} else {
				addToast('error', 'Failed to create profile');
			}
		}

		saving = false;
		profileModal?.close();
		queryClient.refetchQueries(timerProfilesRefetchOptions());
	}

	async function deleteProfile(profile: TimerProfileDB) {
		const res = await api.delete(`timer-profiles/${profile.id}`);
		if (res.status === 204) {
			if (selectedProfileId === profile.id) {
				selectedProfileId = '';
			}
			addToast('success', 'Profile deleted!');
			queryClient.refetchQueries(timerProfilesRefetchOptions());
		} else {
			addToast('error', 'Failed to delete profile');
		}
	}

	let pauseTarget = $state(0);
	let remainingSeconds = $state(0);
	let min: number = $state(0);
	let sec: number = $state(0);
	let started = $state(false);
	let timerInterval: ReturnType<typeof setInterval> | undefined = $state();
	let transitionTimeout: ReturnType<typeof setTimeout> | undefined = $state();

	function doubleDigits(num: number): string {
		if (num < 0) return '00';
		return String(num).padStart(2, '0');
	}

	function getMinSec(total: number) {
		return { min: Math.floor(total / 60), sec: total % 60 };
	}

	function tick() {
		remainingSeconds -= 1;
		if (remainingSeconds > 0 && playSound) {
			audioPlayer?.play();
			play(remainingSeconds, character);
		}

		const ms = getMinSec(remainingSeconds);
		min = ms.min;
		sec = ms.sec;

		if (remainingSeconds === 0) {
			clearInterval(timerInterval);
			timerInterval = undefined;

			if (playSound) {
				new Audio(successSoundPath).play();
			}

			if (mode === 'profile' && currentSegmentIndex < segmentDurations.length - 1) {
				// Auto go to next segment after 1s pause
				transitionTimeout = setTimeout(() => {
					currentSegmentIndex += 1;
					remainingSeconds = segmentDurations[currentSegmentIndex];
					const ms = getMinSec(remainingSeconds);
					min = ms.min;
					sec = ms.sec;
					tick();
					timerInterval = setInterval(tick, 1000);
				}, 1000);
			} else {
				started = false;
			}
		}
	}

	function start() {
		if (mode === 'quick') {
			if (targetTotalSeconds === 0) return;
			remainingSeconds = targetTotalSeconds;
		} else {
			if (!segmentDurations.length || segmentDurations.every((d) => d === 0)) return;
			currentSegmentIndex = 0;
			remainingSeconds = segmentDurations[0];
		}

		started = true;
		tick();
		timerInterval = setInterval(tick, 1000);
	}

	function stop() {
		clearInterval(timerInterval);
		clearTimeout(transitionTimeout);
		timerInterval = undefined;
		transitionTimeout = undefined;
		remainingSeconds = 0;
		pauseTarget = 0;
		currentSegmentIndex = 0;
		started = false;
	}

	function pause() {
		pauseTarget = min * 60 + sec;
		clearInterval(timerInterval);
		clearTimeout(transitionTimeout);
		timerInterval = undefined;
		transitionTimeout = undefined;
		started = false;
	}

	function resume() {
		remainingSeconds = pauseTarget;
		started = true;
		tick();
		timerInterval = setInterval(tick, 1000);
	}

	function skip() {
		if (mode !== 'profile') return;
		clearInterval(timerInterval);
		clearTimeout(transitionTimeout);
		timerInterval = undefined;
		transitionTimeout = undefined;

		if (currentSegmentIndex < segmentDurations.length - 1) {
			if (playSound) {
				new Audio(successSoundPath).play();
			}
			currentSegmentIndex += 1;
			remainingSeconds = segmentDurations[currentSegmentIndex];
			const ms = getMinSec(remainingSeconds);
			min = ms.min;
			sec = ms.sec;
			tick();
			timerInterval = setInterval(tick, 1000);
		} else {
			remainingSeconds = 0;
			min = 0;
			sec = 0;
			started = false;
		}
	}

	let isTimerActive = $derived(started || pauseTarget > 0);

	let displayMin = $derived(
		!isTimerActive && mode === 'quick' ? doubleDigits(targetMinutes) : doubleDigits(min)
	);
	let displaySec = $derived(
		!isTimerActive && mode === 'quick' ? doubleDigits(targetSeconds) : doubleDigits(sec)
	);

	let timerTotal = $derived.by(() => {
		if (mode === 'quick') return targetTotalSeconds;
		return segmentDurations[currentSegmentIndex] ?? 0;
	});
	let progress = $derived(timerTotal > 0 ? remainingSeconds / timerTotal : 1);

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

	async function changeCharacterHandler(char: string) {
		const response = await api.patch('users/me/character', {
			body: JSON.stringify({ preferredCharacter: char })
		});
		if (response.status !== 204 && response.status !== 201) {
			addToast('error', 'Error saving character choice!');
		}
	}

	let canStart = $derived.by(() => {
		if (mode === 'quick') return targetTotalSeconds > 0;
		return segmentDurations.length > 0 && segmentDurations.some((d) => d > 0);
	});
</script>

<PageWrapper title="Count">
	<div
		class="grid min-h-[calc(100vh-3.5rem-6rem)] w-full grid-rows-[auto_1fr_auto] justify-items-center"
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
			<div
				class={[
					'transition-opacity duration-300',
					isTimerActive ? 'pointer-events-none opacity-60' : ''
				]}
			>
				<SegmentedControl items={2}>
					<label>
						<input type="radio" bind:group={mode} value="quick" name="mode" />Quick
					</label>
					<label>
						<input type="radio" bind:group={mode} value="profile" name="mode" />Profiles
					</label>
				</SegmentedControl>
			</div>

			{#if mode === 'quick'}
				<ProgressRing {progress} active={isTimerActive}>
					<QuickTimePicker
						bind:minutes={targetMinutes}
						bind:seconds={targetSeconds}
						disabled={isTimerActive}
						displayMinutes={isTimerActive ? displayMin : undefined}
						displaySeconds={isTimerActive ? displaySec : undefined}
					/>
				</ProgressRing>
			{:else if !isTimerActive}
				{#if profilesQuery.isPending}
					<span class="loading loading-spinner loading-md"></span>
				{:else}
					<div class="flex w-full max-w-lg flex-wrap gap-2">
						{#each profiles as profile (profile.id)}
							<button
								class={[
									'btn grow rounded-full',
									selectedProfileId === profile.id ? 'btn-primary' : 'btn-neutral btn-soft'
								]}
								onclick={() => (selectedProfileId = profile.id)}
							>
								{profile.name}
							</button>
						{/each}
						<button class="btn btn-neutral btn-soft grow rounded-full" onclick={openNewProfile}>
							+ New
						</button>
					</div>

					{#if selectedProfile}
						<SegmentList segments={selectedProfile.segments} bind:durations={segmentDurations} />
						<div class="flex gap-2">
							<button
								class="btn btn-ghost btn-sm"
								onclick={() => selectedProfile && openEditProfile(selectedProfile)}
							>
								Edit Profile
							</button>
							<button
								class="btn btn-ghost btn-sm text-error"
								onclick={() => selectedProfile && deleteProfile(selectedProfile)}
							>
								Delete
							</button>
						</div>
					{:else if profiles.length === 0}
						<p class="text-base-content/50 text-sm">No profiles yet. Create one to get started.</p>
					{/if}
				{/if}
			{:else}
				{#if selectedProfile}
					<SegmentProgress
						segments={selectedProfile.segments}
						currentIndex={currentSegmentIndex}
						isRunning={started}
					/>
					<p class="text-primary text-lg font-semibold">
						{selectedProfile.segments[currentSegmentIndex]?.label}
					</p>
				{/if}
				<ProgressRing {progress} active={isTimerActive}>
					<TimerDisplay minutes={displayMin} seconds={displaySec} />
				</ProgressRing>
			{/if}

			<div class="grid w-full max-w-lg grid-cols-2 gap-2">
				{#if mode === 'profile' && isTimerActive}
					<button
						class="btn btn-lg btn-primary w-full rounded-full p-8"
						onclick={() => {
							if (started) pause();
							else resume();
						}}
					>
						{started ? 'Pause' : 'Resume'}
					</button>
					<div class="grid grid-cols-2 gap-2">
						<button class="btn btn-lg btn-neutral btn-soft w-full rounded-full p-8" onclick={skip}>
							Skip
						</button>
						<button class="btn btn-lg btn-neutral btn-soft w-full rounded-full p-8" onclick={stop}>
							Reset
						</button>
					</div>
				{:else}
					<button
						class="btn btn-lg btn-primary w-full min-w-40 rounded-full p-8"
						onclick={() => {
							if (!started && !pauseTarget) {
								start();
								return;
							}
							if (started) pause();
							else resume();
						}}
						disabled={!isTimerActive && !canStart}
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
						class="btn btn-lg btn-neutral btn-soft w-full min-w-40 rounded-full p-8"
						onclick={stop}
					>
						Reset
					</button>
				{/if}
			</div>

			<div>
				<button
					aria-label="toggle sound"
					class="btn btn-ghost btn-xl"
					onclick={() => (playSound = !playSound)}
				>
					{#if playSound}
						<Icon icon="material-symbols:volume-up" width="1em" height="1em" />
					{:else}
						<Icon icon="material-symbols:volume-off" width="1em" height="1em" />
					{/if}
				</button>
			</div>
		</main>

		<SegmentedControl class="scale-80" items={2}>
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

<dialog bind:this={profileModal} class="modal">
	<div class="modal-box">
		<h3 class="text-lg font-bold">{editingProfile ? 'Edit Profile' : 'New Profile'}</h3>

		<div class="mt-4 grid gap-4">
			<label class="form-control w-full">
				<div class="label"><span class="label-text">Name</span></div>
				<input
					type="text"
					class="input input-bordered w-full"
					placeholder="e.g. Pomodoro"
					bind:value={profileName}
				/>
			</label>

			<div>
				<div class="label"><span class="label-text">Segments</span></div>
				<div class="grid gap-2">
					{#each profileSegments as segment, i (i)}
						<div class="flex items-center gap-2">
							<input
								type="text"
								class="input input-bordered input-sm grow"
								placeholder="Label"
								bind:value={segment.label}
							/>
							<input
								type="number"
								class="input input-bordered input-sm w-20"
								placeholder="Min"
								min="0"
								value={Math.floor(segment.defaultSeconds / 60)}
								oninput={(e) => {
									const min = parseInt(e.currentTarget.value) || 0;
									const sec = segment.defaultSeconds % 60;
									profileSegments[i].defaultSeconds = min * 60 + sec;
								}}
							/>
							<span class="text-base-content/50 text-sm">min</span>
							{#if profileSegments.length > 1}
								<button class="btn btn-ghost btn-sm btn-square" onclick={() => removeSegment(i)}>
									<Icon icon="material-symbols:remove" width="16" height="16" />
								</button>
							{/if}
						</div>
					{/each}
				</div>
				<button class="btn btn-ghost btn-sm mt-2" onclick={addSegment}>+ Add Segment</button>
			</div>

			<label class="label cursor-pointer justify-start gap-3">
				<input type="checkbox" class="toggle toggle-primary" bind:checked={profileIsDefault} />
				<span class="label-text">Set as default</span>
			</label>
		</div>

		<div class="modal-action">
			<form method="dialog">
				<button class="btn btn-ghost">Cancel</button>
			</form>
			<button
				class="btn btn-primary"
				onclick={saveProfile}
				disabled={saving || !profileName.trim() || profileSegments.length === 0}
			>
				{saving ? 'Saving...' : editingProfile ? 'Update' : 'Create'}
			</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop"><button>close</button></form>
</dialog>
