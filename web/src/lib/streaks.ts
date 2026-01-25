import dayjs from 'dayjs';
import { leadTimeHours, dueThresholdDays, dueThresholdWeeks } from './notification';

interface Grace {
	gracePeriodValue: number;
	gracePeriodUnit: 'hour' | 'day' | 'week' | 'hour' | undefined;
}

export function calculateStreak(
	logs: LogsDB[] | undefined,
	tracker: TrackerDB | undefined
): number {
	if (!logs || logs.length === 0 || !tracker) return 0;

	const sortedLogs = [...logs].sort((a, b) => dayjs(b.time).diff(dayjs(a.time)));
	const latestLog = sortedLogs[0];

	const nextDueDate = dayjs(latestLog.time).add(tracker.interval, tracker.intervalUnit);

	const grace = {} as Grace;

	if (tracker.intervalUnit === 'day') {
		grace.gracePeriodValue = leadTimeHours;
		grace.gracePeriodUnit = 'hour';
	} else if (tracker.intervalUnit === 'month') {
		grace.gracePeriodValue = dueThresholdDays;
		grace.gracePeriodUnit = 'day';
	} else if (tracker.intervalUnit === 'year') {
		grace.gracePeriodValue = dueThresholdWeeks;
		grace.gracePeriodUnit = 'week';
	}

	const gracePeriodLimit = nextDueDate.add(grace.gracePeriodValue, grace.gracePeriodUnit);

	if (dayjs().isAfter(gracePeriodLimit)) {
		return 0;
	}

	return checkStreakLength(sortedLogs, tracker, grace);
}

function checkStreakLength(sortedLogs: LogsDB[], tracker: TrackerDB, grace: Grace): number {
	let streak = 1;

	for (let i = 0; i < sortedLogs.length - 1; i++) {
		const currentLog = sortedLogs[i];
		const previousLog = sortedLogs[i + 1];

		const previousNextDueDate = dayjs(previousLog.time).add(tracker.interval, tracker.intervalUnit);
		const previousGraceLimit = previousNextDueDate.add(
			grace.gracePeriodValue,
			grace.gracePeriodUnit
		);

		if (
			dayjs(currentLog.time).isSame(previousGraceLimit) ||
			dayjs(currentLog.time).isBefore(previousGraceLimit)
		) {
			streak++;
		} else {
			break;
		}
	}

	return streak;
}
