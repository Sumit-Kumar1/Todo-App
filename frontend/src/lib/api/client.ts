import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { notifications } from '$lib/stores/notifications';

const baseURL = import.meta.env.VITE_BACKEND_URL || '/api';

interface ApiFetchOptions extends RequestInit {
  skipRedirect?: boolean;
}

export async function apiFetch(path: string, options: ApiFetchOptions = {}) {
  const url = baseURL + path;
  // Default to including credentials
  const fetchOptions: RequestInit = {
    credentials: 'include',
    ...options
  };

  const res = await fetch(url, fetchOptions);

  if (res.status === 401 || res.status === 403) {
    if (!options.skipRedirect && browser) {
      // Try to parse error message if possible
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

  return res;
}
