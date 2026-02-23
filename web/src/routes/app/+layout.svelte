<script lang="ts">
	import { firebaseIconSrc, getFirebaseMessaging } from '$lib/firebase';
	import { onMessage } from 'firebase/messaging';
	import { onMount, type Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	onMount(async () => {
		const messaging = getFirebaseMessaging();

		onMessage(messaging, async (payload) => {
			const registration = await navigator.serviceWorker.ready;
			registration.showNotification(payload.notification?.title ?? '', {
				body: payload.notification?.body,
				icon: firebaseIconSrc,
				data: payload.data
			});
		});
	});
</script>

{@render children()}
