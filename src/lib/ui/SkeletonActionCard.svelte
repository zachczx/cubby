<script lang="ts">
	import Icon from '@iconify/svelte';

	let { size = 'default' }: { size?: 'compact' | 'default' } = $props();
</script>

{#if size === 'compact'}
	<section
		class="skeleton border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 px-2 py-2"
	>
		<div class="flex items-center">
			<div class="flex grow items-center gap-4">
				<div class="skeleton aspect-square size-9 rounded-full opacity-75"></div>
				<div>
					<div class="mb-2 text-xl font-bold">
						<div class="skeleton h-5 w-26"></div>
					</div>
					<div class="text-neutral/70 font-medium tracking-tight">
						<div class="skeleton h-3 w-30"></div>
					</div>
				</div>
			</div>

			<div class="flex h-full items-center">
				<div class="skeleton aspect-square size-12 rounded-2xl"></div>
			</div>
		</div>
	</section>
{:else}
	<section class="border-base-300 skeleton grid min-h-24 gap-4 rounded-3xl border bg-white/70 p-4">
		<div class="flex items-center">
			<div class="flex grow items-center gap-4">
				<div class="skeleton aspect-square size-12 rounded-full opacity-75"></div>
				<div>
					<div class="mb-2 text-xl font-bold">
						<div class="skeleton h-5 w-26"></div>
					</div>
					<div class="text-neutral/70 font-medium tracking-tight">
						<div class="skeleton h-3 w-30"></div>
					</div>
				</div>
			</div>
			<div class="flex h-full items-center">
				<div class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75">
					<Icon icon="material-symbols:chevron-right" class="size-6 opacity-50" />
				</div>
			</div>
		</div>

		<div class="skeleton aspect-square size-12 w-full rounded-full"></div>
	</section>
{/if}
