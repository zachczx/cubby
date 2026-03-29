export const marketCategories = [
	{ value: 'fruit', label: 'Fruit', icon: 'material-symbols:nutrition' },
	{ value: 'vegetable', label: 'Vegetable', icon: 'material-symbols:eco' },
	{ value: 'dairy', label: 'Dairy', icon: 'material-symbols:egg-alt' },
	{ value: 'meat', label: 'Meat', icon: 'material-symbols:restaurant' },
	{ value: 'seafood', label: 'Seafood', icon: 'material-symbols:set-meal' },
	{ value: 'bakery', label: 'Bakery', icon: 'material-symbols:bakery-dining' },
	{ value: 'pantry', label: 'Pantry', icon: 'material-symbols:shelves' },
	{ value: 'frozen', label: 'Frozen', icon: 'material-symbols:ac-unit' },
	{ value: 'beverage', label: 'Beverage', icon: 'material-symbols:local-cafe' },
	{ value: 'snack', label: 'Snack', icon: 'material-symbols:cookie' },
	{ value: 'others', label: 'Others', icon: 'material-symbols:category' }
] as const;

export type MarketCategoryValue = (typeof marketCategories)[number]['value'];

export const marketUnits = [
	'kg',
	'g',
	'each',
	'bunch',
	'punnet',
	'pack',
	'bottle',
	'can',
	'litre',
	'dozen',
	'carton'
] as const;
