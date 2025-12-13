import type { CreateTaskRequest } from '$lib/types/todo';
import { apiFetch } from './client';

export async function AddTask(taskReq: CreateTaskRequest) {
	try {
		const res = await apiFetch('/task', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(taskReq)
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

export async function GetTasks() {
	try {
		const res = await apiFetch('/tasks', {
			method: 'GET'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

export async function DelTask(id: string) {
	try {
		const res = await apiFetch('/tasks/' + id, {
			method: 'DELETE'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return;
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

export async function MarkDone(id: string) {
	try {
		const res = await apiFetch('/tasks/' + id + '/done', {
			method: 'PATCH'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

export async function UpdateTask(id: string, payload: Record<string, unknown>) {
	try {
		const res = await apiFetch('/tasks/' + id, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload)
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		throw new Error((error as Error).message);
	}
}
