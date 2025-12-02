<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		cleanEmail,
		createUserQueryOptions,
		createUserRefetchOptions,
		createVacationQueryOptions,
		createVacationRefetchOptions,
		familyQueryOptions,
		familyRefetchOptions,
		inviteQueryOptions
	} from '$lib/queries';
	import { page } from '$app/state';
	import SegmentedControl from '$lib/ui/SegmentedControl.svelte';
	import type { ChangeEventHandler } from 'svelte/elements';
	import MaterialSymbolsDelete from '$lib/assets/svg/MaterialSymbolsDelete.svelte';
	import { invalidateAll } from '$app/navigation';
	import MaterialSymbolsPerson from '$lib/assets/svg/MaterialSymbolsPerson.svelte';
	import MaterialSymbolsCheck from '$lib/assets/svg/MaterialSymbolsCheck.svelte';
	import MaterialSymbolsFamilyGroup from '$lib/assets/svg/MaterialSymbolsFamilyGroup.svelte';
	import MaterialSymbolsArrowRightAlt from '$lib/assets/svg/MaterialSymbolsArrowRightAlt.svelte';
	import MaterialSymbolsShare from '$lib/assets/svg/MaterialSymbolsShare.svelte';
	import MaterialSymbolsExclamation from '$lib/assets/svg/MaterialSymbolsExclamation.svelte';
	import MdiAlertCircle from '$lib/assets/svg/MdiAlertCircle.svelte';
	import MaterialSymbolsMoreVert from '$lib/assets/svg/MaterialSymbolsMoreVert.svelte';
	import MaterialSymbolsPersonRemove from '$lib/assets/svg/MaterialSymbolsPersonRemove.svelte';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(createUserQueryOptions);
	const family = createQuery(familyQueryOptions);
	const currentInvite = createQuery(inviteQueryOptions);
	const tanstackClient = useQueryClient();

	let spinner = $state(false);
	let invited = $state(false);

	let invitee = $state('');

	let modals = $state<HTMLDialogElement[]>([]);

	async function addHandler() {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			// const start = dayjs.tz(vacationStart, 'Asia/Singapore');
			// const end = dayjs.tz(vacationEnd, 'Asia/Singapore');
			// const result = await pb.collection('vacation').create({
			// 	user: pb.authStore.record?.id,
			// 	startDateTime: start,
			// 	endDateTime: end
			// });
			// if (result.id) {
			// 	addToast('success', 'Added successfully!');
			// 	spinner = false;
			// 	await tanstackClient.refetchQueries(createVacationRefetchOptions());
			// }
		} catch (err) {
			console.log(err);
		}
	}

	async function deleteHandler(userId: string) {
		if (!user.isSuccess || !family.data) return;
		spinner = true;

		try {
			const result = await pb.collection('families').update(family.data.id, {
				'members-': userId
			});

			if (result) {
				addToast('success', 'Removed member!');
				spinner = false;
				await tanstackClient.refetchQueries(familyRefetchOptions());
			}
		} catch (err) {
			console.log(err);
		}
	}

	async function handleInvite() {
		if (!family.data) return;

		invited = true;

		const inviteResponse = await pb.collection('invites').create({
			family: family.data.id,
			familyNameSnapshot: family.data?.name,
			ownerEmailSnapshot: family.data.expand?.owner?.email,
			userEmail: invitee,
			status: 'pending'
		});

		await pb.collection('families').update(family.data.id, {
			'activeInvites+': inviteResponse.id
		});

		if (inviteResponse) {
			addToast('success', 'Invite sent!');
		}

		invitee = '';
		setTimeout(() => {
			invited = false;
		}, 3000);
	}
</script>

<PageWrapper title="Family" {pb} largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl max-lg:h-full max-lg:grid-rows-[1fr_auto] lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Family</h1>

		<div class="grid w-full content-start gap-8">
			{#if currentInvite.isSuccess && currentInvite.data && currentInvite.data.status !== 'completed'}
				<section
					class="border-base-300 bg-info text-info-content grid min-h-18 gap-4 rounded-2xl border p-4 text-lg font-semibold"
				>
					<a
						href="/profile/family/invite?i={currentInvite.data.family}"
						class="flex items-center gap-4"
						><MdiAlertCircle class="size-[2em]" />You received an invite <MaterialSymbolsArrowRightAlt
							class="-ms-3 size-[1.3em]"
						/></a
					>
				</section>
			{/if}

			<section class="border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 p-4">
				<h2 class="flex items-center gap-4 text-xl font-bold">
					<MaterialSymbolsFamilyGroup class="size-[1.3em]" />Your Family
				</h2>
				{#if family.isPending && !family.data}
					<div class="custom-loader"></div>
				{/if}
				{#if family.isSuccess && family.data}
					<div>
						{#if family.data.owner === pb.authStore.record?.id}
							You're the owner of "{family.data.name}".
						{:else}
							You are in
							<span class="text-primary font-bold">
								{#if family.data.expand?.owner?.name}
									{family.data.expand?.owner?.name}
								{:else}
									{@const nameFromEmail = cleanEmail(family.data.expand?.owner?.email)}
									{nameFromEmail}
								{/if}</span
							>'s family.
						{/if}
					</div>
				{/if}
				<div class="grid grid-cols-1 justify-items-start">
					<button class="flex items-center gap-2"
						>Switch to another family<MaterialSymbolsArrowRightAlt /></button
					>
				</div>
			</section>
			<section class="border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 p-4">
				<h2 class="text-xl font-bold">Members ({family.data?.members.length})</h2>
				{#if family.isPending && !family.data}
					<div class="custom-loader"></div>
				{/if}

				<ul class="grid list-disc">
					{#if family.isSuccess && family.data}
						{#each family.data.expand?.members as member, i}
							<li class="flex items-center">
								<div class="flex grow items-center gap-4 py-1">
									<MaterialSymbolsPerson />
									{#if member.name}
										{member.name}
									{:else}
										{member.email}
									{/if}
									<div class="flex items-center gap-2">
										{#if family.data.owner === member.id}
											<span class="btn btn-outline btn-xs border-base-content/50 border">Owner</span
											>
										{/if}
										{#if pb.authStore.record?.id === member.id}
											<span class="btn btn-outline btn-xs border-base-content/50 border">You</span>
										{/if}
									</div>
								</div>

								<div class="dropdown dropdown-end">
									<div tabindex="0" role="button" class="btn btn-ghost drawer-button px-2 py-0">
										<MaterialSymbolsMoreVert class="" />
									</div>
									<ul
										class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 min-w-32 shadow-lg"
									>
										<button
											onclick={() => modals[i].showModal()}
											class="btn btn-ghost flex w-full items-center gap-2 rounded-xl"
											><MaterialSymbolsPersonRemove class="size-[1.3em]" />Remove</button
										>
									</ul>
								</div>
							</li>
						{/each}
					{/if}
				</ul>
			</section>

			<section class="border-base-300 grid min-h-18 rounded-2xl border bg-white/70 p-4">
				<h2 class="mb-4 text-xl font-bold">Invite Someone</h2>
				{#if family.isPending && !family.data}
					<div class="custom-loader"></div>
				{/if}
				{#if family.isSuccess && family.data}
					<form onsubmit={handleInvite} class="join mb-2">
						<input
							type="email"
							placeholder="Enter email address"
							class="input join-item focus-within:border-neutral w-full font-medium focus-within:outline-0"
							bind:value={invitee}
							required
						/>
						<button
							class={[
								'join-item btn btn-neutral btn-0 flex min-w-20 items-center gap-2 border-dashed',
								invited && 'btn-success'
							]}
						>
							{#if !invited}
								Invite
								<!-- <MaterialSymbolsShare class="size-6" /> -->
							{:else}
								<MaterialSymbolsCheck class="size-6" />
							{/if}
						</button>
					</form>
					<p class="text-xs">They'll need to accept your invitation.</p>
				{/if}
			</section>
		</div>
	</div>
</PageWrapper>

{#if family.isSuccess && family.data}
	{#each family.data.expand?.members as member, i}
		<dialog bind:this={modals[i]} class="modal modal-bottom sm:modal-middle">
			<div class="modal-box grid gap-8">
				<div
					class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center overflow-hidden rounded-full"
				>
					<MaterialSymbolsPersonRemove class="ms-2.5 size-12" />
				</div>
				<div class="grid gap-4">
					<h2 class="text-2xl font-bold">Remove Member?</h2>
					<ul class="ms-6 list-disc space-y-2">
						<li>This will revoke their access to your family's logs.</li>
					</ul>
				</div>
				<div class="grid grid-cols-1 gap-4">
					<button class="btn btn-error btn-lg" onclick={() => deleteHandler(member.id)}
						>Remove</button
					>
					<form method="dialog" class="">
						<button class="btn btn-outline btn-lg w-full">Cancel</button>
					</form>
				</div>
			</div>
		</dialog>
	{/each}
{/if}
