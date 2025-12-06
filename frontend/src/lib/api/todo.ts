import { handleApiError } from '$lib/stores/apiUtils';
import type { CreateTaskRequest } from '$lib/types/todo';

// Use environment variable or default to localhost for development
const baseURL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:9003';

export async function AddTask(taskReq: CreateTaskRequest) {
	const url = baseURL + '/task';

	try {
		const res = await fetch(url, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(taskReq),
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		handleApiError(error, 'Add Task')
		throw new Error((error as Error).message);
	}
}

export async function GetTasks() {
	const url = baseURL + '/tasks';

	try {
		const res = await fetch(url, {
			method: 'GET',
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		handleApiError(error, 'Fetch tasks')
		throw new Error((error as Error).message);
	}
}

export async function DelTask(id: string) {
	const url = baseURL + '/tasks/' + id;

	try {
		const res = await fetch(url, {
			method: 'DELETE',
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return;
	} catch (error) {
		handleApiError(error, 'Delete task')
		throw new Error((error as Error).message);
	}
}

export async function MarkDone(id: string) {
	const url = baseURL + '/tasks/' + id + '/done';

	try {
		const res = await fetch(url, {
			method: 'PATCH',
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		handleApiError(error, 'Mark done task')
		throw new Error((error as Error).message);
	}
}

export async function UpdateTask(id: string, payload: Record<string, unknown>) {
	const url = baseURL + '/tasks/' + id;

	try {
		const res = await fetch(url, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload),
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		handleApiError(error, 'Updating Task')
		throw new Error((error as Error).message);
	}
}
