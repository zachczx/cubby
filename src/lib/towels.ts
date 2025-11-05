import { dirtyTowelDays } from './config';

export function getStatusColorFromValue(val: number): string {
	const day = 24;

	if (val > 0 && val <= 2 * day) return 'green';

	if (val > 2 * day && val <= 4 * day) return 'yellow';

	if (val > 4 * day && val <= 6 * day) return 'orange';

	if (val > dirtyTowelDays * day && val <= 999 * day) return 'red';

	return '';
}
