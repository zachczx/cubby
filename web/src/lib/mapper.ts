import type { Component } from 'svelte';
import dayjs from 'dayjs';
import FluentEmojiFlatBed from './assets/expressive-icons/FluentEmojiFlatBed.svelte';
import FluentEmojiFlatBookmarkTabs from './assets/expressive-icons/FluentEmojiFlatBookmarkTabs.svelte';
import FluentEmojiFlatLotionBottle from './assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
import FluentEmojiFlatShield from './assets/expressive-icons/FluentEmojiFlatShield.svelte';
import FluentEmojiFlatShower from './assets/expressive-icons/FluentEmojiFlatShower.svelte';
import StreamlineColorHotelLaundryFlat from './assets/expressive-icons/StreamlineColorHotelLaundryFlat.svelte';
import FluentEmojiFlatTShirt from './assets/expressive-icons/FluentEmojiFlatTShirt.svelte';
import StreamlineColorToothFlat from './assets/expressive-icons/StreamlineColorToothFlat.svelte';
import FluentEmojiFlatZzz from './assets/expressive-icons/FluentEmojiFlatZzz.svelte';
import FluentEmojiFlatElectricPlug from './assets/expressive-icons/FluentEmojiFlatElectricPlug.svelte';
import FluentEmojiFlatSyringe from './assets/expressive-icons/FluentEmojiFlatSyringe.svelte';
import FluentEmojiFlatCreditCard from './assets/expressive-icons/FluentEmojiFlatCreditCard.svelte';
import StreamlineFlexColorSubscriptionCashflowFlat from './assets/expressive-icons/StreamlineFlexColorSubscriptionCashflowFlat.svelte';
import FluentEmojiFlatDvd from './assets/expressive-icons/FluentEmojiFlatDvd.svelte';

const trackerIcons: Record<string, Component> = {
	tshirt: FluentEmojiFlatTShirt,
	bottle: FluentEmojiFlatLotionBottle,
	bed: FluentEmojiFlatBed,
	shower: FluentEmojiFlatShower,
	shield: FluentEmojiFlatShield,
	tooth: StreamlineColorToothFlat,
	washer: StreamlineColorHotelLaundryFlat,
	zzz: FluentEmojiFlatZzz,
	electricPlug: FluentEmojiFlatElectricPlug,
	syringe: FluentEmojiFlatSyringe,
	creditCard: FluentEmojiFlatCreditCard,
	subscription: StreamlineFlexColorSubscriptionCashflowFlat,
	software: FluentEmojiFlatDvd,
	bookmark: FluentEmojiFlatBookmarkTabs
};

export function getTrackerIcon(iconName: string | undefined): Component {
	if (!iconName) return FluentEmojiFlatBookmarkTabs;

	return trackerIcons[iconName] || FluentEmojiFlatBookmarkTabs;
}

export function getAllTrackerIcons() {
	return trackerIcons;
}

export function getFamilyColor(id: string | undefined, familyIds: string[]): TrackerColor {
	if (!id || !familyIds.includes(id)) return 'slate';

	const idx = familyIds.indexOf(id);

	const colors: TrackerColor[] = ['blue', 'purple', 'red', 'teal', 'orange', 'lime'];
	return colors[idx] ?? 'slate';
}

export function getColoredTrackers(trackers: TrackerDB[]): TrackerColored[] {
	const s = new Set<string>();

	for (const t of trackers) {
		if (!t.isOwner && t.familyId) {
			s.add(t.familyId);
		}
	}

	const familyIds = Array.from(s);

	const coloredTrackers: TrackerColored[] = trackers.map((tracker) => {
		if (tracker.isOwner) {
			const color = 'green';
			return { ...tracker, color };
		}

		const color = getFamilyColor(tracker.familyId, familyIds);
		return { ...tracker, color };
	});

	return coloredTrackers;
}

export function generateSubscriptionEntries(tracker: TrackerDB, userId: string): EntryDB[] {
	if (!tracker.startDate) return [];

	const subscriptionStart = tracker.startDate;
	const historicalRecords: EntryDB[] = [];
	const today = dayjs();
	let currentDateTime = dayjs(subscriptionStart);

	// Add the initial start date as a log
	historicalRecords.push({
		id: 'generated_start',
		created_at: subscriptionStart,
		updated_at: subscriptionStart,
		tracker: tracker.id,
		trackerId: tracker.id,
		performedBy: userId,
		performedAt: subscriptionStart,
		remark: '',
		interval: tracker.interval,
		intervalUnit: tracker.intervalUnit
	});

	// Generate subsequent logs based on interval
	while (
		currentDateTime.add(tracker.interval, tracker.intervalUnit).isBefore(today) ||
		currentDateTime.add(tracker.interval, tracker.intervalUnit).isSame(today, 'day')
	) {
		currentDateTime = currentDateTime.add(tracker.interval, tracker.intervalUnit);
		historicalRecords.push({
			id: `generated_${currentDateTime.toISOString()}`,
			remark: '',
			created_at: currentDateTime.toISOString(),
			updated_at: currentDateTime.toISOString(),
			tracker: tracker.id,
			trackerId: tracker.id,
			performedBy: userId,
			performedAt: currentDateTime.toISOString(),
			interval: tracker.interval,
			intervalUnit: tracker.intervalUnit
		});
	}

	// Return reversed so latest is first, similar to DB sort
	return historicalRecords.reverse();
}
