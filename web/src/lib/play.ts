import { Capacitor } from '@capacitor/core';
export async function play(seconds: number, character: Characters = 'robot') {
	if (Capacitor.getPlatform() === 'android') {
		return;
	}

	if (seconds % 60 === 0) {
		const minsRemaining = seconds / 60;

		const audio = new Audio(`/${character}/minutes/${minsRemaining}_remaining.mp3`);
		audio.play();
		return;
	}

	if (seconds <= 10 && seconds > 0) {
		const audio = new Audio(`/${character}/seconds/${seconds}_seconds.mp3`);
		audio.play();
		return;
	}
}
