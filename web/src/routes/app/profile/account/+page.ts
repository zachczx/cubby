export const load = async ({ url }) => {
	const onboarding = url.searchParams.get('onboarding');

	return { onboarding };
};
