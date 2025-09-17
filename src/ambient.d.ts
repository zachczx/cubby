interface TowelDB {
	collectionId: string;
	collectionName: string;
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
	collectionName: string;
	id: string;
	time: string;
	created: string;
	updated: string;
	user: string;
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
	defaultSprayInterval: number;
	created: string;
	updated: string;
}
