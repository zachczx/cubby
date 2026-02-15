<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Icon from '@iconify/svelte';
	// import { inviteQueryOptions, inviteRefetchOptions } from '$lib/queries';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	const destinationFamilyId = page.url.searchParams.get('i') ?? '';
	// const currentInvite = createQuery(inviteQueryOptions);
	const tanstackClient = useQueryClient();

	async function confirmJoinFamily() {
		try {
			/* await pb.collection('families').update(destinationFamilyId, {
				'members+': pb.authStore.record.id,
				'activeInvites-': pb.authStore.record.id
			});

			if (currentInvite.data) {
				const cleanup = await pb.collection('invites').update(currentInvite.data.id, {
					status: 'completed'
				});

				console.log(cleanup);
			}

			// TODO: This part is slow
			await tanstackClient.refetchQueries(inviteRefetchOptions());

			addToast('success', 'Joined family!');
			goto(resolve('/app/profile/family')); */
		} catch (err) {
			console.log(err);
		}
	}
</script>

<PageWrapper title="Join Family" largeScreenCenter={true}>
	<span></span>
</PageWrapper>

<dialog class="modal modal-bottom sm:modal-middle modal-open">
	<div class="modal-box grid gap-8">
		<!-- {#if currentInvite.isSuccess && destinationFamilyId}<div class="grid gap-4">
				<div
					class="bg-primary/10 text-primary flex aspect-square size-20 items-center justify-center justify-self-center overflow-hidden rounded-full"
				>
					<Icon icon="material-symbols:person-add" class="ms-2.5 size-12" />
				</div>
				<h2 class="text-2xl font-bold">
					You've been invited to {currentInvite.data?.ownerEmailSnapshot}'s family!
				</h2>
				<ul class="ms-6 list-disc space-y-2">
					<li>You can see and use the family's trackers.</li>
					<li>You can leave the family any time.</li>
					<li>Your personal trackers will remain private unless you share them.</li>
				</ul>
			</div>
			<div class="grid grid-cols-1 gap-4">
				<button
					class="btn btn-primary btn-lg"
					onclick={() => {
						confirmJoinFamily();
					}}>Accept Invite</button
				>
				<form method="dialog" class="">
					<button class="btn btn-outline btn-primary btn-lg w-full">Cancel</button>
				</form>
			</div>
		{:else}
			<div class="text-error grid content-center justify-items-center gap-2 text-lg font-bold">
				<Icon icon="material-symbols:warning" class="size-12" />
				{#if !destinationFamilyId}
					You're not joining a valid family!
				{:else}
					Something went wrong!
				{/if}
			</div>
		{/if} -->
	</div>
</dialog>
