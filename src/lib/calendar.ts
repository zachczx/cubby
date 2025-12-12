import type { Calendar } from '@event-calendar/core';
import dayjs from 'dayjs';

function hasDateRange(record: CalendarRecord): record is VacationDB {
	return 'startDateTime' in record;
}

export function getCalendarEntries(
	data: LogsDB[] | VacationDB[] | undefined,
	title: string | undefined,
	icon?: string
) {
	if (!data || !title) return [];

	const times: Calendar.EventInput[] = [];

	for (const r of data ?? []) {
		if (hasDateRange(r)) {
			//Adding the timezone here creates a problem for mobile devices, not sure why => .tz('Asia/Singapore');
			const start = dayjs.utc(r.startDateTime);
			const end = dayjs.utc(r.endDateTime);
			times.push({
				start: start.toDate(),
				end: end.toDate(),
				title: icon ? `— ${icon} ${title}` : `— ${title}`,
				backgroundColor: 'var(--color-neutral)' // vacation color
			});
		} else {
			const t = dayjs.utc(r.time);
			times.push({
				start: t.toDate(),
				end: t.toDate(),
				title: icon ? `— ${icon} ${title}` : `— ${title}`
			});
		}
	}

	return times;
}
