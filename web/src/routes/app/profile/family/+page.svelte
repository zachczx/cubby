<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import {
		userQueryOptions,
		familyQueryOptions,
		familyRefetchOptions,
		inviteQueryOptions
	} from '$lib/queries';
	import Icon from '@iconify/svelte';
	import { api } from '$lib/api';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	const user = createQuery(userQueryOptions);
	const families = createQuery(familyQueryOptions);

	$inspect(families.data);
	const invites = createQuery(inviteQueryOptions);
	const tanstackClient = useQueryClient();

	let invited = $state('');

	let invitee = $state('');

	let modals = $state<HTMLDialogElement[]>([]);

	async function deleteHandler(memberId: string, family: Family) {
		if (!user.isSuccess || !family) return;

		try {
			const result = await api.delete(`families/${family.id}/${memberId}`);

			console.log(result);

			if (result) {
				addToast('success', 'Removed member!');

				await tanstackClient.refetchQueries(familyRefetchOptions());
			}
		} catch (err) {
			console.log(err);
		}
	}

	async function handleInvite(family: Family) {
		if (!family) return;

		invited = 'pending';

		try {
			const inviteResponse = await api.post(`families/invites`, {
				body: JSON.stringify({
					inviteeEmail: invitee
				})
			});

			if (inviteResponse.status === 201) {
				addToast('success', 'Invite sent!');
				invited = 'success';
				setTimeout(() => {
					invited = '';
				}, 3000);
			}
		} catch (err) {
			invited = '';
			addToast('error', 'Invite error!');
		}
	}

	async function leaveFamily(family: Family) {
		if (!family) return;

		try {
			const response = await api.post(`families/${family.id}/leave`);

			if (response.status === 204) {
				addToast('success', 'Successfully left family!');
				await tanstackClient.refetchQueries(familyRefetchOptions());
				goto(resolve('/app/profile/family'));
			}
		} catch (err) {
			console.log(err);
		}
	}

	let section = $state('');

	$effect(() => {
		if (families.isSuccess && families.data.length > 0) {
			section = families.data[0].id;
		}
	});
</script>

<PageWrapper title="Manage Family" largeScreenCenter={true}>
	<div
		class="lg:bg-base-200 grid w-full rounded-2xl lg:h-min lg:max-w-lg lg:justify-self-center lg:p-8 lg:shadow-md"
	>
		<h1 class="text-primary mb-4 text-center text-4xl font-bold max-lg:hidden">Family</h1>

		<div class="grid w-full content-start gap-8">
			<nav class="flex items-center gap-2">
				{#if families.isSuccess}
					{#each families.data as family (family.id)}
						<button
							onclick={() => {
								section = family.id;
							}}
							class={[
								'w-full max-w-48 border-b-2 pb-1 text-center',
								section === family.id
									? 'border-b-neutral text-neutral font-bold'
									: 'text-base-content/80 border-b-transparent'
							]}
						>
							{#if family.isOwner}
								My Cubby
							{:else}
								{family.name}
							{/if}
						</button>
					{/each}
				{:else}
					<div
						class="border-b-neutral text-neutral w-full max-w-48 border-b-2 pb-1 text-center font-bold"
					>
						<div class="skeleton h-6 w-24 justify-self-center"></div>
					</div>
				{/if}
			</nav>

			{#if families.isPending}
				<div class="skeleton h-27.5 w-full justify-self-center"></div>
				<div class="skeleton h-27.5 w-full justify-self-center"></div>
			{/if}

			{#each families.data as family (family.id)}
				{@const numberOfMembers = family.members.length + 1}
				{#if family.id === section}
					{#if invites.isSuccess}
						{#each invites.data as invite}
							<section
								class="border-base-300 bg-info text-info-content grid min-h-18 gap-4 rounded-2xl border p-4 text-lg font-semibold"
							>
								<a
									href={resolve(`/app/profile/family/invite/${invite.id}`)}
									class="flex items-center gap-4"
									><Icon icon="mdi:alert-circle" class="size-[2em]" />You received an invite <Icon
										icon="material-symbols:arrow-right-alt"
										class="-ms-3 size-[1.3em]"
									/></a
								>
							</section>
						{/each}
					{/if}

					<section class="border-base-300 grid min-h-18 gap-4 rounded-2xl border bg-white/70 p-4">
						<h2 class="text-xl font-bold">Members ({numberOfMembers})</h2>

						<ul class="grid list-disc">
							{#if families.isSuccess && families.data}
								<li class="flex items-center">
									<div class="flex grow items-center gap-2 py-1">
										<Icon icon="material-symbols:person" class="me-2" />
										{#if family.owner.name}
											{family.owner.name}
										{:else}
											{family.owner.email}
										{/if}

										{#if family.isOwner}
											<span class="text-sm opacity-60">(You)</span>
										{/if}
										<div class="flex items-center gap-2">
											<span class="btn btn-outline border-base-content/50 h-auto border">Owner</span
											>
										</div>
									</div>
								</li>
								{#each family.members as member, i (member.id)}
									<li class="flex items-center">
										<div class="flex grow items-center gap-2 py-1">
											<Icon icon="material-symbols:person" class="me-2" />
											{#if member.name}
												{member.name}
											{:else}
												{member.email}
											{/if}
										</div>

										{#if family.isOwner}
											<div class="dropdown dropdown-end">
												<div
													tabindex="0"
													role="button"
													class="btn btn-ghost drawer-button px-2 py-0"
												>
													<Icon icon="material-symbols:more-vert" class="" />
												</div>
												<ul
													class="dropdown-content menu rounded-box bg-base-100 text-md text-base-content z-1 min-w-32 shadow-lg"
												>
													<button
														onclick={() => modals[i].showModal()}
														class="btn btn-ghost flex w-full items-center gap-2 rounded-xl"
													>
														<Icon
															icon="material-symbols:person-remove"
															class="size-[1.3em]"
														/>Remove
													</button>
												</ul>
											</div>
										{/if}
									</li>
								{/each}
							{/if}
						</ul>

						{#if !family.isOwner}
							<div class="border-t-base-300/70 flex justify-center border-t-2 pt-3">
								<button
									class="btn btn-ghost text-error flex items-center gap-2 py-0"
									onclick={() => leaveFamily(family)}
									>Leave Cubby<Icon icon="material-symbols:arrow-right-alt" /></button
								>
							</div>
						{/if}
					</section>

					{#if family.isOwner}
						<section class="border-base-300 grid min-h-18 rounded-2xl border bg-white/70 p-4">
							<h2 class="mb-4 text-xl font-bold">Invite Someone</h2>
							{#if families.isPending && !families.data}
								<div class="custom-loader"></div>
							{/if}
							{#if families.isSuccess && families.data}
								<form onsubmit={() => handleInvite(family)} class="join mb-2">
									<input
										type="email"
										placeholder="Enter email address"
										class="input join-item focus-within:border-neutral w-full font-medium focus-within:outline-0"
										bind:value={invitee}
										required
									/>
									<button
										class={[
											'join-item btn btn-neutral btn-0 flex min-w-20 items-center gap-2 border-dashed shadow-none',
											invited === 'success' && 'btn-success'
										]}
									>
										{#if invited === 'success'}
											<Icon icon="material-symbols:check" class="size-6" />
										{:else if invited === 'pending'}
											<span class="loading loading-sm loading-spinner"></span>
										{:else}
											Invite
										{/if}
									</button>
								</form>
								<p class="text-xs">They'll need to accept your invitation.</p>
							{/if}
						</section>
					{/if}
				{/if}
			{/each}
		</div>
	</div>
</PageWrapper>

{#if families.isSuccess && families.data}
	{#each families.data as family (family.id)}
		{#each family.members as member, i (member.id)}
			<dialog bind:this={modals[i]} class="modal modal-bottom sm:modal-middle">
				<div class="modal-box grid gap-8">
					<div
						class="bg-error/10 text-error flex aspect-square size-20 items-center justify-center justify-self-center overflow-hidden rounded-full"
					>
						<Icon icon="material-symbols:person-remove" class="ms-2.5 size-12" />
					</div>
					<div class="grid gap-4">
						<h2 class="text-2xl font-bold">Remove Member?</h2>
						<ul class="ms-6 list-disc space-y-2">
							<li>
								This will revoke <span class="font-bold"
									>{member.name ? member.name : member.email}</span
								>'s access to your family's logs.
							</li>
						</ul>
					</div>
					<div class="grid grid-cols-1 gap-4">
						<button class="btn btn-error btn-lg" onclick={() => deleteHandler(member.id, family)}
							>Remove</button
						>
						<form method="dialog" class="">
							<button class="btn btn-outline btn-lg w-full">Cancel</button>
						</form>
					</div>
				</div>
			</dialog>
		{/each}
	{/each}
{/if}
