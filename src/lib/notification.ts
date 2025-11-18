import type { CreateQueryResult } from '@tanstack/svelte-query';
import dayjs from 'dayjs';
import { trackerIdToName } from './queries';

export const defaultNotificationStatus: NotificationStatus = {
	show: false,
	level: 'ok'
};

export function getNotificationStatus(
	query: CreateQueryResult<LogsDB, Error> | CreateQueryResult<LogsDB[], Error>
): NotificationStatus {
	if (!query?.isSuccess) return defaultNotificationStatus;

	let lastRecord;

	if (Array.isArray(query.data)) {
		lastRecord = query.data[0] ?? null;
	} else {
		lastRecord = query.data ?? null;
	}

	if (!lastRecord) return defaultNotificationStatus;

	const notifDetails = getNotificationLabel(lastRecord);

	if (lastRecord.intervalUnit === 'month') {
		const nextDate = dayjs(lastRecord.time).add(lastRecord.interval, 'month');
		const daysTillNextDate = nextDate.diff(dayjs(), 'day', true);

		if (daysTillNextDate > 1) {
			return { ...defaultNotificationStatus, next: nextDate.toString() };
		} else if (daysTillNextDate <= 1 && daysTillNextDate > 0) {
			return {
				show: true,
				level: 'due',
				next: nextDate.toString(),
				...notifDetails
			};
		} else {
			return {
				show: true,
				level: 'overdue',
				next: nextDate.toString(),
				...notifDetails
			};
		}
	} else {
		const now = dayjs();
		const leadTimeHours = 6;

		const intervalHours = lastRecord.interval * 24;

		const nextDate = dayjs(lastRecord.time).add(lastRecord.interval, 'day');

		const hoursSinceLastRecord = now.diff(dayjs(lastRecord.time), 'hour', true);
		if (hoursSinceLastRecord >= intervalHours - leadTimeHours) {
			const overdue = hoursSinceLastRecord > intervalHours;
			return {
				show: true,
				level: overdue ? 'overdue' : 'due',
				next: nextDate.toString(),
				...notifDetails
			};
		}

		return {
			...defaultNotificationStatus,
			next: nextDate.toString(),
			...notifDetails
		};
	}
}

function getNotificationLabel(record: LogsDB) {
	switch (trackerIdToName(record.tracker)) {
		case 'spray':
			return { label: 'Spray your nose!', href: '/personal/spray' };
		case 'towel':
			return { label: 'Wash your towel!', href: '/household/towel' };
		case 'gummy':
			return { label: 'Eat your gummy!', href: '/personal/gummy' };
		case 'bedsheet':
			return { label: 'Change your bedsheet!', href: '/household/bedsheet' };
		case 'doggoBath':
			return { label: 'Bathe your doggo!', href: '/pet/bath' };
		case 'doggoChewable':
			return { label: 'Feed your doggo chewable!', href: '/pet/chewable' };
	}
}
