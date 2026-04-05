type Characters = 'frankenstein' | 'robot' | 'furnando';

interface UserDB {
	id: string;
	email: string;
	name: string;
	soundModeQuick: 'off' | 'end' | 'full';
	soundModeProfile: 'off' | 'end' | 'full';
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

// Gym types

interface WorkoutDB {
	id: string;
	userId: string;
	startTime: string;
	notes: string | null;
	createdAt: string;
	updatedAt: string;
	sets: SetDB[];
}

interface SetDB {
	id: string;
	workoutId: string;
	exerciseId: string;
	weightKg: number | null;
	reps: number | null;
	setType: string;
	isCompleted: boolean;
	position: number;
	createdAt: string;
	updatedAt: string;
}

interface SetInput {
	exerciseId: string;
	weightKg: number | null;
	reps: number | null;
	setType: string;
}

interface FavouriteExercisesDB {
	exerciseIds: string[];
}

interface ExerciseCount {
	exerciseId: string;
	count: number;
}

interface WorkoutSummaryDB {
	totalWorkoutsThisMonth: number;
	totalVolumeThisMonth: number;
	totalSetsThisMonth: number;
	topExercises: ExerciseCount[];
	failureExercises: ExerciseCount[];
}

interface WorkoutCalendarEntryDB {
	workoutId: string;
	startTime: string;
	exerciseCount: number;
	setCount: number;
	exerciseIds: string[];
}

interface ExerciseFailureStatsDB {
	exerciseId: string;
	failureCount: number;
	lastFailureDate: string | null;
	totalSets: number;
}

interface MuscleStatsDB {
	muscle: string;
	failureCount: number;
	lastFailureDate: string | null;
	totalSets: number;
	failureRate: number;
}

interface ExerciseSetStatsDB {
	date: string;
	weightKg: number | null;
	reps: number | null;
	setType: string;
}

interface UserExerciseDB {
	exerciseId: string;
	setCount: number;
}

// Gym routine types

interface RoutineDB {
	id: string;
	userId: string;
	name: string;
	position: number;
	createdAt: string;
	updatedAt: string;
	exercises: RoutineExerciseDB[];
}

interface RoutineExerciseDB {
	id: string;
	routineId: string;
	exerciseId: string;
	sets: number;
	position: number;
	createdAt: string;
	updatedAt: string;
}

// Timer profile types

interface TimerSegmentDef {
	label: string;
	defaultSeconds: number;
}

interface TimerProfileDB {
	id: string;
	userId: string;
	name: string;
	segments: TimerSegmentDef[];
	isDefault: boolean;
	createdAt: string;
	updatedAt: string;
}

interface TimerProfileInput {
	name: string;
	segments: TimerSegmentDef[];
	isDefault: boolean;
}

// Market types

interface MarketPriceDB {
	id: string;
	loggedBy: string | null;
	itemName: string;
	category: string | null;
	country: string | null;
	store: string | null;
	unit: string | null;
	quantity: number | null;
	price: number;
	isPromo: boolean;
	remarks: string | null;
	createdAt: string;
	updatedAt: string;
}

interface MarketInsightDB {
	itemName: string;
	category: string | null;
	country: string | null;
	lowestPrice: number;
	lowestUnit: number | null;
	lowestStore: string | null;
	lowestDate: string;
	latestPrice: number;
	latestUnit: number | null;
	latestStore: string | null;
	latestDate: string;
}

interface MarketPriceInput {
	itemName: string;
	category: string | null;
	country: string | null;
	store: string | null;
	unit: string | null;
	quantity: number | null;
	price: number;
	isPromo: boolean;
	remarks: string | null;
	createdAt?: string;
	updatedAt?: string;
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
