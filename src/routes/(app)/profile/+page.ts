import { pb } from '$lib/pb';
import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	if (!pb.authStore.isValid) {
		redirect(307, '/login');
	}
};
