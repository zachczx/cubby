import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { enhancedImages } from '@sveltejs/enhanced-img';
import { defineConfig } from 'vite';
import { visualizer } from 'rollup-plugin-visualizer';

export default defineConfig({
	plugins: [
		enhancedImages(),
		tailwindcss(),
		sveltekit(),
		visualizer({
			emitFile: true, // or false if you just want to open in browser
			filename: 'bundle-analysis.html',
			open: true // opens the visualization in the browser after build
		})
	]
});
