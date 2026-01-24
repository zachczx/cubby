export const topLevelRoutes = {
	animation: [
		{
			id: 'home',
			href: '/app',
			label: 'Home',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:home'
		},
		{
			id: 'personal',
			href: '/app/personal',
			label: 'Personal',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:person'
		},
		{
			id: 'household',
			href: '/app/household',
			label: 'Household',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:laundry'
		},
		{
			id: 'pet',
			href: '/app/pet',
			label: 'Pet',
			desktopNav: true,
			mobileNav: true,
			icon: 'material-symbols:pets'
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
