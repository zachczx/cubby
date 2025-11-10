import { queryOptions, type RefetchQueryFilters } from '@tanstack/svelte-query';
import { pb } from './pb';

const staleTime = 5 * 60 * 1000;

export function createTowelQueryOptions() {
	return queryOptions({
		queryKey: ['towels', pb.authStore?.record?.id],
		queryFn: async (): Promise<TowelDB[]> =>
			await pb.collection('towel').getFullList({ sort: '-time' }),
		staleTime: staleTime
	});
}

export function createSprayQueryOptions() {
	return queryOptions({
		queryKey: ['sprays', pb.authStore?.record?.id],
		queryFn: async (): Promise<SprayDB[]> =>
			await pb.collection('spray').getFullList({ sort: '-time' }),
		staleTime: staleTime
	});
}

export function createUserQueryOptions() {
	return queryOptions({
		queryKey: ['users', pb.authStore?.record?.id],
		queryFn: async (): Promise<UserDB> =>
			await pb.collection('users').getOne(String(pb.authStore?.record?.id)),
		staleTime: staleTime
	});
}

export function createVacationQueryOptions() {
	return queryOptions({
		queryKey: ['vacations', pb.authStore?.record?.id],
		queryFn: async (): Promise<VacationDB[]> =>
			await pb.collection('vacation').getFullList({ sort: '-startDateTime' }),
		staleTime: staleTime
	});
}

export function createGummyQueryOptions() {
	return queryOptions({
		queryKey: ['gummies', pb.authStore?.record?.id],
		queryFn: async (): Promise<GummyDB[]> =>
			await pb.collection('gummy').getFullList({ sort: '-time' }),
		staleTime: staleTime
	});
}

export function createDoggoChewableQueryOptions() {
	return queryOptions({
		queryKey: ['doggoChewables', pb.authStore?.record?.id],
		queryFn: async (): Promise<DoggoChewableDB[]> =>
			await pb.collection('doggoChewable').getFullList({ sort: '-time' }),
		staleTime: staleTime
	});
}

export function createDoggoBathQueryOptions() {
	return queryOptions({
		queryKey: ['doggoBaths', pb.authStore?.record?.id],
		queryFn: async (): Promise<DoggoBathDB[]> =>
			await pb.collection('doggoBath').getFullList({ sort: '-time' }),
		staleTime: staleTime
	});
}

export function createTowelRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['towels', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createSprayRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['sprays', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createGummyRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['gummies', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createDoggoChewableRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['doggoChewables', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createDoggoBathRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['doggoBaths', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createUserRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['users', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}

export function createVacationRefetchOptions(): RefetchQueryFilters {
	return {
		queryKey: ['vacations', pb.authStore?.record?.id],
		type: 'active',
		exact: true
	};
}
