import dayjs from 'dayjs';
import { leadTimeHours, dueThresholdDays, dueThresholdWeeks } from './notification';

interface Grace {
	gracePeriodValue: number;
	gracePeriodUnit: 'hour' | 'day' | 'week' | 'hour' | undefined;
}

export function calculateStreak(
	entries: EntryDB[] | undefined,
	tracker: TrackerDB | undefined
): number {
	if (!entries || entries.length === 0 || !tracker) return 0;

	const sortedEntries = [...entries].sort((a, b) => dayjs(b.performedAt).diff(dayjs(a.performedAt)));
	const latestEntry = sortedEntries[0];

	const nextDueDate = dayjs(latestEntry.performedAt).add(tracker.interval, tracker.intervalUnit);

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

	return checkStreakLength(sortedEntries, tracker, grace);
}

function checkStreakLength(sortedEntries: EntryDB[], tracker: TrackerDB, grace: Grace): number {
	let streak = 1;

	for (let i = 0; i < sortedEntries.length - 1; i++) {
		const currentEntry = sortedEntries[i];
		const previousEntry = sortedEntries[i + 1];

		const previousNextDueDate = dayjs(previousEntry.performedAt).add(tracker.interval, tracker.intervalUnit);
		const previousGraceLimit = previousNextDueDate.add(
			grace.gracePeriodValue,
			grace.gracePeriodUnit
		);

		if (
			dayjs(currentEntry.performedAt).isSame(previousGraceLimit) ||
			dayjs(currentEntry.performedAt).isBefore(previousGraceLimit)
		) {
			streak++;
		} else {
			break;
		}
	}

	return streak;
}
