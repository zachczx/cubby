<script lang="ts">
	import { PUBLIC_FCM_VAPID_KEY } from '$env/static/public';
	import { api } from '$lib/api';
	import { firebaseIconSrc, getFirebaseMessaging } from '$lib/firebase';
	import { getToken, onMessage } from 'firebase/messaging';
	import { onMount } from 'svelte';

	let permission = $state();
	let token = $state();

	onMount(async () => {
		permission = await Notification.requestPermission();
		if (permission !== 'granted') {
			console.log('Permission denied');
			return;
		}

		const registration = await navigator.serviceWorker.ready;

		const messaging = getFirebaseMessaging();
		console.log('messaging instance:', messaging);

		token = await getToken(messaging, {
			vapidKey: PUBLIC_FCM_VAPID_KEY,
			serviceWorkerRegistration: registration
		});

		const response = await api.post('tokens', {
			body: JSON.stringify({ token: token, platform: 'web' })
		});

		console.log(response);

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

<h1 class="text-4xl">Test Firebase</h1>

<p>Token: <b>{token}</b></p>
<p>Permission: <b>{permission}</b></p>
