import { writable } from 'svelte/store';

export type NotificationLevel = 'info' | 'success' | 'warning' | 'error';

export interface ToastNotification {
	id: string;
	message: string;
	level: NotificationLevel;
	durationMs: number;
}

function createId(): string {
	return Math.random().toString(36).slice(2) + Date.now().toString(36);
}

function createNotificationsStore() {
	const { subscribe, update } = writable<ToastNotification[]>([]);

	function remove(id: string) {
		update((items) => items.filter((n) => n.id !== id));
	}

	function push(message: string, level: NotificationLevel = 'info', durationMs = 3500) {
		const id = createId();
		const item: ToastNotification = { id, message, level, durationMs };
		update((items) => [item, ...items]);
		if (durationMs > 0) {
			setTimeout(() => remove(id), durationMs);
		}
		return id;
	}

	return {
		subscribe,
		push,
		remove,
		info: (message: string, durationMs?: number) => push(message, 'info', durationMs ?? 3000),
		success: (message: string, durationMs?: number) => push(message, 'success', durationMs ?? 2500),
		warning: (message: string, durationMs?: number) => push(message, 'warning', durationMs ?? 4000),
		error: (message: string, durationMs?: number) => push(message, 'error', durationMs ?? 5000)
	};
}

export const notifications = createNotificationsStore();
