import { marketCategories } from '$lib/market';

export function match(val: string): boolean {
	const result = marketCategories.find((m) => m.value === val);

	if (result) return true;

	return false;
}
