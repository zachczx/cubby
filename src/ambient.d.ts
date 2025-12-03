type Characters = 'frankenstein' | 'robot' | 'furnando';

interface UserDB {
	collectionId: string;
	collectionName: 'users';
	id: string;
	email: string;
	emailVisibility: boolean;
	verified: boolean;
	name: string;
	avatar: string;
	sound: boolean;
	created: string;
	updated: string;
}

interface VacationDB {
	collectionId: string;
	collectionName: 'vacation';
	id: string;
	created: string;
	updated: string;
	user: string;
	startDateTime: string;
	endDateTime: string;
}

type IntervalUnit = 'day' | 'month' | 'year';

interface LogsDB {
	collectionId: string;
	collectionName: 'logs';
	id: string;
	user: string;
	tracker: string;
	interval: number;
	intervalUnit: IntervalUnit;
	time: string;
	oldId?: string | null;
	created: string;
	updated: string;
	expand?: { tracker?: TrackerDB };
}

interface LogsRecord extends LogsDB {
	gap: number;
}

interface TrackerColored extends TrackerDB {
	color: string;
}

interface TrackerDB {
	id: string;
	user: string;
	family: string;
	name: string;
	display: string;
	interval: number;
	intervalUnit: IntervalUnit;
	category: string;
	actionLabel: string;
	pinned: boolean;
	show: boolean;
	created: string;
	updated: string;
	expand?: { family?: FamilyDB };
}

interface FamilyDB {
	id: string;
	name: string;
	members: string[];
	owner: string;
	created: string;
	updated: string;
	expand?: { members?: UserDB[]; owner?: UserDB };
}

interface InviteDB {
	id: string;
	family: string;
	code: string;
	familyNameSnapshot: string;
	ownerEmailSnapshot: string;
	status?: 'pending' | 'completed' | null;
	created: string;
	updated: string;
	expand?: { family?: FamilyDB };
}

interface NotificationStatus {
	show: boolean;
	level: 'ok' | 'due' | 'overdue' | null;
	next?: string;
	label?: string | undefined;
	href?: string | undefined;
}

type ButtonState = 'default' | 'loading' | 'success' | 'error';

// type CollectionName =
// 	| 'towel'
// 	| 'spray'
// 	| 'gummy'
// 	| 'user'
// 	| 'bedsheet'
// 	| 'vacation'
// 	| 'bath'
// 	| 'chewable';

// Types/Interfaces for TrackerPage component abstraction

interface TrackerPageOptions {
	tracker: TrackerDB | undefined;
	collectionName: string | undefined;
	labels: {
		pageTitle: string | undefined;
		ctaButtonText: string | undefined;
		noun: string | undefined;
	};
	calculateGaps?: (records: LogsDB[], vacations: VacationDB[]) => LogsRecord[];
}

type Collections = LogsDB | VacationDB;

type CalendarRecord = Collections;

interface ActionCardOptions {
	title: string | undefined;
	size?: 'compact' | 'default' | 'list';
	tracker: TrackerDB | TrackerColored;
	icon: Component;
	route: string;
	lastChild?: boolean;
	button: {
		status?: ButtonState;
		text: string | undefined;
	};
}
