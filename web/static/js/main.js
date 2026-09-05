import * as auth from "./auth.js";
import * as store from "./store.js";
import * as router from "./router.js";

function updateNavActive() {
	const hash = location.hash || "#/inbox";
	document.querySelectorAll("[data-nav]").forEach((a) => {
		const key = a.getAttribute("data-nav");
		const active =
			(key === "inbox" && (hash.startsWith("#/inbox") || hash.startsWith("#/mail/") || hash === "#" || hash === "")) ||
			(key === "saved-searches" && hash.startsWith("#/saved-searches")) ||
			(key === "settings" && (hash.startsWith("#/settings") || hash.startsWith("#/admin")));
		a.classList.toggle("is-active", active);
	});
}

async function updateNav() {
	const logoutLink = document.getElementById("nav-logout");
	const needsAuth = await auth.authRequired().catch(() => false);
	logoutLink.style.display = needsAuth && auth.isLoggedIn() ? "" : "none";
	updateNavActive();
}

async function boot() {
	store.applyThemeEarly();

	const logoutLink = document.getElementById("nav-logout");
	logoutLink.addEventListener("click", async (e) => {
		e.preventDefault();
		await auth.logout();
		location.hash = "#/login";
		updateNav();
	});

	store.on("token", updateNav);
	window.addEventListener("hashchange", updateNavActive);
	await updateNav();

	router.init(document.getElementById("app-root"));
}

boot();
