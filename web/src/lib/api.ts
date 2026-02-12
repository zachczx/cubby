import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import ky from 'ky';
import { addToast } from './ui/ArkToaster.svelte';

export const api = ky.create({
	prefixUrl: PUBLIC_API_URL,
	credentials: 'include',
	hooks: {
		afterResponse: [
			(_request, _options, response) => {
				if (response.status === 403 || response.status === 401) {
					if (browser) {
						console.error('403 Not auth');
						addToast('error', "You're not logged in!");
						redirect(307, '/login');
					}

					return new Response(JSON.stringify(null), {
						status: 403,
						headers: { 'Content-Type': 'application/json' }
					});
				}
			}
		]
	}
});
