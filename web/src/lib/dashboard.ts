import dayjs from 'dayjs';
import { getTrackerStatus } from './notification';
import { calculateStreak } from './streaks';

export interface ClassifiedTracker {
	trackerName: string;
	trackerData: TrackerColored;
	entries: EntryDB[];
	notification: NotificationStatus;
	streak: number;
}

export function getDashboardTasks(
	trackers: TrackerColored[],
	entries: EntryDB[]
): ClassifiedTracker[] {
	const data: ClassifiedTracker[] = [];

	for (const t of trackers) {
		const entryData = entries.filter((entry) => t.id === entry.trackerId);
		const notification = getTrackerStatus(entryData);
		const streak = calculateStreak(entryData, t);

		// Only show overdue or due within 24 hours
		if (notification.level === 'overdue') {
			data.push({ trackerName: t.name, trackerData: t, entries: entryData, notification, streak });
			continue;
		}

		if (notification.level === 'due' && notification.next) {
			const hoursUntilDue = dayjs(notification.next).diff(dayjs(), 'hour', true);
			if (hoursUntilDue <= 24) {
				data.push({
					trackerName: t.name,
					trackerData: t,
					entries: entryData,
					notification,
					streak
				});
			}
		}
	}

	return data;
}

export function getActiveWorkout(workouts: WorkoutDB[] | undefined): WorkoutDB | undefined {
	if (!workouts || workouts.length === 0) return undefined;

	const today = dayjs().startOf('day');
	return workouts.find((w) => dayjs(w.startTime).isAfter(today));
}
