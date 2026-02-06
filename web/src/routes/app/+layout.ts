import { api } from '$lib/api';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from '../$types';

export const load: LayoutLoad = async () => {
	try {
		console.log('[Layout] Checking authentication...');
		const user = await api.get('check').json();
		console.log('[Layout] Auth check passed:', user);
		return { user };
	} catch (error) {
		console.error('[Layout] Auth check failed:', error);
		redirect(307, '/login');
	}
};
