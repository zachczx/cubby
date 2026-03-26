import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const exerciseId = params.exerciseId;
	return { exerciseId };
};