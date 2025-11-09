import type { Calendar } from '@event-calendar/core';
import type { CreateQueryResult } from '@tanstack/svelte-query';
import dayjs from 'dayjs';

function hasDateRange(record: CalendarRecord): record is VacationDB {
	return 'startDateTime' in record;
}

export function getCalendarEntries(
	query: CreateQueryResult<CalendarRecord[], Error>,
	title: string,
	icon?: string
) {
	const times: Calendar.EventInput[] = [];
	if (query.isSuccess) {
		for (const r of query.data ?? []) {
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
	}
	return times;
}
