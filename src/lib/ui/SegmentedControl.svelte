<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children, items }: { children: Snippet; items: number } = $props();
</script>

<div class="segmented-control grid-cols-{items}">
	{@render children()}
</div>

<style>
	:global {
		.segmented-control {
			display: grid;
			/* grid-template-columns: repeat(2, minmax(0, 1fr)); */
			cursor: pointer;
			text-align: center;

			label {
				width: 100%;
				font-size: 1.125rem;

				&:first-child {
					border-start-start-radius: var(--join-ss, var(--radius-field) /* var(--radius-field) */);
					border-end-start-radius: var(--join-es, var(--radius-field) /* var(--radius-field) */);
				}

				&:last-child {
					border-start-end-radius: var(--join-se, var(--radius-field) /* var(--radius-field) */);
					border-end-end-radius: var(--join-ee, var(--radius-field) /* var(--radius-field) */);
				}
				border: 1px solid rgba();
				height: var(--size);
				outline: var(--border) solid #0000;
				position: relative;
				display: inline-flex;
				align-items: center;
				justify-content: center;

				outline-color: #d6cecb;
				/* var(--input-color);
				--input-color: color-mix(
					in oklab,
					var(--color-base-content) 40%,
					#0000
				); */
				--size: calc(var(--size-field, 0.25rem) * 10);

				&:has(input:checked) {
					color: var(--color-neutral-content);
					background-color: var(--color-neutral);
					font-weight: 600;
				}
				&:has(input:not(checked)) {
					color: var(--color-base-content);
					background-color: var(--color-base-100);
				}

				input {
					display: none;
				}
			}
		}
	}
</style>
