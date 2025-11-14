<script lang="ts">
	import MaterialSymbolsChevronRight from '$lib/assets/svg/MaterialSymbolsChevronRight.svelte';
	import type { RecordModel } from 'pocketbase';
	import type { Component } from 'svelte';
	import ActionButton from './ActionButton.svelte';
	import MaterialSymbolsCheckCircle from '$lib/assets/svg/MaterialSymbolsCheckCircle.svelte';
	import MaterialSymbolsAdd from '$lib/assets/svg/MaterialSymbolsAdd.svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';

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
		'border-base-300 grid min-h-18 gap-4 rounded-3xl border px-4 py-2',
		options.notification.show ? 'bg-error/15 outline-error/30 outline' : 'bg-base-100'
	]}
>
	<div class="flex items-center">
		<a href={options.route} class="flex grow items-center gap-4">
			<options.icon class="size-9 opacity-75" />
			<div>
				<p class="flex items-center gap-2 text-xl font-bold">
					{options.title}
					<MaterialSymbolsArrowRightAlt class="size-[1em]" />
				</p>
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
		</a>
		<div class="flex h-full items-center">
			<ActionButton
				query={options.button.query}
				refetch={options.button.refetch}
				text={options.button.text}
				compact={true}
				color={options.notification.level === 'overdue' || options.notification.level === 'due'
					? 'primary'
					: 'neutral'}
				icon={MaterialSymbolsAdd}
			/>
		</div>
	</div>
</section>
