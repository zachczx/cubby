import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';
import ky from 'ky';

export const api = ky.create({
	prefixUrl: PUBLIC_API_URL,
	credentials: 'include',
	hooks: {
		afterResponse: [
			(_request, options, response) => {
				if (response.status === 403 || response.status === 401) {
					if (browser) {
						console.error('403 Not auth');
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
