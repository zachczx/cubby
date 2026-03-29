<script lang="ts">
	interface Props {
		points: { date: string; price: number }[];
		width?: number;
		height?: number;
		strokeColor?: string;
		fillColor?: string;
	}

	let {
		points,
		width = 280,
		height = 64,
		strokeColor = 'var(--color-primary)',
		fillColor = 'var(--color-primary)'
	}: Props = $props();

	const padX = 4;
	const padY = 8;

	let stats = $derived.by(() => {
		if (points.length < 2) return null;
		const prices = points.map((p) => p.price);
		const min = Math.min(...prices);
		const max = Math.max(...prices);
		const range = max - min || 1;
		const innerW = width - padX * 2;
		const innerH = height - padY * 2;
		return { min, max, range, innerW, innerH };
	});

	let linePath = $derived.by(() => {
		if (!stats) return '';
		const { min, range, innerW, innerH } = stats;
		return points
			.map((p, i) => {
				const x = padX + (i / (points.length - 1)) * innerW;
				const y = padY + innerH - ((p.price - min) / range) * innerH;
				return `${i === 0 ? 'M' : 'L'}${x.toFixed(1)},${y.toFixed(1)}`;
			})
			.join(' ');
	});

	let areaPath = $derived.by(() => {
		if (!linePath || !stats) return '';
		const firstX = padX;
		const lastX = padX + stats.innerW;
		return `${linePath} L${lastX.toFixed(1)},${height} L${firstX.toFixed(1)},${height} Z`;
	});

	let latestDot = $derived.by(() => {
		if (!stats) return null;
		const { min, range, innerW, innerH } = stats;
		const last = points[points.length - 1];
		const x = padX + innerW;
		const y = padY + innerH - ((last.price - min) / range) * innerH;
		return { x, y };
	});

	let minMaxLabels = $derived.by(() => {
		if (!stats) return null;
		const { min, max, innerW, innerH } = stats;
		if (min === max) return null;

		let minX = padX,
			minY = padY + innerH;
		let maxX = padX,
			maxY = padY + innerH;

		for (let i = 0; i < points.length; i++) {
			const x = padX + (i / (points.length - 1)) * innerW;
			const y = padY + innerH - ((points[i].price - min) / (max - min)) * innerH;
			if (points[i].price === min) {
				minX = x;
				minY = y;
			}
			if (points[i].price === max) {
				maxX = x;
				maxY = y;
			}
		}

		return {
			min: { x: minX, y: minY, label: `$${min.toFixed(2)}` },
			max: { x: maxX, y: maxY, label: `$${max.toFixed(2)}` }
		};
	});
</script>

{#if points.length >= 2}
	<svg viewBox="0 0 {width} {height}" class="w-full" style="height: {height}px;">
		<defs>
			<linearGradient id="sparkFill" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color={fillColor} stop-opacity="0.15" />
				<stop offset="100%" stop-color={fillColor} stop-opacity="0.02" />
			</linearGradient>
		</defs>

		<path d={areaPath} fill="url(#sparkFill)" />

		<path
			d={linePath}
			fill="none"
			stroke={strokeColor}
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>

		{#if latestDot}
			<circle cx={latestDot.x} cy={latestDot.y} r="4" fill={strokeColor} />
			<circle cx={latestDot.x} cy={latestDot.y} r="7" fill={strokeColor} opacity="0.15" />
		{/if}

		{#if minMaxLabels}
			<text
				x={minMaxLabels.min.x}
				y={minMaxLabels.min.y + 12}
				text-anchor="middle"
				class="fill-base-content/40 text-[10px]">{minMaxLabels.min.label}</text
			>
			<text
				x={minMaxLabels.max.x}
				y={minMaxLabels.max.y - 5}
				text-anchor="middle"
				class="fill-base-content/40 text-[10px]">{minMaxLabels.max.label}</text
			>
		{/if}
	</svg>
{/if}
