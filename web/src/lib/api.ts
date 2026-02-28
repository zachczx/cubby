import { PUBLIC_API_URL } from '$env/static/public';
import ky from 'ky';
import { browser } from '$app/environment';
import { Capacitor } from '@capacitor/core';

let apiUrl = PUBLIC_API_URL;
if (Capacitor.getPlatform() === 'android' && apiUrl.includes('localhost')) {
	apiUrl = apiUrl.replace('localhost', '10.0.2.2');
}

export const api = ky.create({
	prefixUrl: apiUrl,
	throwHttpErrors: false,
	credentials: 'include',
	hooks: {
		afterResponse: [
			(_request, _options, response) => {
				if (response.status === 403 || response.status === 401) {
					if (browser && window.location.pathname !== '/login') {
						console.error('Not authenticated');
						window.location.href = '/login';
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
