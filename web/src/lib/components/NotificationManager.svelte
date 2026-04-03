<script lang="ts">
	import { onMount } from 'svelte';
	import { Capacitor } from '@capacitor/core';
	import { PushNotifications } from '@capacitor/push-notifications';
	import { api } from '$lib/api';
	import { addToast } from '$lib/ui/ArkToaster.svelte';

	onMount(async () => {
		if (!Capacitor.isNativePlatform()) {
			console.log('not a native platform');
			return;
		}

		let permStatus = await PushNotifications.checkPermissions();
		if (permStatus.receive === 'prompt') {
			permStatus = await PushNotifications.requestPermissions();
		}

		if (permStatus.receive !== 'granted') {
			console.log('user rejected push notification permission');
			return;
		}

		await PushNotifications.register();

		PushNotifications.addListener('registration', async (token) => {
			console.log('Push registration success, token: ' + token.value);
			try {
				await api.post('tokens', {
					body: JSON.stringify({
						token: token.value,
						platform: Capacitor.getPlatform() // 'ios' or 'android'
					})
				});
			} catch (err) {
				console.error('Failed to save push token to backend', err);
			}
		});

		PushNotifications.addListener('registrationError', (error) => {
			console.error('Error on registration: ' + JSON.stringify(error));
			addToast('error', 'Failed to register for notifications.');
		});

		PushNotifications.addListener('pushNotificationReceived', (notification) => {
			console.log('Push received: ' + JSON.stringify(notification));
			// addToast('info', notification.body ?? notification.title ?? 'New notification');
		});

		PushNotifications.addListener('pushNotificationActionPerformed', (notification) => {
			console.log('Push action performed: ' + JSON.stringify(notification));
			// Handle routing when the user taps the notification
			// TODO: if (notification.notification.data.trackerId) goto tracker page
		});
	});
</script>
