import { pb } from '$lib/pb';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from '../$types';

export const load: LayoutLoad = async () => {
	if (!pb.authStore.isValid) {
		redirect(307, '/login');
	}
};
