import { PUBLIC_API_URL } from '$env/static/public';
import ky from 'ky';
import { browser } from '$app/environment';
import { Capacitor } from '@capacitor/core';
import { CapacitorCookies } from '@capacitor/core';

let apiUrl = PUBLIC_API_URL;
const cookieUrl = PUBLIC_API_URL;
if (Capacitor.getPlatform() === 'android' && apiUrl.includes('localhost')) {
	apiUrl = apiUrl.replace('localhost', '10.0.2.2');
}

let headers: Record<string, string> = {};
if (Capacitor.getPlatform() === 'android') {
	headers['x-capacitor-app'] = 'true';
}

export const api = ky.create({
	prefixUrl: apiUrl,
	throwHttpErrors: false,
	credentials: 'include',
	headers: headers,
	hooks: {
		beforeRequest: [
			async (request) => {
				if (Capacitor.getPlatform() === 'android') {
					const cookies = await CapacitorCookies.getCookies({ url: cookieUrl });
					if (Object.keys(cookies).length > 0) {
						const cookieString = Object.entries(cookies)
							.map(([key, value]) => `${key}=${value}`)
							.join('; ');
						request.headers.set('Cookie', cookieString);
					}
				}
			}
		],
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
