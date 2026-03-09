<script lang="ts">
	import NotificationManager from '$lib/components/NotificationManager.svelte';
	import { type Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<NotificationManager />
{@render children()}
