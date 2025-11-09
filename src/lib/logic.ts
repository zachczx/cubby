export const dirtyTowelDays = 6;

export function getTowelStatusColor(hoursDiff: number): string {
	const day = 24;

	if (hoursDiff > 0 && hoursDiff <= 2 * day) return 'green';

	if (hoursDiff > 2 * day && hoursDiff <= 4 * day) return 'yellow';

	if (hoursDiff > 4 * day && hoursDiff <= 6 * day) return 'orange';

	if (hoursDiff > dirtyTowelDays * day && hoursDiff <= 999 * day) return 'red';

	return '';
}

export function getSprayStatusColor(hoursDiff: number, daysToNext: number | undefined): string {
	if (!daysToNext) return '';

	if (hoursDiff === 0) return '';

	const day = 24;

	if (hoursDiff > 0 && hoursDiff <= daysToNext * day) return 'green';

	if (hoursDiff > 2 * day && hoursDiff <= 4 * day) return 'yellow';

	if (hoursDiff > 4 * day && hoursDiff <= 6 * day) return 'orange';

	if (hoursDiff > daysToNext * day && hoursDiff <= 999 * day) return 'red';

	return '';
}

export function getGummyStatusColor(hoursDiff: number, daysToNext: number | undefined): string {
	if (!daysToNext) return '';

	if (hoursDiff === 0) return 'empty';

	const day = 24;

	if (hoursDiff > 0 && hoursDiff <= daysToNext * day) return 'green';

	if (hoursDiff > 2 * day && hoursDiff <= 4 * day) return 'yellow';

	if (hoursDiff > 4 * day && hoursDiff <= 6 * day) return 'orange';

	if (hoursDiff > daysToNext * day && hoursDiff <= 999 * day) return 'red';

	return '';
}

export function getDoggoChewableStatusColor(
	daysDiff: number,
	monthsToNext: number | undefined
): string {
	if (!monthsToNext) return '';

	if (daysDiff === 0) return 'empty';

	if (daysDiff < 0) return 'green';

	if (daysDiff > 0 && daysDiff < 1) return 'orange';

	if (daysDiff >= 1) return 'red';

	return '';
}

export function getDoggoBathStatusColor(hoursDiff: number, daysToNext: number | undefined): string {
	if (!daysToNext) return '';

	if (hoursDiff === 0) return 'empty';

	if (hoursDiff < 0) return 'green';

	if (hoursDiff > 0 && hoursDiff < 1) return 'orange';

	if (hoursDiff >= 1) return 'red';

	return '';
}
