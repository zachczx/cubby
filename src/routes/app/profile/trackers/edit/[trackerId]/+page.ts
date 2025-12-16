import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const trackerId = params.trackerId;

	return { trackerId };
};

export const prerender = false;
