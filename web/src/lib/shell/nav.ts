import { router } from '$lib/routes';

export const topLevelRoutes = {
	animation: [
		{
			id: 'home',
			href: '/app',
			label: 'Dashboard',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:home'
		},
		{
			id: 'tasks',
			href: '/app/trackers',
			label: 'Tasks',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:checklist'
		},
		{
			id: 'gym',
			href: router.app('gym'),
			label: 'Gym',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:exercise'
		},
		{
			id: 'tools',
			href: router.app('count'),
			label: 'Tools',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:widgets'
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
