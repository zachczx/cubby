import MaterialSymbolsHome from '$lib/assets/svg/MaterialSymbolsHome.svelte';
import MaterialSymbolsLaundry from '$lib/assets/svg/MaterialSymbolsLaundry.svelte';
import MaterialSymbolsPerson from '$lib/assets/svg/MaterialSymbolsPerson.svelte';
import MaterialSymbolsPets from '$lib/assets/svg/MaterialSymbolsPets.svelte';

export const topLevelRoutes = {
	animation: [
		{
			id: 'home',
			href: '/app',
			label: 'Home',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsHome
		},
		{
			id: 'personal',
			href: '/app/personal',
			label: 'Personal',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsPerson
		},
		{
			id: 'household',
			href: '/app/household',
			label: 'Household',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsLaundry
		},
		{
			id: 'pet',
			href: '/app/pet',
			label: 'Pet',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsPets
		}
	],
	topNavAnimation: [
		{
			id: 'profile',
			href: '/app/profile',
			label: 'Profile',
			desktopNav: false,
			mobileNav: false
		}
	],
	noAnimation: [
		{
			id: 'login',
			href: '/login',
			label: 'Login',
			desktopNav: false,
			mobileNav: false
		},
		{
			id: 'register',
			href: '/register',
			label: 'Register',
			desktopNav: false,
			mobileNav: false
		}
	]
} as const;
