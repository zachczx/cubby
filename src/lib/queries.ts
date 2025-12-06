import { QueryClient, queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import { pb } from './pb';
import dayjs from 'dayjs';
import type { ListResult } from 'pocketbase';

const staleTime = 5 * 60 * 1000;
const rootKey = [pb.authStore?.record?.id];
export const queryClient = new QueryClient();

function createQueryFactory<T>(key: string[], queryFn: () => Promise<T>) {
	return {
		options: () =>
			queryOptions<T>({
				queryKey: [...rootKey, ...key],
				queryFn,
				staleTime: staleTime
			}),

		refetch: (): RefetchQueryFilters => ({
			queryKey: [...rootKey, ...key],
			type: 'active',
			exact: true
		})
	};
}

const allLogsQuery = createQueryFactory(
	['logs-all'],
	async (): Promise<LogsDB[]> =>
		await pb.collection('logs').getFullList({
			sort: '-time'
		})
);
export const allLogsQueryOptions = allLogsQuery.options;
export const allLogsRefetchOptions = allLogsQuery.refetch;

const allTrackersQuery = createQueryFactory(
	['allTrackers'],
	async (): Promise<TrackerDB[]> =>
		await pb.collection('trackers').getFullList({ sort: 'display', expand: 'family.owner' })
);
export const allTrackersQueryOptions = allTrackersQuery.options;
export const allTrackersRefetchOptions = allTrackersQuery.refetch;

const notificationQuery = createQueryFactory(
	['notification'],
	async (): Promise<LogsDB[]> =>
		await pb.collection('latest_logs').getFullList({ expand: 'tracker' })
);
export const notificationQueryOptions = notificationQuery.options;
export const notificationRefetchOptions = notificationQuery.refetch;

const feedQuery = createQueryFactory(
	['feed'],
	async (): Promise<ListResult<LogsDB>> =>
		await pb.collection('logs').getList(1, 5, { expand: 'tracker', sort: '-time' })
);
export const feedQueryOptions = feedQuery.options;
export const feedRefetchOptions = feedQuery.refetch;

/** 
requestKey prevents "Auto-cancellation" errors by pb. Spray was getting loaded twice, with first being cancelled.
Cause: pb treats simultaneous requests to same endpoint as duplicates, so cancels it.
Solution: 'requestKey' needed so only the same tracker fetch cancels the previous.
*/
const getLogsFactory = (name: string | undefined) => {
	const safeName = name ?? '';

	return createQueryFactory(['logs', safeName], async (): Promise<LogsDB[] | []> => {
		if (!name) return [];

		return await pb.collection('logs').getFullList({
			filter: `tracker.name="${name}"`,
			sort: '-time',
			requestKey: `${name}`
		});
	});
};

export const logsQueryOptions = (name: string | undefined) => {
	const factory = getLogsFactory(name);

	return {
		...factory.options(),
		enabled: !!name
	};
};

export const logsRefetchOptions = (name: string | undefined) => getLogsFactory(name).refetch();

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

const getTrackerFactory = (name: string | undefined) => {
	const safeName = name ?? '';

	return createQueryFactory(['trackers', safeName], async (): Promise<TrackerDB | null> => {
		if (!name) return null;

		return await pb.collection('trackers').getFirstListItem(`name="${name}"`, {
			requestKey: `${name}-tracker-details`
		});
	});
};
export const trackerQueryOptions = (name: string | undefined) => {
	const factory = getTrackerFactory(name);

	return { ...factory.options(), enabled: !!name };
};
export const trackerRefetchOptions = (name: string) => getTrackerFactory(name).refetch();

const familyQuery = createQueryFactory(['family'], async (): Promise<FamilyDB[]> => {
	const resp: FamilyDB[] = await pb.collection('families').getFullList({
		filter: `members.id?="${pb.authStore?.record?.id}" && enabled=true`,
		expand: 'members,owner'
	});
	return resp ?? null;
});
export const familyQueryOptions = familyQuery.options;
export const familyRefetchOptions = familyQuery.refetch;

const inviteQuery = createQueryFactory(['invite'], async (): Promise<InviteDB> => {
	const res: InviteDB = await pb
		.collection('invites')
		.getFirstListItem(`userEmail="${pb.authStore.record?.email}" && status="pending"`, {
			expand: 'family'
		});
	return res ?? null;
});
export const inviteQueryOptions = inviteQuery.options;
export const inviteRefetchOptions = inviteQuery.refetch;

const userQuery = createQueryFactory(['users'], async (): Promise<UserDB> => {
	return await pb.collection('users').getOne(String(pb.authStore?.record?.id));
});
export const userQueryOptions = userQuery.options;
export const userRefetchOptions = userQuery.refetch;

const vacationQuery = createQueryFactory(
	['vacations'],
	async (): Promise<VacationDB[]> =>
		await pb.collection('vacation').getFullList({ sort: '-startDateTime' })
);
export const vacationQueryOptions = vacationQuery.options;
export const vacationRefetchOptions = vacationQuery.refetch;

export function cleanEmail(email: string | undefined): string {
	if (!email) return '';

	const name = email.split('@')?.[0];
	const maxLength = name.length > 11 ? 11 : name.length;
	const clean = name.slice(0, maxLength);

	return clean;
}

export function trackerNameToId(name: string, trackers: TrackerDB[] | undefined): string | null {
	if (!trackers) return null;

	return trackers.find((item) => item.name === name)?.id ?? null;
}

export function trackerIdToName(id: string, trackers: TrackerDB[] | undefined): string | null {
	if (!trackers) return null;

	return trackers.find((item) => id === item.id)?.name ?? null;
}
