import dayjs from 'dayjs';

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
