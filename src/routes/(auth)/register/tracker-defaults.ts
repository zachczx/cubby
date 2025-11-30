export const trackerDefaults = [
	{
		name: 'bedsheet',
		display: 'Bedsheet',
		interval: 14,
		intervalUnit: 'day',
		category: 'household',
		actionLabel: 'Changed',
		pinned: false
	},
	{
		name: 'bath',
		display: 'Bath',
		interval: 14,
		intervalUnit: 'day',
		category: 'pet',
		actionLabel: 'Bathed',
		pinned: false
	},
	{
		name: 'chewable',
		display: 'Nexgard',
		interval: 1,
		intervalUnit: 'month',
		category: 'pet',
		actionLabel: 'Fed',
		pinned: false
	},
	{
		name: 'gummy',
		display: 'Gummy',
		interval: 2,
		intervalUnit: 'day',
		category: 'personal',
		actionLabel: 'Ate',
		pinned: true
	},
	{
		name: 'spray',
		display: 'Nasal Spray',
		interval: 3,
		intervalUnit: 'day',
		category: 'personal',
		actionLabel: 'Sprayed',
		pinned: true
	},
	{
		name: 'towel',
		display: 'Towel Wash',
		interval: 5,
		intervalUnit: 'day',
		category: 'household',
		actionLabel: 'Washed',
		pinned: true
	}
] as const;
