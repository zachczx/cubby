export function cleanEmail(email: string | undefined): string {
	if (!email) return '';

	const name = email.split('@')?.[0];
	const maxLength = name.length > 11 ? 11 : name.length;
	const clean = name.slice(0, maxLength);

	return clean;
}
