<script lang="ts">
	import { NumberInput } from '@ark-ui/svelte/number-input';

	let { value = $bindable('1') } = $props();
</script>

<NumberInput.Root name="interval" bind:value min={1}>
	<!-- <NumberInput.Label>Count</NumberInput.Label> -->
	<NumberInput.Control class="input input-lg w-full p-0">
		<NumberInput.DecrementTrigger
			class="text-base-content/70 h-full min-w-16 rounded-l-md bg-(--color-base-input) px-8 text-3xl font-medium"
			>–</NumberInput.DecrementTrigger
		>
		<NumberInput.Input class="text-center" />
		<NumberInput.IncrementTrigger
			class="text-base-content/70 h-full min-w-16 rounded-r-md bg-(--color-base-input) px-2 text-3xl font-medium"
			>+</NumberInput.IncrementTrigger
		>
	</NumberInput.Control>
</NumberInput.Root>
