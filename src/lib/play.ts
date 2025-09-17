import { Capacitor } from '@capacitor/core';
export async function play(count: number, character: Characters = 'robot') {
	if (Capacitor.getPlatform() === 'android') {
		return;
	}

	// let minuteAudio = '';

	// if (count === 1) {
	// 	minuteAudio = `/${character}/minute.mp3`;
	// } else {
	// 	minuteAudio = `/${character}/minutes.mp3`;
	// }

	// const audioFiles = [`/${character}/${count}.mp3`, minuteAudio, `/${character}/remaining.mp3`];

	// const audioElements = audioFiles.map((file) => {
	// 	const audio = new Audio(file);
	// 	audio.preload = 'auto';

	// 	return audio;
	// });

	// const playSequentially = async (audioElements: HTMLAudioElement[]) => {
	// 	for (const audio of audioElements) {
	// 		await new Promise((resolve) => {
	// 			audio.addEventListener('ended', resolve, { once: true });
	// 			audio.addEventListener('error', resolve, { once: true });
	// 			audio.playbackRate = 1.2;
	// 			audio.play();
	// 		});
	// 	}
	// };

	// await playSequentially(audioElements);

	const audio = new Audio(`/${character}/${count}_remaining.mp3`);
	audio.play();
}
