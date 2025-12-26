<script lang="ts">
	import hero from '$lib/assets/hero.webp?w=250&enhanced';
	import heroSubscription from '$lib/assets/hero-subscription.webp?w=250&enhanced';
	import AntDesignExclamationCircleFilled from '$lib/assets/svg/AntDesignExclamationCircleFilled.svelte';
	import EmptyState from '$lib/assets/svg/EmptyState.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	let { notification, kind = 'task' }: { notification: NotificationStatus; kind?: Kind } = $props();
</script>

{#if notification}
	{#if kind === 'task'}
		<div class="avatar relative mt-2 mb-4">
			<div class="w-40 rounded-full shadow-md">
				<enhanced:img
					src={hero}
					alt="hero"
					class={[
						notification.level === 'due' && 'saturate-75',
						notification.level === 'overdue' && 'saturate-50'
					]}
					fetchpriority="high"
				/>
			</div>
			<div
				class="bg-base-100 absolute top-2 -right-4 flex size-14 items-center justify-center rounded-full"
			>
				{#if notification.level === 'ok'}
					<MaterialSymbolsCheckCircle class="text-success size-13" />
				{:else if notification.level === 'due'}
					<AntDesignExclamationCircleFilled class="text-warning size-13" />
				{:else if notification.level === 'overdue'}
					<AntDesignExclamationCircleFilled class="text-error size-13" />
				{/if}
			</div>
		</div>
	{:else}
		<div class="avatar relative mt-2 mb-4">
			<div class="w-40 rounded-full shadow-md">
				<enhanced:img src={heroSubscription} alt="hero" class="scale-110" fetchpriority="high" />
			</div>
		</div>
	{/if}
{:else}
	<EmptyState class="text-primary/30 my-4" />
{/if}
