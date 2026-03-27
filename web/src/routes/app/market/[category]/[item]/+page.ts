import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	return {
		category: decodeURIComponent(params.category),
		item: decodeURIComponent(params.item)
	};
};
