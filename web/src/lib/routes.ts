export const router = {
	tracker: (id: string) => `/app/trackers/${id}`,
	app: (route: string) => `/app/${route}`
} as const;
