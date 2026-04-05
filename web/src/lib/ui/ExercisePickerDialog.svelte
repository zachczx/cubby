<script lang="ts">
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { favouriteExercisesQueryOptions, getFavouriteExercisesQueryKey } from '$lib/queries';
	import { exercises } from '$lib/exercises';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import Icon from '@iconify/svelte';
	import Dialog from '$lib/ui/Dialog.svelte';

	let {
		open = $bindable(false),
		onpick
	}: {
		open?: boolean;
		onpick: (exerciseId: string) => void;
	} = $props();

	const queryClient = useQueryClient();
	const favouritesDb = createQuery(favouriteExercisesQueryOptions);

	let exerciseSearch = $state('');

	let filteredExercises = $derived(
		exerciseSearch
			? exercises.filter((e) => e.name.toLowerCase().includes(exerciseSearch.toLowerCase()))
			: exercises
	);

	let favouriteIds = $derived(favouritesDb.data?.exerciseIds ?? []);
	let favouriteExercises = $derived(exercises.filter((e) => favouriteIds.includes(e.id)));

	async function toggleFavourite(exerciseId: string) {
		const response = await api.post('gym/favourites', {
			json: { exerciseId }
		});
		if (response.ok) {
			const data: FavouriteExercisesDB = await response.json();
			queryClient.setQueryData(getFavouriteExercisesQueryKey(), data);
		} else {
			addToast('error', 'Failed to update favourite');
		}
	}

	function pick(id: string) {
		exerciseSearch = '';
		onpick(id);
	}
</script>

<Dialog bind:open title="Select Exercise">
	<div class="grid gap-3">
		<input
			type="text"
			class="input input-lg w-full"
			placeholder="Search exercises..."
			bind:value={exerciseSearch}
		/>
		<div class="max-h-[60vh] overflow-y-auto">
			{#if favouriteExercises.length > 0}
				<p
					class="text-base-content/50 px-3 pt-2 pb-1 text-xs font-semibold tracking-wider uppercase"
				>
					Favourites
				</p>
				{#each favouriteExercises as ex (ex.id)}
					<div class="hover:bg-base-200 flex w-full items-center gap-1 rounded-lg px-3 py-2.5">
						<button class="flex-1 cursor-pointer text-left" onclick={() => pick(ex.id)}>
							{ex.name}
						</button>
						<button
							class="btn btn-ghost btn-xs btn-square"
							onclick={() => toggleFavourite(ex.id)}
						>
							<Icon icon="material-symbols:star" class="text-warning size-4" />
						</button>
					</div>
				{/each}
				<div class="divider my-0"></div>
			{/if}
			{#each filteredExercises as ex (ex.id)}
				<div class="hover:bg-base-200 flex w-full items-center gap-1 rounded-lg px-3 py-2.5">
					<button class="flex-1 cursor-pointer text-left" onclick={() => pick(ex.id)}>
						{ex.name}
					</button>
					<button
						class="btn btn-ghost btn-xs btn-square"
						onclick={() => toggleFavourite(ex.id)}
					>
						<Icon
							icon={favouriteIds.includes(ex.id)
								? 'material-symbols:star'
								: 'material-symbols:star-outline'}
							class={favouriteIds.includes(ex.id)
								? 'text-warning size-4'
								: 'text-base-content/30 size-4'}
						/>
					</button>
				</div>
			{/each}
		</div>
	</div>
</Dialog>
