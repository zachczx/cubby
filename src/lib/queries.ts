import { QueryClient, queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import { pb } from './pb';
import dayjs from 'dayjs';
import type { ListResult } from 'pocketbase';

const staleTime = 5 * 60 * 1000;

export const queryClient = new QueryClient();

export function allLogsQueryOptions() {
	return queryOptions({
		queryKey: ['logs-all', pb.authStore?.record?.id],
		queryFn: async (): Promise<LogsDB[]> =>
			await pb.collection('logs').getFullList({
				sort: '-time'
			}),
		staleTime: staleTime
	});
}

export function allLogsRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['logs-all', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function allTrackersQueryOptions() {
	return queryOptions({
		queryKey: ['trackers-all', pb.authStore?.record?.id],
		queryFn: async (): Promise<TrackerDB[]> =>
			await pb.collection('trackers').getFullList({ sort: 'display', expand: 'family.owner' }),
		staleTime: staleTime
	});
}

export function allTrackersRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['trackers-all', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

/** 
requestKey prevents "Auto-cancellation" errors by pb. Spray was getting loaded twice, with first being cancelled.
Cause: pb treats simultaneous requests to same endpoint as duplicates, so cancels it.
Solution: 'requestKey' needed so only the same tracker fetch cancels the previous.
*/
export function logsQueryOptions(name: string | undefined) {
	return queryOptions({
		queryKey: ['log-' + name, pb.authStore?.record?.id],
		queryFn: async (): Promise<LogsDB[]> => {
			if (!name) return [];

			return await pb.collection('logs').getFullList({
				filter: `tracker.name="${name}"`,
				sort: '-time',
				requestKey: `${name}-logs`
			});
		},
		staleTime: staleTime
	});
}

export function logsRefetchOptions(name: string | undefined): RefetchQueryFilters {
	return {
		queryKey: ['log-' + name, pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export async function createLogsQuery(options: {
	trackerId: string;
	interval: number | undefined;
	intervalUnit: IntervalUnit | undefined;
}) {
	const response = await pb.collection('logs').create({
		tracker: options.trackerId,
		time: dayjs.tz(new Date(), 'Asia/Singapore'),
		interval: options.interval,
		intervalUnit: options.intervalUnit
	});

	return response;
}

export function notificationQueryOptions() {
	return queryOptions({
		queryKey: ['notif', pb.authStore?.record?.id],
		queryFn: async (): Promise<LogsDB[]> =>
			await pb.collection('latest_logs').getFullList({ expand: 'tracker' }),
		staleTime: staleTime
	});
}

export function notificationRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['notif', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function feedQueryOptions() {
	return queryOptions({
		queryKey: ['feed', pb.authStore?.record?.id],
		queryFn: async (): Promise<ListResult<LogsDB>> =>
			await pb.collection('logs').getList(1, 5, { expand: 'tracker', sort: '-time' }),
		staleTime: staleTime
	});
}

export function feedRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['feed', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function trackerQueryOptions(name: string | undefined) {
	return queryOptions({
		queryKey: ['tracker' + name, pb.authStore?.record?.id],
		queryFn: async (): Promise<TrackerDB | null> => {
			if (!name) return null;

			return await pb.collection('trackers').getFirstListItem(`name="${name}"`, {
				requestKey: `${name}-tracker-details`
			});
		},
		staleTime: staleTime
	});
}

export function trackerRefetchOptions(name: string): RefetchQueryFilters {
	return {
		queryKey: ['tracker' + name, pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function familyQueryOptions() {
	return queryOptions({
		queryKey: ['family', pb.authStore?.record?.id],
		queryFn: async (): Promise<FamilyDB[]> => {
			const resp: FamilyDB[] = await pb.collection('families').getFullList({
				filter: `members.id?="${pb.authStore?.record?.id}" && enabled=true`,
				expand: 'members,owner'
			});
			return resp ?? null;
		},
		staleTime: staleTime
	});
}

export function familyRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['family', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function inviteQueryOptions() {
	return queryOptions({
		queryKey: ['invite', pb.authStore?.record?.id],
		queryFn: async (): Promise<InviteDB> => {
			const res: InviteDB = await pb
				.collection('invites')
				.getFirstListItem(`userEmail="${pb.authStore.record?.email}" && status="pending"`, {
					expand: 'family'
				});
			return res ?? null;
		},
		staleTime: staleTime
	});
}

export function inviteRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['invite', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function cleanEmail(email: string | undefined): string {
	if (!email) return '';

	const name = email.split('@')?.[0];
	const maxLength = name.length > 11 ? 11 : name.length;
	const clean = name.slice(0, maxLength);

	return clean;
}

export function createUserQueryOptions() {
	return queryOptions({
		queryKey: ['users', pb.authStore?.record?.id],
		queryFn: async (): Promise<UserDB> => {
			return await pb.collection('users').getOne(String(pb.authStore?.record?.id));
		},
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

export function trackerNameToId(name: string, trackers: TrackerDB[] | undefined): string | null {
	if (!trackers) return null;

	return trackers.find((item) => item.name === name)?.id ?? null;
}

export function trackerIdToName(id: string, trackers: TrackerDB[] | undefined): string | null {
	if (!trackers) return null;

	return trackers.find((item) => id === item.id)?.name ?? null;
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
// 		queryKey: ['chewables', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<DoggoChewableDB[]> =>
// 			await pb.collection('chewable').getFullList({ sort: '-time' }),
// 		staleTime: staleTime
// 	});
// }

// export function createDoggoBathQueryOptions() {
// 	return queryOptions({
// 		queryKey: ['baths', pb.authStore?.record?.id],
// 		queryFn: async (): Promise<DoggoBathDB[]> =>
// 			await pb.collection('bath').getFullList({ sort: '-time' }),
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
// 		queryKey: ['chewables', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }

// export function createDoggoBathRefetchOptions(): RefetchQueryFilters {
// 	return {
// 		queryKey: ['baths', pb.authStore?.record?.id],
// 		type: 'active',
// 		exact: true
// 	};
// }
