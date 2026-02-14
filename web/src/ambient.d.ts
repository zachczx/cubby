type Characters = 'frankenstein' | 'robot' | 'furnando';

interface UserDB {
	id: string;
	email: string;
	name: string;
	soundOn: boolean;
	taskLookaheadDays: number;
	created_at: string;
	updated_at: string;
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

interface EntryDB {
	id: string;
	trackerId: string;
	performedBy: string;
	tracker: string;
	interval: number;
	intervalUnit: IntervalUnit;
	performedAt: string;
	remark: string;
	created_at: string;
	updated_at: string;
	expand?: { tracker?: TrackerDB };
}

interface EntryRecord extends EntryDB {
	gap: number;
}

type TrackerColor = 'blue' | 'purple' | 'red' | 'teal' | 'orange' | 'lime' | 'slate' | 'green';

interface TrackerColored extends TrackerDB {
	color: TrackerColor;
}

type Kind = 'task' | 'subscription';

type TrackerIcon =
	| 'tshirt'
	| 'bottle'
	| 'bed'
	| 'shower'
	| 'shield'
	| 'tooth'
	| 'washer'
	| 'zzz'
	| 'electricPlug'
	| 'syringe'
	| 'creditCard'
	| 'subscription'
	| 'software'
	| 'bookmark';

type TrackerCategory = 'household' | 'personal' | 'pet';

interface TrackerInput {
	familyId: string;
	name: string;
	display: string;
	interval: string;
	intervalUnit: IntervalUnit;
	category: TrackerCategory;
	kind: Kind;
	actionLabel: string;
	pinned: boolean;
	show: boolean;
	icon: TrackerIcon;
	startDate?: string;
	cost?: number | undefined;
}

interface TrackerDB extends TrackerInput {
	id: string;
	interval: number;
	created_at: string;
	updated_at: string;
	familyId: string;
	familyName: string;
	isOwner: boolean;
	startDate?: string;
	cost?: number;
}

interface Family {
	id: string;
	name: string;
	isOwner: boolean;
	created: string;
	updated: string;
	owner: FamilyOwnerMember;
	members: FamilyOwnerMember[];
}

interface FamilyOwnerMember {
	id: string;
	email: string;
	name: string | null;
	created_at: string;
	updated_at: string;
}

type InviteStatus = 'pending' | 'completed' | null;

interface InviteDB {
	id: string;
	family: string;
	code: string;
	familyNameSnapshot: string;
	ownerEmailSnapshot: string;
	status?: InviteStatus;
	created: string;
	updated: string;
	expand?: { family?: FamilyDB };
}

type NotificationLevel = 'ok' | 'due' | 'overdue';

interface NotificationStatus {
	show: boolean;
	level: NotificationLevel | null;
	next?: string;
	label?: string | undefined;
	href?: string | undefined;
}

type ButtonState = 'default' | 'loading' | 'success' | 'error';

// Types/Interfaces for TrackerPage component abstraction

interface TrackerPageOptions {
	tracker: TrackerDB | undefined;
	labels: {
		pageTitle: string | undefined;
		ctaButtonText: string | undefined;
		noun: string | undefined;
	};
}

type Collections = EntryDB | VacationDB;

type CalendarRecord = Collections;

interface ActionCardOptions {
	title: string | undefined;
	size?: 'compact' | 'default' | 'list';
	tracker: TrackerDB | TrackerColored;
	entries: EntryDB[] | undefined;
	icon: Component;
	route: string;
	lastChild?: boolean;
	button: {
		status?: ButtonState;
		text: string | undefined;
	};
	streak?: number;
}
