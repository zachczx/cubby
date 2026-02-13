<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { PUBLIC_API_URL } from '$env/static/public';
	import { addToast } from '$lib/ui/ArkToaster.svelte';
	import ky from 'ky';
	import { onMount } from 'svelte';

	onMount(async () => {
		await ky.get(`${PUBLIC_API_URL}/logout`, {
			credentials: 'include'
		});

		addToast('success', 'Successfully logged out!');
		goto(resolve('/login'));
	});
</script>
