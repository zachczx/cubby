import type { Calendar } from '@event-calendar/core';
import dayjs from 'dayjs';

function hasDateRange(record: CalendarRecord): record is VacationDB {
	return 'startDateTime' in record;
}

export function getCalendarEntries(data: LogsDB[] | VacationDB[] | undefined, icon?: string) {
	if (!data) return [];

	const times: Calendar.EventInput[] = [];

	for (const r of data ?? []) {
		if (hasDateRange(r)) {
			//Adding the timezone here creates a problem for mobile devices, not sure why => .tz('Asia/Singapore');
			const start = dayjs.utc(r.startDateTime);
			const end = dayjs.utc(r.endDateTime);
			times.push({
				start: start.toDate(),
				/**
				 * Add 1 second duration to ensure events display in calendar grid.
				 * Zero-duration events at midnight (12:00am) fail to render due to day boundary ambiguity.
				 */
				end: end.add(1, 'second').toDate(),
				title: icon ? `— ${icon}` : ``,
				backgroundColor: 'var(--color-neutral)' // vacation color
			});
		} else {
			const t = dayjs.utc(r.time);
			times.push({
				start: t.toDate(),
				/**
				 * Ditto above
				 */
				end: t.add(1, 'second').toDate(),
				title: icon ? `— ${icon}` : ``
			});
		}
	}

	return times;
}
