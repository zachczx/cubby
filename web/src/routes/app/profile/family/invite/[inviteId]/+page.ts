import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const inviteId = params.inviteId;

	return { inviteId };
};
