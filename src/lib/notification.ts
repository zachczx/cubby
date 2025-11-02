import dayjs from 'dayjs';
import { dirtyTowelDays } from './config';
import type { CreateQueryResult } from '@tanstack/svelte-query';

export function availableTowelNotification(towels: CreateQueryResult<TowelDB[], Error>): boolean {
	if (!towels?.isSuccess || towels.data.length === 0) return false;

	const lastWash = towels.data?.[0] ?? null;

	if (!lastWash) {
		return false;
	}

	const now = dayjs();
	const leadTimeHours = 12;
	const intervalHours = dirtyTowelDays * 24;

	const hoursSinceLastWash = now.diff(dayjs(lastWash.time), 'hour', true);

	if (hoursSinceLastWash >= intervalHours - leadTimeHours) {
		return true;
	}

	return false;
}

export function availableSprayNotification(sprays: CreateQueryResult<SprayDB[], Error>): boolean {
	if (!sprays?.isSuccess || sprays.data.length === 0) return false;

	const lastSpray = sprays.data?.[0] ?? null;

	if (!lastSpray) {
		return false;
	}

	const now = dayjs();
	const leadTimeHours = 12;
	const intervalHours = lastSpray.daysToNext * 24;

	const hoursSinceLastSpray = now.diff(dayjs(lastSpray.time), 'hour', true);
	if (hoursSinceLastSpray >= intervalHours - leadTimeHours) {
		return true;
	}

	return false;
}
