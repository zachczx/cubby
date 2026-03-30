export function capitalize(text: string) {
	if (text.length == 0) return text;
	return text[0].toUpperCase() + text.substr(1);
}

export function titleCase(text: string) {
	return text.split(' ').map(capitalize).join(' ');
}
