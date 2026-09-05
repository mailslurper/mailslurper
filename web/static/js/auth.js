import * as store from "./store.js";

let cachedSettings = null;

// fetchSettings hits the one unauthenticated endpoint the SPA needs at boot
// to know whether a login gate is required at all.
export async function fetchSettings() {
	if (cachedSettings) return cachedSettings;
	const res = await fetch("/api/settings");
	if (!res.ok) throw new Error(`GET /api/settings: ${res.status}`);
	cachedSettings = await res.json();
	return cachedSettings;
}

export async function authRequired() {
	const settings = await fetchSettings();
	return settings.authenticationScheme === "basic";
}

export async function login(userName, password) {
	const res = await fetch("/api/login", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ userName, password }),
	});
	if (!res.ok) {
		const body = await res.json().catch(() => ({}));
		throw new Error(body.error || `login failed (${res.status})`);
	}
	const { token } = await res.json();
	store.setToken(token);
}

export async function logout() {
	try {
		await fetch("/api/logout", {
			method: "POST",
			headers: authHeader(),
		});
	} finally {
		store.clearToken();
	}
}

export function authHeader() {
	const token = store.getToken();
	return token ? { Authorization: `Bearer ${token}` } : {};
}

export function isLoggedIn() {
	return Boolean(store.getToken());
}
