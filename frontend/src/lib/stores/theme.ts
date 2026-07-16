import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

function getInitialTheme(): Theme {
	if (browser) {
		const stored = localStorage.getItem('theme');
		if (stored === 'light' || stored === 'dark') return stored;
		if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark';
	}
	return 'light';
}

export const theme = writable<Theme>(getInitialTheme());

if (browser) {
	theme.subscribe((value) => {
		document.documentElement.setAttribute('data-theme', value);
		localStorage.setItem('theme', value);
	});
}

export function toggleTheme(): void {
	theme.update((current) => (current === 'light' ? 'dark' : 'light'));
}
