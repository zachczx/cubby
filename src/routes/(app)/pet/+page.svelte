<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import FluentEmojiFlatShower from '$lib/assets/expressive-icons/FluentEmojiFlatShower.svelte';
	import FluentEmojiFlatShield from '$lib/assets/expressive-icons/FluentEmojiFlatShield.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import { doggoBath, doggoChewable } from '$lib/ui/ActionCardDefaults';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $state<Record<string, ButtonState>>({
		doggoBath: 'default',
		doggoChewable: 'default'
	});
</script>

<PageWrapper title="Pet" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<ActionCard
					options={{
						...doggoBath,
						button: { ...doggoBath.button, status: buttonStatuses.doggoBath }
					}}
				></ActionCard>

				<ActionCard
					options={{
						...doggoChewable,
						button: { ...doggoChewable.button, status: buttonStatuses.doggoChewable }
					}}
				></ActionCard>
			</div>
		</div>
	</main>
</PageWrapper>
