<script lang="ts">
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import type { RecordModel } from 'pocketbase';
	import { addToast } from './ArkToaster.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';
	import { getAllLogsQueryKey } from '$lib/queries';

	// @ts-ignore
	import confetti from 'canvas-confetti';

	let {
		text,
		query,
		compact = false,
		color = 'primary',
		rounded = 'full',
		class: className,
		icon: CustomIcon,
		celebrate = true
	}: {
		text: string | undefined;
		query: () => Promise<RecordModel>;
		compact?: boolean;
		color?: 'primary' | 'neutral';
		rounded?: 'full' | '3xl' | '2xl' | 'xl' | 'lg' | 'md' | 'sm' | 'xs';
		class?: string;
		icon?: Component;
		celebrate?: boolean;
	} = $props();

	let status: ButtonState = $state('default');
	let isCelebrating = $state(false);

	const tanstackClient = useQueryClient();
	export const insertNewLogToCache = (newLog: RecordModel) =>
		tanstackClient.setQueryData(getAllLogsQueryKey(), (oldLogs: LogsDB[] | undefined) => {
			if (!oldLogs) return [newLog];
			return [newLog, ...oldLogs];
		});

	/**
	 * Confetti
	 */

	let buttonRef = $state<HTMLButtonElement>();

	interface Options {
		spread: number;
		startVelocity?: number;
		decay?: number;
		scalar?: number;
		zIndex?: number;
		origin: Point;
		gravity?: number;
	}

	type Point = { x: number; y: number };

	const count = 150;

	function triggerConfetti() {
		if (!buttonRef) return;

		const rect = buttonRef.getBoundingClientRect();
		const x = (rect.left + rect.width / 2) / window.innerWidth;
		const y = (rect.top + rect.height / 2) / window.innerHeight;

		const origin = { x, y };

		if (compact) {
			fire(0.25, { spread: 26, startVelocity: 25, origin, gravity: 0.31 });
			fire(0.2, { spread: 60, startVelocity: 15, origin, gravity: 0.3 });
			fire(0.35, {
				spread: 100,
				startVelocity: 10,
				decay: 0.91,
				scalar: 0.8,
				origin,
				gravity: 0.5
			});
			fire(0.1, { spread: 120, startVelocity: 15, decay: 0.92, scalar: 1.2, origin, gravity: 0.6 });
			fire(0.1, { spread: 120, startVelocity: 20, origin, gravity: 0.7 });
		} else {
			fire(0.25, { spread: 26, startVelocity: 55, origin });
			fire(0.2, { spread: 60, origin });
			fire(0.35, { spread: 100, decay: 0.91, scalar: 0.8, origin });
			fire(0.1, { spread: 120, startVelocity: 25, decay: 0.92, scalar: 1.2, origin });
			fire(0.1, { spread: 120, startVelocity: 45, origin });
		}
	}

	function fire(particleRatio: number, opts: Options) {
		confetti({
			...opts,
			colors: [
				'#ff7f50', // Coral (bright accent)
				'#ffb347', // Pastel orange
				'#ffd700', // Gold
				'#98fb98', // Pale green
				'#87ceeb' // Sky blue
			],
			particleCount: Math.floor(count * particleRatio),
			zIndex: 1000
		});
	}

	async function addHandler() {
		status = 'loading';

		try {
			const result = await query();
			if (result) {
				if (celebrate) {
					triggerConfetti();
				}
				isCelebrating = true;
				addToast('success', 'Added successfully!');
				status = 'success';
				await insertNewLogToCache(result);

				setTimeout(() => {
					status = 'default';
					isCelebrating = false;
				}, 3000);
			}
		} catch (err) {
			console.log(err);
			status = 'error';
			addToast('error', 'Error creating!');
			setTimeout(() => {
				status = 'default';
			}, 3000);
		}
	}

	function getButtonClasses() {
		return [
			'btn btn-lg flex items-center gap-2',
			`rounded-${rounded}`,
			(status === 'default' || status === 'loading') && `btn-${color}`,
			status === 'success' && 'btn-success isCelebrating',
			status === 'error' && 'btn-error',
			isCelebrating && 'scale-105 shadow-lg',
			className
		];
	}
</script>

{#if !compact}
	<button
		class={['btn btn-lg flex w-full items-center gap-2', getButtonClasses()]}
		onclick={addHandler}
		bind:this={buttonRef}
	>
		{#if status === 'success'}
			<MaterialSymbolsCheck class="size-6" />Added!
		{/if}
		{#if status === 'loading'}
			<span class="loading loading-spinner loading-md"></span>
		{/if}
		{#if status === 'error'}
			Error!
		{/if}
		{#if status !== 'success' && status !== 'loading' && status !== 'error'}
			{#if text}
				{text}
			{:else}
				Submit
			{/if}
		{/if}
	</button>
{:else}
	<button
		class={['btn btn-lg flex aspect-square w-full items-center gap-2 p-0', getButtonClasses()]}
		onclick={addHandler}
		bind:this={buttonRef}
	>
		{#if status === 'loading'}
			<span class="loading loading-spinner loading-md"></span>
		{:else if CustomIcon && status !== 'success'}
			<CustomIcon class="size-8" />
		{:else}
			<MaterialSymbolsCheck class="size-8" />
		{/if}
	</button>
{/if}

<style>
	@keyframes glow {
		0%,
		100% {
			box-shadow: 0 0 5px rgba(124, 179, 66, 0.3);
		}
		50% {
			box-shadow: 0 0 20px rgba(124, 179, 66, 0.6);
		}
	}

	.isCelebrating {
		animation: glow 0.5s ease-in-out 3;
	}
</style>
