<script lang="ts">
	import Icon from '@iconify/svelte';

	let { notification }: { notification: NotificationStatus } = $props();

	interface StatusLabels {
		ok: string;
		due: string;
		overdue: string;
	}

	const statusLabels: StatusLabels = {
		ok: 'Done',
		due: 'Due',
		overdue: 'Overdue'
	};
</script>

<div class="flex items-center gap-2">
	{#if notification.level === 'ok'}
		<Icon icon="material-symbols:check-circle" class="size-8 text-lime-600" />
		<span class="text-lime-600">{statusLabels.ok}</span>
	{:else if notification.level === 'due'}
		<Icon icon="ant-design:exclamation-circle-filled" class="text-warning size-8" />

		<span class="text-orange-400">{statusLabels.due}</span>
	{:else if notification.level === 'overdue'}
		<Icon icon="ant-design:exclamation-circle-filled" class="text-error size-8" />
		<span class="text-red-700">{statusLabels.overdue}</span>
	{:else}
		<span>Nil</span>
	{/if}
</div>
