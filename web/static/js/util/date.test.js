import assert from "node:assert/strict";
import test from "node:test";
import { formatCountdown, formatDateTime } from "./date.js";

test("formatCountdown handles zero and negative", () => {
	assert.equal(formatCountdown(0), "0:00");
	assert.equal(formatCountdown(-500), "0:00");
});

test("formatCountdown pads seconds", () => {
	assert.equal(formatCountdown(65000), "1:05");
});

test("formatDateTime returns empty string for empty input", () => {
	assert.equal(formatDateTime("", "iso"), "");
});

test("formatDateTime iso format round-trips a known date", () => {
	const iso = "2026-08-18T09:12:03.000Z";
	assert.equal(formatDateTime(iso, "iso"), iso);
});
