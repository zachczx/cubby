// import adapter from '@sveltejs/adapter-cloudflare';
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: { adapter: adapter() },
	compilerOptions: {
		experimental: {
			async: true
		}
	}
};

export default config;
