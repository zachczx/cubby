<script lang="ts">
	import NotificationManager from '$lib/components/NotificationManager.svelte';
	import Refresher from '$lib/components/Refresher.svelte';
	import { invalidateAll } from '$app/navigation';
	import { type Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<NotificationManager />
<Refresher onRefresh={() => invalidateAll()}>
	{@render children()}
</Refresher>
