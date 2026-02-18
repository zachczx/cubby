<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Icon from '@iconify/svelte';
	import { singleInviteQueryOptions, singleInviteRefetchOptions } from '$lib/queries';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api } from '$lib/api';

	let { data } = $props();

	const invite = $derived(createQuery(singleInviteQueryOptions(data.inviteId)));
	const tanstackClient = useQueryClient();

	async function confirmJoinFamily() {
		try {
			const response = await api.post(`families/invites/${data.inviteId}/accept`);

			await tanstackClient.refetchQueries(singleInviteRefetchOptions(data.inviteId));

			if (response.status === 204) {
				addToast('success', 'Joined family!');
				goto(resolve('/app/profile/family'));
			}
		} catch (err) {
			addToast('error', 'Error joining family!');
		}
	}

	async function declineJoinFamily() {
		try {
			const response = await api.post(`families/invites/${data.inviteId}/decline`);

			await tanstackClient.refetchQueries(singleInviteRefetchOptions(data.inviteId));

			if (response.status === 204) {
				goto(resolve('/app/profile/family'));
			}
		} catch (err) {
			addToast('error', 'Error declining invite!');
		}
	}
</script>

<PageWrapper title="Join Family" largeScreenCenter={true}>
	<span></span>
</PageWrapper>

<dialog class="modal modal-bottom sm:modal-middle modal-open">
	<div class="modal-box grid gap-8">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<div class="grid gap-4">
			<div
				class="bg-primary/10 text-primary flex aspect-square size-20 items-center justify-center justify-self-center overflow-hidden rounded-full"
			>
				<Icon icon="material-symbols:person-add" class="ms-2.5 size-12" />
			</div>
			<h2 class="text-2xl font-bold">You've been invited!</h2>

			<ul class="ms-6 list-disc space-y-2">
				{#if invite.isPending}
					<li><div class="skeleton h-5 w-68 rounded-full"></div></li>
					<li><div class="skeleton h-5 w-52 rounded-full"></div></li>
					<li>
						<div class="skeleton mb-2 h-5 w-full rounded-full"></div>
						<div class="skeleton h-5 w-26 rounded-full"></div>
					</li>
				{:else if invite.isSuccess}
					<li>
						You can see and use <span class="font-bold">{invite.data.familyName}</span>'s trackers.
					</li>
					<li>You can leave <span class="font-bold">{invite.data.familyName}</span> any time.</li>
					<li>Your personal trackers will remain private unless you share them.</li>
				{/if}
			</ul>
		</div>
		<div class="grid grid-cols-1 gap-4">
			<button
				class="btn btn-primary btn-lg"
				onclick={() => {
					confirmJoinFamily();
				}}>Accept Invite</button
			>
			<button
				class="btn btn-outline btn-neutral btn-lg w-full"
				onclick={() => {
					declineJoinFamily();
				}}>Decline</button
			>
		</div>
	</div>
</dialog>
