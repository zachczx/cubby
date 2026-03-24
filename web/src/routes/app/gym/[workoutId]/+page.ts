import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const workoutId = params.workoutId;

	return { workoutId };
};
