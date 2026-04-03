export const router = {
	tracker: (id: string) => `/app/trackers/${id}`,
	marketCategory: (category: string) => `/app/market/${encodeURIComponent(category)}`,
	marketItem: (category: string, item: string) =>
		`/app/market/${encodeURIComponent(category)}/${encodeURIComponent(item)}`,
	marketAdd: (opts?: { category?: string; item?: string }) => {
		const params = new URLSearchParams();
		if (opts?.category) params.set('category', opts.category);
		if (opts?.item) params.set('item', opts.item);
		const qs = params.toString();
		return `/app/market/add${qs ? `?${qs}` : ''}`;
	},
	marketEdit: (id: string) => `/app/market/edit/${encodeURIComponent(id)}`,
	marketCategories: '/app/market/categories',
	gym: (workoutId?: string) => (workoutId ? `/app/gym/${workoutId}` : '/app/gym'),
	gymStats: '/app/gym/stats',
	app: (route: string) => `/app/${route}`
} as const;
