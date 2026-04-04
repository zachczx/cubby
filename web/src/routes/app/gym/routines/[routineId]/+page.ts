import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const routineId = params.routineId;

	return { routineId };
};
