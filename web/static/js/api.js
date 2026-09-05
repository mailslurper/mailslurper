import { authHeader } from "./auth.js";
import * as store from "./store.js";
import { stateToQueryString } from "./util/query.js";

class UnauthorizedError extends Error {}

async function request(path, options = {}) {
	const res = await fetch(path, {
		...options,
		headers: { ...authHeader(), ...(options.headers || {}) },
	});

	if (res.status === 401 || res.status === 403) {
		store.clearToken();
		throw new UnauthorizedError("unauthorized");
	}

	if (!res.ok) {
		const body = await res.json().catch(() => ({}));
		throw new Error(body.error || `request failed (${res.status})`);
	}

	return res;
}

export { UnauthorizedError };

export async function getMails(searchState) {
	const qs = stateToQueryString(searchState);
	const res = await request(`/api/mail${qs ? "?" + qs : ""}`);
	return res.json();
}

export async function getMailByID(id) {
	const res = await request(`/api/mail/${encodeURIComponent(id)}`);
	return res.json();
}

export function getMailMessageURL(id) {
	return `/api/mail/${encodeURIComponent(id)}/message`;
}

export function getMailMessageRawURL(id) {
	return `/api/mail/${encodeURIComponent(id)}/messageraw`;
}

export async function getMailCount() {
	const res = await request("/api/mail/count");
	return res.json();
}

export async function pruneMail(code) {
	const res = await request("/api/mail/prune", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ code }),
	});
	return res.json();
}

export async function getPruneOptions() {
	const res = await request("/api/prune-options");
	return res.json();
}

export async function getVersion() {
	const res = await request("/api/version");
	return res.json();
}
