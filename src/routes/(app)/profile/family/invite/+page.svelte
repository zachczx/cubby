<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import MaterialSymbolsWarning from '$lib/assets/svg/MaterialSymbolsWarning.svelte';
	import { pb } from '$lib/pb';
	import {
		allLogsQueryOptions,
		familyQueryOptions,
		familyRefetchOptions,
		inviteQueryOptions,
		inviteRefetchOptions
	} from '$lib/queries';
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	const destinationFamilyId = page.url.searchParams.get('i') ?? '';
	const userAllLogs = createQuery(allLogsQueryOptions);
	const currentFamily = createQuery(familyQueryOptions);
	const currentInvite = createQuery(inviteQueryOptions);
	const tanstackClient = useQueryClient();

	async function confirmJoinFamily() {
		if (!pb.authStore.record) return;

		try {
			if (currentFamily.data?.owner === pb.authStore.record?.id) {
				await pb.collection('families').delete(currentFamily.data.id);
			}

			if (
				currentFamily.data &&
				currentFamily.data.owner !== pb.authStore.record.id &&
				currentFamily.data.members.includes(pb.authStore.record.id ?? '')
			) {
				await pb
					.collection('families')
					.update(currentFamily.data.id, { 'members-': pb.authStore.record.id });
			}

			console.log(destinationFamilyId);
			const resp = await pb.collection('families').update(destinationFamilyId, {
				'members+': pb.authStore.record.id,
				'activeInvites-': pb.authStore.record.id
			});

			await pb.collection('users').update(pb.authStore.record.id, {
				family: destinationFamilyId
			});

			if (currentInvite.data) {
				const cleanup = await pb.collection('invites').update(currentInvite.data.id, {
					status: 'completed'
				});

				console.log(cleanup);
			}

			await tanstackClient.refetchQueries(familyRefetchOptions());
			await tanstackClient.refetchQueries(inviteRefetchOptions());

			addToast('success', 'Joined family!');
			goto('/profile/family');
		} catch (err) {
			console.log(err);
		}
	}

	let modal = $state<HTMLDialogElement>();

	let currentUserTrackerCounts = $derived.by(() => {
		if (!userAllLogs.isSuccess) return;

		return userAllLogs.data.length;
	});
</script>

<PageWrapper title="Join Family" {pb} largeScreenCenter={true}>
	<span></span>
</PageWrapper>

<dialog bind:this={modal} class="modal modal-bottom sm:modal-middle modal-open">
	<div class="modal-box grid gap-8">
		{#if currentInvite.isSuccess && destinationFamilyId}<div class="grid gap-4">
				<h2 class="text-2xl font-bold">
					You've been invited to "{currentInvite.data?.familyNameSnapshot}" Family!
				</h2>
				<ul class="ms-6 list-disc space-y-2">
					{#if currentUserTrackerCounts && currentUserTrackerCounts > 0}
						<li>All of your existing log data will be deleted.</li>
						<li>
							You have <span class="text-error font-bold">{currentUserTrackerCounts} records</span> to
							be deleted.
						</li>
						<li>This action cannot be undone.</li>
						<li>Future activity will be saved to the new family.</li>
					{:else}
						<li>Future activity will be saved to the new family.</li>
						<li>Your account is empty, so no data will be lost.</li>
					{/if}
				</ul>
			</div>
			<div class="grid grid-cols-1 gap-4">
				<button class="btn btn-primary btn-lg" onclick={() => confirmJoinFamily()}
					>Join New Family</button
				>
				<form method="dialog" class="">
					<button class="btn btn-outline btn-primary btn-lg w-full">Cancel</button>
				</form>
			</div>
		{:else}
			<div class="text-error grid content-center justify-items-center gap-2 text-lg font-bold">
				<MaterialSymbolsWarning class="size-12" />
				{#if !destinationFamilyId}
					You're not joining a valid family!
				{:else}
					Something went wrong!
				{/if}
			</div>
		{/if}
	</div>
</dialog>
