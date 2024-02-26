import { defineConfig } from 'vite';

export default defineConfig({
	test: {
		clearMocks: true,
		globals: true,
	},
});
