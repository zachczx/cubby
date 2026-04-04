import { QueryClient, queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import dayjs from 'dayjs';
import { api } from './api';

const defaultTz = 'Asia/Singapore';
const staleTime = 5 * 60 * 1000;
const rootKey = ['cubby'];
export const queryClient = new QueryClient();

// Helper to get the current rootKey (useful for manual cache updates)
export const getRootKey = () => rootKey;

// Helper to get the allLogs query key
export const getAllEntriesQueryKey = () => [...getRootKey(), 'entries-all'];
export const getAllTrackersQueryKey = () => [...getRootKey(), 'trackers-all'];

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
	['trackers-all'],
	async (): Promise<TrackerDB[]> => await api.get(`trackers`).json()
);
export const allTrackersQueryOptions = allTrackersQuery.options;
export const allTrackersRefetchOptions = allTrackersQuery.refetch;

export async function createEntryQuery(options: {
	trackerId: string;
	interval: number | undefined;
	intervalUnit: IntervalUnit | undefined;
}) {
	return await api
		.post(`trackers/${options.trackerId}/entries`, {
			json: {
				trackerId: options.trackerId,
				interval: options.interval,
				intervalUnit: options.intervalUnit,
				performedAt: dayjs.tz(new Date(), defaultTz)
			}
		})
		.json<EntryDB>();
}

const familyQuery = createQueryFactory(['family'], async (): Promise<Family[]> => {
	return await api.get('users/me/families').json();
});
export const familyQueryOptions = familyQuery.options;
export const familyRefetchOptions = familyQuery.refetch;

const userQuery = createQueryFactory(['users'], async (): Promise<UserDB> => {
	return await api.get('users').json<UserDB>();
});
export const userQueryOptions = userQuery.options;
export const userRefetchOptions = userQuery.refetch;
export const getUserQueryKey = () => [...getRootKey(), 'users'];

const vacationQuery = createQueryFactory(
	['vacations'],
	async (): Promise<VacationDB[]> => await api.get('vacations').json()
);
export const vacationQueryOptions = vacationQuery.options;
export const vacationRefetchOptions = vacationQuery.refetch;

const inviteQuery = createQueryFactory(['invites'], async (): Promise<InviteDB[]> => {
	return await api.get('families/invites').json();
});
export const inviteQueryOptions = inviteQuery.options;
export const inviteRefetchOptions = inviteQuery.refetch;

// Doing this instead of adjusting the factory, since this requires 1 param
export const singleInviteQueryOptions = (inviteId: string) => () =>
	queryOptions<InviteDB>({
		queryKey: [...rootKey, 'invites', inviteId],
		queryFn: async () => await api.get(`families/invites/${inviteId}`).json(),
		staleTime: staleTime
	});

export const singleInviteRefetchOptions = (inviteId: string): RefetchQueryFilters => ({
	queryKey: [...rootKey, 'invites', inviteId],
	exact: true
});

const timerProfilesQuery = createQueryFactory(
	['timer-profiles'],
	async (): Promise<TimerProfileDB[]> => await api.get('timer-profiles').json()
);
export const timerProfilesQueryOptions = timerProfilesQuery.options;
export const timerProfilesRefetchOptions = timerProfilesQuery.refetch;

const allWorkoutsQuery = createQueryFactory(
	['gym-workouts'],
	async (): Promise<WorkoutDB[]> => await api.get('gym/workouts').json()
);
export const allWorkoutsQueryOptions = allWorkoutsQuery.options;
export const allWorkoutsRefetchOptions = allWorkoutsQuery.refetch;
export const getAllWorkoutsQueryKey = () => [...getRootKey(), 'gym-workouts'];

const favouriteExercisesQuery = createQueryFactory(
	['gym-favourites'],
	async (): Promise<FavouriteExercisesDB> => await api.get('gym/favourites').json()
);
export const favouriteExercisesQueryOptions = favouriteExercisesQuery.options;
export const favouriteExercisesRefetchOptions = favouriteExercisesQuery.refetch;
export const getFavouriteExercisesQueryKey = () => [...getRootKey(), 'gym-favourites'];

const routinesQuery = createQueryFactory(
	['gym-routines'],
	async (): Promise<RoutineDB[]> => await api.get('gym/routines').json()
);
export const routinesQueryOptions = routinesQuery.options;
export const routinesRefetchOptions = routinesQuery.refetch;
export const getRoutinesQueryKey = () => [...getRootKey(), 'gym-routines'];

const gymSummaryQuery = createQueryFactory(
	['gym-stats-summary'],
	async (): Promise<WorkoutSummaryDB> => await api.get('gym/stats/summary').json()
);
export const gymSummaryQueryOptions = gymSummaryQuery.options;
export const gymSummaryRefetchOptions = gymSummaryQuery.refetch;

const gymCalendarQuery = createQueryFactory(
	['gym-stats-calendar'],
	async (): Promise<WorkoutCalendarEntryDB[]> => await api.get('gym/stats/calendar').json()
);
export const gymCalendarQueryOptions = gymCalendarQuery.options;
export const gymCalendarRefetchOptions = gymCalendarQuery.refetch;

export const gymMusclesQueryOptions = (weeks: number) =>
	queryOptions<ExerciseFailureStatsDB[]>({
		queryKey: ['gym-stats-muscles', weeks],
		queryFn: async () => await api.get(`gym/stats/muscles?weeks=${weeks}`).json(),
		staleTime: staleTime
	});

export const gymExerciseQueryOptions = (exerciseId: string) =>
	queryOptions<ExerciseSetStatsDB[]>({
		queryKey: ['gym-stats-exercise', exerciseId],
		queryFn: async () => await api.get(`gym/stats/exercises/${encodeURIComponent(exerciseId)}`).json(),
		staleTime: staleTime
	});

export const gymUserExercisesQueryOptions = () =>
	queryOptions<UserExerciseDB[]>({
		queryKey: ['gym-stats-user-exercises'],
		queryFn: async () => await api.get('gym/stats/exercises').json(),
		staleTime: staleTime
	});

// Market queries
const marketPricesQuery = createQueryFactory(
	['market-prices'],
	async (): Promise<MarketPriceDB[]> => await api.get('market/prices').json()
);
export const marketPricesQueryOptions = marketPricesQuery.options;
export const marketPricesRefetchOptions = marketPricesQuery.refetch;

const marketInsightsQuery = createQueryFactory(
	['market-insights'],
	async (): Promise<MarketInsightDB[]> => await api.get('market/insights').json()
);
export const marketInsightsQueryOptions = marketInsightsQuery.options;
export const marketInsightsRefetchOptions = marketInsightsQuery.refetch;

export function marketPriceQueryOptions(id: string) {
	return queryOptions({
		queryKey: [...rootKey, 'market-price', id],
		queryFn: async (): Promise<MarketPriceDB> => await api.get(`market/prices/${id}`).json(),
		staleTime
	});
}

export function filteredMarketPricesQueryOptions(filter: { category?: string; item?: string }) {
	const params = new URLSearchParams();
	if (filter.category) params.set('category', filter.category);
	if (filter.item) params.set('item', filter.item);
	const qs = params.toString();
	return queryOptions({
		queryKey: [...rootKey, 'market-prices', filter.category ?? '', filter.item ?? ''],
		queryFn: async (): Promise<MarketPriceDB[]> =>
			await api.get(`market/prices${qs ? `?${qs}` : ''}`).json(),
		staleTime
	});
}

export function filteredMarketInsightsQueryOptions(category: string) {
	return queryOptions({
		queryKey: [...rootKey, 'market-insights', category],
		queryFn: async (): Promise<MarketInsightDB[]> =>
			await api.get(`market/insights?category=${encodeURIComponent(category)}`).json(),
		staleTime
	});
}

export interface MarketPriceUpsertResult {
	id: string;
	isUpdate: boolean;
}

function refetchMarketQueries() {
	queryClient.refetchQueries({ queryKey: [...rootKey, 'market-prices'] });
	queryClient.refetchQueries({ queryKey: [...rootKey, 'market-insights'] });
}

export async function createMarketPriceMutation(input: MarketPriceInput) {
	const result = await api
		.post('market/prices', {
			json: input
		})
		.json<MarketPriceUpsertResult>();
	refetchMarketQueries();
	return result;
}

export async function updateMarketPriceMutation(id: string, input: MarketPriceInput) {
	await api
		.patch(`market/prices/${id}`, {
			json: input
		})
		.json<void>();
	refetchMarketQueries();
}
