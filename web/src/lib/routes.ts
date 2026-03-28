export const router = {
	tracker: (id: string) => `/app/trackers/${id}`,
	marketCategory: (category: string) => `/app/market/${encodeURIComponent(category)}`,
	marketItem: (category: string, item: string) =>
		`/app/market/${encodeURIComponent(category)}/${encodeURIComponent(item)}`,
	marketCategories: '/app/market/categories',
	app: (route: string) => `/app/${route}`
} as const;
