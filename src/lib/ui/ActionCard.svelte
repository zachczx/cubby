<script lang="ts">
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import type { RecordModel } from 'pocketbase';
	import type { Component } from 'svelte';
	import ActionButton from './ActionButton.svelte';

	let { options }: { options: Options } = $props();

	interface Options {
		title: string;
		query: Query;
		notification: NotificationStatus;
		icon: Component;
		route: string;
		last: string;
		button: {
			query: () => Promise<RecordModel>;
			refetch: () => Promise<void>;
			status: ButtonState;
			text: string;
		};
	}
</script>

<section
	class={[
		'border-base-300 grid min-h-24 gap-4 rounded-3xl border p-4',
		options.notification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-base-100'
	]}
>
	<a href={options.route} class="flex items-center">
		<div class="flex grow items-center gap-4">
			<options.icon class="size-12 opacity-75" />
			<div>
				<p class="text-xl font-bold">{options.title}</p>
				{#if options.query.isPending && !options.query.data}
					<div class="custom-loader"></div>
				{/if}
				{#if options.query.error}
					An error has occurred:
					{options.query.error.message}
				{/if}
				{#if options.query.isSuccess}
					{#if options.notification.show}
						<span class="text-error font-medium tracking-tight">
							{#if options.notification.level === 'overdue'}
								Overdue
							{:else if options.notification.level === 'due'}
								Due
							{/if}
						</span>&nbsp;&nbsp;•&nbsp;&nbsp;
					{/if}<span>{options.last}</span>
				{/if}
			</div>
		</div>
		<div class="flex h-full items-center">
			<button class="active:bg-neutral/10 cursor-pointer rounded-lg p-1 opacity-75"
				><MaterialSymbolsChevronRight class="size-6" /></button
			>
		</div>
	</a>
	<ActionButton
		query={options.button.query}
		refetch={options.button.refetch}
		text={options.button.text}
	/>
</section>
