export const trackerDefaults = [
	{
		name: 'bedsheet',
		display: 'Bedsheet',
		interval: 14,
		intervalUnit: 'day',
		category: 'household',
		actionLabel: 'Changed',
		pinned: false,
		show: true
	},
	{
		name: 'bath',
		display: 'Bath',
		interval: 14,
		intervalUnit: 'day',
		category: 'pet',
		actionLabel: 'Bathed',
		pinned: false,
		show: true
	},
	{
		name: 'chewable',
		display: 'Nexgard',
		interval: 1,
		intervalUnit: 'month',
		category: 'pet',
		actionLabel: 'Fed',
		pinned: false,
		show: true
	},
	{
		name: 'gummy',
		display: 'Gummy',
		interval: 2,
		intervalUnit: 'day',
		category: 'personal',
		actionLabel: 'Ate',
		pinned: true,
		show: true
	},
	{
		name: 'spray',
		display: 'Nasal Spray',
		interval: 3,
		intervalUnit: 'day',
		category: 'personal',
		actionLabel: 'Sprayed',
		pinned: true,
		show: true
	},
	{
		name: 'towel',
		display: 'Towel Wash',
		interval: 5,
		intervalUnit: 'day',
		category: 'household',
		actionLabel: 'Washed',
		pinned: true,
		show: true
	}
] as const;
