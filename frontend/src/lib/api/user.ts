import { handleApiError } from '$lib/stores/apiUtils';

const baseURL ='/api';

async function Login(email: string, password: string) {
	const url = baseURL + '/login';
	const reqData = { email: email, password: password };

	try {
		const res = await fetch(url, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reqData),
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error || 'Login failed'}`);
		}

		return await res;
	} catch (error) {
		handleApiError(error, 'Login')
		throw new Error((error as Error).message);
	}
}

async function Register(email: string, password: string) {
	const url = baseURL + '/register';
	const reqData = { email: email, password: password };

	try {
		const res = await fetch(url, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reqData)
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res;
	} catch (error) {
		handleApiError(error, 'Registeration')
		throw new Error((error as Error).message);
	}
}

async function Logout() {
	const url = baseURL + '/logout';

	try {
		const res = await fetch(url, {
			method: 'POST',
			credentials: 'include'
		});

		if (!res.ok) {
			const error = await res.json();
			throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
		}

		return await res.json();
	} catch (error) {
		handleApiError(error, 'Logout')
		throw new Error((error as Error).message);
	}
}

export { Login, Register, Logout };
	