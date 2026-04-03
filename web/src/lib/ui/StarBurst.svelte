<script lang="ts">
	let {
		show = false,
		count = 10,
		size = '20px'
	}: {
		show: boolean;
		count?: number;
		size?: string;
	} = $props();

	const colors = ['#ff7f50', '#ffb347', '#ffd700', '#98fb98', '#87ceeb'];

	function generatePositions(count: number): { x: number; y: number }[] {
		const centerX = 50;
		const centerY = 50;
		const slice = 360 / count;
		const globalOffset = Math.random() * slice;

		return Array.from({ length: count }, (_, i) => {
			let angleDeg = i * slice + globalOffset;
			const maxWobble = slice * 0.5;
			angleDeg += (Math.random() - 0.5) * maxWobble;
			const angleRad = angleDeg * (Math.PI / 180);
			const distance = 45 + Math.random() * 20;

			return {
				x: centerX + Math.cos(angleRad) * distance,
				y: centerY + Math.sin(angleRad) * distance
			};
		});
	}

	let positions = $derived(show ? generatePositions(count) : []);
</script>

{#if show}
	<div class="star-container" style="--star-size: {size}">
		{#each positions as pos, i}
			<div
				class="star-pop"
				style="
					left: {pos.x}%;
					top: {pos.y}%;
					background: {colors[i % colors.length]};
					animation-delay: {i * 0.08}s;
				"
			></div>
		{/each}
	</div>
{/if}

<style>
	.star-container {
		position: absolute;
		inset: 0;
		pointer-events: none;
		z-index: 20;
	}

	.star-pop {
		position: absolute;
		width: var(--star-size);
		height: var(--star-size);
		opacity: 0;
		animation: appear-stay-fade 1.5s ease-out forwards;
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
