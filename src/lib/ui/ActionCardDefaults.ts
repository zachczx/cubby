import FluentEmojiFlatBed from '$lib/assets/expressive-icons/FluentEmojiFlatBed.svelte';
import FluentEmojiFlatLotionBottle from '$lib/assets/expressive-icons/FluentEmojiFlatLotionBottle.svelte';
import FluentEmojiFlatShield from '$lib/assets/expressive-icons/FluentEmojiFlatShield.svelte';
import FluentEmojiFlatShower from '$lib/assets/expressive-icons/FluentEmojiFlatShower.svelte';
import StreamlineColorHotelLaundryFlat from '$lib/assets/expressive-icons/StreamlineColorHotelLaundryFlat.svelte';

export const towel: ActionCardOptions = {
	collectionName: 'towel',
	title: 'Wash Towel',
	route: '/household/towel',
	icon: StreamlineColorHotelLaundryFlat,
	button: {
		text: 'Washed'
	}
};

export const spray: ActionCardOptions = {
	collectionName: 'spray',
	title: 'Spray Nose',
	route: '/personal/spray',
	icon: FluentEmojiFlatLotionBottle,
	button: {
		text: 'Sprayed'
	}
};

export const gummy: ActionCardOptions = {
	collectionName: 'gummy',
	title: 'Elderberry Gummy',
	route: '/personal/gummy',
	icon: FluentEmojiFlatShield,
	button: {
		text: 'Ate'
	}
};

export const bedsheet: ActionCardOptions = {
	collectionName: 'bedsheet',
	title: 'Bedsheet',
	route: '/household/bedsheet',
	icon: FluentEmojiFlatBed,
	button: {
		text: 'Changed'
	}
};

export const bath: ActionCardOptions = {
	collectionName: 'bath',
	title: 'Bath',
	route: '/pet/bath',
	icon: FluentEmojiFlatShower,
	button: {
		text: 'Bathed'
	}
};

export const chewable: ActionCardOptions = {
	collectionName: 'chewable',
	title: 'Chewable',
	route: '/pet/chewable',
	icon: FluentEmojiFlatShield,
	button: {
		text: 'Fed'
	}
};
