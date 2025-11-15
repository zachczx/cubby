import MaterialSymbolsCleaningServices from '$lib/assets/svg/MaterialSymbolsCleaningServices.svelte';
import MaterialSymbolsHealthAndSafety from '$lib/assets/svg/MaterialSymbolsHealthAndSafety.svelte';
import MaterialSymbolsHome from '$lib/assets/svg/MaterialSymbolsHome.svelte';
import MaterialSymbolsPets from '$lib/assets/svg/MaterialSymbolsPets.svelte';

export const topLevelRoutes = {
	animation: [
		{
			id: 'home',
			href: '/',
			label: 'Home',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsHome
		},
		{
			id: 'personal',
			href: '/personal',
			label: 'Personal',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsHealthAndSafety
		},
		{
			id: 'household',
			href: '/household',
			label: 'Household',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsCleaningServices
		},
		{
			id: 'pet',
			href: '/pet',
			label: 'Pet',
			desktopNav: true,
			mobileNav: true,
			icon: MaterialSymbolsPets
		}
	],
	noAnimation: [
		{
			id: 'profile',
			href: '/profile',
			label: 'Profile',
			desktopNav: false,
			mobileNav: false
		},
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
