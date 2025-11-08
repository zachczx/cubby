interface TowelDB {
	collectionId: string;
	collectionName: 'towel';
	id: string;
	time: string;
	created: string;
	updated: string;
	user: string;
}

interface TowelRecord extends TowelDB {
	gap: number;
}

interface SprayDB {
	collectionId: string;
	collectionName: 'spray';
	id: string;
	time: string;
	created: string;
	updated: string;
	user: string;
	daysToNext: number;
}

interface SprayRecord extends SprayDB {
	gap: number;
}

interface GummyDB {
	collectionId: string;
	collectionName: 'gummy';
	id: string;
	time: string;
	created: string;
	updated: string;
	user: string;
	daysToNext: number;
}

interface GummyRecord extends GummyDB {
	gap: number;
}

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
	mute: boolean;
	sprayInterval: number;
	gummyInterval: number;
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

interface NotificationStatus {
	show: boolean;
	level: 'ok' | 'due' | 'overdue';
}

type ButtonState = 'default' | 'loading' | 'success';

interface StatusDescriptions {
	green: string;
	yellow: string;
	orange: string;
	red: string;
}

type Query =
	| CreateQueryResult<SprayDB[], Error>
	| CreateQueryResult<TowelDB[], Error>
	| CreateQueryResult<GummyDB[], Error>;
