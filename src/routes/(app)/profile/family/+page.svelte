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

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(createUserQueryOptions);
	const family = createQuery(familyQueryOptions);
	const currentInvite = createQuery(inviteQueryOptions);

	$inspect(family);
	let spinner = $state(false);
	let invited = $state(false);

	let invitee = $state('');

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

	async function deleteHandler(deleteId: string) {
		if (!user.isSuccess) return;
		spinner = true;

		try {
			// const result = await pb.collection('vacation').delete(deleteId);
			// if (result) {
			// 	addToast('success', 'Deleted successfully!');
			// 	spinner = false;
			// 	await tanstackClient.refetchQueries(createVacationRefetchOptions());
			// 	invalidateAll();
			// }
		} catch (err) {
			console.log(err);
		}
	}

	async function handleInvite() {
		if (!family.data) return;
		// const shareData = {
		// 	title: 'Join my Sundry Family',
		// 	text: `Join my family on Sundry! Use code: ${inviteId}`,
		// 	url: `http://localhost:5173/profile/family/invite?i=${inviteId}`
		// };
		// if (navigator.share) {
		// 	try {
		// 		await navigator.share(shareData);
		// 		return;
		// 	} catch (err) {
		// 		console.log(err);
		// 	}
		// }
		// await navigator.clipboard.writeText(code);
		// addToast('success', 'Code copied to clipboard!');

		invited = true;

		const inviteResponse = await pb.collection('invites').create({
			family: family.data.id,
			familyNameSnapshot: family.data?.name,
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
					{@const f = family.data}

					<div>
						{#if f.owner === pb.authStore.record?.id}
							You're the owner of "{f.name}".
						{:else}
							You are in
							<span class="text-primary font-bold">
								{#if f.expand?.owner?.name}
									{f.expand?.owner?.name}
								{:else}
									{@const nameFromEmail = cleanEmail(f.expand?.owner?.email)}
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
						{@const f = family.data}
						{#each f.expand?.members as member}
							<li class="flex items-center gap-4 py-1">
								<MaterialSymbolsPerson />

								{#if member.name}
									{member.name}
								{:else}
									{member.email}
								{/if}
								<div class="flex items-center gap-2">
									{#if f.owner === member.id}
										<span class="btn btn-outline btn-xs border-base-content/50 border">Owner</span>
									{/if}
									{#if pb.authStore.record?.id === member.id}
										<span class="btn btn-outline btn-xs border-base-content/50 border">You</span>
									{/if}
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
					<p class="text-xs">They need to accept your invitation.</p>
				{/if}
			</section>

			<!-- <section class="border-base-300 grid min-h-18 rounded-2xl border bg-white/70 p-4">
				<h2 class="mb-4 text-xl font-bold">Join Someone's Family</h2>
				<div class="join mb-2">
					<input type="text" disabled class="input join-item w-full font-medium" value="" />
					<button
						type="button"
						class={[
							'join-item btn btn-neutral btn-0 flex min-w-20 items-center gap-2',
							copied && 'btn-success'
						]}
					>
						Join
					</button>
				</div>
			</section> -->
		</div>
	</div>
</PageWrapper>

<!-- <dialog bind:this={vacationsModal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box">
		<form method="dialog">
			<button class="btn btn-sm btn-circle btn-ghost absolute top-2 right-2">✕</button>
		</form>
		<h3 class="mb-4 text-lg font-bold uppercase">Recent Vacations</h3>
		<ul class="list-disc">
			{#if vacations.isSuccess}
				{#each vacations.data as v}
					{@const dateTime = formatTime(v.startDateTime, v.endDateTime)}
					<li class="border-b-base-300 ms-6 border-b py-4">
						<div class="flex items-center gap-4">
							<div class="grow">{dateTime}</div>
							<button
								class="btn btn-error btn-sm"
								onclick={() => deleteHandler(v.id)}
								aria-label="delete"><MaterialSymbolsDelete class="size-5" /></button
							>
						</div>
					</li>
				{/each}
			{/if}
		</ul>
	</div>
</dialog> -->
