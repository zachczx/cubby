<script lang="ts">
	import Icon from '@iconify/svelte';
	import { addToast } from './ArkToaster.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import type { Component } from 'svelte';
	import { getAllEntriesQueryKey } from '$lib/queries';

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
		query: () => Promise<EntryDB>;
		compact?: boolean;
		color?: 'primary' | 'neutral';
		rounded?: 'full' | '3xl' | '2xl' | 'xl' | 'lg' | 'md' | 'sm' | 'xs';
		class?: string;
		icon?: Component | string;
		celebrate?: boolean;
	} = $props();

	let status: ButtonState = $state('default');
	let isCelebrating = $state(false);
	let showStars = $state(false);
	let starPositions: { x: number; y: number }[] = $state([]);

	const tanstackClient = useQueryClient();
	export const insertNewEntryToCache = (newEntry: EntryDB) =>
		tanstackClient.setQueryData(getAllEntriesQueryKey(), (oldEntries: EntryDB[] | undefined) => {
			if (!oldEntries) return [newEntry];
			return [newEntry, ...oldEntries];
		});

	const starColors = ['#ff7f50', '#ffb347', '#ffd700', '#98fb98', '#87ceeb'];

	function generateStarPositions(count: number): { x: number; y: number }[] {
		const centerX = 50;
		const centerY = 50;
		const starCount = count;

		// Calculate size of one slice
		const slice = 360 / starCount;

		// Global rotation to make it not fixed where 1st always starts at 3pm position.
		const globalOffset = Math.random() * slice;

		return Array.from({ length: starCount }, (_, i) => {
			let angleDeg = i * slice + globalOffset;

			const maxWobble = slice * 0.5;
			const randomShift = (Math.random() - 0.5) * maxWobble;

			angleDeg += randomShift;

			const angleRad = angleDeg * (Math.PI / 180);

			const distance = 45 + Math.random() * 20;

			return {
				x: centerX + Math.cos(angleRad) * distance,
				y: centerY + Math.sin(angleRad) * distance
			};
		});
	}

	async function addHandler() {
		status = 'loading';

		try {
			const result = await query();
			if (result) {
				isCelebrating = true;
				if (celebrate) {
					starPositions = generateStarPositions(compact ? 6 : 10);
					showStars = true;

					setTimeout(() => {
						showStars = false;
					}, 2000);
				}

				addToast('success', 'Added successfully!');
				status = 'success';
				await insertNewEntryToCache(result);

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

<div class="relative w-full" style="--star-size: {compact ? '12px' : '20px'}">
	{#if showStars}
		{#each starPositions as pos, i}
			<div
				class="star-pop"
				style="
					left: {pos.x}%;
					top: {pos.y}%;
					background: {starColors[i % starColors.length]};
					animation-delay: {i * 0.08}s;
				"
			></div>
		{/each}
	{/if}

	{#if !compact}
		<button
			class={['btn btn-lg flex w-full items-center gap-2', getButtonClasses()]}
			onclick={addHandler}
		>
			{#if status === 'success'}
				<Icon icon="material-symbols:check" class="size-6" />Added!
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
		>
			{#if status === 'loading'}
				<span class="loading loading-spinner loading-md"></span>
			{:else if typeof CustomIcon === 'string'}
				<Icon icon={CustomIcon} class="size-8" />
			{:else if CustomIcon && status !== 'success'}
				<CustomIcon class="size-8" />
			{:else}
				<Icon icon="material-symbols:check" class="size-8" />
			{/if}
		</button>
	{/if}
</div>

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

	.star-pop {
		position: absolute;
		width: var(--star-size);
		height: var(--star-size);
		opacity: 0;
		animation: appear-stay-fade 1.5s ease-out forwards;
		z-index: 20;
		clip-path: polygon(
			50% 0%,
			61% 35%,
			98% 35%,
			68% 57%,
			79% 91%,
			50% 70%,
			21% 91%,
			32% 57%,
			2% 35%,
			39% 35%
		);
	}

	@keyframes appear-stay-fade {
		0% {
			transform: translate(-50%, -50%) scale(0) rotate(0deg);
			opacity: 0;
		}
		20% {
			transform: translate(-50%, -50%) scale(1) rotate(45deg);
			opacity: 1;
		}
		80% {
			transform: translate(-50%, -50%) scale(1) rotate(90deg);
			opacity: 1;
		}
		100% {
			transform: translate(-50%, -50%) scale(1) rotate(90deg);
			opacity: 0;
		}
	}
</style>
