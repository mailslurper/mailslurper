import assert from "node:assert/strict";
import test from "node:test";
import { stateToQueryString, queryStringToState } from "./query.js";

test("stateToQueryString omits empty fields and default page", () => {
	const qs = stateToQueryString({ q: "hello", from: "", to: "", start: "", end: "", sort: "subject", dir: "asc", page: 1 });
	assert.equal(qs, "q=hello&sort=subject&dir=asc");
});

test("stateToQueryString includes page when not 1", () => {
	const qs = stateToQueryString({ page: 3 });
	assert.equal(qs, "page=3");
});

test("queryStringToState fills defaults", () => {
	const state = queryStringToState("");
	assert.equal(state.sort, "date");
	assert.equal(state.dir, "desc");
	assert.equal(state.page, 1);
});

test("queryStringToState round-trips a full query", () => {
	const original = { q: "invoice", from: "a@b.com", to: "c@d.com", start: "2026-01-01", end: "2026-02-01", sort: "from", dir: "asc", page: 2 };
	const roundTripped = queryStringToState(stateToQueryString(original));
	assert.deepEqual(roundTripped, original);
});

test("stateToQueryString handles empty state", () => {
	assert.equal(stateToQueryString({}), "");
});

test("queryStringToState preserves encoded values", () => {
	const state = queryStringToState("q=hello%20world&from=a%40b.com");
	assert.equal(state.q, "hello world");
	assert.equal(state.from, "a@b.com");
});
