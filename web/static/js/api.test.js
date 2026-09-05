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
let api;

test.before(async () => {
	installBrowserMocks();
	store = await import("./store.js");
	api = await import("./api.js");
});

test("getMailCount sends Authorization header when token exists", async (t) => {
	const originalFetch = globalThis.fetch;
	t.after(() => { globalThis.fetch = originalFetch; });
	globalThis.fetch = async (url, options) => {
		assert.equal(url, "/api/mail/count");
		assert.equal(options.headers.Authorization, "Bearer secret-token");
		return new Response(JSON.stringify({ count: 3 }), {
			status: 200,
			headers: { "Content-Type": "application/json" },
		});
	};

	store.setToken("secret-token");
	const result = await api.getMailCount();
	assert.equal(result.count, 3);
	store.clearToken();
});

test("getMailCount clears token and throws on 401", async (t) => {
	const originalFetch = globalThis.fetch;
	t.after(() => { globalThis.fetch = originalFetch; });
	globalThis.fetch = async () => new Response(JSON.stringify({ error: "unauthorized" }), {
		status: 401,
		headers: { "Content-Type": "application/json" },
	});

	store.setToken("expired-token");
	await assert.rejects(() => api.getMailCount(), api.UnauthorizedError);
	assert.equal(store.getToken(), "");
});
