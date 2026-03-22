<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		items,
		class: className
	}: { children: Snippet; items: number; class?: string } = $props();
</script>

<form
	class={['segmented-control', className]}
	style="grid-template-columns: repeat({items}, minmax(0, 1fr))"
>
	{@render children()}
</form>

<style>
	:global {
		.segmented-control {
			display: grid;
			cursor: pointer;
			text-align: center;

			background: rgba(0, 0, 0, 0.02);
			padding: 0.25rem;
			margin-bottom: 0.25rem;
			border-radius: calc(infinity * 1px);

			label {
				width: 100%;
				font-size: 1.125rem;
				padding: 2rem;
				cursor: pointer;
				height: var(--size);
				position: relative;
				display: inline-flex;
				align-items: center;
				justify-content: center;
				outline-color: #d6cecb;
				--size: calc(var(--size-field, 0.25rem) * 10);

				&:first-child {
					border-start-start-radius: var(--join-ss, var(--radius-field) /* var(--radius-field) */);
					border-end-start-radius: var(--join-es, var(--radius-field) /* var(--radius-field) */);
				}

				&:last-child {
					border-start-end-radius: var(--join-se, var(--radius-field) /* var(--radius-field) */);
					border-end-end-radius: var(--join-ee, var(--radius-field) /* var(--radius-field) */);
				}

				&:not(:first-child) {
					border-left: 0;
				}

				&:has(input:checked) {
					color: var(--color-primary);
					background-color: var(--color-bg-checked);
					font-weight: 700;
					border-radius: calc(infinity * 1px);
					--tw-shadow:
						0 1px 3px 0 var(--tw-shadow-color, rgb(0 0 0 / 0.1)),
						0 1px 2px -1px var(--tw-shadow-color, rgb(0 0 0 / 0.1));
					box-shadow:
						var(--tw-inset-shadow), var(--tw-inset-ring-shadow), var(--tw-ring-offset-shadow),
						var(--tw-ring-shadow), var(--tw-shadow);
				}
				&:has(input:not(checked)) {
					color: rgba(45, 37, 32, 0.7);
				}

				input {
					display: none;
				}
			}
		}
	}
</style>
