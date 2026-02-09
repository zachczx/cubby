import { QueryClient, queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import { pb } from './pb';
import dayjs from 'dayjs';
import { api } from './api';

const staleTime = 5 * 60 * 1000;
const rootKey = [pb.authStore?.record?.id];
export const queryClient = new QueryClient();

// Helper to get the current rootKey (useful for manual cache updates)
export const getRootKey = () => [pb.authStore?.record?.id];

// Helper to get the allLogs query key
export const getAllEntriesQueryKey = () => [...getRootKey(), 'entries-all'];

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

const allEntriesQuery = createQueryFactory(
	['entries-all'],
	async (): Promise<EntryDB[]> => await api.get(`entries`).json()
);
export const allEntriesQueryOptions = allEntriesQuery.options;
export const allEntriesRefetchOptions = allEntriesQuery.refetch;

const allTrackersQuery = createQueryFactory(
	['allTrackers'],
	async (): Promise<TrackerDB[]> => await api.get(`trackers`).json()
);
export const allTrackersQueryOptions = allTrackersQuery.options;
export const allTrackersRefetchOptions = allTrackersQuery.refetch;

export async function createEntryQuery(options: {
	trackerId: string;
	interval: number | undefined;
	intervalUnit: IntervalUnit | undefined;
}) {
	return await api.post(`trackers/${options.trackerId}/entries`, {
		body: JSON.stringify({
			trackerId: options.trackerId,
			interval: options.interval,
			intervalUnit: options.intervalUnit,
			performedAt: dayjs.tz(new Date(), 'Asia/Singapore')
		})
	});
}

const familyQuery = createQueryFactory(['family'], async (): Promise<Family[]> => {
	return await api.get('users/me/families').json();
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
export const getUserQueryKey = () => [...getRootKey(), 'users'];

const vacationQuery = createQueryFactory(
	['vacations'],
	async (): Promise<VacationDB[]> =>
		await pb.collection('vacation').getFullList({ sort: '-startDateTime' })
);
export const vacationQueryOptions = vacationQuery.options;
export const vacationRefetchOptions = vacationQuery.refetch;
