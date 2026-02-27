import type { Task, TaskResponse } from '$lib/types/todo';
import { writable } from 'svelte/store';
import { GetTasks, DeleteCompleted as DeleteCompletedApi } from '$lib/api/todo';

export const tasks = writable<Task[]>([]);
export const deletedTasks = writable<Task[]>([]);

export async function loadTasks(): Promise<void> {
	try {
		const res = await GetTasks();
		const list: TaskResponse[] = Array.isArray(res)
			? (res as TaskResponse[])
			: res
				? ([res] as TaskResponse[])
				: [];
		tasks.set(list.map(mapResponseToTask));
	} catch (err) {
		console.error((err as Error).message);
	}
}

export function addTaskToStore(task: Task): void {
	tasks.update((current) => [task, ...current]);
}

export function mapResponseToTask(resp: TaskResponse): Task {
	return {
		Id: resp.id,
		Title: resp.title,
		Description: resp.description,
		DueDate: resp.dueDate,
		IsDone: resp.status === 'DONE',
		Priority: (resp.priority as Task['Priority']) || 'MEDIUM',
		Category: resp.category || '',
		DueWarning: resp.dueWarning || '',
		ChildTasks: (resp.childTasks || []).map(mapResponseToTask),
		ParentId: resp.parentId
	};
}

export function markDoneInStore(id: string): void {
	tasks.update((current) => current.map((t) => (t.Id === id ? { ...t, IsDone: true } : t)));
}

export function updateInStore(updated: Task): void {
	tasks.update((current) => current.map((t) => (t.Id === updated.Id ? updated : t)));
}

export function moveToDeleted(id: string): void {
	let removed: Task | undefined;
	tasks.update((current) => {
		const idx = current.findIndex((t) => t.Id === id);
		if (idx === -1) return current;
		removed = current[idx];
		const next = current.slice();
		next.splice(idx, 1);
		return next;
	});
	if (removed) {
		deletedTasks.update((cur) => [removed as Task, ...cur]);
	}
}

export async function clearCompleted(): Promise<void> {
	await DeleteCompletedApi();
	tasks.update((current) => current.filter((t) => !t.IsDone));
}
