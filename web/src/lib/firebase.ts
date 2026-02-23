import { initializeApp } from 'firebase/app';
import { getMessaging } from 'firebase/messaging';
import { getAnalytics } from 'firebase/analytics';

const firebaseConfig = {
	apiKey: 'AIzaSyBJVeYpH1mvHBNdYhslq9Zan3p_hT0deyc',
	authDomain: 'cubbydotdev.firebaseapp.com',
	projectId: 'cubbydotdev',
	storageBucket: 'cubbydotdev.firebasestorage.app',
	messagingSenderId: '466629353659',
	appId: '1:466629353659:web:4257b35546b449a8eb099e',
	measurementId: 'G-PETLZ8V716'
};

const app = initializeApp(firebaseConfig);
export const getFirebaseMessaging = () => getMessaging(app);
export const analytics = () => getAnalytics(app);

export const firebaseIconSrc = '/icons/icon-192x192.webp';
