import type { Component } from 'svelte';
import FluentEmojiFlatBed from './assets/expressive-icons/FluentEmojiFlatBed.svelte';
import FluentEmojiFlatBookmarkTabs from './assets/expressive-icons/FluentEmojiFlatBookmarkTabs.svelte';
import FluentEmojiFlatLotionBottle from './assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
import FluentEmojiFlatShield from './assets/expressive-icons/FluentEmojiFlatShield.svelte';
import FluentEmojiFlatShower from './assets/expressive-icons/FluentEmojiFlatShower.svelte';
import StreamlineColorHotelLaundryFlat from './assets/expressive-icons/StreamlineColorHotelLaundryFlat.svelte';

export function getTrackerIcon(trackerName: string): Component {
	switch (trackerName) {
		case 'towel':
			return StreamlineColorHotelLaundryFlat;

		case 'spray':
			return FluentEmojiFlatLotionBottle;

		case 'gummy':
			return FluentEmojiFlatShield;

		case 'bedsheet':
			return FluentEmojiFlatBed;

		case 'bath':
			return FluentEmojiFlatShower;

		case 'chewable':
			return FluentEmojiFlatShield;

		default:
			return FluentEmojiFlatBookmarkTabs;
	}
}
