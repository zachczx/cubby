// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {

		// }

		declare module '*?enhanced' {
			import type { Picture } from 'vite-imagetools';

			const value: Picture;
			export default value;
		}

		declare module '*&enhanced' {
			import type { Picture } from 'vite-imagetools';

			const value: Picture;
			export default value;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
