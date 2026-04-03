/**
 * Personal Record detection for gym sets.
 *
 * A set is a PR if it has the heaviest weight ever for that exercise,
 * or the heaviest weight at that specific rep count.
 */

export type PrType = 'weight' | 'rep-weight' | 'both';

export interface PrResult {
	type: PrType;
	label: string;
}

/**
 * Check if a set is a PR given all historical sets for the same exercise.
 * `historicalSets` should include all sets BEFORE the current one (not including it).
 */
export function detectPr(
	set: { weightKg: number | null; reps: number | null },
	historicalSets: { weightKg: number | null; reps: number | null }[]
): PrResult | null {
	if (set.weightKg == null) return null;

	const validHistory = historicalSets.filter((s) => s.weightKg != null);
	if (validHistory.length === 0) return null;

	const maxWeight = Math.max(...validHistory.map((s) => s.weightKg!));
	const isWeightPr = set.weightKg > maxWeight;

	let isRepWeightPr = false;
	if (set.reps != null) {
		const sameRepSets = validHistory.filter((s) => s.reps === set.reps && s.weightKg != null);
		if (sameRepSets.length > 0) {
			const maxAtReps = Math.max(...sameRepSets.map((s) => s.weightKg!));
			isRepWeightPr = set.weightKg > maxAtReps;
		}
	}

	if (isWeightPr && isRepWeightPr) {
		return { type: 'both', label: 'PR' };
	} else if (isWeightPr) {
		return { type: 'weight', label: 'PR' };
	} else if (isRepWeightPr) {
		return {
			type: 'rep-weight',
			label: `${set.reps}-rep PR`
		};
	}

	return null;
}

/**
 * For a list of sets in chronological order (oldest first), compute which ones are PRs.
 * Returns a Map from set index to PrResult.
 */
export function detectAllPrs(
	sets: { weightKg: number | null; reps: number | null }[]
): Map<number, PrResult> {
	const prs = new Map<number, PrResult>();

	for (let i = 0; i < sets.length; i++) {
		const historical = sets.slice(0, i);
		const pr = detectPr(sets[i], historical);
		if (pr) prs.set(i, pr);
	}

	return prs;
}
