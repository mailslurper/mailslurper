import assert from "node:assert/strict";
import test from "node:test";

function installBrowserMocks() {
	const backing = new Map();
	globalThis.localStorage = {
		getItem: (key) => (backing.has(key) ? backing.get(key) : null),
		setItem: (key, value) => { backing.set(key, String(value)); },
		removeItem: (key) => { backing.delete(key); },
	};
	globalThis.window = {
		matchMedia: () => ({ matches: false, addEventListener() {} }),
	};
	return backing;
}

let store;
let auth;

test.before(async () => {
	installBrowserMocks();
	store = await import("./store.js");
	auth = await import("./auth.js");
});

test("authHeader includes bearer when token is set", () => {
	store.setToken("jwt-token");
	assert.deepEqual(auth.authHeader(), { Authorization: "Bearer jwt-token" });
	store.clearToken();
});

test("authHeader is empty without token", () => {
	store.clearToken();
	assert.deepEqual(auth.authHeader(), {});
});

test("isLoggedIn reflects stored token", () => {
	store.clearToken();
	assert.equal(auth.isLoggedIn(), false);
	store.setToken("jwt-token");
	assert.equal(auth.isLoggedIn(), true);
	store.clearToken();
	assert.equal(auth.isLoggedIn(), false);
});
