import dayjs from 'dayjs';

const defaultNotificationStatus: NotificationStatus = {
	show: false,
	level: 'ok'
};

export function getNotificationStatus(query: Query): NotificationStatus {
	if (!query?.isSuccess || query.data?.length === 0) return defaultNotificationStatus;

	const lastRecord = query.data?.[0] ?? null;
	if (!lastRecord) return defaultNotificationStatus;

	if (lastRecord.collectionName === 'doggoChewable') {
		const nextDate = dayjs(lastRecord.time).add(lastRecord.monthsToNext, 'month');
		const daysTillNextDate = nextDate.diff(dayjs(), 'day', true);

		if (daysTillNextDate > 1) {
			return defaultNotificationStatus;
		} else if (daysTillNextDate <= 1 && daysTillNextDate > 0) {
			return { show: true, level: 'due' };
		} else {
			return { show: true, level: 'overdue' };
		}
	}

	const now = dayjs();
	const leadTimeHours = 6;

	const intervalHours = lastRecord.daysToNext * 24;

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

export function getNotificationCount(
	sprayNotification: NotificationStatus,
	towelNotification: NotificationStatus,
	gummyNotification: NotificationStatus
): number {
	let count = 0;

	if (sprayNotification.show) {
		count += 1;
	}

	if (towelNotification.show) {
		count += 1;
	}

	if (gummyNotification.show) {
		count += 1;
	}

	return count;
}
