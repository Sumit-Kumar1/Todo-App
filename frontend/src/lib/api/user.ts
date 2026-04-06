import { apiFetch } from './client';

async function Login(email: string, password: string) {
	const reqData = { email: email, password: password };

	try {
		const res = await apiFetch('/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reqData),
			skipRedirect: true
		});

		if (res === undefined) {
			throw new Error(`undefined response from POST/login`)
		}

		if (res.status == 404) {
			throw new Error(`user doesn't exist, please register first!`);
		}

		if (!res.ok) {
			throw new Error(`Login error! status: ${res.status}`);
		}

		return res;
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

async function Register(email: string, password: string) {
	const reqData = { email: email, password: password };

	try {
		const res = await apiFetch('/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(reqData),
			skipRedirect: true
		});

		if (res === undefined) {
			throw new Error(`undefined response from POST/register`)
		}

		if (res.status === 409) {
			throw new Error(`user already, please use login!`);
		}

		if (!res.ok) {
			throw new Error(`register error! status: ${res.status}`);
		}

		return res;
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

async function Logout() {
	try {
		const res = await apiFetch('/logout', {
			method: 'POST'
		});

		if (res === undefined) {
			throw new Error(`undefined response from POST/logout`)
		}

		if (!res.ok) {
			const err = await res.clone().json();
			throw new Error(`logout error! status: ${res.status}, error: ${err.message}`);
		}

		return res;
	} catch (error) {
		throw new Error((error as Error).message);
	}
}

export { Login, Register, Logout };
