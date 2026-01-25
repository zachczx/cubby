export const trackerDefaults = [
	{
		name: 'bedsheet',
		display: 'Bedsheet',
		interval: 14,
		intervalUnit: 'day',
		category: 'household',
		kind: 'task',
		actionLabel: 'Changed',
		pinned: false,
		show: true
	},
	{
		name: 'petBath',
		display: 'petBath',
		interval: 14,
		intervalUnit: 'day',
		category: 'pet',
		kind: 'task',
		actionLabel: 'Bathed',
		pinned: false,
		show: true
	},
	{
		name: 'petChewable',
		display: 'Nexgard',
		interval: 1,
		intervalUnit: 'month',
		category: 'pet',
		kind: 'task',
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
		kind: 'task',
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
		kind: 'task',
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
		kind: 'task',
		actionLabel: 'Washed',
		pinned: true,
		show: true
	}
] as const;
