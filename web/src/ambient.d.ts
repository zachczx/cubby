type Characters = 'frankenstein' | 'robot' | 'furnando';

interface UserDB {
	id: string;
	email: string;
	name: string;
	soundOn: boolean;
	taskLookaheadDays: number;
	preferredCharacter: 'default' | Characters;
	createdAt: string;
	updatedAt: string;
}

interface VacationDB {
	id: string;
	familyId: string;
	createdBy: string;
	startDateTime: string;
	endDateTime: string;
	label?: string | undefined;
	createdAt: string;
	updatedAt: string;
}

type IntervalUnit = 'day' | 'month' | 'year';

interface EntryDB {
	id: string;
	trackerId: string;
	interval: number;
	intervalUnit: IntervalUnit;
	performedBy: string;
	performedAt: string;
	remark: string;
	createdAt: string;
	updatedAt: string;
}

interface EntryWithTracker extends EntryDB {
	tracker: TrackerDB;
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
	isMuted: boolean;
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
	createdAt: string;
	updatedAt: string;
}

interface InviteDB {
	id: string;
	familyId: string;
	inviteeId: string;
	status: InviteStatus;
	familyName: string;
	createdAt: string;
	updatedAt: string;
}

type InviteStatus = 'pending' | 'accepted' | 'declined';

interface InviteRequest {
	familyId: string;
	inviteeEmail: string;
}

interface InviteActionRequest {
	inviteId: string;
	status: InviteStatus;
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
