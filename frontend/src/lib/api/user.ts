import { goto } from "$app/navigation";

const baseURL = 'http://localhost:9003'

async function Login(email: string, password: string) {
  const url = baseURL + '/login'
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
      throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
    }

    return await res.json();
  } catch (error) {
    console.error('Error posting data:', (error as Error).message);
    throw new Error((error as Error).message);
  }
}

async function Register(email: string, password: string) {
  const url = baseURL + '/register'
  const reqData = { email: email, password: password };

  try {
    const res = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(reqData),
    });

    if (!res.ok) {
      const error = await res.json();
      throw new Error(`HTTP error! status: ${res.status}, message: ${error.error}`);
    }

    return await res.json();
  } catch (error) {
    console.error('Error posting data:', (error as Error).message);
    throw new Error((error as Error).message);
  }
}

async function Logout() {
  const url = baseURL + '/logout'

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
		console.error('Error during logout:', (error as Error).message);
		throw new Error((error as Error).message);
	}
}

export { Login, Register, Logout };
