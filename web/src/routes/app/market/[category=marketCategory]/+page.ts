import type { MarketCategoryValue } from '$lib/market';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const category = params.category as MarketCategoryValue;
	return { category };
};
