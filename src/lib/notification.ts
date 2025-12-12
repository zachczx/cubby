import dayjs from 'dayjs';

export const emptyNotificationStatus: NotificationStatus = {
	show: false,
	level: null
};

export const defaultNotificationStatus: NotificationStatus = {
	show: false,
	level: 'ok'
};

const leadTimeHours = 6;
const dueThresholdDays = 1;

function createStatus(
	level: NotificationStatus['level'],
	nextDate: dayjs.Dayjs,
	show = false
): NotificationStatus {
	const baseStatus = show ? { show: true, level } : defaultNotificationStatus;
	return {
		...baseStatus,
		next: nextDate.toString()
	};
}

function getLatestRecord(data: LogsDB | LogsDB[] | undefined): LogsDB | null {
	if (!data) return null;
	return Array.isArray(data) ? data[0] : data;
}

function calculateNextDate(record: LogsDB): dayjs.Dayjs {
	return dayjs(record.time).add(record.interval, record.intervalUnit);
}

function getMonthlyStatus(record: LogsDB) {
	const nextDate = calculateNextDate(record);
	const daysRemaining = nextDate.diff(dayjs(), 'day', true);

	if (daysRemaining > dueThresholdDays) {
		return createStatus('ok', nextDate);
	}

	if (daysRemaining > 0) {
		return createStatus('due', nextDate, true);
	}

	return createStatus('overdue', nextDate, true);
}

function getDailyStatus(record: LogsDB) {
	const nextDate = calculateNextDate(record);
	const intervalHours = record.interval * 24;

	const hoursSinceLastRecord = dayjs().diff(dayjs(record.time), 'hour', true);

	const windowToNotify = intervalHours - leadTimeHours;
	if (hoursSinceLastRecord < windowToNotify) {
		return createStatus('ok', nextDate);
	}

	const isOverdue = hoursSinceLastRecord > intervalHours;
	return createStatus(isOverdue ? 'overdue' : 'due', nextDate, true);
}

export function getTrackerStatus(data: LogsDB[] | undefined): NotificationStatus {
	if (!data) return emptyNotificationStatus;

	const record = getLatestRecord(data);

	if (!record) {
		return emptyNotificationStatus;
	}

	return record.intervalUnit === 'month' ? getMonthlyStatus(record) : getDailyStatus(record);
}
