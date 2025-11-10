import dayjs from 'dayjs';

/**
 * Function for TrackerPage to calculate towel gaps
 */
export function calculateTowelGapsWithVacation(records: TrackerDB[], vacations: VacationDB[]) {
	return records.map((record, i, allRecords) => {
		const nextRecord = allRecords[i + 1];
		if (!nextRecord) {
			return { ...record, gap: 0 };
		}

		// 1. Calculate the raw gap
		const rawGap = dayjs(record.time).diff(nextRecord.time, 'day', true);

		// 2. Calculate vacation days within that gap
		const vacationDays = calculateVacationOverlap(
			record.time,
			nextRecord.time,
			vacations ?? [],
			'day'
		);

		// 3. Subtract vacation time from the raw gap
		const finalGap = rawGap - vacationDays;

		// Ensure gap is never negative
		return { ...record, gap: Math.max(0, finalGap) };
	});
}

/**
 * Calculates the total vacation time (in specified units)
 * that overlaps with a given time interval.
 */
export function calculateVacationOverlap(
	intervalEnd: string | dayjs.Dayjs,
	intervalStart: string | dayjs.Dayjs,
	vacations: VacationDB[],
	unit: dayjs.UnitType = 'day' // Default to 'day'
): number {
	let totalVacationTime = 0;
	const gapStart = dayjs(intervalStart);
	const gapEnd = dayjs(intervalEnd);

	// Ensure start is actually before end
	if (gapStart.isAfter(gapEnd)) return 0;

	for (const vacation of vacations) {
		const vacationStart = dayjs(vacation.startDateTime);
		const vacationEnd = dayjs(vacation.endDateTime);

		// Find the latest start time (either gap start or vacation start)
		const overlapStart = gapStart.isAfter(vacationStart) ? gapStart : vacationStart;

		// Find the earliest end time (either gap end or vacation end)
		const overlapEnd = gapEnd.isBefore(vacationEnd) ? gapEnd : vacationEnd;

		// If overlapEnd is after overlapStart, there is a valid overlap
		if (overlapEnd.isAfter(overlapStart)) {
			totalVacationTime += overlapEnd.diff(overlapStart, unit, true);
		}
	}
	return totalVacationTime;
}
