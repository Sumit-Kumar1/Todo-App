import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { notifications } from '$lib/stores/notifications';

const baseURL = import.meta.env.TODO_BACKEND_URL || '/api';

interface ApiFetchOptions extends RequestInit {
	skipRedirect?: boolean;
}

export async function apiFetch(path: string, options: ApiFetchOptions = {}) {
	const url = baseURL + path;
	const fetchOptions: RequestInit = {
		credentials: 'include',
		...options
	};

	const res = await fetch(url, fetchOptions);

	switch (res.status) {
		case 200:
		case 201:
			return res
		case 401:
		case 403:
			if (!options.skipRedirect && browser) {
				let errorMessage = 'Session expired or unauthorized. Please login again.';
				try {
					const errorBody = await res.clone().json();
					if (errorBody.error) {
						errorMessage = errorBody.error;
					}
				} catch {
					// ignore parsing error
				}

				notifications.error(errorMessage);
				await goto('/');

				// Throwing error after redirect to ensure calling code doesn't proceed
				const error = await res.json().catch(() => ({ error: 'Unauthorized' }));
				throw new Error(
					`apiFetch error! status: ${res.status}, message: ${error.error || 'Unauthorized'}`
				);
			}
	}
}
