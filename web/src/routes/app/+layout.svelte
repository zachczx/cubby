<script lang="ts">
	import { page } from '$app/state';
	import { PUBLIC_FCM_VAPID_KEY } from '$env/static/public';
	import { api } from '$lib/api';
	import { firebaseIconSrc, getFirebaseMessaging } from '$lib/firebase';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { getToken, onMessage } from 'firebase/messaging';
	import { onMount, type Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	onMount(async () => {
		const permission = await Notification.requestPermission();
		if (permission !== 'granted') return;

		const registration = await navigator.serviceWorker.ready;

		const messaging = getFirebaseMessaging();

		const newToken = await getToken(messaging, {
			vapidKey: PUBLIC_FCM_VAPID_KEY,
			serviceWorkerRegistration: registration
		});

		try {
			await api.post('tokens', {
				body: JSON.stringify({ token: newToken, platform: 'web' })
			});
		} catch (err) {
			console.log(err);
			addToast('error', 'Notification error!');
			return;
		}

		if (page.url.searchParams.get('login') === 'true') {
			const url = new URL(window.location.href);
			url.searchParams.delete('login');
			window.history.replaceState({}, '', url);
		}

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
