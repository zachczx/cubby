import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
	appId: 'dev.cubby.app',
	appName: 'Cubby',
	webDir: 'build',
	server: {
		androidScheme: 'http',
		cleartext: true
	}
};

export default config;
