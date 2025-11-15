<script lang="ts">
	import AntDesignExclamationCircleFilled from '$lib/assets/svg/AntDesignExclamationCircleFilled.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';

	let { notification }: { notification: NotificationStatus } = $props();

	interface Statuses {
		ok: string;
		due: string;
		overdue: string;
	}

	const statusLabels: Statuses = {
		ok: 'Done',
		due: 'Due',
		overdue: 'Overdue'
	};
</script>

<div class="flex items-center gap-2">
	{#if notification.level === 'ok'}
		<MaterialSymbolsCheckCircle class="size-8 text-lime-600" />
		<span class="text-lime-600">{statusLabels.ok}</span>
	{:else if notification.level === 'due'}
		<AntDesignExclamationCircleFilled class="text-warning size-13" />

		<span class="text-orange-400">{statusLabels.due}</span>
	{:else if notification.level === 'overdue'}
		<AntDesignExclamationCircleFilled class="text-error size-13" />
		<span class="text-red-700">{statusLabels.overdue}</span>
	{:else}
		<span>Nil</span>
	{/if}
</div>
