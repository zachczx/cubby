<script lang="ts">
	import { PinInput, usePinInput } from '@ark-ui/svelte/pin-input';

	let { otp = $bindable() } = $props();

	const id = $props.id();

	const pinInput = usePinInput({
		id,
		otp: true,
		required: true,
		type: 'alphanumeric',
		onValueComplete: (evt) => (otp = evt.valueAsString)
	});
</script>

<!-- styles.Root styles.Control styles.Input-->
<PinInput.RootProvider value={pinInput} class="w-full justify-self-center py-4">
	<!-- <PinInput.Label class={styles.Label}>Label</PinInput.Label> -->
	<PinInput.Control class={['grid grid-cols-6 gap-1']}>
		{#each [0, 1, 2, 3, 4, 5] as id, index (id)}
			<PinInput.Input
				{index}
				class={[
					'border-primary/10 flex h-18 items-center justify-center rounded-xl border-2 bg-white/70 text-center'
				]}
			/>
		{/each}
	</PinInput.Control>
	<PinInput.HiddenInput />
</PinInput.RootProvider>

<style>
</style>
