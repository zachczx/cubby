<script lang="ts">
	import hero from '$lib/assets/hero.webp?w=600&enhanced';
	import AntDesignExclamationCircleFilled from '$lib/assets/svg/AntDesignExclamationCircleFilled.svelte';
	import EmptyState from '$lib/assets/svg/EmptyState.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	let { status }: { status: string } = $props();
</script>

{#if status !== 'empty'}
	<div class="avatar relative mb-4">
		<div class="w-40 rounded-full shadow-md">
			<enhanced:img
				src={hero}
				alt="hero"
				class={[status === 'orange' && 'saturate-75', status === 'red' && 'saturate-50']}
			/>
		</div>
		<div
			class="bg-base-100 absolute top-2 -right-4 flex size-14 items-center justify-center rounded-full"
		>
			{#if status === 'green' || status === 'yellow'}
				<MaterialSymbolsCheckCircle class="text-success size-13" />
			{:else if status === 'orange'}
				<AntDesignExclamationCircleFilled class="text-warning size-13" />
			{:else if status === 'red'}
				<AntDesignExclamationCircleFilled class="text-error size-13" />
			{/if}
		</div>
	</div>
{:else}
	<EmptyState class="text-primary/30 my-4" />
{/if}
