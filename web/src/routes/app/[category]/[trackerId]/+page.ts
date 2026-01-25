import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const category = params.category;
	const trackerId = params.trackerId;

	return { category, trackerId };
};
