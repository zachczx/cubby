import { api } from '$lib/api';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from '../$types';

export const load: LayoutLoad = async () => {
	try {
		const user = await api.get('check').json();
		return { user };
	} catch (error) {
		redirect(307, '/login');
	}
};
