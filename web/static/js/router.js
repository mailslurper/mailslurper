import * as auth from "./auth.js";
import { clear } from "./util/dom.js";
import * as inbox from "./views/inbox.js";
import * as settings from "./views/settings.js";
import * as login from "./views/login.js";
import * as savedSearches from "./views/savedSearches.js";

const routes = [
	{
		match: (h) => h.startsWith("#/mail/"),
		view: inbox,
		parse: (h) => {
			const [idPart, queryString] = h.slice("#/mail/".length).split("?");
			return { selectedId: decodeURIComponent(idPart), queryString: queryString || "" };
		},
	},
	{ match: (h) => h.startsWith("#/settings"), view: settings, parse: () => ({}) },
	{
		match: (h) => h.startsWith("#/login"),
		view: login,
		parse: (h) => ({ next: new URLSearchParams(h.split("?")[1] || "").get("next") || "#/inbox" }),
	},
	{ match: (h) => h.startsWith("#/saved-searches"), view: savedSearches, parse: () => ({}) },
	{
		match: (h) => h.startsWith("#/inbox"),
		view: inbox,
		parse: (h) => ({ queryString: h.split("?")[1] || "" }),
	},
];

let container = null;
let currentView = null;

export function init(rootEl) {
	container = rootEl;
	window.addEventListener("hashchange", render);
	render();
}

export function navigate(hash) {
	location.hash = hash;
}

function resolveRoute(hash) {
	for (const route of routes) {
		if (route.match(hash)) return { view: route.view, params: route.parse(hash) };
	}
	return { view: inbox, params: {} };
}

async function render() {
	const hash = location.hash || "#/inbox";

	if (hash.startsWith("#/admin")) {
		location.replace("#/settings");
		return;
	}

	if (!hash.startsWith("#/login")) {
		const needsAuth = await auth.authRequired().catch(() => false);
		if (needsAuth && !auth.isLoggedIn()) {
			location.hash = `#/login?next=${encodeURIComponent(hash)}`;
			return;
		}
	}

	const { view, params } = resolveRoute(hash);

	if (currentView && typeof currentView.unmount === "function") {
		currentView.unmount();
	}

	clear(container);
	currentView = view;
	await view.mount(container, params);
}
