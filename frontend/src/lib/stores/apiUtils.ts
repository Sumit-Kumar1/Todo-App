import { notifications, type ToastNotification } from '$lib/stores/notifications';

export function handleApiError(error: unknown, context: string) {
    let message = 'An unknown error occurred';
    let statusCode = 500;

    if (error instanceof Error) {
        try {
            // Try to parse error message that might contain status code
            const errorParts = error.message.match(/HTTP error! status: (\d+), message: (.+)/);
            if (errorParts) {
                statusCode = parseInt(errorParts[1]);
                message = errorParts[2];
            } else {
                message = error.message;
            }
        } catch {
            message = error.message;
        }
    }

    // Show appropriate notification based on status code
    showNotification(message, statusCode, context);
    throw new Error(`[${context}] ${message} (Status: ${statusCode})`);
}

function showNotification(message: string, statusCode: number, context: string) {
    let type: 'error' | 'warning' | 'info' = 'error';
    let fullMessage = `${context}: ${message}`;

    switch (statusCode) {
        case 400:
            type = 'warning';
            fullMessage = `${context}: Bad request - ${message}`;
            break;
        case 401:
            type = 'warning';
            fullMessage = `${context}: Unauthorized - Please login again`;
            break;
        case 403:
            type = 'warning';
            fullMessage = `${context}: Forbidden - You don't have permission`;
            break;
        case 404:
            type = 'warning';
            fullMessage = `${context}: Not found - ${message}`;
            break;
        case 409:
            type = 'warning';
            fullMessage = `${context}: Conflict - ${message}`;
            break;
        case 500:
        default:
            type = 'error';
            fullMessage = `${context}: Server error - ${message}`;
            break;
    }

    notifications.push(fullMessage, type, 5000)
}