import { queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import { pb } from './pb';
import dayjs from 'dayjs';

const staleTime = 5 * 60 * 1000;

export function logsQueryOptions(name: CollectionName) {
	return queryOptions({
		queryKey: ['log-' + name, pb.authStore?.record?.id],
		queryFn: async (): Promise<LogsDB[]> =>
			await pb.collection('logs').getFullList({
				filter: `user="${pb.authStore?.record?.id}" && tracker.name="${name}"`,
				sort: '-time'
			}),
		staleTime: staleTime
	});
}

export function logsRefetchOptions(name: CollectionName): RefetchQueryFilters {
	return {
		queryKey: ['log-' + name, pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export async function createLogsQuery(options: {
	collectionName: CollectionName;
	interval: number | undefined;
	intervalUnit: IntervalUnit | undefined;
}) {
	const response = await pb.collection('logs').create({
		tracker: trackerNameToId(options.collectionName),
		user: pb.authStore.record?.id,
		time: dayjs.tz(new Date(), 'Asia/Singapore'),
		interval: options.interval,
		intervalUnit: options.intervalUnit
	});

	return response;
}

export function notificationQueryOptions(name: CollectionName) {
	return queryOptions({
		queryKey: ['notif-' + name, pb.authStore?.record?.id],
		queryFn: async (): Promise<LogsDB> =>
			await pb
				.collection('logs')
				.getFirstListItem(`user="${pb.authStore?.record?.id}" && tracker.name="${name}"`, {
					sort: '-time'
				}),
		staleTime: staleTime
	});
}

export function notificationRefetchOptions(name: CollectionName): RefetchQueryFilters {
	return {
		queryKey: ['notif-' + name, pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function trackerQueryOptions(name: CollectionName) {
	return queryOptions({
		queryKey: ['tracker-' + name, pb.authStore?.record?.id],
		queryFn: async (): Promise<TrackerDB> =>
			await pb
				.collection('trackers')
				.getFirstListItem(`user="${pb.authStore?.record?.id}" && name="${name}"`),
		staleTime: staleTime
	});
}

export function trackerRefetchOptions(name: CollectionName): RefetchQueryFilters {
	return {
		queryKey: ['tracker-' + name, pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createUserQueryOptions() {
	return queryOptions({
		queryKey: ['users', pb.authStore?.record?.id],
		queryFn: async (): Promise<UserDB> =>
			await pb.collection('users').getOne(String(pb.authStore?.record?.id)),
		staleTime: staleTime
	});
}

export function createUserRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['users', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

const trackers = {
	bedsheet: 'vk58159wczyxmus',
	doggoBath: 'hhz09lsfj1o5mbp',
	doggoChewable: 'h3e3xkbmoxma6dv',
	gummy: '8t9hsvqah63rs7h',
	spray: '381t91o03thrvyd',
	towel: 'vvd9jnl0uw8qnie'
};

export function trackerNameToId(name: CollectionName): string | null {
	switch (name) {
		case 'bedsheet':
			return trackers.bedsheet;
		case 'doggoBath':
			return trackers.doggoBath;
		case 'doggoChewable':
			return trackers.doggoChewable;
		case 'gummy':
			return trackers.gummy;
		case 'spray':
			return trackers.spray;
		case 'towel':
			return trackers.towel;
		default:
			return null;
	}
}

export function trackerIdToName(id: string): string | undefined {
	const finder = Object.entries(trackers).find((item) => id === item[1]);

	if (finder) return finder[0];

	return undefined;
}

export function createVacationQueryOptions() {
	return queryOptions({
		queryKey: ['vacations', pb.authStore?.record?.id],
		queryFn: async (): Promise<VacationDB[]> =>
			await pb.collection('vacation').getFullList({ sort: '-startDateTime' }),
		staleTime: staleTime
	});
}

export function createVacationRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['vacations', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

// Legacy

// export function createTowelQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['towels', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<TowelDB[]> =>
// 			await pb.collection('towel').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createSprayQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['sprays', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<SprayDB[]> =>
// 			await pb.collection('spray').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createGummyQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['gummies', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<GummyDB[]> =>
// 			await pb.collection('gummy').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createDoggoChewableQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['doggoChewables', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<DoggoChewableDB[]> =>
// 			await pb.collection('doggoChewable').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createDoggoBathQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['doggoBaths', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<DoggoBathDB[]> =>
// 			await pb.collection('doggoBath').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createTowelRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['towels', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }

// export function createSprayRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['sprays', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }

// export function createGummyRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['gummies', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }

// export function createDoggoChewableRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['doggoChewables', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }

// export function createDoggoBathRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['doggoBaths', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }
