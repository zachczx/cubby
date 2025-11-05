<script lang="ts" module>
	import { createToaster } from '@ark-ui/svelte/toast';

	const toaster = createToaster({
		placement: 'top',
		overlap: true,
		gap: 24,
		duration: 7000
	});

	export const addToast = (type: string = 'success', description: string) => {
		toaster.create({
			title: 'title',
			description: description,
			type: type
		});
	};
</script>

<script lang="ts">
	import { Toast, Toaster } from '@ark-ui/svelte/toast';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsClose from '$lib/assets/svg/MaterialSymbolsClose.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import { SafeArea } from '@capacitor-community/safe-area';
</script>

<div>
	<Toaster {toaster}>
		{#snippet children(toast)}
			<Toast.Root>
				<div class="content">
					<div><MaterialSymbolsCheckCircle class="size-[1.3em]" /></div>
					<div>
						<!-- {#if toast().title}
						<Toast.Title>{toast().title}</Toast.Title>
					{/if} -->
						<Toast.Description>{toast().description}</Toast.Description>
					</div>
				</div>
				<div class="flex items-center">
					<Toast.CloseTrigger>
						<MaterialSymbolsClose class="size-[1em] cursor-pointer" />
					</Toast.CloseTrigger>
				</div>
			</Toast.Root>
		{/snippet}
	</Toaster>
</div>
<span class="placeholder"></span>

<style>
	.placeholder {
		display: none;
	}

	:global {
		[data-scope='toast'][data-part='root'] {
			.content {
				display: flex;
				align-items: center;
				column-gap: 0.75rem;
			}

			display: flex;
			align-items: center;
			column-gap: 2rem;
			padding: 1rem;
			margin-top: calc(var(--safe-area-inset-top) + 1rem) !important;
			border-radius: 0.5rem;
			box-shadow:
				0 4px 6px -1px rgb(0 0 0 / 0.3),
				0 2px 4px -2px rgb(0 0 0 / 0.3);

			translate: var(--x) var(--y);
			scale: var(--scale);
			z-index: var(--z-index);
			height: var(--height);
			opacity: var(--opacity);
			will-change: translate, opacity, scale;

			transition:
				translate 400ms,
				scale 400ms,
				opacity 400ms,
				height 400ms,
				box-shadow 200ms;
			transition-timing-function: cubic-bezier(0.21, 1.02, 0.73, 1);

			&[data-state='closed'] {
				transition:
					translate 400ms,
					scale 400ms,
					opacity 200ms;
				transition-timing-function: cubic-bezier(0.06, 0.71, 0.55, 1);
			}

			&[data-type='error'] {
				background: var(--color-error);
				color: white;
			}

			&[data-type='info'] {
				background: var(--color-secondary);
				color: white;
			}

			&[data-type='warning'] {
				background: orange;
			}

			&[data-type='success'] {
				background: var(--color-success);
				color: white;
			}
		}
		/* [data-scope='toast'][data-part='close-trigger'] {
			cursor: pointer;
			&:hover {
				border-radius: 50%;
				background-color: --var(color-success-content);
				color: --var(color-success);
			}
		} */
	}
</style>
