<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import { api } from '$lib/api';
	import { onMount } from 'svelte';
	import { queryClient } from '$lib/queries';

	onMount(async () => {
		await api.post('logout');

		// Explicitly clear cookies on the frontend to prevent /check from refreshing logic
		document.cookie = 'stytch_session_jwt=; Max-Age=0; path=/; domain=' + window.location.hostname;
		document.cookie =
			'stytch_session_token=; Max-Age=0; path=/; domain=' + window.location.hostname;

		// Fallback without domain just in case
		document.cookie = 'stytch_session_jwt=; Max-Age=0; path=/';
		document.cookie = 'stytch_session_token=; Max-Age=0; path=/';

		queryClient.clear();
		addToast('success', 'Successfully logged out!');
		goto(resolve('/login'));
	});
</script>
