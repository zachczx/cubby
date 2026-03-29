<!-- svelte-ignore state_referenced_locally -->
<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import {
		marketPricesQueryOptions,
		createMarketPriceMutation,
		updateMarketPriceMutation
	} from '$lib/queries';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { marketCategories, marketUnits, type MarketCategoryValue } from '$lib/market';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';

	dayjs.extend(utc);
	dayjs.extend(timezone);

	let {
		onClose,
		editPrice = null,
		paramCategory
	}: {
		onClose: () => void;
		editPrice?: MarketPriceDB | null;
		paramCategory?: MarketCategoryValue;
	} = $props();

	const pricesQuery = createQuery(marketPricesQueryOptions);

	let itemName = $state(editPrice?.itemName ?? '');
	let category = $state(editPrice?.category ?? paramCategory ?? 'fruit');
	let country = $state(editPrice?.country ?? '');
	let store = $state(editPrice?.store ?? '');
	let unit = $state(editPrice?.unit ?? 'pack');
	let quantity = $state(editPrice?.quantity?.toString() ?? '1');
	let price = $state(editPrice?.price?.toString() ?? '');
	let isPromo = $state(editPrice?.isPromo ?? false);
	let remarks = $state(editPrice?.remarks ?? '');
	let customDateEnabled = $state(false);
	let date = $state(dayjs().format('YYYY-MM-DD'));
	let time = $state(dayjs().format('HH:mm'));
	let timestamp = $derived.by(() => {
		const ts = dayjs(date + 'T' + time);
		const tzTime = dayjs.tz(ts, 'Asia/Singapore');

		if (tzTime.diff(dayjs(), 'hours') > 0) {
			return dayjs();
		}

		return tzTime;
	});

	let isSubmitting = $state(false);
	let showSuggestions = $state(false);

	function toTitleCase(s: string): string {
		return s.replace(/\S+/g, (w) => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase());
	}

	let existingNames = $derived.by(() => {
		if (!pricesQuery.isSuccess || !pricesQuery.data) return [];
		const seen = new Set<string>();
		const names: string[] = [];
		for (const p of pricesQuery.data) {
			const lower = p.itemName.toLowerCase();
			if (!seen.has(lower)) {
				seen.add(lower);
				names.push(toTitleCase(p.itemName));
			}
		}
		return names.sort((a, b) => a.localeCompare(b));
	});

	let suggestions = $derived.by(() => {
		if (!itemName || !showSuggestions) return [];
		const lower = itemName.toLowerCase();
		return existingNames
			.filter((n) => n.toLowerCase() !== lower && n.toLowerCase().includes(lower))
			.slice(0, 6);
	});

	function selectSuggestion(name: string) {
		itemName = name;
		showSuggestions = false;
	}

	function handleNameInput() {
		showSuggestions = true;
	}

	function handleNameBlur() {
		setTimeout(() => {
			showSuggestions = false;
		}, 150);
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const trimmedName = itemName.trim();
		if (!trimmedName || !price) return;

		isSubmitting = true;
		try {
			const payload: MarketPriceInput = {
				itemName: trimmedName,
				category: category || null,
				country: country || null,
				store: store || null,
				unit: unit || null,
				quantity: quantity ? parseFloat(quantity) : null,
				price: parseFloat(price),
				isPromo,
				remarks: remarks.trim() || null
			};

			if (customDateEnabled) {
				payload.createdAt = timestamp.toISOString();
			}

			if (editPrice) {
				await updateMarketPriceMutation(editPrice.id, payload);
				addToast('success', 'Price updated successfully!');
				onClose();
			} else {
				const res = await createMarketPriceMutation(payload);
				if (res) {
					if (res.isUpdate) {
						addToast('info', 'Duplicate found — updated existing entry for today.');
					} else {
						addToast('success', 'Price logged successfully!');
					}
					onClose();
				} else {
					addToast('error', 'Failed to log price');
				}
			}
		} catch (err) {
			addToast('error', editPrice ? 'Failed to update price' : 'Failed to log price');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<dialog class="modal modal-open modal-bottom sm:modal-middle bg-base-300/50 backdrop-blur-[2px]">
	<div class="modal-box outline-base-300/50 relative outline">
		<button
			class="btn btn-sm btn-circle btn-ghost absolute top-4 right-4 focus:outline-none"
			onclick={onClose}
		>
			<Icon icon="material-symbols:close" class="size-5" />
		</button>

		<h3 class="mb-6 text-xl font-bold">{editPrice ? 'Edit Price Log' : 'Add Price Log'}</h3>

		<form onsubmit={handleSubmit} class="grid gap-4">
			<div class="form-control relative w-full">
				<label for="item-name" class="label py-1"
					><span class="label-text text-base-content/80 font-medium">Item Name *</span></label
				>
				<input
					id="item-name"
					type="text"
					bind:value={itemName}
					oninput={handleNameInput}
					onfocus={() => (showSuggestions = true)}
					onblur={handleNameBlur}
					placeholder="e.g. Fuji Apple"
					class="input input-bordered focus:outline-primary w-full transition-all"
					required
					autocomplete="off"
				/>
				{#if suggestions.length > 0}
					<div
						class="bg-base-100 border-base-300/50 absolute top-full right-0 left-0 z-50 mt-1 overflow-hidden rounded-xl border shadow-lg"
					>
						{#each suggestions as suggestion}
							<button
								type="button"
								class="hover:bg-base-200/70 w-full px-4 py-2.5 text-left text-sm transition-colors"
								onmousedown={() => selectSuggestion(suggestion)}
							>
								{suggestion}
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<div class="grid grid-cols-3 gap-4">
				<div class="form-control w-full">
					<label for="price" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Price ($) *</span></label
					>
					<input
						id="price"
						type="number"
						step="0.01"
						min="0"
						bind:value={price}
						placeholder="0.00"
						class="input input-bordered focus:outline-primary w-full transition-all"
						required
					/>
				</div>
				<div class="form-control w-full">
					<label for="quantity" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Qty *</span></label
					>
					<input
						id="quantity"
						type="number"
						step="0.01"
						min="0"
						bind:value={quantity}
						placeholder="1"
						required
						class="input input-bordered focus:outline-primary w-full transition-all"
					/>
				</div>
				<div class="form-control w-full">
					<label for="unit" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Unit *</span></label
					>
					<select
						id="unit"
						bind:value={unit}
						class="select select-bordered focus:outline-primary w-full transition-all"
					>
						{#each marketUnits as u}
							<option value={u}>{u}</option>
						{/each}
					</select>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="form-control w-full">
					<label for="category" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Category *</span></label
					>
					<select
						id="category"
						bind:value={category}
						class="select select-bordered focus:outline-primary w-full transition-all"
					>
						{#each marketCategories as cat}
							<option value={cat.value}>{cat.label}</option>
						{/each}
					</select>
				</div>
				<div class="form-control w-full">
					<label for="store" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Store</span></label
					>
					<input
						id="store"
						type="text"
						bind:value={store}
						placeholder="e.g. NTUC"
						class="input input-bordered focus:outline-primary w-full transition-all"
					/>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="form-control w-full">
					<label for="country" class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Country</span></label
					>
					<input
						id="country"
						type="text"
						bind:value={country}
						placeholder="e.g. Japan"
						class="input input-bordered focus:outline-primary w-full transition-all"
					/>
				</div>
				<fieldset class="form-control w-full">
					<legend class="label py-1"
						><span class="label-text text-base-content/80 font-medium">Price Type *</span></legend
					>
					<div class="join w-full">
						<input
							type="radio"
							name="price-type"
							class="btn join-item checked:bg-segmented checked:text-primary-content flex-1"
							aria-label="Regular"
							checked={!isPromo}
							onchange={() => (isPromo = false)}
						/>
						<input
							type="radio"
							name="price-type"
							class="btn join-item checked:bg-segmented checked:text-primary-content flex-1"
							aria-label="Promo"
							checked={isPromo}
							onchange={() => (isPromo = true)}
						/>
					</div>
				</fieldset>
			</div>

			<div class="form-control w-full">
				<label for="remarks" class="label py-1"
					><span class="label-text text-base-content/80 font-medium">Remarks</span></label
				>
				<input
					id="remarks"
					type="text"
					bind:value={remarks}
					placeholder="e.g. Members promo, bulk buy"
					class="input input-bordered focus:outline-primary w-full transition-all"
				/>
			</div>

			<div class={['form-control w-full']}>
				<div class="flex items-center">
					<button
						type="button"
						class={[
							'btn btn-sm my-1',
							!customDateEnabled && 'btn-ghost',
							customDateEnabled && 'btn-soft btn-primary'
						]}
						onclick={() => (customDateEnabled = !customDateEnabled)}
						>Set Date/Time<Icon icon="material-symbols:arrow-right-alt" class="size-[1.3em]" />
					</button>
				</div>
				{#if customDateEnabled}
					<div class="form-control w-full">
						<div class="grid grid-cols-2 gap-2">
							<input type="date" bind:value={date} class="input input-bordered focus:outline-primary w-full transition-all" />
							<input type="time" bind:value={time} class="input input-bordered focus:outline-primary w-full transition-all" />
						</div>
					</div>
				{/if}
			</div>

			<div class="mt-6 grid grid-cols-2 gap-3">
				<button
					type="button"
					class="btn btn-ghost bg-base-200/50 hover:bg-base-300/50 rounded-2xl border border-transparent"
					onclick={onClose}
					disabled={isSubmitting}>Cancel</button
				>
				<button type="submit" class="btn btn-primary rounded-2xl" disabled={isSubmitting}>
					{#if isSubmitting}
						<span class="loading loading-spinner"></span>
					{/if}
					Save Price
				</button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button onclick={onClose}>close</button>
	</form>
</dialog>
