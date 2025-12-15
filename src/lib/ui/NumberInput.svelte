<script lang="ts">
	import { NumberInput } from '@ark-ui/svelte/number-input';

	let { value = $bindable('1') } = $props();
</script>

<NumberInput.Root name="interval" bind:value min={1}>
	<!-- <NumberInput.Label>Count</NumberInput.Label> -->
	<NumberInput.Control class="input input-lg w-full p-0">
		<NumberInput.DecrementTrigger
			class="bg-neutral/40 text-primary-content h-full min-w-16 rounded-l-md px-8 text-3xl"
			>-</NumberInput.DecrementTrigger
		>
		<NumberInput.Input class="text-center" />
		<NumberInput.IncrementTrigger
			class="text-primary-content bg-neutral/40 h-full min-w-16 rounded-r-md px-2 text-3xl"
			>+</NumberInput.IncrementTrigger
		>
	</NumberInput.Control>
</NumberInput.Root>
