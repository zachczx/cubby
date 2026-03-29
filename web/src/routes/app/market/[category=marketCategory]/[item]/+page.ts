import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const category = params.category;
	const item = params.item;
	return { category, item };
};
