<script lang="ts">
	let {
		count = 12,
		size = 'md',
		class: className
	}: {
		count?: number;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		class?: string;
	} = $props();

	const sizeRanges = {
		sm: { base: 4, range: 5 },
		md: { base: 6, range: 8 },
		lg: { base: 10, range: 12 },
		xl: { base: 14, range: 16 }
	};

	const colors = ['#ff7f50', '#ffb347', '#ffd700', '#98fb98', '#87ceeb'];

	const { base, range } = $derived(sizeRanges[size]);
	const minGap = $derived(base * 1.5);

	function generateStars(count: number) {
		const placed: { x: number; y: number; size: number }[] = [];
		const stars: {
			x: number;
			y: number;
			size: number;
			rotation: number;
			color: string;
			delay: number;
			duration: number;
		}[] = [];

		for (let i = 0; i < count; i++) {
			let x: number, y: number, starSize: number;
			let attempts = 0;

			do {
				x = Math.random() * 90 + 5;
				y = Math.random() * 80 + 10;
				starSize = base + Math.random() * range;
				attempts++;
			} while (attempts < 50 && placed.some((p) => Math.hypot(p.x - x, p.y - y) < minGap));

			placed.push({ x, y, size: starSize });
			stars.push({
				x,
				y,
				size: starSize,
				rotation: Math.random() * 360,
				color: colors[i % colors.length],
				delay: Math.random() * 3,
				duration: 2 + Math.random() * 2
			});
		}

		return stars;
	}

	const stars = $derived(generateStars(count));
</script>

<div class={['pointer-events-none absolute inset-0 z-0 overflow-hidden', className]}>
	{#each stars as star}
		<div
			class="star-twinkle absolute"
			style="
				left: {star.x}%;
				top: {star.y}%;
				width: {star.size}px;
				height: {star.size}px;
				background: {star.color};
				rotate: {star.rotation}deg;
				animation-delay: {star.delay}s;
				animation-duration: {star.duration}s;
			"
		></div>
	{/each}
</div>

<style>
	.star-twinkle {
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
		opacity: 0.4;
		animation: twinkle ease-in-out infinite;
	}

	@keyframes twinkle {
		0%,
		100% {
			opacity: 0.15;
			scale: 0.85;
		}
		50% {
			opacity: 0.5;
			scale: 1;
		}
	}
</style>
