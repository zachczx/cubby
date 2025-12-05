import type { Component } from 'svelte';
import FluentEmojiFlatBed from './assets/expressive-icons/FluentEmojiFlatBed.svelte';
import FluentEmojiFlatBookmarkTabs from './assets/expressive-icons/FluentEmojiFlatBookmarkTabs.svelte';
import FluentEmojiFlatLotionBottle from './assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
import FluentEmojiFlatShield from './assets/expressive-icons/FluentEmojiFlatShield.svelte';
import FluentEmojiFlatShower from './assets/expressive-icons/FluentEmojiFlatShower.svelte';
import StreamlineColorHotelLaundryFlat from './assets/expressive-icons/StreamlineColorHotelLaundryFlat.svelte';
import { pb } from './pb';

const trackerIcons: Record<string, Component> = {
	towel: StreamlineColorHotelLaundryFlat,
	spray: FluentEmojiFlatLotionBottle,
	gummy: FluentEmojiFlatShield,
	bedsheet: FluentEmojiFlatBed,
	bath: FluentEmojiFlatShower,
	chewable: FluentEmojiFlatShield
};

export function getTrackerIcon(trackerName: string): Component {
	return trackerIcons[trackerName] || FluentEmojiFlatBookmarkTabs;
}

export function getFamilyColor(id: string | undefined, familyIds: string[]): string {
	if (!id || !familyIds.includes(id)) return '';

	const idx = familyIds.indexOf(id);

	switch (idx) {
		case 0:
			return 'blue';

		case 1:
			return 'purple';

		case 2:
			return 'red';

		case 3:
			return 'teal';

		case 4:
			return 'orange';

		case 5:
			return 'lime';

		default:
			return 'slate';
	}
}

export function getColoredTrackers(trackers: TrackerDB[]): TrackerColored[] {
	const s = new Set<string>();

	for (const t of trackers) {
		const owner = t.expand?.family?.owner;
		const familyId = t.expand?.family?.id;

		if (owner !== pb.authStore.record?.id && familyId) {
			s.add(familyId);
		}
	}

	const familyIds = Array.from(s);

	const coloredTrackers: TrackerColored[] = trackers.map((tracker) => {
		if (tracker.expand?.family?.owner === pb.authStore.record?.id) {
			const color = 'green';
			return { ...tracker, color };
		}

		const color = getFamilyColor(tracker.expand?.family?.id, familyIds);
		return { ...tracker, color };
	});

	return coloredTrackers;
}
