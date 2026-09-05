// localStorage-backed settings, saved searches, and auth token, with a tiny
// pub/sub so views can react to changes (e.g. live theme switching).

const KEYS = {
	theme: "mylslurper.theme",
	dateFormat: "mylslurper.dateFormat",
	autoRefreshMinutes: "mylslurper.autoRefreshMinutes",
	savedSearches: "mylslurper.savedSearches",
	inboxListWidth: "mylslurper.inboxListWidth",
	inboxColumnWidths: "mylslurper.inboxColumnWidths",
	token: "mylslurper.jwt",
};

function migrateLegacyKeys() {
	for (const newKey of Object.values(KEYS)) {
		if (localStorage.getItem(newKey) !== null) continue;
		const oldKey = newKey.replace(/^mylslurper\./, "mailslurper.");
		if (oldKey === newKey) continue;
		const old = localStorage.getItem(oldKey);
		if (old === null) continue;
		localStorage.setItem(newKey, old);
		localStorage.removeItem(oldKey);
	}
}

migrateLegacyKeys();

const listeners = new Map(); // event name -> Set<fn>

function emit(event, value) {
	for (const fn of listeners.get(event) || []) fn(value);
}

export function on(event, fn) {
	if (!listeners.has(event)) listeners.set(event, new Set());
	listeners.get(event).add(fn);
	return () => listeners.get(event)?.delete(fn);
}

function prefersDark() {
	return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

function normalizeThemePref(value) {
	if (value === "dark" || value === "slate") return "dark";
	if (value === "system") return "system";
	return "light";
}

function resolveTheme(pref) {
	if (pref === "system") return prefersDark() ? "dark" : "light";
	return pref;
}

export function getTheme() {
	return normalizeThemePref(localStorage.getItem(KEYS.theme));
}

export function setTheme(theme) {
	const pref = normalizeThemePref(theme);
	localStorage.setItem(KEYS.theme, pref);
	document.documentElement.dataset.theme = resolveTheme(pref);
	emit("theme", pref);
}

export function getDateFormat() {
	return localStorage.getItem(KEYS.dateFormat) || "locale";
}

export function setDateFormat(format) {
	localStorage.setItem(KEYS.dateFormat, format);
	emit("dateFormat", format);
}

export function getAutoRefreshMinutes() {
	const raw = localStorage.getItem(KEYS.autoRefreshMinutes);
	return raw === null ? 1 : Number(raw);
}

export function setAutoRefreshMinutes(minutes) {
	localStorage.setItem(KEYS.autoRefreshMinutes, String(minutes));
	emit("autoRefreshMinutes", minutes);
}

export function getSavedSearches() {
	try {
		return JSON.parse(localStorage.getItem(KEYS.savedSearches) || "[]");
	} catch {
		return [];
	}
}

export function saveSavedSearch(search) {
	const all = getSavedSearches();
	all.push(search);
	localStorage.setItem(KEYS.savedSearches, JSON.stringify(all));
	emit("savedSearches", all);
}

export function deleteSavedSearch(name) {
	const all = getSavedSearches().filter((s) => s.name !== name);
	localStorage.setItem(KEYS.savedSearches, JSON.stringify(all));
	emit("savedSearches", all);
}

export function getInboxListWidth() {
	const n = Number(localStorage.getItem(KEYS.inboxListWidth));
	return Number.isFinite(n) && n > 0 ? n : 0;
}

export function setInboxListWidth(px) {
	if (!px) localStorage.removeItem(KEYS.inboxListWidth);
	else localStorage.setItem(KEYS.inboxListWidth, String(Math.round(px)));
}

export const DEFAULT_INBOX_COLUMNS = { date: 96, from: 110, att: 24 };

export function getInboxColumnWidths() {
	try {
		const raw = JSON.parse(localStorage.getItem(KEYS.inboxColumnWidths) || "null");
		if (!raw || typeof raw !== "object") return { ...DEFAULT_INBOX_COLUMNS };
		const num = (v, fallback) => {
			const n = Number(v);
			return Number.isFinite(n) && n > 0 ? n : fallback;
		};
		return {
			date: num(raw.date, DEFAULT_INBOX_COLUMNS.date),
			from: num(raw.from, DEFAULT_INBOX_COLUMNS.from),
			att: num(raw.att, DEFAULT_INBOX_COLUMNS.att),
		};
	} catch {
		return { ...DEFAULT_INBOX_COLUMNS };
	}
}

export function setInboxColumnWidths(widths) {
	if (!widths) {
		localStorage.removeItem(KEYS.inboxColumnWidths);
		return;
	}
	localStorage.setItem(KEYS.inboxColumnWidths, JSON.stringify({
		date: Math.round(widths.date),
		from: Math.round(widths.from),
		att: Math.round(widths.att),
	}));
}

export function getToken() {
	return localStorage.getItem(KEYS.token) || "";
}

export function setToken(token) {
	localStorage.setItem(KEYS.token, token);
	emit("token", token);
}

export function clearToken() {
	localStorage.removeItem(KEYS.token);
	emit("token", "");
}

// Applies the persisted theme synchronously; called inline in <head> so
// there's no flash of the wrong theme before main.js loads.
let themeMedia = null;

export function applyThemeEarly() {
	const pref = getTheme();
	document.documentElement.dataset.theme = resolveTheme(pref);
	const stored = localStorage.getItem(KEYS.theme);
	if (stored && stored !== pref) {
		localStorage.setItem(KEYS.theme, pref);
	}
	if (!themeMedia) {
		themeMedia = window.matchMedia("(prefers-color-scheme: dark)");
		themeMedia.addEventListener("change", () => {
			if (getTheme() === "system") {
				document.documentElement.dataset.theme = resolveTheme("system");
				emit("theme", "system");
			}
		});
	}
}
