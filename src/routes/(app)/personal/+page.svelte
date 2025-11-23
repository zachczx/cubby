<script lang="ts">
	import PageWrapper from '$lib/shell/PageWrapper.svelte';
	import { pb } from '$lib/pb';
	import dayjs from 'dayjs';
	import utc from 'dayjs/plugin/utc';
	import timezone from 'dayjs/plugin/timezone';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import FluentEmojiFlatShield from '$lib/assets/expressive-icons/FluentEmojiFlatShield.svelte';
	import FluentEmojiFlatLotionBottle from '$lib/assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
	import ActionCard from '$lib/ui/ActionCard.svelte';
	import { gummy, spray } from '$lib/ui/ActionCardDefaults';

	dayjs.extend(relativeTime);
	dayjs.extend(utc);
	dayjs.extend(timezone);

	let buttonStatuses = $state<Record<string, ButtonState>>({
		spray: 'default',
		gummy: 'default'
	});
</script>

<PageWrapper title="Personal" {pb}>
	<main class="h-full">
		<div id="mobile" class="grid w-full max-w-lg gap-8 justify-self-center lg:text-base">
			<div class="grid gap-8 py-4">
				<ActionCard
					options={{
						...spray,
						button: { ...spray.button, status: buttonStatuses.spray }
					}}
				></ActionCard>

				<ActionCard
					options={{
						...gummy,
						button: { ...gummy.button, status: buttonStatuses.gummy }
					}}
				></ActionCard>
			</div>
		</div>
	</main>
</PageWrapper>
