import dayjs from 'dayjs';
import { dirtyTowelDays } from './config';
import type { CreateQueryResult } from '@tanstack/svelte-query';

export function getNotificationStatus(
	query: CreateQueryResult<SprayDB[], Error> | CreateQueryResult<TowelDB[], Error>
): NotificationStatus {
	if (!query?.isSuccess || query.data.length === 0) return defaultNotificationStatus;

	const lastRecord = query.data?.[0] ?? null;
	if (!lastRecord) return defaultNotificationStatus;

	const now = dayjs();
	const leadTimeHours = 12;

	let intervalHours = 0;

	if (lastRecord.collectionName === 'spray') {
		intervalHours = lastRecord.daysToNext * 24;
	} else if (lastRecord.collectionName === 'towel') {
		intervalHours = dirtyTowelDays * 24;
	}

	const hoursSinceLastRecord = now.diff(dayjs(lastRecord.time), 'hour', true);
	if (hoursSinceLastRecord >= intervalHours - leadTimeHours) {
		const overdue = hoursSinceLastRecord > intervalHours;
		return {
			show: true,
			level: overdue ? 'overdue' : 'due'
		};
	}

	return defaultNotificationStatus;
}

const defaultNotificationStatus: NotificationStatus = {
	show: false,
	level: 'ok'
};
