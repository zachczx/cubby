<script lang="ts">
	import { browser } from '$app/environment';
	import { PUBLIC_API_URL } from '$env/static/public';
	import ky from 'ky';
	import { onMount } from 'svelte';

	let message = $state('');

	const api = ky.create({
		prefixUrl: PUBLIC_API_URL + '/',
		credentials: 'include',
		hooks: {
			afterResponse: [
				(_request, options, response) => {
					if (response.status === 403 || response.status === 401) {
						if (browser) {
							message = '403 Not auth';
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

	onMount(async () => {
		message = await api.get('check').json();
	});
</script>

{message}
