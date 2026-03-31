import type { MarketCategoryValue } from '$lib/market';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url }) => {
	const category = url.searchParams.get('category') as MarketCategoryValue | null;
	const item = url.searchParams.get('item');
	return { prefillCategory: category ?? undefined, prefillItemName: item ?? undefined };
};
